package network

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gdiamond/common/namesrv"
	"gdiamond/namesrv/routeinfo"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"log"
)

//NameServer name server
type NameServer struct {
	*gnet.EventServer
	async      bool
	workerPool *goroutine.Pool
	addr       string
	routeInfo  *routeinfo.RouteInfo
}

//New setup name server
func New(addr string, routeInfo *routeinfo.RouteInfo) error {
	codec := getCodec()
	ns := &NameServer{addr: addr, async: true, workerPool: goroutine.Default(), routeInfo: routeInfo}
	err := gnet.Serve(ns, addr, gnet.WithMulticore(true),
		/*gnet.WithTCPKeepAlive(time.Second*20)*/ gnet.WithCodec(codec), gnet.WithReusePort(true))
	return err
}

//OnInitComplete implement gnet method
func (ns *NameServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("NameServer is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

//OnShutdown implement gnet method
func (ns *NameServer) OnShutdown(srv gnet.Server) {
	log.Printf("NameServer OnShutdown on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

//OnOpened implement gnet method
func (ns *NameServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("NameServer OnOpened")
	return
}

//OnClosed implement gnet method
func (ns *NameServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("NameServer OnClosed")
	return
}

//React  implement gnet method, handle logic
func (ns *NameServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Println("React frame:", string(frame))

	if ns.async {
		_ = ns.workerPool.Submit(func() {
			response := ns.processRequest(frame, c)
			//TODO Marshal err?
			result, _ := json.Marshal(response)
			c.AsyncWrite(result)
		})
		return
	}
	out = frame
	return
}

func (ns *NameServer) processRequest(data []byte, c gnet.Conn) namesrv.Response {
	response := namesrv.Response{}
	response.Code = namesrv.NotSupported

	request := &namesrv.Request{}
	err := json.Unmarshal(data, request)
	if err != nil {
		response.Code = namesrv.SystemError
		return response
	}

	body := request.Body
	rr := &namesrv.RegisterRequest{}
	err = json.Unmarshal(body, rr)
	if err != nil {
		response.Code = namesrv.SystemError
		return response
	}

	switch request.ActionType {
	case namesrv.ActionRegister:
		rresult := ns.routeInfo.RegisterServer(rr.ClusterName, rr.ServerAddr, rr.ServerName, rr.HaServerAddr, rr.ServerId, c)
		resBody, _ := json.Marshal(rresult)
		response.Body = resBody
		response.Code = namesrv.Success
	case namesrv.ActionUnRegister:
		ns.routeInfo.UnregisterServer(rr.ClusterName, rr.ServerAddr, rr.ServerName, rr.ServerId)
		response.Code = namesrv.Success
	}
	return response
}

func getCodec() gnet.ICodec {
	encoderConfig := gnet.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}
	decoderConfig := gnet.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 4,
	}
	return gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)
}
