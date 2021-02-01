package namesrv

//Request general request
type Request struct {
	ActionType uint16
	//store different
	Body []byte
}

//RegisterRequest RegisterRequest
type RegisterRequest struct {
	ServerName   string
	ServerAddr   string
	ClusterName  string
	HaServerAddr string
	ServerId     int
}

const (
	//ActionRegister register server
	ActionRegister = 0x0001
	//ActionUnRegister unregister server
	ActionUnRegister = 0x0002
)
