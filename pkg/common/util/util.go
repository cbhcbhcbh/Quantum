package util

import (
	"strings"
)

func GetServerAddrs(addrs string) []string {
	return strings.Split(addrs, ",")
}

