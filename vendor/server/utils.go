package server

const (

	SYSTEM_CHAT				int32	= 0
	SINGLE_CHAT				int32	= 1
	GROUP_CHAT				int32	= 2


	SOCKET_KEY_PREFIX   	string	= "socket_key:"
	CONVERSATION_PREFIX 	string	= "conversation:"
	UNREAD_PREFIX			string	= "unread:"

	ERROR_IM_TOKEN_MSG		string	= "IM Token錯誤"
	ERROR_IM_TOKEN_CODE		int32	= -1
	ERROR_NONE_USER_MSG		string	= "無此用戶"
	ERROR_NONE_USER_CODE	int32	= -2
	ERROR_UNAUTHORIZED_MSG	string	= "用戶未登錄"
	ERROR_UNAUTHORIZED_CODE	int32	= -3

	ERROR_MYSQL_ERROR_MSG				string	= "MySQL錯誤"
	ERROR_MYSQL_ERROR_CODE				int32	= 40001
	ERROR_REDIS_ERROR_MSG				string	= "Redis錯誤"
	ERROR_REDIS_ERROR_CODE				int32	= 40002
	ERROR_MYSQL_TRANSACTION_ERROR_MSG	string	= "提交處理出錯"
	ERROR_MYSQL_TRANSACTION_ERROR_CODE	int32	= 40003

	ERROR_NONE_GROUP_MSG	string	= "無此群組"
	ERROR_NONE_GROUP_CODE	int32	= 41001

)