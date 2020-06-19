package main

import (
	"math/rand"
	"myserver/config"
	"myserver/util"
	"net"
	"sync"
	"time"
)

func CreateNewClient(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	conn, err := net.Dial("tcp", "localhost"+config.ServerConf.GameIp)
	util.HandleErr(err, "连接服务器失败")
	stopChan := make(chan struct{}, 1)
	//读协程
	go ReadLoop(stopChan, conn)
	for {
		//if !HandleInput(stopChan, conn) {
		if !AutoInput(stopChan, conn) {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(500)))
	}
	err = conn.Close()
	util.HandleErr(err, "关闭连接错误")
}
