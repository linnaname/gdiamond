package protocol

import "github.com/panjf2000/gnet"

type LiveServer struct {
	LastUpdateTimestamp int64
	HaServerAddr        string
	DataVersion         *DataVersion
	Conn                gnet.Conn
}

func NewLiveServer(lastUpdateTimestamp int64, haServerAddr string, dataVersion *DataVersion, conn gnet.Conn) *LiveServer {
	ls := &LiveServer{LastUpdateTimestamp: lastUpdateTimestamp, HaServerAddr: haServerAddr,
		DataVersion: dataVersion, Conn: conn}
	return ls
}
