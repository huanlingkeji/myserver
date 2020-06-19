package main

import (
	"math/rand"
	"myserver/core/log"
	"net"
	"strings"
)

var SendTestData = []string{
	`测试数据1`,
	`测试数据2`,
	`测试数据3`,
	`测试数据4`,
}

const TestDataNum = 4

func AutoInput(stopChan chan struct{}, conn net.Conn) (isContinue bool) {
	log.Info("发送随机内容:")
	sendData := SendTestData[rand.Intn(TestDataNum)]
	select {
	case _ = <-stopChan:
		return false
	default:
	}
	n, err := conn.Write([]byte(sendData))
	_ = n
	//TODO 发送错误了该怎么做 什么类型的错误？
	if err != nil {
		log.Info("发送数据错误:", err)
		return false
	}
	return true
}

func HandleInput(stopChan chan struct{}, conn net.Conn) (isContinue bool) {
	log.Info("请输入内容:")
	input, _ := inputReader.ReadString('\n')
	input = strings.Trim(input, "\n")
	if input == "quit" {
		return false
	}
	select {
	case _ = <-stopChan:
		return false
	default:
	}
	n, err := conn.Write([]byte(input))
	_ = n
	//TODO 发送错误了该怎么做 什么类型的错误？
	if err != nil {
		log.Info("发送数据错误:", err)
		return false
	}
	return true
}
