package restfulapi

import(
	"net/http"

	"github.com/gin-gonic/gin"
)

const (

	SUCCESS_CODE						int32	= 0

	SUCCESS_SMS_SEND_MSG         		string = "簡訊發送成功"
	SUCCESS_SMS_VERIFICATION_MSG 		string = "簡訊驗證成功"

	ERROR_MISSING_PARAMETER_MSG  		string = "缺少參數"
	ERROR_MISSING_PARAMETER_CODE 		int32  = -1
	ERROR_UNKNOW_MSG					string = "未知錯誤"
	ERROR_UNKNOW_CODE					int32  = -2

	ERROR_MYSQL_ERROR_MSG				string	= "MySQL錯誤"
	ERROR_MYSQL_ERROR_CODE				int32	= 40001
	ERROR_REDIS_ERROR_MSG				string	= "Redis錯誤"
	ERROR_REDIS_ERROR_CODE				int32	= 40002
	ERROR_MYSQL_TRANSACTION_ERROR_MSG	string	= "提交處理出錯"
	ERROR_MYSQL_TRANSACTION_ERROR_CODE	int32	= 40003

	ERROR_NONE_VERIFICATION_CODE_MSG	string = "無驗證碼"
	ERROR_NONE_VERIFICATION_CODE_CODE	int32  = 41001
	ERROR_VERIFICATION_CODE_MSG			string = "驗證碼錯誤"
	ERROR_VERIFICATION_CODE_CODE		int32  = 41002
	ERROR_SEND_SMS_MSG					string = "簡訊發送失敗"
	ERROR_SEND_SMS_CODE					int32  = 41003
	ERROR_MISSING_PHONE_MSG				string = "手機不能為空白"
	ERROR_MISSING_PHONE_CODE			int32  = 41004
	ERROR_PHONE_ISEXIST_MSG				string = "該電話已被申請"
	ERROR_PHONE_ISEXIST_CODE			int32  = 41005
	ERROR_USER_LOCK_MSG					string = "用戶已鎖定"
	ERROR_USER_LOCK_CODE				int32  = 41006

	ERROR_WRONG_INPUT_MSG     			string = "輸入錯誤"
	ERROR_WRONG_INPUT_CODE    			int32  = 42001
	ERROR_WRONG_PASSWORD_MSG  			string = "密碼錯誤"
	ERROR_WRONG_PASSWORD_CODE 			int32  = 42002
	ERROR_WRONG_OPEN_ID_MSG   			string = "Open ID錯誤"
	ERROR_WRONG_OPEN_ID_CODE  			int32  = 42003

	ERROR_NONE_FRIEND_REQUEST_MSG		string = "無好友邀請資訊"
	ERROR_NONE_FRIEND_REQUEST_CODE		int32  = 43001
	ERROR_SENDED_FRIEND_REQUEST_MSG		string = "已發送好友邀請"
	ERROR_SENDED_FRIEND_REQUEST_CODE	int32  = 43002
	ERROR_FRIEND_ALREADY_MSG			string = "已成為好友"
	ERROR_FRIEND_ALREADY_CODE			int32  = 43003
	ERROR_UNHANDLED_FRIEND_REQUEST_MSG	string = "對方已加您為好友"
	ERROR_UNHANDLED_FRIEND_REQUEST_CODE	int32  = 43004

)

type RestfulAPIResponse struct {
	Status	int32
	Data	interface{}
	Msg		string
}

func Response(c *gin.Context, code int32, data interface{}, msg string) {
	c.JSON(http.StatusOK, RestfulAPIResponse{code, data, msg})
}