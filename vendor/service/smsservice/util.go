package smsservice

const (
	VERIFICATION_CODE_REDIS_PREFIX      string = "verification_code:"

	ERROR_MISSING_PARAMETER_STRING      string = "缺少參數"

	ERROR_MISSING_PHONE_STRING			string = "手機不能為空白"
	ERROR_NONE_VERIFICATION_CODE_STRING string = "無驗證碼"
	ERROR_VERIFICATION_CODE_STRING      string = "驗證碼錯誤"
	ERROR_SEND_STRING					string = "簡訊發送失敗！"

	SUCCESS_SMS_SEND_STRING				string = "發送成功"
	SUCCESS_SMS_VERIFICATION_STRING		string = "驗證成功"
)

type PreCheckVerificationInput struct {
	Zone		string	`json:"zone"`
	Mobile		string	`json:"mobile"`
}

type SMS interface {
	// 發送
	Send(phone, content string) error
	// 發送指定模板ID
	SendVerifyCodeForTemplate(phone, templateId string, datas interface{}) error
}

type VerifyCodeSMS interface {
	// 發送驗證碼
	SendVerifyCode(zone, mobile, code string) error
}