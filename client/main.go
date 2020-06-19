/*
模拟客户端 自动手动输入内容
*/
package main

import (
	"bufio"
	_ "myserver/client/cli"
	"myserver/client/ui"
	_ "myserver/client/ui"
	"myserver/core/log"
	"os"
	"sync"
)

var inputReader *bufio.Reader

func init() {
	inputReader = bufio.NewReader(os.Stdin)
}

const clientNum = 100

var wg = &sync.WaitGroup{}

func main() {
	defer func() {
		if x := recover(); x != nil {
			log.Info("客户端出现错误:", x)
		}
	}()
	//for i := 0; i < clientNum; i++ {
	//	go CreateNewClient(wg)
	//}
	//wg.Wait()

	ui.CreateMainUI()
}

//TODO 统一编码为proto3
//TODO 数据的封包 解包 传输管理 源生包
