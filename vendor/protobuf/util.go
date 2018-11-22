package protobuf

const (
	GET_CONN             uint32 = 10001	// 建立TCP長連接
	GET_CONN_RETURN      uint32 = 10002	// 獲得連接返回

	SEND_MSG             uint32 = 20001	// 發送訊息
	SEND_MSG_RETURN      uint32 = 20002	// 發送訊息返回

	UNAUTHORIZED         uint32 = 90001	// 未授權
)