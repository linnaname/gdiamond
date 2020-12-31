package service

import (
	"fmt"
	"gdiamond/server/common"
	"github.com/stretchr/testify/suite"
	"testing"
)

type _S struct {
	suite.Suite
}

func (s *_S) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	common.InitConfig()
	common.InitDBConn()
}

func (s *_S) TearDownSuite() {
	common.CloseConn()
}

func (s *_S) TestInit() {
	Init()
	for {

	}
}

func TestS(t *testing.T) {
	suite.Run(t, new(_S))
}
