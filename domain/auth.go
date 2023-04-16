package domain

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/elisalimli/go_graphql_template/graphql/models"
)

func (d *Domain) Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
	user, err := d.UsersRepo.GetUserByEmail(input.Email)
	if err != nil {
		return nil, ErrBadCredentials
	}

	err = user.ComparePassword(input.Password)
	if err != nil {
		return nil, ErrBadCredentials
	}

	token, err := user.GenToken()
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return &models.AuthResponse{
		AuthToken: token,
		User:      user,
	}, nil
}

func (d *Domain) Register(ctx context.Context, input models.RegisterInput) (*models.AuthResponse, error) {
	_, err := d.UsersRepo.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("email already in used")
	}

	_, err = d.UsersRepo.GetUserByUsername(input.Username)
	if err == nil {
		return nil, errors.New("username already in used")
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

	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	if err := tx.Create(user).Error; err != nil {
		log.Printf("error creating a user: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("error while commiting: %v", err)
		return nil, err
	}

	token, err := user.GenToken()
	if err != nil {
		log.Printf("error while generating the token: %v", err)
		return nil, errors.New("something went wrong")
	}

	return &models.AuthResponse{
		AuthToken: token,
		User:      user,
	}, nil
}
