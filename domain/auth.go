package domain

import (
	"context"
	"errors"
	"log"
	"net/http"

	// "net/http"

	"github.com/elisalimli/go_graphql_template/graphql/models"
	"github.com/elisalimli/go_graphql_template/middleware"
	"github.com/elisalimli/go_graphql_template/validator"
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

	rtCookie := http.Cookie{
		Name:    "refresh_token",
		Path:    "/", // <--- add this line
		Value:   refreshToken.RefreshToken,
		Expires: refreshToken.ExpiredAt,
	}

	writer, _ := ctx.Value(middleware.HttpWriterKey).(http.ResponseWriter)
	http.SetCookie(writer, &rtCookie)

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
