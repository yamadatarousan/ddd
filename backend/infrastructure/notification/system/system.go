package system

import (
	"fmt"
	"sync/atomic"
	"time"
)

// SequenceIDGeneratorは通知IDを連番で生成する。
type SequenceIDGenerator struct {
	sequence uint64
}

func NewSequenceIDGenerator() *SequenceIDGenerator {
	return &SequenceIDGenerator{}
}

func (g *SequenceIDGenerator) NextID() string {
	id := atomic.AddUint64(&g.sequence, 1)
	return fmt.Sprintf("notification-%d", id)
}

// RealtimeClockはシステム時刻を返す。
type RealtimeClock struct{}

func NewRealtimeClock() RealtimeClock {
	return RealtimeClock{}
}

func (c RealtimeClock) Now() time.Time {
	return time.Now()
}
