package smsservice

import (
	"errors"
	"net/http"
	"rz/util"
	"rz/redis"
	"rz/restfulapi"
	"github.com/gin-gonic/gin"
)

func GetVerificationCode(c *gin.Context) {

	zone := c.Query("zone")
	mobile := c.Query("mobile")

	if zone == "" || mobile == "" {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, ERROR_MISSING_PHONE_STRING))
		return
	}

	phone := zone + mobile

	code := util.GetRandomCode(4)
	key := VERIFICATION_CODE_REDIS_PREFIX + phone
	_, err := redis.Instance().SetAndExpire(key, code, 5*60)
	if (err != nil) {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	// 傳送簡訊
	err = send(zone, mobile, code);
	if err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	c.JSON(http.StatusOK, restfulapi.SuccessMsg(SUCCESS_SMS_SEND_STRING))
}



func PreCheckVerificationCode(c *gin.Context) {

	code := c.Param("code")

	// 綁定輸入參數
	var input PreCheckVerificationInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, ERROR_MISSING_PARAMETER_STRING))
		return
	}

	// 驗證碼確認
	err := preCheck(input.Zone + input.Mobile, code)
	if (err != nil) {
		c.JSON(http.StatusBadRequest, restfulapi.Error(1, err.Error()))
		return
	}

	c.JSON(http.StatusOK, restfulapi.SuccessMsg(SUCCESS_SMS_VERIFICATION_STRING))

}

func Check(phone, code string) error {
	key := VERIFICATION_CODE_REDIS_PREFIX + phone
	if err := preCheck(phone, code); err != nil {
		return err
	}

	redis.Instance().Del(key)
	return nil
}

func send(zone, mobile, code string) error {

	var verifyCodeSms VerifyCodeSMS

	switch zone {
	case "0086":
		break
	default:
		verifyCodeSms = &NexmoVerifyCode{}
	}
	if err := verifyCodeSms.SendVerifyCode(zone, mobile, code); err != nil {
		return err
	}
	return nil
}

func preCheck(phone, code string) error {
	key := VERIFICATION_CODE_REDIS_PREFIX + phone
	v, err := redis.Instance().Get(key)
	if (err != nil) {
		return err
	}

	if (v == "") {
		return errors.New(ERROR_NONE_VERIFICATION_CODE_STRING)
	} else if (v != code) {
		return errors.New(ERROR_VERIFICATION_CODE_STRING)
	}

	return nil
}

func combinePhoneFormat(zone, mobile string) string {
	//去除前面00
	z := zone[2:]
	zeroIndex := 0
	for zeroIndex = 0; zeroIndex < len(mobile); zeroIndex++ {
		if mobile[zeroIndex] != '0' {
			break
		}
	}
	//取出mobile裡前面沒有0的部分
	p := mobile[zeroIndex:]
	return z + p
}