package netx

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"myserver/core/log"
	"myserver/core/rpc"
	"myserver/util"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

const ReadDeadLine = 60

//TODO 一个链接/Session打印的日志添加标识

//TODO 不能直接建立session 而应该区分net层和应用层
//TODO net层数据会处理数据 根据数据类型来创建session或者发送数据给session处理
//TODO 连接和session的对应关系 session包含一个连接 连接不可用时可以替换掉 session数据不改动

type Connection struct {
	Conn net.Conn
	Out  chan []byte
	In   chan []byte
}

func HandleConn(conn net.Conn, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		_ = conn.Close()
	}()
	//构建一个session并处理它
	c := NewConnection(conn)
	wg.Add(1)
	log.Info("创建一个新的连接 远程地址:", c.Conn.RemoteAddr())
	c.ReadLoop()
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		Conn: conn,
		Out:  make(chan []byte, 1),
		In:   make(chan []byte, 1),
	}
}

func (c *Connection) ReadLoop() {
	defer func() {
		log.Info("关闭session的数据读取循环")
		if x := recover(); x != nil {
			log.Info("读取循环出错:", x)
			log.Info(string(debug.Stack()))
		}
	}()
	var sess *Session
	for {
		select {
		case <-AppDie:
			return
		default:

		}
		_ = c.Conn.SetReadDeadline(time.Now().Add(time.Second * ReadDeadLine))
		var bytes = make([]byte, 4)
		_, err := io.ReadFull(c.Conn, bytes)
		util.HandleErr(err, "连接读取数据错误")
		size := binary.BigEndian.Uint32(bytes)
		payload := make([]byte, size)
		_ = c.Conn.SetReadDeadline(time.Now().Add(time.Second * ReadDeadLine))
		_, err = io.ReadFull(c.Conn, payload)
		util.HandleErr(err, "连接读取数据错误")

		//解包处理
		//数据响应 TODO 一定会需要进行响应吗?
		_, _, _, protoId, reqData := util.Unpack(payload)
		_ = protoId
		//var ack proto.Message //TODO 不用使用反射获取协议实例 使用代码生成获取实例的代码
		var ack = &rpc.Msg{}
		err = proto.UnmarshalMerge(reqData, ack)
		util.HandleErr(err, "反序列化数据错误")
		log.Info(fmt.Sprintf("收到数据:%v", ack))

		if ack.Data == "连接" {
			//TODO 处理session注册相关
			sess = NewSession(c.Conn)
			log.Info("创建一个新的session")
			go c.WriteLoop(sess)
			go sess.Handle()
			sess.Send(payload)
		} else if ack.Data == "重连" {
			//TODO 让旧的session的连接断开 然后使用此连接
		} else {
			//TODO 连接部分已经把数据处理好了 直接使用就行
			if sess == nil {
				//TODO 把操作失败的原因发送给前端
				log.Info("session尚未建立")
			} else {
				sess.In <- payload
			}
		}
	}
}

func (c *Connection) WriteLoop(sess *Session) {
	defer func() {
		log.Info("关闭session的数据写循环")
		if x := recover(); x != nil {
			log.Info("发送数据循环错误:", x)
		}
	}()
	for {
		select {
		case data := <-sess.Out:
			//util.ShowBytes(data)
			n, err := c.Conn.Write(util.FinalPkg(data))
			_ = n
			util.HandleErr(err, "发送数据错误")
			log.Info("发送数据成功 字节数量:", n)
		case <-AppDie:
			return
		}
	}
}
