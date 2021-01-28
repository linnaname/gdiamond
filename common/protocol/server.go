package protocol

type Server struct {
	cluster     string
	ServerName  string
	serverAddrs map[int] /* serverId */ string /* server address */
}

func NewServer(cluster, serverName string, serverAddrs map[int]string) *Server {
	s := &Server{cluster: cluster, ServerName: serverName, serverAddrs: serverAddrs}
	return s
}

func (s *Server) GetServerAddrs() map[int]string {
	return s.serverAddrs
}
