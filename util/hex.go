package util

func TrimHex(hexStr string) string {
	if len(hexStr) == 0 {
		return hexStr
	} else if len(hexStr) >= 2 && (hexStr[:2] == "0x" || hexStr[:2] == "0X") {
		hexStr = hexStr[2:]
	}

	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	return hexStr
}
