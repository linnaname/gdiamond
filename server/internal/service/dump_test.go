package service

import (
	"fmt"
	"gdiamond/server/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type _S struct {
	suite.Suite
}

func (s *_S) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	common.ParseCmdAndInitConfig()
	common.InitDBConn()
}

func (s *_S) TearDownSuite() {
	common.CloseConn()
}

func (s *_S) TestInit() {
	SetupDumpTask()
	for {

	}
}

func (s *_S) TestDumpAll2Disk() {
	err := DumpAll2Disk()
	assert.NoError(s.T(), err)
}

func TestS(t *testing.T) {
	suite.Run(t, new(_S))
}
