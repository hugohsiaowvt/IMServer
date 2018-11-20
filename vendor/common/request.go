package common

import (
	"github.com/golang/protobuf/proto"
	pb "protobuf"
)

/*
 輸入訊息通道
*/
type IMRequest  struct {
	Client  *Client
	Request *pb.IMRequest
}
type InMessage chan IMRequest

/*
 轉成protobuf資料
*/
func (this *IMRequest) Encode() []byte {
	//s, _ := json.Marshal(*this)
	s, _ := proto.Marshal(this.Request)
	return s
}

/*
 解析protobuf資料
*/
func (this *IMRequest) Decode(data []byte) error {
	//err := json.Unmarshal(data, this)
	obj := &pb.IMRequest{}
	err := proto.Unmarshal(data, obj)
	this.Request = obj
	return err
}

/*
 解析protobuf資料
*/
func DecodeIMRequest(data []byte) (*IMRequest, error) {
	req := &IMRequest{}
	err := req.Decode(data)
	if err != nil {
		return nil, err
	}
	return req, nil
}
