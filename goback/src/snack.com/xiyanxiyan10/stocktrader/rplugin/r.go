package rplugin

import (
	"fmt"
	"github.com/senseyeio/roger"
)

// NewRClient ...
func NewRClient(host string, port int64) (roger.RClient, error) {
	rClient, err := roger.NewRClient(host, port)
	if err != nil {
		fmt.Println("Failed to connect")
		return nil, err
	}
	return rClient, nil
}
