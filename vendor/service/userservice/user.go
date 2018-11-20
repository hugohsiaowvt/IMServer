package userservice

import (
	"fmt"
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"rz/mysql"
	"rz/util"
	"rz/redis"
	"rz/restfulapi"
	"service/smsservice"
	"github.com/pborman/uuid"
)

// PUBLIC METHOD

func Test(c *gin.Context) {
	//user, query := &User{}, &User{}
	//var count int
	//query.OpenId = "875c5080f1f7e5eb7cbd4c893290994f"
	//query.Password = "ss"
	v, err := redis.Instance().LRange("mylist", 0, -1)
	fmt.Println(v, err)

	c.JSON(http.StatusOK, restfulapi.Success(v))
}

func CreateUser(c *gin.Context) {

	res := LoginResponse{}
	var loginInput LoginInput

	code := c.Param("code")

	// 綁定輸入參數
	var input RegisterInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, ERROR_MISSING_PARAMETER_STRING))
		return
	}

	// 輸入參數確認
	if err := input.CheckInput(); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	// 驗證碼確認
	err := smsservice.Check(input.Zone + input.Mobile, code)
	if (err != nil) {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	// 將輸入參數轉成User
	user := input.ToModel()

	// 新增一個Open id
	user.OpenId = util.GenerateKey("open_id:")

	// 開啟Transactions
	tx := mysql.Instance().Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()

	// 新增用戶
	if err := InsertUsers(tx, user); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	// 提交
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, "提交處理出錯！"))
		return
	}

	// 登入
	loginInput.Zone = input.Zone
	loginInput.Mobile = input.Mobile
	u := User{}
	token, im_token, err := login(&u, loginInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	res.Token = token
	res.IMToken = im_token
	res.User = u

	c.JSON(http.StatusOK, restfulapi.Success(res))
}

func Login(c *gin.Context) {

	res := LoginResponse{}
	user := User{}

	// 綁定輸入參數
	var input LoginInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, "缺少參數"))
		return
	}

	// 輸入參數確認
	if err := input.CheckInput(); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	// 登入
	token, im_token, err := login(&user, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	res.Token = token
	res.IMToken = im_token
	res.User = user

	c.JSON(http.StatusOK, restfulapi.Success(res))
}

func SetPassword(c *gin.Context) {

	user := &User{}

	// 綁定輸入參數
	var input UpdatePasswordInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, "缺少參數"))
		return
	}

	// 輸入參數確認
	if err := input.CheckInput(); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	// 查詢用戶
	if err := QueryUserByOpenId(user, input.OpenId); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	// 檢查密碼
	if user.Password != "" {
		if user.Password != input.OldPassword {
			c.JSON(http.StatusBadRequest, restfulapi.Error(1, ERROR_MESSAGE_WRONG_PASSWORD_STRING))
			return
		}
	}

	// 更新密碼
	if err := UpdatePassword(input.OpenId, input.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	c.JSON(http.StatusOK, restfulapi.Success(input))
}

// PRIVATE METHOD

// 確認用戶存不存在
func QueryUserExistByMobile(zone, mobile string) (User, bool) {
	user := User{}
	var count int
	query := &User{
		Zone: zone,
		Mobile: mobile,
	}
	QueryUser(&user, query, &count)
	if count != 0 {
		return user, true
	}
	return user, false
}

func QueryUserExistByOpenId(openId string) (User, bool) {
	user := User{}
	var count int
	query := &User{
		OpenId: openId,
	}
	QueryUser(&user, query, &count)
	if count != 0 {
		return user, true
	}
	return user, false
}

//登入
func login(user *User, input LoginInput) (string, string, error) {

	// 查詢用戶
	var count int
	query := &User{}
	query.Zone = input.Zone
	query.Mobile = input.Mobile
	query.Password = input.Password
	if err := QueryUser(user, query, &count); err != nil {
		if (count != 0) {
			return "", "", err
		} else {
			return "", "", errors.New(ERROR_MESSAGE_WRONG_INPUT_STRING)
		}
	}

	// 製作登入token
	key := LOGIN_TOKEN_REDIS_PREFIX + user.OpenId
	token := uuid.New()
	if _, err := redis.Instance().Set(key, token); err != nil {
		return "", "", err
	}

	// 取得im token
	userId := user.OpenId
	if imToken, tokenError := getIMToken(userId); tokenError != nil {
		return "", "", tokenError
	} else {
		return token, imToken, nil
	}

}

func getIMToken(openId string) (string, error) {
	token := uuid.New()
	key := IM_TOKEN_REDIS_PREFIX + openId
	if _, err := redis.Instance().Set(key, token); err != nil {
		return "", err
	}
	return token, nil
}

func CheckIMToken(openId, token string) bool {
	key := IM_TOKEN_REDIS_PREFIX + openId
	if v, err := redis.Instance().Get(key); err != nil {
		return false
	} else {
		if v == token {
			return true
		} else {
			return false
		}
	}
}