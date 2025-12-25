package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/cloud/cloudflare"
	"github.com/stretchr/testify/assert"
)

func TestCloudflare(t *testing.T) {
	client := cloudflare.NewFromConfigure()
	flog.Infof("AccountId: %s", client.AccountId)
	flog.Infof("ApiToken: %s\n", client.ApiToken)

	flog.Infof("=======================验证账户=======================")
	success, code, err := client.Verify()
	flog.Infof("验证结果: %v, 状态码: %d, 错误信息: %v\n", success, code, err)

	var testZoneId string
	var testDomain string
	flog.Infof("=======================列表所有区域=======================")
	zones, _ := client.ZoneList(100, 1)
	for _, zone := range zones.Result {
		flog.Infof("区域：%s，域名：%s", zone.ID, zone.Name)

		// 找到测试域名，用来模拟操作DNS记录
		if strings.Contains(zone.Name, "test") {
			testZoneId = zone.ID
			testDomain = zone.Name
		}
	}

	flog.Infof("=======================查询特定dns=======================")
	dnsClient := client.NewDnsClient(testZoneId)
	dnss, _ := dnsClient.List(100, 1)
	for _, dns := range dnss.Result {
		flog.Infof(fmt.Sprintf("[%s] %s %s", dns.Type, dns.Name, dns.Content))
	}

	flog.Infof("=======================添加dns记录=======================")
	success1, id1, _, err1 := dnsClient.Create("A", "test1."+testDomain, "127.0.0.1", false, 3600, "测试添加记录"+time.Now().Format("2006-01-02 15:04:05"))
	assert.Truef(t, success1, "添加dns记录失败:%+v", err1)

	dns1, _ := dnsClient.Info(id1)
	flog.Infof(fmt.Sprintf("%s [%s] %s %s", dns1.Result.ID, dns1.Result.Type, dns1.Result.Name, dns1.Result.Content))

	flog.Infof("=======================修改dns记录=======================")
	success2, _, err2 := dnsClient.Update(id1, "CNAME", "test2."+testDomain, "www."+testDomain, true, 1800, "测试修改记录"+time.Now().Format("2006-01-02 15:04:05"))
	assert.Truef(t, success2, "修改dns记录失败:%+v", err2)

	dns2, _ := dnsClient.Info(id1)
	flog.Infof(fmt.Sprintf("%s [%s] %s %s", dns2.Result.ID, dns2.Result.Type, dns2.Result.Name, dns2.Result.Content))

	flog.Infof("=======================删除dns记录=======================")
	oprId, _ := dnsClient.Delete(id1)
	flog.Infof("删除操作ID: %s", oprId)
	dns3, _ := dnsClient.Info(id1)
	flog.Infof(fmt.Sprintf("%s [%s] %s %s", dns3.Result.ID, dns3.Result.Type, dns3.Result.Name, dns3.Result.Content))
}
