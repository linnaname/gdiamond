package protocol

import "github.com/emirpasic/gods/sets/hashset"

//ClusterInfo cluster info
type ClusterInfo struct {
	ServerAddrTable  map[string] /* serverName */ *Server
	ClusterAddrTable map[string] /* clusterName */ *hashset.Set /* serverrName */
}
