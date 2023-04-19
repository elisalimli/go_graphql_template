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
)

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
		return nil, errors.New("something went wrong")
	}

	refreshToken, err := user.GenRefreshToken()
	if err != nil {
		return nil, errors.New("something went wrong")
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
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.HashPassword(input.Password)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return nil, errors.New("something went wrong")
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

	token, err := user.GenAccessToken()
	if err != nil {
		log.Printf("error while generating the token: %v", err)
		return nil, errors.New("something went wrong")
	}

	return &models.AuthResponse{
		AuthToken: token,
		User:      user,
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
	t, ok := userId.(string)
	user, err := d.UsersRepo.GetUserByID(t)
	fmt.Println(t, ok, err)
	if err != nil {
		return &models.AuthResponse{Ok: false, Errors: []*validator.FieldError{{Message: "User not found", Field: "general"}}}, nil
	}

	newRefreshToken, err := user.GenRefreshToken()
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	newAccessToken, err := user.GenAccessToken()
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	user.SaveRefreshToken(ctx, newRefreshToken)

	return &models.AuthResponse{Ok: true, AuthToken: newAccessToken}, nil
}
