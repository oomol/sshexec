package os

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/shirou/gopsutil/v3/process"
)

// IsProcessAliveV4 if return err != nil, process not found
func IsProcessAliveV4(pid int) (bool, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return false, fmt.Errorf("failed to find process: %w", err)
	}
	running, err := proc.IsRunning()
	if err != nil {
		return false, fmt.Errorf("failed to check if process is running: %w", err)
	}
	if !running {
		return false, fmt.Errorf("process %d not found", pid)
	}
	return true, nil
}

const tickerInterval = 300 * time.Millisecond

func WatchPPID(g *errgroup.Group, ctx context.Context, parentPid int) {
	logrus.Info("listen parent pid: ", parentPid)
	g.Go(func() error {
		ticker := time.NewTicker(tickerInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				isRunning, err := IsProcessAliveV4(parentPid)
				if !isRunning {
					return fmt.Errorf("PPID %d exited, possible error: %w", parentPid, err)
				}
			}
		}
	})
}

func ListenSignal(g *errgroup.Group, ctx context.Context) {
	g.Go(func() error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-sigChan:
			return fmt.Errorf("catch signal: %v", s.String())
		}
	})
}
