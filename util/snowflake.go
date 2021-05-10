package util

import (
	"errors"
	"sync"
	"time"
)

const (
	machineIDBits = int64(10)
	sequenceBits  = int64(12)
	maxSequence   = int64(1<<12 - 1)
)

type Snowflake struct {
	mu        *sync.Mutex
	startTime int64
	machineID int64
	sequence  int64
	lastTime  int64
}

func NewSnowflake(machineID int64) *Snowflake {
	return &Snowflake{new(sync.Mutex), time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC).UnixNano() / 1e6, machineID, 0, 0}
}

func (sf *Snowflake) NextID() (int64, error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	now := currentEpoch()
	duration := now - sf.lastTime
	if duration < 0 {
		return 0, errors.New("the machine time scroll back")
	}
	if duration == 0 {
		sf.sequence = (sf.sequence + 1) & maxSequence
		if sf.sequence == 0 {
			for duration <= 0 {
				now = currentEpoch()
				duration = now - sf.lastTime
			}
		}
	} else {
		sf.sequence = 0
	}
	sf.lastTime = now
	return sf.generateID(now), nil
}

func (sf *Snowflake) generateID(now int64) int64 {
	return (now-sf.startTime)<<(machineIDBits+sequenceBits) | sf.machineID<<sequenceBits | sf.sequence
}

func currentEpoch() int64 {
	return time.Now().UnixNano() / 1e6
}
