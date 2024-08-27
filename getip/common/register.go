package common

type GetIpFunc func(conf map[string]interface{}) (string, error)

var GetIpMap = make(map[string]GetIpFunc)

func RegisterGetIp(provider string, getIp GetIpFunc) {
	GetIpMap[provider] = getIp
}
