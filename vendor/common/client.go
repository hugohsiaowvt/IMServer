package common

import (
	"log"
	"net"
	"service/userservice"
)

/*
 Client端結構體
 */
type Client struct {
	// 連接訊息
	Key    string        // Client端連接的key
	Conn   net.Conn      // Socket連線
	In     InMessage     // 輸入訊息
	Out    OutMessage    // 輸出訊息
	Quit   chan *Client  // 退出
	// 登入訊息
	User	*userservice.User
}

/*
 Client端列表
 */
type ClientTable map[string]*Client

/*
 獲取輸入訊息
 */
func (this *Client) GetIn() IMRequest {
	return <-this.In
}

/*
 設置輸出訊息
 */
func (this *Client) PutOut(resp *IMResponse) {
	this.Out <- *resp
}

/*
 創建Client端
 */
func CreateClient(key string, conn net.Conn) *Client {
	client := &Client{
		Key:    key,
		Conn:   conn,
		In:     make(InMessage),
		Out:    make(OutMessage),
		Quit:   make(chan *Client),
	}
	client.Listen()
	return client
}

/*
 自動讀入或寫出訊息
 */
func (this *Client) Listen() {
	go this.read()
	go this.write()
}

/*
 退出連線
 */
func (this *Client) Quiting() {
	this.Quit <- this
}

/*
 關閉連接
 */
func (this *Client) Close() {
	this.Conn.Close()
}

/*
 讀取訊息
 */
func (this *Client) read() {
	for {
		data := make([]byte, 32768)
		n, err := this.Conn.Read(data)
		if err != nil {
			log.Printf("讀取錯誤: %s\n", err)
			this.Quiting()
			return
		}
		req, err := DecodeIMRequest(data[0:n])
		if err == nil {
			req.Client = this
			this.In <- *req
		} else {
			log.Printf("解析封包錯誤: %s", data)
		}
	}
}

/*
 輸出訊息
 */
func (this *Client) write() {
	for resp := range this.Out {
		if _, err := this.Conn.Write(resp.Encode()); err != nil {
			this.Quiting()
			return
		}
	}
}
