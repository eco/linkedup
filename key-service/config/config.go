package config

type cfg struct {
	longyRestURL string
}

var globalCfg = cfg{}

// SetLongyRestURL -
func SetLongyRestURL(url string) {
	globalCfg.longyRestURL = url
}

// LongyRestURL -
func LongyRestURL() string {
	return globalCfg.longyRestURL
}
