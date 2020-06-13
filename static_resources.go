package pithy

var staticResouces = make(map[string][]byte)

// add static resouces to cache
func SetSRC(key string, bytes []byte) {
	staticResouces[key] = bytes
}

// add static resouces to cache
func SetSRCByFunc(key string, f func() []byte) {
	bytes := f()
	staticResouces[key] = bytes
}

// get resouces from cache
func GetResource(key string) []byte {
	if r, ok := staticResouces[key]; ok {
		return r
	}
	return nil
}
