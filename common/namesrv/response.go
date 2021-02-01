package namesrv

//Response general register response
type Response struct {
	Code int
	Body []byte
}

//RegisterResponse RegisterResponse
type RegisterResponse struct {
	HaServerAddr string
	MasterAddr   string
	KvTable      map[string]string
}

const (
	//Success success
	Success = 0
	//SystemError system internal error
	SystemError = 1
	//SystemBusy system busy can't response
	SystemBusy = 2
	//NotSupported not supported op
	NotSupported = 3
)
