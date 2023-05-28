package user

import (
	"github.com/Degoke/doc-verify/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSerializer struct {
	User UserModel
}

type UsersSerializer struct {
	Users []UserModel
}

type UserResponse struct {
	ID primitive.ObjectID `json:"id"`
	Email string `json:"email"`
}

type LoginSerializer struct {
	User UserModel
}

type LoginResponse struct {
	ID primitive.ObjectID `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (u *UserSerializer) Response() UserResponse {
	return UserResponse{
		ID: u.User.ID,
		Email: u.User.Email,
	}
}

func (u *UsersSerializer) Response() []UserResponse {
	var users []UserResponse
	for _, user := range u.Users {
		users = append(users, UserResponse{
			ID: user.ID,
			Email: user.Email,
		})
	}
	return users
}

func (u *LoginSerializer) Response() LoginResponse {
	return LoginResponse{
		ID: u.User.ID,
		Email: u.User.Email,
		Token: common.GenToken(u.User.ID),
	}
}