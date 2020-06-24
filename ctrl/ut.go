package ctrl

import "strings"

func IsLocalNetwork(host string) bool {
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "192.") {
		return true
	}
	return false
}
