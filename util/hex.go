package util

import "strings"

// TrimHex will remove the "0x" hex prefix if present and also 0-pad the
// hex string to an even length
func TrimHex(hexStr string) string {
	if len(hexStr) == 0 {
		return hexStr
	} else if strings.EqualFold(hexStr[:2], "0x") {
		hexStr = hexStr[2:]
	}

	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	return hexStr
}
