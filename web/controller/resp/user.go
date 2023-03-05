package resp

import "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"

type User struct {
	Token string
	*model.User
}
