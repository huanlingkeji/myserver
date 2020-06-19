package main

import (
	"myserver/core/log"
	"myserver/util"
	"net"
)

func ReadLoop(stopChan chan struct{}, conn net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			//TODO 细分错误类型
			stopChan <- struct{}{}
		}
	}()
	for {
		var bytes = make([]byte, 512)
		_, err := conn.Read(bytes)
		util.HandleErr(err, "连接读取数据错误")
		log.Info("收到的数据:", string(bytes))
		//TODO 转交业务逻辑层
	}
}
