package util

import (
	"errors"
	"strconv"
	"strings"
)

func HostAndPort(server string) (string, int, error) {
	if len(server) == 0 {
		return "", -1, errors.New("empty server destination")
	}

	if !strings.Contains(server, ":") {
		return "", -1, errors.New("expected host:port")
	}

	split := strings.Split(server, ":")
	if len(split) != 2 {
		return "", -1, errors.New("only expected a single `:` seperator")
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
