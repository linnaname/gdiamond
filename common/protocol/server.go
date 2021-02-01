package protocol

//Server server info
type Server struct {
	cluster     string
	ServerName  string
	serverAddrs map[int] /* serverId */ string /* server address */
}

//NewServer new
func NewServer(cluster, serverName string, serverAddrs map[int]string) *Server {
	s := &Server{cluster: cluster, ServerName: serverName, serverAddrs: serverAddrs}
	return s
}

//GetServerAddrs getter
func (s *Server) GetServerAddrs() map[int]string {
	return s.serverAddrs
}
