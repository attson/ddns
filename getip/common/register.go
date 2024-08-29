package common

type GetIpFunc func(conf map[string]interface{}) (string, error)

var GetIpProviders = make(map[string]GetIpFunc)

func RegisterGetIp(provider string, getIp GetIpFunc) {
	GetIpProviders[provider] = getIp
}
