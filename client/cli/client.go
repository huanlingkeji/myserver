package cli

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"myserver/core/log"
	"myserver/util"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

//chan是否关闭 无法判断
//调用close conn控制client结构的生命周期结束

//定义一个客户端的网络底层操作结构

//数据如何发送

//数据如何接收

//网络协议格式

//数据包字节的处理

//这个结构的生命周期控制

//发送结构
type SendFuture struct {
	Future    chan struct{} //通知结构
	Cli       *NetClient    //客户端
	StartTime time.Time     //开启时间
	ProtoId   int           //协议id
	Req, Ack  proto.Message //发送接收数据
}

func NewSendFuture(cli *NetClient, protoId int, req, ack proto.Message) *SendFuture {
	return &SendFuture{
		Future:    make(chan struct{}),
		Cli:       cli,
		StartTime: time.Now(),
		ProtoId:   protoId,
		Req:       req,
		Ack:       ack,
	}
}

func (sf *SendFuture) Get(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-sf.Future:
			if ok {
				return nil
			}
		}
	}
}

func (sf *SendFuture) Send() error {
	req := sf.Req
	conn := sf.Cli.Conn
	client := sf.Cli
	protoId := sf.ProtoId
	ack := sf.Ack
	bytes, err := proto.Marshal(req)
	util.HandleErr(err, "解析数据错误")
	_ = protoId
	finalData := util.FinalPkg(util.PacketPayload(0, 0, client.ReqId, 0, bytes))
	if _, err = conn.Write(finalData); err == nil {
		reqId := client.ReqId
		atomic.AddInt32(&client.ReqId, 1)
		if nil != ack {
			client.SendMap[reqId] = sf
		}
	} else {
		log.Info("写数据错误:", err)
		return err
	}
	return nil
}

type NetClient struct {
	Conn       net.Conn              //连接
	Id         string                //标识
	OriginAddr string                //远程地址
	ConTimeOut time.Duration         //连接超时时间
	GetTimeOut time.Duration         //获取数据超时时间
	ReqId      int32                 //请求id号
	Out        chan *SendFuture      //发送数据
	In         chan []byte           //收取数据
	PushChan   chan string           //推送数据
	SendMap    map[int32]*SendFuture //发送结构
	StopChan   chan struct{}         //关闭的chan
	Die        bool                  //是否死亡
	//HeartBeatReq struct{}            //TODO 定义的心跳包数据
}

func NewNetClient(network string, id string, addr string, pushChan chan string) *NetClient {
	client := &NetClient{
		Id:         id,
		OriginAddr: addr,
		Die:        false,
		ConTimeOut: 5 * time.Second,
		GetTimeOut: 2 * time.Second,
		ReqId:      1,
		PushChan:   pushChan,
		Out:        make(chan *SendFuture),
		In:         make(chan []byte, 10), //接收有推送数据
		SendMap:    make(map[int32]*SendFuture),
		StopChan:   make(chan struct{}),
	}
	var err error

	switch network {
	//TODO udp类型的服务以后再处理
	//case "udp", "udp4", "udp6":
	//	client.con, err = dialKcpServer(client.targetAddr)
	default:
		client.Conn, err = net.DialTimeout("tcp", client.OriginAddr, client.ConTimeOut)
	}
	if nil != err {
		tips := fmt.Sprintf(`客户端:%v 连接服务器错误:%v`, client.Id, err)
		log.Info(tips)
		panic(tips)
	}
	return client
}

func (client *NetClient) Close() {
	if client.Die {
		return
	}
	err := client.Conn.Close()
	util.HandleErr(err, "关闭连接错误")
}

const heartBeatTime = 60 * time.Second

func (client *NetClient) Start() {
	//遇到一个bug wait没起效果  原因：在go协程内开启wg.add可能来不及  不要在协程内add
	defer func() {
		_ = client.Conn.Close()
		client.Die = true
		log.Info("关闭socket")
	}()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	//开启写
	go func() {
		defer func() {
			wg.Done()
			close(client.Out)
			log.Info("主/写循环结束")
		}()
		for {
			select {
			case out, ok := <-client.Out:
				if ok {
					if err := out.Send(); err != nil {
						close(client.StopChan)
						return
					}
				}
			case in, ok := <-client.In:
				if ok {
					client.HandleRecvMsg(in)
				} else {
					//close(client.StopChan) //In的关闭说明读协程已经关闭 这里是不需要关闭StopChan
					return
				}
				//TODO 还要处理心跳
			}
		}
	}()
	//如何主动关闭阻塞的读 关闭socket即可 如果需要等待读取数据可以关闭StopChan（是否还有意义呢？）
	go func() {
		defer func() {
			wg.Done()
			close(client.In)
			log.Info("读取循环结束")
		}()
		for {
			//读取超时退出时 写协程会因为收到in关闭而退出
			recvMsg := client.ReadResponse(heartBeatTime)
			if nil == recvMsg {
				return
			}
			//检测是否出现写错误导致的关闭
			select {
			case <-client.StopChan:
			default:
			}
			select {
			case client.In <- recvMsg: //这么写 跟 client.In<-recvMsg写法差别不大 当有阻塞情况发生时 添加case来做超时判断
			}
		}
	}()
	wg.Wait()
}

func (client *NetClient) ReadResponse(readTime time.Duration) []byte {
	//怎么读取数据
	//先读取包体长度 然后在读取完整的数据包
	sizeBuf := make([]byte, 4)
	_ = client.Conn.SetReadDeadline(time.Now().Add(readTime))
	_, err := io.ReadFull(client.Conn, sizeBuf)
	if err != nil {
		if err == io.EOF {
			//正常关闭流程
			log.Info("读取到EOF")
		} else {
			//异常流程
			log.Info("读取数据错误:", err)
			//出现的错误类型 use of closed network connection
		}
		return nil
	}
	size := binary.BigEndian.Uint32(sizeBuf)
	log.Info("读取到服务器数据 数据长度为:", size)
	recvMsg := make([]byte, size)
	_ = client.Conn.SetReadDeadline(time.Now().Add(readTime))
	_, err = io.ReadFull(client.Conn, recvMsg)
	util.HandleErr(err, "读取数据错误")
	log.Info("读取完整个包数据")
	return recvMsg
}

//处理获取数据
func (client *NetClient) HandleRecvMsg(in []byte) {
	//数据解包
	//反序列化成具体类型结构
	//分发处理
	srcId, mode, seqId, protoId, reqData := util.Unpack(in)
	if ft, ok := client.SendMap[int32(seqId)]; ok {
		delete(client.SendMap, int32(seqId))
		err := proto.Unmarshal(reqData, ft.Ack)
		if err == nil {
			ft.Future <- struct{}{}
		} else {
			panic("反序列化数据错误:" + err.Error())
		}
		log.Info("收到服务器的数据为:", fmt.Sprintf("%v", ft.Ack))
		//统计耗时长的调用
		cost := time.Now().Sub(ft.StartTime)
		if cost > 2*time.Second {
			log.Info("调用超过2s 客户端id:", client.Id, " 获取请求id:", seqId, " 协议名:", "xxxx", " 耗时:", cost, " 响应:", fmt.Sprintf("%v", ft.Ack))
		}
	} else { //收到服务器推送数据
		_, _, _ = srcId, mode, protoId
		//TODO 一定协议id范围的/特定协议id的为服务器推送数据
		var ack proto.Message
		_ = proto.Unmarshal(reqData, ack)
		client.PushChan <- "收到服务器数据推送:" + fmt.Sprintf("%v", ack)
	}
}

func (client *NetClient) Send(protoId int, req, ack proto.Message) (err error) {
	//创建一个发生结构
	//把数据封包然后发送
	defer func() {
		if x := recover(); x != nil {
			log.Info("发送数据错误:", x)
			err = errors.New(fmt.Sprintf("%v", x))
		}
	}()
	future := NewSendFuture(client, protoId, req, ack)
	client.Out <- future
	if nil != ack {
		if err := future.Get(client.GetTimeOut); err != nil {
			log.Info("发送结果结构获取错误:", err)
			return err
		} else {
			return nil
		}
	}
	return nil
}
