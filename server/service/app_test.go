package service

import (
	"container/list"
	"testing"
)

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
