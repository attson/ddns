package getip

import (
	"github.com/attson/ddns/getip/common"
	_ "github.com/attson/ddns/getip/tygw"
)

func IpProvider(provider string) common.GetIpFunc {
	return common.GetIpProviders[provider]
}
