package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	// "net/http"

	myContext "github.com/elisalimli/go_graphql_template/context"
	"github.com/elisalimli/go_graphql_template/graphql/models"
	"github.com/elisalimli/go_graphql_template/validator"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"

	"github.com/twilio/twilio-go"
)

func envACCOUNTSID() string {
	println(godotenv.Unmarshal(".env"))
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("TWILIO_ACCOUNT_SID")
}

func envAUTHTOKEN() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("TWILIO_AUTH_TOKEN")
}

func envSERVICESID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("VERIFY_SERVICE_SID")
}

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: envACCOUNTSID(),
	Password: envAUTHTOKEN(),
})

func (d *Domain) Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
	user, err := d.UsersRepo.GetUserByEmail(input.Email)
	if err != nil {
		return NewFieldError(validator.FieldError{Message: ErrBadCredentials, Field: "general"}), nil
	}

	err = user.ComparePassword(input.Password)
	if err != nil {
		return NewFieldError(validator.FieldError{Message: ErrBadCredentials, Field: "general"}), nil
	}

	accessToken, err := user.GenAccessToken()
	if err != nil {
		return nil, errors.New(ErrSomethingWentWrong)
	}

	refreshToken, err := user.GenRefreshToken()
	if err != nil {
		return nil, errors.New(ErrSomethingWentWrong)
	}
	user.SaveRefreshToken(ctx, refreshToken)

	return &models.AuthResponse{
		AuthToken: accessToken,
		User:      user,
	}, nil
}

func (d *Domain) Register(ctx context.Context, input models.RegisterInput) (*models.AuthResponse, error) {
	_, err := d.UsersRepo.GetUserByEmail(input.Email)
	if err == nil {
		return NewFieldError(validator.FieldError{Message: "Email already in used", Field: "email"}), nil
	}

	_, err = d.UsersRepo.GetUserByUsername(input.Username)
	if err == nil {
		return NewFieldError(validator.FieldError{Message: "Username already in used", Field: "username"}), nil
	}

	user := &models.User{
		Username:    input.Username,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
	}

	err = user.HashPassword(input.Password)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return nil, errors.New(ErrSomethingWentWrong)
	}

	// TODO: create verification code

	tx := d.UsersRepo.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		log.Printf("error creating a user: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("error while commiting: %v", err)
		return nil, err
	}

	// expiredAt := time.Minute * 10 // 10 mins

	// err = d.UsersRepo.RedisClient.Set(ctx, user.PhoneNumber, user.ID, expiredAt).Err()
	// if err != nil {
	// 	return nil, errors.New(ErrSomethingWentWrong)
	// }

	// token, err := user.GenAccessToken()
	// if err != nil {
	// 	log.Printf("error while generating the token: %v", err)
	// 	return nil, errors.New(ErrSomethingWentWrong)
	// }

	return &models.AuthResponse{
		// AuthToken: token,
		User: user,
	}, nil
}

func (d *Domain) RefreshToken(ctx context.Context) (*models.AuthResponse, error) {
	refreshTokenCookie, ok := ctx.Value(myContext.CookieRefreshTokenKey).(*http.Cookie)

	if !ok {
		return &models.AuthResponse{Ok: false}, nil
	}

	// Verify that the refresh token is valid
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(refreshTokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return &models.AuthResponse{Ok: false}, nil
	}
	userId, ok := claims["jti"]

	if !ok {
		return &models.AuthResponse{Ok: false}, nil
	}

	user, err := d.UsersRepo.GetUserByID(userId.(string))

	if err != nil {
		return &models.AuthResponse{Ok: false, Errors: []*validator.FieldError{{Message: "User not found", Field: "general"}}}, nil
	}

	if !user.Verified {
		return &models.AuthResponse{Ok: false, Errors: []*validator.FieldError{{Message: "You need to verify your account.", Field: "general"}}}, nil
	}

	newRefreshToken, err := user.GenRefreshToken()
	if err != nil {
		return nil, errors.New(ErrSomethingWentWrong)
	}

	newAccessToken, err := user.GenAccessToken()
	if err != nil {
		return nil, errors.New(ErrSomethingWentWrong)
	}
	user.SaveRefreshToken(ctx, newRefreshToken)

	return &models.AuthResponse{Ok: true, AuthToken: newAccessToken}, nil
}

func (d *Domain) SendOtp(ctx context.Context, input models.SendOtpInput) (*models.FormResponse, error) {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(input.To)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(envSERVICESID(), params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
	}

	return &models.FormResponse{Ok: true}, nil

}

func (d *Domain) VerifyOtp(ctx context.Context, input models.VerifyOtpInput) (*models.FormResponse, error) {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(input.To)
	params.SetCode(input.Code)

	resp, err := client.VerifyV2.CreateVerificationCheck(envSERVICESID(), params)
	fmt.Println("code", resp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New(ErrSomethingWentWrong)

	} else if *resp.Status == "approved" {
		fmt.Println("Correct!")
		if err != nil {
			panic(err)
		}

		d.UsersRepo.DB.Model(&models.User{}).Where("phone_number = ?", input.To).Update("verified", true)

		return &models.FormResponse{Ok: true}, nil
	} else {
		fmt.Println("Incorrect!")
		return &models.FormResponse{Ok: false}, nil
	}
}
