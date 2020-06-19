package ui

import (
	"fmt"
	"github.com/lxn/walk"
	"myserver/client/cli"
	"myserver/core/log"
	"myserver/core/rpc"
	"runtime/debug"
)

import (
	. "github.com/lxn/walk/declarative"
)

type ClientWindow struct {
	*walk.MainWindow
	OutPrintTe *walk.TextEdit //结果输出区
	Client     *cli.NetClient //网络底层客户端
	SendLt     *walk.LineEdit //发送的数据输入框
	PushChan   chan string
	//TODO 状态数据存储 模拟游戏客户端数据
	//TODO 玩家的状态数据如何发送到应用层进行处理呢？ 处理函数集合 可以访问到状态数据对象  可以提供给网络层做回调
}

func CreateClientUI(id string, addr string) (exitCode int) {
	defer func() {
		if x := recover(); x != nil {
			log.Info("创建客户端窗口失败:", x)
			log.Info(string(debug.Stack()))
			exitCode = -1
		}
	}()
	pushChan := make(chan string)
	mw := &ClientWindow{
		PushChan: pushChan,
		Client:   cli.NewNetClient("tcp", id, addr, pushChan),
	}
	go mw.Client.Start()
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "模拟客户端",
		Layout:   VBox{},
		Size:     Size{Width: 500, Height: 800},
		Children: []Widget{
			PushButton{
				Text: "关闭客户端",
				OnClicked: func() {
					mw.Client.Close()
					_ = mw.Close()
				},
			},
			PushButton{
				Text: "关闭连接",
				OnClicked: func() {
					mw.Client.Close()
				},
			},
			PushButton{
				Text:      "重新连接",
				OnClicked: mw.ReConnected,
			},
			HSplitter{
				Children: []Widget{
					Label{
						Text: "请输入内容",
					},
					LineEdit{
						AssignTo: &mw.SendLt,
						MaxSize:  Size{Width: 200, Height: 50},
					},
					PushButton{
						Text:      "发送",
						OnClicked: mw.HandleSend,
					},
				},
			},
			TextEdit{
				VScroll:  true,
				HScroll:  true,
				AssignTo: &mw.OutPrintTe,
				Text:     "",
				ReadOnly: true,
			},
		},
	}.Create()); err != nil {
		log.Info("创建客户端窗口失败:", err)
	}
	go ShowPushMsg(mw)
	return mw.Run()
}

func ShowPushMsg(mw *ClientWindow) {
	for {
		select {
		case data := <-mw.PushChan:
			_ = mw.OutPrintTe.SetText(string(data) + mw.OutPrintTe.Text())
		}
	}
}

func (mw *ClientWindow) HandleSend() {
	data := mw.SendLt.Text()
	var req, ack = &rpc.Msg{Data: data}, &rpc.Msg{}
	if err := mw.Client.Send(0, req, ack); err == nil {
		log.Info("收到的数据为：", fmt.Sprintf("%v", ack))
		_ = mw.OutPrintTe.SetText("收到数据为:" + fmt.Sprintf("%v", ack) + "\r\n" + mw.OutPrintTe.Text())
	} else {
		_ = mw.OutPrintTe.SetText("发送数据失败 请重新连接 原因:" + err.Error() + "\r\n" + mw.OutPrintTe.Text())
		walk.MsgBox(mw, "失败", "发送数据失败 请重新连接", walk.MsgBoxOK)
	}
}

func (mw *ClientWindow) ReConnected() {
	id, addr, pushChan := mw.Client.Id, mw.Client.OriginAddr, mw.PushChan
	mw.Client.Close()
	mw.Client = cli.NewNetClient("tcp", id, addr, pushChan)
	go mw.Client.Start()
	walk.MsgBox(mw, "成功", "重新连接成功", walk.MsgBoxOK)
}
