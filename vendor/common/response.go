package common

import (
	"github.com/golang/protobuf/proto"
	pb "protobuf"
)

/*
 輸出訊息通道
*/
type IMResponse  struct {
	Response *pb.IMResponse
}
type OutMessage chan IMResponse

/*
 錯誤訊息構造方法
*/
func NewIMResponseSimple(status int32, msg string, refer string) *IMResponse {
	res := &pb.IMResponse{
		Status: status,
		Msg: msg,
		Data: nil,
		Refer: refer,
	}
	return &IMResponse{res}
}

/*
 成功消息構造方法
*/
func NewIMResponseData(data []byte, refer string) *IMResponse {
	res := &pb.IMResponse{
		Status: 0,
		Msg: "",
		Data: data,
		Refer: refer,
	}
	return &IMResponse{res}
}

/*
 將返回訊息轉成protobuf
*/
func (this *IMResponse) Encode() []byte {
	s, _ := proto.Marshal(this.Response)
	return s
}

/*
 將protobuf轉成返回訊息
*/
func (this *IMResponse) Decode(data []byte) error {
	err := proto.Unmarshal (data, this.Response)
	return err
}

