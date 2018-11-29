package server

import (
	"common"
	"net"
	"strings"
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"log"
	"github.com/pborman/uuid"
	"service/userservice"
	"rz/config"
	"rz/redis"

	"github.com/golang/protobuf/proto"
	pb "protobuf"
)



type Server struct {
	listener    net.Listener        // Server監聽器
	clients     common.ClientTable  // 客戶端列表，抽象出來單獨維護和入參，更方便管理連接
	joinsniffer chan net.Conn       // 觸發創建Client連接處理方法
	quitsniffer chan *common.Client // 觸發連接退出處理方法
	insniffer   common.InMessage    // 觸發接收訊息處理方法
}

/*
 IM Server啟動方法
*/
func StartIMServer() {
	log.Println("	Server啟動中...")
	// 初始化Server
	server := &Server{
		clients:     make(common.ClientTable, config.MAX_CLIENTS),
		joinsniffer: make(chan net.Conn),
		quitsniffer: make(chan *common.Client),
		insniffer:   make(common.InMessage),
	}
	// 添加關閉構子，當關Server時執行
	server.interruptHandler()
	// 啟動監聽方法(包含各種探測器)
	server.listen()
	// 啟動Server端口監聽(等待連接)
	server.start()
}

/*
 監聽方法
*/
func (this *Server) listen() {
	go func() {
		for {
			select {
			// 接收到了訊息
			case message := <-this.insniffer:
				this.receivedHandler(message)
			// 新來了一個連線
			case conn := <-this.joinsniffer:
				this.joinHandler(conn)
			// 退出了一個連線
			case client := <-this.quitsniffer:
				this.quitHandler(client)
			}
		}
	}()
}


/*
 新Client端請求處理方法
*/
func (this *Server) joinHandler(conn net.Conn) {
	// 獲取UUID作為Client端的key
	key := uuid.New()
	// 創建一個Client端
	client := common.CreateClient(key, conn)
	// 给客户端指定key
	this.clients[key] = client
	log.Printf("新客户端Key:[%s] Online:%d", client.Key, len(this.clients))
	// 開始協程不斷地接收訊息
	go func() {
		for {
			msg := <-client.In
			this.insniffer <- msg
		}
	}()
	// 開始協程一直等待斷開
	go func() {
		for {
			conn := <-client.Quit
			this.quitsniffer <- conn
		}
	}()
}


/*
 Client端退出處理方法
*/
func (this *Server) quitHandler(client *common.Client) {
	if client != nil {
		// 調用Client端關閉方法
		client.Close()
		delete(this.clients, client.Key)

		log.Printf("Client端退出: %s Online:%d", client.Key, len(this.clients))
	}
}


/*
 接收訊息處理方法
*/
func (this *Server) receivedHandler(request common.IMRequest) {
	log.Println("開始讀取資料")

	// 獲得請求的Client
	client := request.Client
	// 獲得請求的資料
	reqData := request.Request
	log.Printf("Client:[%s]發送命令:[%s] 訊息內容:[%s]", client.Key, reqData.Command, reqData.Data)

	// 未授權業務處理部分
	switch reqData.Command {
	case pb.GET_CONN:
		req := &pb.ReqConn{}
		res := &pb.ResConn{}
		if err := proto.Unmarshal(reqData.Data, req); err != nil {
			return
		}
		openId := req.OpenId
		token := req.Token
		// 判斷 IM token 存不存在
		if userservice.CheckIMToken(openId, token) {
			// 判斷用戶存不存在
			user, exist := userservice.QueryUserExistByOpenId(openId)
			if exist {
				// 將key寫進redis
				key := SOCKET_KEY_PREFIX + openId
				if _, err := redis.Instance().HSet(key, "key", client.Key); err != nil {
					client.PutOut(common.NewIMResponseSimple(1, err.Error(), pb.GET_CONN_RETURN))
					return
				}
				// 將這條socket的用戶設為剛查詢到的用戶
				client.User = &user
				res.Status = 1
				data, _ := proto.Marshal(res)
				client.PutOut(common.NewIMResponseData(data, pb.GET_CONN_RETURN))
				return
			} else {
				client.PutOut(common.NewIMResponseSimple(ERROR_NONE_USER_CODE, ERROR_NONE_USER_MSG, pb.GET_CONN_RETURN))
				return
			}
		} else {
			client.PutOut(common.NewIMResponseSimple(ERROR_IM_TOKEN_CODE, ERROR_IM_TOKEN_MSG, pb.GET_CONN_RETURN))
			return
		}
	}

	// 驗證連線是否已授權
	if client.User == nil {
		client.PutOut(common.NewIMResponseSimple(ERROR_UNAUTHORIZED_CODE, ERROR_UNAUTHORIZED_MSG, pb.UNAUTHORIZED))
		return
	}

	// 已授權業務處理部分
	switch reqData.Command {
	case pb.SEND_MSG:
		req := &pb.ReqMsg{}
		if err := proto.Unmarshal(reqData.Data, req); err != nil {
			return
		}
		// 判斷溝通通道存不存在
		conversationId := req.ToConversationId
		smallId := req.FromId
		bigId := req.ToConversationId
		if (req.RoomType == SINGLE_CHAT) {
			log.Println("單人訊息")
			if (strings.Compare(smallId, bigId) == 1) {
				smallId, bigId = bigId, smallId
			}
			conversationId = smallId + "_" + bigId
		}

		list, err := redis.Instance().SMembers(CONVERSATION_PREFIX + conversationId)
		if err != nil {
			client.PutOut(common.NewIMResponseSimple(ERROR_REDIS_ERROR_CODE, err.Error(), pb.SEND_MSG_RETURN))
		}
		if len(list) == 0 {
			if (req.RoomType == SINGLE_CHAT) {
				log.Println("建立conversation")
				if _, err1 := redis.Instance().SAdd(CONVERSATION_PREFIX + conversationId, smallId); err1 != nil {
					client.PutOut(common.NewIMResponseSimple(ERROR_REDIS_ERROR_CODE, err1.Error(), pb.SEND_MSG_RETURN))
				}
				if _, err1 := redis.Instance().SAdd(CONVERSATION_PREFIX + conversationId, bigId); err1 != nil {
					client.PutOut(common.NewIMResponseSimple(ERROR_REDIS_ERROR_CODE, err1.Error(), pb.SEND_MSG_RETURN))
				}
				this.sendMsgByKeys([]string {smallId, bigId}, req)
				return
			} else {
				client.PutOut(common.NewIMResponseSimple(ERROR_NONE_GROUP_CODE, ERROR_NONE_GROUP_MSG, pb.SEND_MSG_RETURN))
			}
		} else {
			this.sendMsgByKeys(list, req)
			return
		}

	}
	
}

/*
 發送訊息
*/
func (this *Server) sendMsgByKeys(list []string, req *pb.ReqMsg) {
	res := &pb.ResMsg{}
	res.FromId = req.FromId
	res.ToConversationId = req.ToConversationId
	res.RoomType = req.RoomType
	res.MsgType = req.MsgType
	res.Data = req.Data
	res.Time = time.Now().Unix()
	data, _ := proto.Marshal(res)
	for _, openId := range list {
		key := SOCKET_KEY_PREFIX + openId
		v, err := redis.Instance().HGet(key, "key")
		if err != nil {

		}
		if this.clients[v] == nil {
			//推送、離線訊息處理
			continue
		}
		this.clients[v].PutOut(common.NewIMResponseData(data, pb.SEND_MSG_RETURN))
	}
}


func (this *Server) start() {
	// 設置監聽地址及端口
	addr := fmt.Sprintf("0.0.0.0:%s", config.IM_PORT)
	this.listener, _ = net.Listen("tcp", addr)
	log.Printf("開始監聽IM Server[%s]端口", config.IM_PORT)
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			log.Fatalln(err)
			return
		}
		// log.Printf("新連接地址為:[%s]", conn.RemoteAddr())
		this.joinsniffer <- conn
	}
}

// FIXME: need to figure out if this is the correct approach to gracefully
// terminate a server.
func (this *Server) Stop() {
	this.listener.Close()
}

// Server關閉時執行
func (this *Server) interruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		sig := <-c
		log.Printf("captured %v, stopping profiler and exiting..", sig)
		// 清除客户端连接
		for _, v := range this.clients {
			this.quitHandler(v)
		}
		// 退出
		os.Exit(1)
	}()
}