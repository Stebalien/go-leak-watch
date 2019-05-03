package watch

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func Watch(threshold uint64, interval, maxFreq time.Duration, cb func()) {
	var lastNotif time.Time
	for range time.Tick(interval) {
		var mstats runtime.MemStats
		runtime.ReadMemStats(&mstats)
		if mstats.Alloc > threshold && time.Since(lastNotif) > maxFreq {
			cb()
			lastNotif = time.Now()
		}
	}
}

func init() {
	go Watch(1<<30, time.Second*10, time.Minute*5, func() {
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
