package namesrv

type Request struct {
	ActionType uint16
	Body       []byte
}

type RegisterRequest struct {
	ServerName   string
	ServerAddr   string
	ClusterName  string
	HaServerAddr string
	ServerId     int
}

const (
	ActionRegister   = 0x0001 // register server
	ActionUnRegister = 0x0002 // unregister server

)
