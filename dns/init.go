package dns

import (
	_ "github.com/attson/ddns/dns/aliyun"
	"github.com/attson/ddns/dns/common"
)

func GetDNSProvider(provider string) common.DNSProvider {
	return common.DNSProviders[provider]
}
