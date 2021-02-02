package service

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gdiamond/common/namesrv"
	"gdiamond/server/common"
	"gdiamond/util/netutil"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/smallnest/goframe"
	"log"
	"net"
	"sync"
	"time"
)

//Register  register gdiamond-server to nameserver
type Register struct {
}

//RegisterServerAll register gdiamond-server to all nameserver
func (r *Register) RegisterServerAll(nameServerAddressList *singlylinkedlist.List) {
	if nameServerAddressList != nil && !nameServerAddressList.Empty() {
		request := namesrv.Request{}
		request.ActionType = namesrv.ActionRegister
		rreq := common.RegisterRequestConfig
		if rreq.ServerAddr == "" {
			rreq.ServerAddr = netutil.GetLocalIP()
		}
		if rreq.HaServerAddr == "" {
			rreq.HaServerAddr = netutil.GetLocalIP()
		}
		body, _ := json.Marshal(rreq)
		request.Body = body
		wg := sync.WaitGroup{}
		wg.Add(nameServerAddressList.Size())
		nameServerAddressList.Each(func(index int, value interface{}) {
			namesrvAddr, _ := value.(string)
			go registerServer(namesrvAddr, 1000, request, &wg)
		})
		wg.Wait()
	}
}

func registerServer(namesrvAddr string, timeoutMills int, request namesrv.Request, wg *sync.WaitGroup) error {
	defer wg.Done()

	address := fmt.Sprintf("%s:%v", namesrvAddr, 9000)
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeoutMills)*time.Millisecond)
	//conn.(*net.TCPConn).SetKeepAlive(true)
	//conn.(*net.TCPConn).SetKeepAlivePeriod(time.Second * 10)
	if err != nil {
		log.Println("DialTimeout err:", err)
		return err
	}
	defer conn.Close()

	fc := getFrameConn(conn)

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	err = fc.WriteFrame(data)
	if err != nil {
		return err
	}

	buf, err := fc.ReadFrame()
	if err != nil {
		return err
	}
	resp := namesrv.Response{}
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return err
	}

	if resp.Code != namesrv.Success {
		return fmt.Errorf("register wrong, res code: %v", resp.Code)
	}
	return nil
}

func getFrameConn(conn net.Conn) goframe.FrameConn {
	encoderConfig := goframe.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}

	decoderConfig := goframe.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 4,
	}

	return goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)
}
