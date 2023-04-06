package cache

import (
	"fmt"
	"time"

	"github.com/mdaliyan/icache/v2"
)

var monitorCheckPot icache.Pot[int64]

func SetupLocalCache() {
	monitorCheckPot = icache.NewPot[int64](icache.WithTTL(8 * time.Second))
}

// SetMonitorToRun caches for 8 seconds
func SetMonitorToRun(monitorID, regionID int64) {
	monitorCheckPot.Set(getMonitorCacheKey(monitorID), regionID)
}

func GetMonitorToRun(monitorID int64) *int64 {
	value, err := monitorCheckPot.Get(getMonitorCacheKey(monitorID))
	if err != nil {
		return nil
	}
	return &value
}

func getMonitorCacheKey(monitorID int64) string {
	return fmt.Sprintf("monitor_to_run_%d", monitorID)
}
