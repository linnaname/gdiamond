package kvconfig

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type _S struct {
	suite.Suite
	kc *KVConfig
}

func (s *_S) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	s.kc = New("./test/kv.json")
	s.kc.Load()
}

func (s *_S) TearDownSuite() {
	fmt.Printf("TearDownSuite() ...\n")
}

func TestS(t *testing.T) {
	suite.Run(t, new(_S))
}

func (s *_S) TestPutKVConfig() {
	s.kc.PutKVConfig("ns", "linnana", "me")
	s.kc.PutKVConfig("ns", "aa", "11")

}

func (s *_S) TestGetKVConfig() {
	assert.Equal(s.T(), s.kc.GetKVConfig("ns", "linnana"), "me")
}

func (s *_S) TestGetKVListByNamespace() {
	b := s.kc.GetKVListByNamespace("ns")
	assert.NotNil(s.T(), b)
	fmt.Println(b)
}

func (s *_S) TestDeleteKVConfig() {
	s.kc.DeleteKVConfig("ns", "linnana")
	assert.Empty(s.T(), s.kc.GetKVConfig("ns", "linnana"))

}
