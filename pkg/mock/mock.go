package mock

import (
	"fmt"
	"time"
)

type MockTimer struct {
	NowTime       time.Time
	SleepCalled   bool
	SleepDuration time.Duration
}

func (m *MockTimer) Now() time.Time {
	return m.NowTime
}

func (m *MockTimer) Sleep(d time.Duration) {
	m.SleepCalled = true
	m.SleepDuration = d
}

type MockLogger struct {
	Messages []string
}

func (m *MockLogger) Printf(format string, args ...interface{}) {
	m.Messages = append(m.Messages, fmt.Sprintf(format, args...))
}
