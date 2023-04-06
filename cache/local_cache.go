package cache

import (
	"fmt"
	"time"

	"github.com/mdaliyan/icache/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
)

var monitorCheckPot icache.Pot[int64]
var monitorRegionCheckPot icache.Pot[string]

func SetupLocalCache() {
	monitorCheckPot = icache.NewPot[int64](icache.WithTTL(8 * time.Second))
}

// SetMonitorToRun caches for 8 seconds
func SetMonitorToRun(monitorID, regionID int64) {
	monitorCheckPot.Set(getMonitorToRunCacheKey(monitorID), regionID)
}

func GetMonitorToRun(monitorID int64) *int64 {
	value, err := monitorCheckPot.Get(getMonitorToRunCacheKey(monitorID))
	if err != nil {
		return nil
	}
	return &value
}

func getMonitorToRunCacheKey(monitorID int64) string {
	return fmt.Sprintf("monitor_to_run_%d", monitorID)
}

// SetMonitorRegionRunning caches for 8 seconds
func SetMonitorRegionRunning(monitorRegionID int64) {
	monitorRegionCheckPot.Set(getMonitorRegionRunningCacheKey(monitorRegionID), config.App.FlyRegion)
}

func GetMonitorRegionRunning(monitorRegionID int64) *string {
	value, err := monitorRegionCheckPot.Get(getMonitorRegionRunningCacheKey(monitorRegionID))
	if err != nil {
		return nil
	}
	return &value
}

func getMonitorRegionRunningCacheKey(monitorRegionID int64) string {
	return fmt.Sprintf("monitor_region_running_%d", monitorRegionID)
}
