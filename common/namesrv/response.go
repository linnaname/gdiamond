package namesrv

type Response struct {
	Code int
	Body []byte
}

type RegisterResponse struct {
	HaServerAddr string
	MasterAddr   string
	KvTable      map[string]string
}

const (
	Success      = 0
	SystemError  = 1
	SystemBusy   = 2
	NotSupported = 3
)
