package routeinfo

import (
	"gdiamond/common/namesrv"
	"gdiamond/common/protocol"
	logger "gdiamond/namesrv/internal/log"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/panjf2000/gnet"
	"github.com/sirupsen/logrus"
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
	MasterId                 = 0
	ServerChannelExpiredTime = 60
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
		logger.Logger.WithFields(logrus.Fields{
			"serverLiveTable":  r.serverLiveTable,
			"serverAddr":       serverAddr,
			"clusterAddrTable": r.clusterAddrTable,
		}).Warn("prevServerLiveInfo is nil")
	}
	if MasterId != serverId {
		masterAddr := serverData.GetServerAddrs()[MasterId]
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

	logger.Logger.WithFields(logrus.Fields{
		"serverlive": serverlive,
		"serverName": serverName,
	}).Info("unregisterBroker, remove from serverLiveTable")

	removeServerName := false
	server := r.serverAddrTable[serverName]
	if server != nil {
		serverAddress := server.GetServerAddrs()
		addr := serverAddress[serverId]
		delete(serverAddress, serverId)

		logger.Logger.WithFields(logrus.Fields{
			"addr": addr,
		}).Info("unregisterServer, remove addr from serverAddrTable")

		if len(serverAddress) == 0 {
			delete(r.serverAddrTable, serverName)
			logger.Logger.WithFields(logrus.Fields{
				"serverName": serverName,
			}).Info("unregisterServer, remove name from serverAddrTable")

			removeServerName = true
		}
	}
	if removeServerName {
		nameSet := r.clusterAddrTable[clusterName]
		if nameSet != nil {
			nameSet.Remove(serverName)
			logger.Logger.WithFields(logrus.Fields{
				"serverName": serverName,
			}).Info("unregisterServer, remove name from clusterAddrTable")

			if nameSet.Size() == 0 {
				delete(r.clusterAddrTable, clusterName)
				logger.Logger.WithFields(logrus.Fields{
					"clusterName": clusterName,
				}).Info("unregisterServer, remove name from clusterAddrTable")
			}
		}
	}
}

//ScanNotActiveServer namesrv will scheduled invoke this method to delete unlive diamond server
func (r *RouteInfo) ScanNotActiveServer() {
	for k, v := range r.serverLiveTable {
		last := v.LastUpdateTimestamp
		if last+ServerChannelExpiredTime < time.Now().Unix() {
			if v.Conn != nil {
				err := v.Conn.Close()
				logger.Logger.WithFields(logrus.Fields{
					"address": k,
					"err":     err,
				}).Info("close the connection to remote address")
			}
			delete(r.serverLiveTable, k)
			logger.Logger.WithFields(logrus.Fields{
				"address":                  k,
				"ServerChannelExpiredTime": ServerChannelExpiredTime,
			}).Info("The server channel expired")

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
		logger.Logger.WithFields(logrus.Fields{
			"serverAddrFound": serverAddrFound,
		}).Info("the broker's channel destroyed, clean it's data structure at once")
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
					logger.Logger.WithFields(logrus.Fields{
						"serverId":   serverId,
						"serverAddr": serverAddr,
					}).Info("remove serverAddr from serverAddrTable, because channel destroyed")
					break
				}
			}

			if len(server.GetServerAddrs()) == 0 {
				removeServerName = true
				delete(r.serverAddrTable, serverName)
				logger.Logger.WithFields(logrus.Fields{
					"ServerName": server.ServerName,
				}).Info("remove serverAddr from serverAddrTable, because channel destroyed")
			}
		}

		if serverNameFound != "" && removeServerName {
			for clusterName, serverNames := range r.clusterAddrTable {
				serverNames.Remove(serverNameFound)
				logger.Logger.WithFields(logrus.Fields{
					"serverNameFound": serverNameFound,
					"clusterName":     clusterName,
				}).Info("remove serverAddr,clusterName from serverAddrTable, because channel destroyed")

				if serverNames.Empty() {
					logger.Logger.WithFields(logrus.Fields{
						"clusterName": clusterName,
					}).Info("remove the clusterName[{}] from clusterAddrTable, because channel destroyed and no broker in this cluster")
					delete(r.clusterAddrTable, clusterName)
				}
				break
			}
		}
		r.Unlock()
	}
}
