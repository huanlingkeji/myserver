package backup

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var inputReader *bufio.Reader

func init() {
	inputReader = bufio.NewReader(os.Stdin)
}

func _main() {
	wg := &sync.WaitGroup{}
	//wg.Add(1)
	wg.Add(3)
	kill := make(chan struct{})
	go server55(wg, kill)
	go func() {
		defer wg.Done()
		input, _ := inputReader.ReadString('\n')
		input = strings.Trim(input, "\n")
		if input == "quit" {
			kill <- struct{}{}
		}
	}()
	time.Sleep(time.Millisecond * 200)
	go client66(wg)
	wg.Wait()
}

func client1(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
	//回应：connectex: No connection could be made because the target machine actively refused it.
}

func server2(wg *sync.WaitGroup) {
	defer wg.Done()
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("error listen:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

	var i int
	for {
		time.Sleep(time.Second * 10)
		if _, err := l.Accept(); err != nil {
			log.Println("accept error:", err)
			break
		}
		i++
		log.Printf("%d: accept a new connection\n", i)
	}
}

func client2(wg *sync.WaitGroup) {
	defer wg.Done()
	var sl []net.Conn
	for i := 1; i < 1000; i++ {
		conn := establishConn(i)
		if conn != nil {
			sl = append(sl, conn)
		}
	}

	time.Sleep(time.Second * 10000)
	//windows下测试发现连接会被阻塞一段时间（2s）后失败 连接数200
}

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}
	log.Println(i, ":connect to server ok")
	return conn
}

func client3(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("begin dial...")
	conn, err := net.DialTimeout("tcp", "104.236.176.96:80", 2*time.Second)
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
	//使用clumsy做模拟延时
}

func client22(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	time.Sleep(time.Second * 2)
	conn.Write([]byte("test data342432"))

	time.Sleep(time.Second * 10000)
}

func server22(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	l, _ := net.Listen("tcp", ":8888")
	defer l.Close()
	c, _ := l.Accept()
	defer c.Close()
	for {
		// read from the connection
		var buf = make([]byte, 10)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error:", err)
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
	//如果收到的数据超过slice大小，重复收取
}

func server33(wg *sync.WaitGroup) {
	defer wg.Done()
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		log.Println("accept a new connection")
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		time.Sleep(10 * time.Second)
		var buf = make([]byte, 10)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error:", err)
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
	//数据会放置在缓冲区中 即使对方早已关闭了连接 数据依然可以读取处理
}

func client33(wg *sync.WaitGroup) {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	time.Sleep(time.Second * 2)
	conn.Write([]byte("test datadgesgrgregas543645754532"))
}

func client44(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	data := make([]byte, 65536)
	conn.Write(data)

	time.Sleep(time.Second * 10000)
}

func server44(wg *sync.WaitGroup) {
	defer wg.Done()
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		log.Println("accept a new connection")
		go handleConn2(c)
	}
}

func handleConn2(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		time.Sleep(1 * time.Second)
		var buf = make([]byte, 65536)
		log.Println("start to read from conn")
		c.SetReadDeadline(time.Now().Add(time.Microsecond * 1))
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
	//模拟不出效果
}

func client55(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	data := make([]byte, 65536)
	var total int
	for {
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error:%s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n", n, total)
	}

	log.Printf("write %d bytes in total\n", total)
	time.Sleep(time.Second * 10000)
	//关闭服务器 之后的write调用会写入0个字节（这个可能不是一定的）
	//数据写到一定之后会阻塞 即使对方接收了数据也不能发送（不应该是这样子的 为什么呢）
}

func server55(wg *sync.WaitGroup, kill chan struct{}) {
	defer wg.Done()
	l, _ := net.Listen("tcp", ":8888")
	defer l.Close()
	c, _ := l.Accept()
	defer c.Close()
	for {
		// read from the connection
		select {
		case <-kill:
			return
		default:

		}
		time.Sleep(5 * time.Second)
		var buf = make([]byte, 60000)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

func client66(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	data := make([]byte, 65536)
	var total int
	for {
		conn.SetWriteDeadline(time.Now().Add(time.Microsecond * 10))
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error:%s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n", n, total)
	}

	log.Printf("write %d bytes in total\n", total)
	time.Sleep(time.Second * 10000)
}

/*
从client1的结果来看，在己方已经关闭的socket上再进行read和write操作，会得到”use of closed network connection” error；
从server1的执行结果来看，在对方关闭的socket上执行read操作会得到EOF error，但write操作会成功，因为数据会成功写入己方的
内核socket缓冲区中，即便最终发不到对方socket缓冲区了，因为己方socket并未关闭。因此当发现对方socket关闭后，己方应该正确
合理处理自己的socket，再继续write已经无任何意义了。
*/
