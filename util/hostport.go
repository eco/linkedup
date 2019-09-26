package util

import (
	"errors"
	"strconv"
	"strings"
)

// HostAndPort will parse the `address` into it's host and port. An error is
// returned if address is in the incorrect format, expected; "host:port"
func HostAndPort(address string) (string, int, error) {
	if len(address) == 0 {
		return "", -1, errors.New("empty address destination")
	}

	if !strings.Contains(address, ":") {
		return "", -1, errors.New("expected host:port")
	}

	split := strings.Split(address, ":")
	if len(split) != 2 {
		return "", -1, errors.New("only expected a single `:` separator")
	}

	host := split[0]
	if len(host) == 0 {
		return "", -1, errors.New("expected host")
	}
	port, err := strconv.Atoi(split[1])
	if err != nil || port < 0 {
		return "", -1, errors.New("port must be in decimal format and positive")
	}

	return host, port, nil
}
