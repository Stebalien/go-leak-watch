package watch

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func Watch(ctx context.Context, threshold uint64, interval, maxFreq time.Duration, cb func()) {
	var lastNotif time.Time
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			var mstats runtime.MemStats
			runtime.ReadMemStats(&mstats)
			if mstats.Alloc <= threshold {
				continue
			}
			if now.Before(lastNotif.Add(maxFreq)) {
				continue
			}
			lastNotif = now
			cb()
		case <-ctx.Done():
			return
		}
	}
}

func Auto(ctx context.Context) {
	go Watch(ctx, 1<<30, time.Second*10, time.Minute*5, func() {
		fmt.Println("memory threshold triggered! writing profile!")
		fi, err := os.Create(fmt.Sprintf("high-mem-profile-%d", time.Now().UnixNano()))
		if err != nil {
			fmt.Println("Error: failed to open file for memory profile: ", err)
			return
		}
		defer fi.Close()
		if err := pprof.WriteHeapProfile(fi); err != nil {
			fmt.Println("failed to write memory profile: ", err)
		}
	})
}
