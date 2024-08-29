package aliyun

import (
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/attson/ddns/dns/common"
)

func init() {
	common.RegisterDNS("aliyun", AliDNS{})
}

type config struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	RR              string `json:"rr"`
	DomainName      string `json:"domain_name"`
}

func (c *config) fromMap(m map[string]interface{}) *config {
	c.AccessKeyId = m["access_key_id"].(string)
	c.AccessKeySecret = m["access_key_secret"].(string)
	c.RR = m["rr"].(string)
	c.DomainName = m["domain_name"].(string)

	return c
}

type Record struct {
	id    string
	value string
}

func (r Record) Value() string {
	return r.value
}

type AliDNS struct {
}

func (a AliDNS) AddRecord(conf map[string]interface{}, value string) error {
	c := (&config{}).fromMap(conf)
	client, err := createClient(c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return err
	}

	_, err = client.AddDomainRecord(&alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(c.DomainName),
		Type:       tea.String("A"),
		RR:         tea.String(c.RR),
		Value:      tea.String(value),
	})

	return err
}

func (a AliDNS) GetRecord(conf map[string]interface{}) (record common.Record, err error) {
	c := (&config{}).fromMap(conf)
	client, err := createClient(c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	result, err := client.DescribeSubDomainRecordsWithOptions(&alidns20150109.DescribeSubDomainRecordsRequest{
		SubDomain: tea.String(c.RR + "." + c.DomainName),
	}, &util.RuntimeOptions{})

	if err != nil {
		return nil, err
	}

	for _, v := range result.Body.DomainRecords.Record {
		return &Record{
			id:    *v.RecordId,
			value: *v.Value,
		}, nil
	}

	return nil, common.RecordNotFound
}

func (a AliDNS) UpdateRecord(conf map[string]interface{}, record common.Record, value string) error {
	c := (&config{}).fromMap(conf)
	client, err := createClient(c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return err
	}

	_, err = client.UpdateDomainRecord(&alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(record.(*Record).id),
		Value:    tea.String(value),
		Type:     tea.String("A"),
		RR:       tea.String(c.RR),
	})

	return err
}

func createClient(accessKey string, accessKeySecret string) (_result *alidns20150109.Client, _err error) {
	c := &openapi.Config{
		AccessKeyId:     tea.String(accessKey),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	c.Endpoint = tea.String("dns.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(c)
	return _result, _err
}
