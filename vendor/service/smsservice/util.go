package smsservice

const (
	VERIFICATION_CODE_REDIS_PREFIX      string = "verification_code:"

	ERROR_MISSING_PARAMETER_STRING      string = "缺少參數"

	ERROR_MISSING_PHONE_STRING			string = "手機不能為空白"
	ERROR_NONE_VERIFICATION_CODE_STRING string = "無驗證碼"
	ERROR_VERIFICATION_CODE_STRING      string = "驗證碼錯誤"

	SUCCESS_SMS_SEND_STRING				string = "發送成功"
	SUCCESS_SMS_VERIFICATION_STRING		string = "驗證成功"
)

type PreCheckVerificationInput struct {
	Zone		string	`json:"zone"`
	Mobile		string	`json:"mobile"`
}