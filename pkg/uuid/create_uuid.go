package uuid

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type SnowFlake struct {
	mu sync.Mutex
	twepoch int64

	workerIdBits     int64 
	datacenterIdBits int64 
	sequenceBits     int64 

	maxWorkerId     int64
	maxDatacenterId int64
	maxSequence     int64

	workerIdShift     int64
	datacenterIdShift int64
	timestampShift    int64

	datacenterId int64
	workerId int64
	sequence int64
	lastTimestamp int64
}

func (s *SnowFlake) timeGen() int64 {
	return time.Now().Unix()
}

func (s *SnowFlake) tilNextMills() int64 {
	timeStampMill := s.timeGen()
	for timeStampMill <= s.lastTimestamp {
		timeStampMill = s.timeGen()
	}
	return timeStampMill
}
func (s *SnowFlake) NextId() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nowTimestamp := s.timeGen() 
	if nowTimestamp < s.lastTimestamp {
		return -1, errors.New(fmt.Sprintf("clock moved backwards, Refusing to generate id for %d milliseconds", s.lastTimestamp-nowTimestamp))
	}
	if nowTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & s.maxSequence
		if s.sequence == 0 {
			nowTimestamp = s.tilNextMills()
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = nowTimestamp
	return (nowTimestamp-s.twepoch)<<s.timestampShift | 
			s.datacenterId<<s.datacenterIdShift | 
			s.workerId<<s.workerIdShift |
			s.sequence,
		nil
}

func NewSnowFlake(workerId int64, datacenterId int64) (*SnowFlake, error) {
	mySnow := new(SnowFlake)
	mySnow.twepoch = time.Now().Unix()
	if workerId < 0 || datacenterId < 0 {
		return nil, errors.New("workerId or datacenterId must not lower than 0 ")
	}

	mySnow.workerIdBits = 5
	mySnow.datacenterIdBits = 5
	mySnow.sequenceBits = 12

	mySnow.maxWorkerId = -1 ^ (-1 << mySnow.workerIdBits)         
	mySnow.maxDatacenterId = -1 ^ (-1 << mySnow.datacenterIdBits) 
	mySnow.maxSequence = -1 ^ (-1 << mySnow.sequenceBits)         

	if workerId >= mySnow.maxWorkerId || datacenterId >= mySnow.maxDatacenterId {
		return nil, errors.New("workerId or datacenterId must not higher than max value ")
	}
	mySnow.workerIdShift = mySnow.sequenceBits
	mySnow.datacenterIdShift = mySnow.sequenceBits + mySnow.workerIdBits
	mySnow.timestampShift = mySnow.sequenceBits + mySnow.workerIdBits + mySnow.datacenterIdBits

	mySnow.lastTimestamp = -1
	mySnow.workerId = workerId
	mySnow.datacenterId = datacenterId

	return mySnow, nil
}
