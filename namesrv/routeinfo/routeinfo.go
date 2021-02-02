package routeinfo

import (
	"gdiamond/common/namesrv"
	"gdiamond/common/protocol"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/panjf2000/gnet"
	"log"
	"reflect"
	"sync"
	"time"
)

type RouteInfo struct {
	sync.RWMutex
	serverAddrTable  map[string] /* serverName */ *protocol.Server
	serverLiveTable  map[string] /* serverAddr */ *protocol.LiveServer
	clusterAddrTable map[string] /* clusterName */ *hashset.Set /* brokerName */
}

const (
	MASTER_ID                   = 0
	SERVER_CHANNEL_EXPIRED_TIME = 60
)

func New() *RouteInfo {
	r := &RouteInfo{serverAddrTable: make(map[string]*protocol.Server), serverLiveTable: make(map[string]*protocol.LiveServer), clusterAddrTable: make(map[string]*hashset.Set)}
	return r
}

//GetAllClusterInfo
func (r *RouteInfo) GetAllClusterInfo() protocol.ClusterInfo {
	cInfo := protocol.ClusterInfo{}
	cInfo.ClusterAddrTable = r.clusterAddrTable
	cInfo.ServerAddrTable = r.serverAddrTable
	return cInfo
}

//RegisterServer  diamond server  register it self,namesrv will keep it in mem
//too many args?
func (r *RouteInfo) RegisterServer(clusterName, serverAddr, serverName, haServerAddr string, serverId int, conn gnet.Conn) *namesrv.RegisterResponse {
	r.Lock()
	defer r.Unlock()

	result := &namesrv.RegisterResponse{}

	serverNames := r.clusterAddrTable[clusterName]
	if serverNames == nil {
		serverNames = hashset.New()
		r.clusterAddrTable[clusterName] = serverNames
	}
	serverNames.Add(serverName)
	registerFirst := false
	serverData := r.serverAddrTable[serverName]
	if serverData == nil {
		registerFirst = true
		serverData = protocol.NewServer(clusterName, serverName, make(map[int]string))
		r.serverAddrTable[serverName] = serverData
	}
	serverAddrsMap := serverData.GetServerAddrs()
	//Switch slave to master: first remove <1, IP:PORT> in namesrv, then add <0, IP:PORT>
	//The same IP:PORT must only have one record in serverAddrTable
	for k, v := range serverAddrsMap {
		if serverAddr != "" && serverAddr == v && serverId != k {
			delete(serverAddrsMap, k)
		}
	}

	oldAddr := serverData.GetServerAddrs()[serverId]
	serverData.GetServerAddrs()[serverId] = serverAddr
	registerFirst = registerFirst || "" == oldAddr

	prevServerLiveInfo := r.serverLiveTable[serverAddr]
	ls := protocol.NewLiveServer(time.Now().Unix(), haServerAddr, nil, conn)
	r.serverLiveTable[serverAddr] = ls
	if prevServerLiveInfo == nil {
		log.Println("new broker registered,  HAServer: ", serverAddr, haServerAddr)
	}
	if MASTER_ID != serverId {
		masterAddr := serverData.GetServerAddrs()[MASTER_ID]
		if masterAddr != "" {
			serverLiveInfo := r.serverLiveTable[masterAddr]
			if serverLiveInfo != nil {
				result.HaServerAddr = serverLiveInfo.HaServerAddr
				result.MasterAddr = masterAddr
			}
		}
	}
	return result
}

//UnregisterServer  unregister diamond server
//too many args???
func (r *RouteInfo) UnregisterServer(clusterName, serverAddr, serverName string, serverId int) {
	r.Lock()
	defer r.Unlock()
	serverlive := r.serverLiveTable[serverAddr]
	delete(r.serverLiveTable, serverAddr)
	log.Println("unregisterBroker, remove from serverLiveTable", serverlive, serverName)

	removeServerName := false
	server := r.serverAddrTable[serverName]
	if server != nil {
		serverAddress := server.GetServerAddrs()
		addr := serverAddress[serverId]
		delete(serverAddress, serverId)
		log.Println("unregisterServer, remove addr from serverAddrTable:", addr)

		if len(serverAddress) == 0 {
			delete(r.serverAddrTable, serverName)
			log.Println("unregisterServer, remove name from serverAddrTable:", serverName)
			removeServerName = true
		}
	}
	if removeServerName {
		nameSet := r.clusterAddrTable[clusterName]
		if nameSet != nil {
			nameSet.Remove(serverName)
			log.Println("unregisterServer, remove name from clusterAddrTable:", serverName)
			if nameSet.Size() == 0 {
				delete(r.clusterAddrTable, clusterName)
				log.Println("unregisterServer, remove cluster from clusterAddrTable:", clusterName)
			}
		}
	}
}

//ScanNotActiveServer namesrv will scheduled invoke this method to delete unlive diamond server
func (r *RouteInfo) ScanNotActiveServer() {
	for k, v := range r.serverLiveTable {
		last := v.LastUpdateTimestamp
		if last+SERVER_CHANNEL_EXPIRED_TIME < time.Now().Unix() {
			if v.Conn != nil {
				err := v.Conn.Close()
				log.Printf("close the connection to remote address:%s,err:%v,", k, err)
			}
			delete(r.serverLiveTable, k)
			log.Println("The server channel expired,", k, SERVER_CHANNEL_EXPIRED_TIME)
			r.onChannelDestroy(k, v.Conn)
		}
	}
}

func (r *RouteInfo) onChannelDestroy(remoteAddr string, conn gnet.Conn) {
	serverAddrFound := ""

	if conn != nil {
		r.RLock()
		for serverAddr, liveServerInfo := range r.serverLiveTable {
			if reflect.DeepEqual(liveServerInfo.Conn, conn) {
				serverAddrFound = serverAddr
				break
			}
		}
		r.RUnlock()
	}

	if serverAddrFound == "" {
		serverAddrFound = remoteAddr
	} else {
		log.Println("the broker's channel destroyed, {}, clean it's data structure at once", serverAddrFound)
	}
	if serverAddrFound != "" {
		r.Lock()
		delete(r.serverLiveTable, serverAddrFound)
		serverNameFound := ""
		removeServerName := false

		for serverName, server := range r.serverAddrTable {
			if serverNameFound != "" {
				break
			}

			serverAddrs := server.GetServerAddrs()
			for serverId, serverAddr := range serverAddrs {
				if serverAddr == serverAddrFound {
					serverNameFound = server.ServerName
					delete(serverAddrs, serverId)
					log.Println("remove serverAddr[{}, {}] from serverAddrTable, because channel destroyed", serverId, serverAddr)
					break
				}
			}

			if len(server.GetServerAddrs()) == 0 {
				removeServerName = true
				delete(r.serverAddrTable, serverName)
				log.Println("remove serverName  from serverAddrTable, because channel destroyed", server.ServerName)
			}
		}

		if serverNameFound != "" && removeServerName {
			for clusterName, serverNames := range r.clusterAddrTable {
				serverNames.Remove(serverNameFound)
				log.Println("remove serverName, clusterName  from clusterAddrTable, because channel destroyed", serverNameFound, clusterName)
				if serverNames.Empty() {
					log.Println("remove the clusterName[{}] from clusterAddrTable, because channel destroyed and no broker in this cluster", clusterName)
					delete(r.clusterAddrTable, clusterName)
				}
				break
			}
		}
		r.Unlock()
	}
}
