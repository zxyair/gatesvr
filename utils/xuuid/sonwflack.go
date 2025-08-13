package xuuid

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch         = int64(1609459200000) // 2021-01-01 00:00:00 UTC
	machineIDBits = uint(10)
	sequenceBits  = uint(12)

	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits

	maxMachineID = -1 ^ (-1 << machineIDBits)
	maxSequence  = -1 ^ (-1 << sequenceBits)
)

type Snowflake struct {
	machineID int64
	sequence  int64
	lastStamp int64
	mu        sync.Mutex
}

func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, errors.New("machine ID out of range")
	}
	return &Snowflake{
		machineID: machineID,
	}, nil
}

func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if now < s.lastStamp {
		return 0, errors.New("clock moved backwards")
	}

	if now == s.lastStamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastStamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastStamp = now

	id := (now-epoch)<<timestampShift |
		(s.machineID << machineIDShift) |
		s.sequence

	return id, nil
}
