package helpers

import (
	"os"
)

func GetHostName() string {
	host, _ := os.Hostname()
	return host
}
