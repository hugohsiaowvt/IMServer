package userservice

import (
	"errors"
)

const (

	TABLE_USERS					string	= "users"
	LOGIN_TOKEN_REDIS_PREFIX	string	= "login_token:"
	IM_TOKEN_REDIS_PREFIX		string	= "im_token:"
	OPEN_ID_PREFIX				string	= "open_id:"
	OPEN_ID_KEY					string	= "f?KQk#KZqvu.BvAx"

	USER_STATUS_OK				int		= 1
	USER_STATUS_LOCKED			int		= -1
	
)

type User struct {
	OpenId     	string	`json:"open_id" gorm:"Column:open_id"`
	Username   	string	`json:"user_name" gorm:"Column:user_name"`
	Nickname   	string	`json:"nick_name" gorm:"Column:nick_name"`
	UserId   	string	`json:"user_id" gorm:"Column:user_id"`
	Zone       	string	`json:"zone" gorm:"Column:zone"`
	Mobile     	string	`json:"mobile" gorm:"Column:mobile"`
	Sex        	int		`json:"sex" gorm:"Column:sex"`
	Password   	string	`json:"password" gorm:"Column:password"`
	Email		string	`json:"email" gorm:"Column:email"`
	Area		string	`json:"area" gorm:"Column:area"`
	Status     	int		`json:"status" gorm:"Column:status"`
}

type GetTokenInput struct {
	OpenId	string `json:"open_id"`
}

type GetTokenRes struct {
	Token	string
}

type RegisterInput struct {
	Username   	string	`json:"user_name"`
	Zone       	string	`json:"zone"`
	Mobile     	string	`json:"mobile"`
	Password	string	`json:"password"`
}

type PreCheckVerificationInput struct {
	Zone		string	`json:"zone"`
	Mobile		string	`json:"mobile"`
}

type LoginInput struct {
	Zone		string	`json:"zone"`
	Mobile		string	`json:"mobile"`
	Password	string	`json:"password"`
}

type LoginResponse struct {
	Token       string
	IMToken		string
	User		User
}

type UpdatePasswordInput struct {
	OpenId		string	`json:"open_id"`
	OldPassword	string	`json:"old_password"`
	NewPassword	string	`json:"new_password"`
}

type UpdateUserInfoInput struct {
	Username   	string	`json:"user_name"`
	Sex        	int		`json:"sex"`
	Email      	string	`json:"email"`
}

func (self RegisterInput) ToModel() *User {
	model := &User{}
	model.Username = self.Username
	model.Mobile = self.Mobile
	model.Zone = self.Zone
	model.Password = self.Password
	return model
}

func (input RegisterInput) CheckInput() error {
	if input.Username == "" {
		return errors.New("名字不得為空！")
	}
	if input.Zone == "" {
		return errors.New("區碼不得為空！")
	}
	if input.Mobile == "" {
		return errors.New("手機不得為空！")
	}
	if input.Password == "" {
		return errors.New("密碼不得為空！")
	}
	return nil
}

func (input LoginInput) CheckInput() error {
	if input.Zone == "" {
		return errors.New("區碼不得為空！")
	}
	if input.Mobile == "" {
		return errors.New("手機不得為空！")
	}
	if input.Password == "" {
		return errors.New("密碼不得為空！")
	}
	return nil
}

func (input UpdatePasswordInput) CheckInput() error {
	if input.OpenId == "" {
		return errors.New("Open id不得為空！")
	}
	if input.NewPassword == "" {
		return errors.New("密碼不得為空！")
	}
	return nil
}