package main

import (
	"server"
	"log"
)

const (
	Name    string = "IM System"
	Version string = "1.0"
)

/*
 啟動Server方法
*/
func Start() {
	go func() {
		err := server.StartHttpServer()
		log.Fatalf("HttpServer", err)
	}()
	// 啟動IM服務
	server.StartIMServer()
}

func main() {
	log.Println("*********************************************")
	log.Printf("           系統:[%s]版本:[%s]", Name, Version)
	log.Println("*********************************************")
	Start()
}