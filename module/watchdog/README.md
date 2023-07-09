So this needs quick explanation -

1. We first run the monitor - `run.go`
2. Which hits the route - `hit.go`
3. We assert the response back in the run - `assert.go`
4. We also update the check - `run.go`
5. We then send for verification - `monitor.go`
6. We match against the configured alarm policy to realize if a monitor is passing or failing - `monitor.go`
7. Then we send for alarm check - `alarm.go`
8. We check if we should raise or resolve an alarm - `alarm.go`
9. We then proceed to send notification - `notification.go`
10. We look at the configured alarm channels - `notification.go`
11. We maintain the configured reminder policy - `notification.go`
12. Alarm channels are in - `channel folder`
