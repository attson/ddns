package getip

import (
	"github.com/attson/ddns/getip/common"
	_ "github.com/attson/ddns/getip/gw"
)

func IpProvider(provider string) common.GetIpFunc {
	return common.GetIpMap[provider]
}
