package routeinfo

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRouteInfo_RegisterServer(t *testing.T) {
	r := New()
	r.RegisterServer("testCluster", "127.0.0.2", "testServer2", "127.0.0.2", 0, nil)
	resp := r.RegisterServer("testCluster", "127.0.0.1", "testServer2", "127.0.0.1", 1, nil)
	assert.Equal(t, resp.MasterAddr, "127.0.0.2")
	assert.Equal(t, resp.HaServerAddr, "127.0.0.2")

	r.RegisterServer("testCluster2", "127.0.0.3", "testServer", "127.0.0.3", 1, nil)
	assert.Equal(t, len(r.clusterAddrTable), 2)
	assert.Equal(t, len(r.serverAddrTable), 2)
	assert.Equal(t, len(r.serverLiveTable), 3)
}

func TestRouteInfo_UnregisterServer(t *testing.T) {
	r := New()
	r.RegisterServer("testCluster", "127.0.0.2", "testServer2", "127.0.0.2", 0, nil)
	r.RegisterServer("testCluster", "127.0.0.1", "testServer2", "127.0.0.1", 1, nil)
	r.RegisterServer("testCluster2", "127.0.0.3", "testServer", "127.0.0.3", 1, nil)

	r.UnregisterServer("testCluster", "127.0.0.1", "testServer2", 1)
	assert.Equal(t, len(r.clusterAddrTable), 2)
	assert.Equal(t, len(r.serverAddrTable), 2)
	assert.Equal(t, len(r.serverLiveTable), 2)

	r.UnregisterServer("testCluster2", "127.0.0.3", "testServer", 1)
	assert.Equal(t, len(r.clusterAddrTable), 1)
	assert.Equal(t, len(r.serverAddrTable), 1)
	assert.Equal(t, len(r.serverLiveTable), 1)
}

func TestRouteInfo_ScanNotActiveServer(t *testing.T) {
	r := New()
	r.RegisterServer("testCluster", "127.0.0.2", "testServer2", "127.0.0.2", 0, nil)
	r.RegisterServer("testCluster2", "127.0.0.3", "testServer", "127.0.0.3", 1, nil)
	r.ScanNotActiveServer()
	assert.Equal(t, len(r.clusterAddrTable), 2)
	assert.Equal(t, len(r.serverAddrTable), 2)
	assert.Equal(t, len(r.serverLiveTable), 2)

	time.Sleep(15 * time.Second)

	r.ScanNotActiveServer()
	assert.Equal(t, len(r.clusterAddrTable), 0)
	assert.Equal(t, len(r.serverAddrTable), 0)
	assert.Equal(t, len(r.serverLiveTable), 0)
}

func TestRouteInfo_GetAllClusterInfo(t *testing.T) {
	r := New()
	r.RegisterServer("testCluster", "127.0.0.2", "testServer2", "127.0.0.2", 0, nil)
	r.RegisterServer("testCluster", "127.0.0.1", "testServer2", "127.0.0.1", 1, nil)
	r.RegisterServer("testCluster2", "127.0.0.3", "testServer", "127.0.0.3", 1, nil)
	cInfo := r.GetAllClusterInfo()
	assert.NotNil(t, cInfo)
}
