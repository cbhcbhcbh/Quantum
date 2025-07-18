package uuid

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestUUid(T *testing.T) {
	mySnow, _ := NewSnowFlake(829323232, 9999999999999999)
	group := sync.WaitGroup{}
	startTime := time.Now()
	generateId := func(s SnowFlake, requestNumber int) {
		for i := 0; i < requestNumber; i++ {
			uuids, _ := s.NextId()
			fmt.Println(uuids)
			group.Done()
		}
	}
	group.Add(400)
	currentThreadNum := 4
	for i := 0; i < currentThreadNum; i++ {
		generateId(*mySnow, 100)
	}
	group.Wait()
	fmt.Printf("time: %v\n", time.Now().Sub(startTime))

}
