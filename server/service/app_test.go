package service

import (
	"container/list"
	"fmt"
	"gdiamond/server/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type _SS struct {
	suite.Suite
}

func (s *_SS) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	common.InitConfig()
	common.InitDBConn()
}

func (s *_SS) TearDownSuite() {
	common.CloseConn()
}

func TestSS(t *testing.T) {
	suite.Run(t, new(_SS))
}

func (s *_SS) TestAddConfigInfo() {
	err := AddConfigInfo("gdiamond.test.vv", "GDIADMOND", "whether true")
	assert.NoError(s.T(), err)
}

func (s *_SS) TestUpdateConfigInfo() {
	err := UpdateConfigInfo("gdiamond.test.vv", "GDIADMOND", "whether update ok")
	assert.NoError(s.T(), err)
}

func (s *_SS) TestFindConfigInfo() {
	cInfo, err := FindConfigInfo("gdiamond.test.vv", "GDIADMOND")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cInfo)
	assert.NotEmpty(s.T(), cInfo.Content)
}

func (s *_SS) TestFindConfigInfoPage() {
	page, err := FindConfigInfoPage(1, 10, "GDIADMOND", "gdiamond.test.vv")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.Equal(s.T(), page.TotalCount, 1)

	pEmptyDataId, err := FindConfigInfoPage(1, 10, "GDIADMOND", "")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), pEmptyDataId)
	assert.Greater(s.T(), pEmptyDataId.TotalCount, 0)

	pEmptyGroup, err := FindConfigInfoPage(1, 10, "", "gdiamond.test.vv")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), pEmptyGroup)
	assert.Greater(s.T(), pEmptyGroup.TotalCount, 0)

	pEmptyGroupAndDataId, err := FindConfigInfoPage(1, 10, "", "")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), pEmptyGroupAndDataId)
	assert.Greater(s.T(), pEmptyGroupAndDataId.TotalCount, 0)
}

func (s *_SS) TestFindConfigInfoLike() {
	page, err := FindConfigInfoLike(1, 2, "lin", "DEFAULT_GROUP")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.NotNil(s.T(), page.PageItems)
	assert.Greater(s.T(), page.TotalCount, 1)
}

func TestName(t *testing.T) {
	nums := []int{-2, -1, 0, 3, 5, 6, 7, 9, 13, 14}
	result := pairSums(nums, 12)
	println(result)
}

func pairSums(nums []int, target int) [][]int {
	totalPair := 0
	counter := make(map[int]int, len(nums))
	l := list.New()
	for _, num := range nums {
		pairNum := target - num
		cnt, ok := counter[pairNum]

		if !ok {
			val, ok := counter[num]
			if !ok {
				val = 1
			} else {
				val += 1
			}
			counter[num] = val
		} else {
			pair := make([]int, 2)
			pair[0] = num
			pair[1] = pairNum
			l.PushBack(pair)
			totalPair++
			if cnt == 1 {
				delete(counter, pairNum)
			} else {
				cnt--
				counter[pairNum] = cnt
			}
		}
	}

	pairs := make([][]int, totalPair)
	for i := l.Front(); i != nil; i = i.Next() {
		val := i.Value
		pair, _ := val.([]int)
		pairs[totalPair] = pair
		totalPair--
	}

	return pairs
}
