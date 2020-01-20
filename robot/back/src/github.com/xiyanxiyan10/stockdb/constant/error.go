package constant

import (
	"fmt"
)

var (
	ErrHTTPUnauthorized     = fmt.Errorf("Unauthorized")
	ErrInfluxdbNotConnected = fmt.Errorf("Influxdb is not connected")
)
