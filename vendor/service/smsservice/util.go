package smsservice

const (
	
	VERIFICATION_CODE_REDIS_PREFIX      string = "verification_code:"

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