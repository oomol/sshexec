package flags

import "sync"

var (
	OomolStudioPID int = -1
	once           sync.Once
)

func SetOomolStudioPID(pid int) {
	once.Do(func() {
		OomolStudioPID = pid
	})
}
