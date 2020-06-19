package ui

import (
	"github.com/lxn/walk"
	"myserver/config"
	"myserver/core/log"
)

import (
	. "github.com/lxn/walk/declarative"
)

type MyMainWindow struct {
	*walk.MainWindow
	IdLt   *walk.LineEdit //id的输入框
	AddrLt *walk.LineEdit //地址的输入框
}

func CreateMainUI() {
	mw := &MyMainWindow{}
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "客户端控制器",
		Layout:   VBox{},
		Size:     Size{Width: 400, Height: 300},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Label{
						Text: "客户端标识",
					},
					LineEdit{
						AssignTo: &mw.IdLt,
						MaxSize:  Size{Width: 200, Height: 50},
					},
				},
			},
			HSplitter{
				Children: []Widget{
					Label{
						Text: "服务器地址",
					},
					LineEdit{
						AssignTo: &mw.AddrLt,
						Text:     "localhost" + config.ServerConf.GameIp,
						MaxSize:  Size{Width: 200, Height: 50},
					},
				},
			},
			PushButton{
				Text:      "创建一个客户端",
				OnClicked: mw.CreateClientUI,
			},
		},
	}.Create()); err != nil {
		log.Info("创建客户端窗口失败:", err)
	}
	mw.Run()
}

func (mw *MyMainWindow) CreateClientUI() {
	id := mw.IdLt.Text()
	addr := mw.AddrLt.Text()
	if CreateClientUI(id, addr) != 0 {
		walk.MsgBox(mw, "错误", "创建客户端失败或者客户端运行错误", walk.MsgBoxOK)
	}
}
