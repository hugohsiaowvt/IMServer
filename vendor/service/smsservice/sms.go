package smsservice

import (

	"errors"

	"rz/util"
	"rz/redis"
	"rz/restfulapi"

	"github.com/gin-gonic/gin"

)

func GetVerificationCode(c *gin.Context) {

	zone := c.Query("zone")
	mobile := c.Query("mobile")

	if zone == "" || mobile == "" {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PHONE_CODE, struct{}{}, restfulapi.ERROR_MISSING_PHONE_MSG)
		return
	}

	phone := zone + mobile

	// 產生驗證碼
	code := util.GetRandomCode(6)
	key := VERIFICATION_CODE_REDIS_PREFIX + phone
	_, err := redis.Instance().SetAndExpire(key, code, 5*60)
	if (err != nil) {
		restfulapi.Response(c, restfulapi.ERROR_REDIS_ERROR_CODE, struct{}{}, err.Error())
		return
	}

	// 傳送簡訊(目前先不花錢發簡訊)
	//err = send(zone, mobile, code);
	//if err != nil {
	//	restfulapi.Response(c, restfulapi.ERROR_SEND_SMS_CODE, struct{}{}, restfulapi.ERROR_SEND_SMS_MSG)
	//	return
	//}
	//
	//restfulapi.Response(c, restfulapi.SUCCESS_CODE, struct{}{}, restfulapi.SUCCESS_SMS_SEND_MSG)
	restfulapi.Response(c, restfulapi.SUCCESS_CODE, struct{}{}, code)
}



func PreCheckVerificationCode(c *gin.Context) {

	vcode := c.Param("code")

	// 綁定輸入參數
	var input PreCheckVerificationInput
	if err := c.BindJSON(&input); err != nil {
		restfulapi.Response(c, restfulapi.ERROR_MISSING_PARAMETER_CODE, struct{}{}, restfulapi.ERROR_MISSING_PARAMETER_MSG)
		return
	}

	// 驗證碼確認
	code, err := preCheck(input.Zone + input.Mobile, vcode)
	if (err != nil) {
		restfulapi.Response(c, code, struct{}{}, err.Error())
		return
	}

	restfulapi.Response(c, restfulapi.SUCCESS_CODE, struct{}{}, restfulapi.SUCCESS_SMS_VERIFICATION_MSG)

}

func Check(phone, code string) (int32, error) {
	key := VERIFICATION_CODE_REDIS_PREFIX + phone
	if code, err := preCheck(phone, code); err != nil {
		return code, err
	}

	redis.Instance().Del(key)
	return 0, nil
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

func preCheck(phone, code string) (int32, error) {
	key := VERIFICATION_CODE_REDIS_PREFIX + phone
	v, err := redis.Instance().Get(key)
	if (err != nil) {
		return restfulapi.ERROR_REDIS_ERROR_CODE, err
	}

	if (v == "") {
		return restfulapi.ERROR_NONE_VERIFICATION_CODE_CODE, errors.New(restfulapi.ERROR_NONE_VERIFICATION_CODE_MSG)
	} else if (v != code) {
		return restfulapi.ERROR_VERIFICATION_CODE_CODE, errors.New(restfulapi.ERROR_VERIFICATION_CODE_MSG)
	}

	return restfulapi.SUCCESS_CODE, nil
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