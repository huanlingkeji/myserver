package netx

import (
	"myserver/core/log"
	"net"
)

//TODO 链接和session的关系 session生命周期的管理
//TODO 重连
//TODO 链接的断开  如何优雅主动断开 如何被动断开

type Session struct {
	Conn       net.Conn
	LocalAddr  string
	RemoteAddr string
	Out        chan []byte
	In         chan []byte
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		Conn:       conn,
		LocalAddr:  conn.LocalAddr().String(),
		RemoteAddr: conn.RemoteAddr().String(),
		Out:        make(chan []byte, 1),
		In:         make(chan []byte, 1),
	}
}

func (sess *Session) Handle() {
	defer func() {
		log.Info("session的处理循环关闭")
	}()
	for {
		select {
		case data, ok := <-sess.In:
			if ok {
				route(sess, data)
			} else {
				return
			}
		case <-AppDie:
			return
		}
	}
}

//TODO 应该是处理成了好的数据包
func route(sess *Session, frame []byte) {
	sess.Send(frame)
}

func (sess *Session) Send(payload []byte) {
	sess.Out <- payload
}
