package common

type DNSProvider interface {
	GetRecord(conf map[string]interface{}) (r Record, err error)
	UpdateRecord(conf map[string]interface{}, r Record, value string) error
	AddRecord(conf map[string]interface{}, value string) error
}

type Record interface {
	Value() string
}

var DNSProviders = make(map[string]DNSProvider)

func RegisterDNS(provider string, dns DNSProvider) {
	DNSProviders[provider] = dns
}
