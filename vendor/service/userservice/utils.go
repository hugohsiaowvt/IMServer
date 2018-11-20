package userservice

import (
	"errors"
)

const (
	TABLE_USERS					string = "users"
	LOGIN_TOKEN_REDIS_PREFIX	string = "login_token:"
	IM_TOKEN_REDIS_PREFIX		string = "im_token:"
	PASSWORD_PREFIX				string = "hUGO%pTT"

	
	ERROR_PHONE_ISEXIST_STRING			string = "該電話已被申請"
	ERROR_MISSING_PARAMETER_STRING      string = "缺少參數"
	ERROR_MESSAGE_WRONG_INPUT_STRING    string = "輸入錯誤"
	ERROR_MESSAGE_WRONG_PASSWORD_STRING string = "密碼錯誤"

	SUCCESS_SMS_SEND_STRING				string = "發送成功"
	SUCCESS_SMS_VERIFICATION_STRING		string = "驗證成功"
)

type User struct {
	OpenId     	string	`json:"open_id"`
	Username   	string	`json:"username"`
	UserId   	string	`json:"user_id"`
	Zone       	string	`json:"zone"`
	Mobile     	string	`json:"mobile"`
	Avatar		string	`json:"avatar"`
	Sex        	int		`json:"sex"`
	Password   	string	`json:"password"`
	Email		string	`json:"email"`
	Area		string	`json:"area"`
	Status     	int		`json:"status"`
}

type GetTokenInput struct {
	OpenId	string `json:"open_id"`
}

type GetTokenRes struct {
	Token	string
}

type RegisterInput struct {
	Username   	string	`json:"username"`
	Zone       	string	`json:"zone"`
	Mobile     	string	`json:"mobile"`
	Avatar		string	`json:"avatar"`
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
	Username   	string	`json:"username"`
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