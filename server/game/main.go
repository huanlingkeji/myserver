package main

import (
	"bufio"
	"myserver/config"
	"myserver/core/log"
	"myserver/server/netx"
	"myserver/util"
	"net"
	"os"
	"strings"
	"sync"
)

var inputReader *bufio.Reader

func init() {
	inputReader = bufio.NewReader(os.Stdin)
}

func main() {
	defer func() {
		if x := recover(); x != nil {
			log.Info("服务器出现错误:", x)
		}
	}()
	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.ServerConf.GameIp)
	util.HandleErr(err, "获取tcpAddr")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	util.HandleErr(err, "获取listener错误")
	log.Info("开启监听:", listener.Addr().String())
	//TODO 处理信号 让服务器迅速处理剩余连接的数据后退出
	go SolveInput(listener)
	wg := &sync.WaitGroup{}
	for {
		newConn, err := listener.AcceptTCP()
		//util.HandleErr(err, "获取新连接失败")
		if err != nil {
			log.Info("获取新连接失败:", err)
			break
		}
		err = newConn.SetReadBuffer(config.ServerConf.SockBuf)
		util.HandleErr(err, "设置读缓存区大小失败")
		err = newConn.SetWriteBuffer(config.ServerConf.SockBuf)
		util.HandleErr(err, "设置写缓存区大小失败")

		go netx.HandleConn(newConn, wg)
	}
	wg.Wait()
}

func SolveInput(listener *net.TCPListener) {
	for {
		input, _ := inputReader.ReadString('\n')
		input = strings.Trim(input, "\n")
		//TODO 什么情况下 服务器才能真正关闭?
		if input == "quit" {
			err := listener.Close()
			close(netx.AppDie)
			util.HandleErr(err, "关闭服务器错误")
			return
		}
		log.Info("你的输入为:", input)
	}
}

//怎么自己做一个框架
//可以处理多连接
//单连接的处理 获取数据 写数据 异常关闭
//网络协议设计
//日志怎么做?
//服务交互
//心跳机制
//服务组网
//连接管理
//服务应用
//构建运行发布
