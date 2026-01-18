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
	assert.Equal(t, true, success)

	var testZoneId string
	var testDomain string
	flog.Infof("=======================查询所有区域=======================")
	zones, _ := client.ZoneList(100, 1)
	for _, zone := range zones.Result {
		flog.Infof("区域：%s，域名：%s", zone.ID, zone.Name)

		// 找到测试域名，用来模拟操作DNS记录
		if strings.Contains(zone.Name, "test") {
			testZoneId = zone.ID
			testDomain = zone.Name
		}
	}
	assert.NotEmpty(t, testZoneId, "未找到测试域名的ZoneId，请确保在Cloudflare中有包含test字样的域名")

	flog.Infof("=======================查询所有dns记录=======================")
	dnsClient := client.NewDnsClient(testZoneId)
	dnss, _ := dnsClient.List("", 100, 1)
	for _, dns := range dnss.Result {
		flog.Infof(fmt.Sprintf("[%s] %s %s", dns.Type, dns.Name, dns.Content))
	}
	assert.Greaterf(t, len(dnss.Result), 0, "测试域名下未找到任何DNS记录，请确保在Cloudflare中有对应域名的DNS记录")

	flog.Infof("=======================添加自定义主机=======================")
	hostName := "testdcv3." + testDomain
	success, dnsId, dcvDnsId, customHostnameId, err := client.CreateCustomHostnameAndVerify(testZoneId, testZoneId, hostName, "ddns."+testDomain, "CNAME")
	assert.Truef(t, success, "添加自定义主机失败:%+v", err)

	// 验证dns解析记录
	dnsDetail, _ := dnsClient.Info(dnsId)
	flog.Infof("CNAME %s -> %s", dnsDetail.Result.Name, dnsDetail.Result.Content)
	assert.Equal(t, hostName, dnsDetail.Result.Name)
	assert.Equal(t, "ddns."+testDomain, dnsDetail.Result.Content)

	// 验证dcv记录
	dcvDnsDetail, _ := dnsClient.Info(dcvDnsId)
	flog.Infof("CNAME %s -> %s", dcvDnsDetail.Result.Name, dcvDnsDetail.Result.Content)
	assert.Equal(t, "_acme-challenge."+hostName, dcvDnsDetail.Result.Name)

	flog.Infof("=======================查询自定义主机=======================")
	customHostnameClient := client.NewCustomHostnameClient(testZoneId)
	customHostnames, err := customHostnameClient.List("", 100, 1)
	for _, customHostname := range customHostnames.Result {
		flog.Infof(fmt.Sprintf("%s [%s] %s %s", customHostname.ID, customHostname.Status, customHostname.Hostname, customHostname.Ssl.Status))
	}

	// 验证自定义主机记录
	flog.Infof("=======================等待20秒后验证自定义主机记录=======================")
	time.Sleep(20 * time.Second)
	customHostnameDetial, err := customHostnameClient.Info(customHostnameId)
	flog.Infof(fmt.Sprintf("%s [%s] 主机状态：%s 证书状态：%s", customHostnameId, customHostnameDetial.Result.Hostname, customHostnameDetial.Result.Status, customHostnameDetial.Result.Ssl.Status))

	// 删除自定义主机及对应的验证DNS记录
	flog.Infof("=======================删除自定义主机=======================")
	success, err = client.DeleteCustomHostnameAndDns(testZoneId, testZoneId, hostName)
	assert.Truef(t, success, "删除自定义主机失败:%+v", err)

	// 验证dns解析记录被删除
	dnsDetail, _ = dnsClient.Info(dnsId)
	assert.Equal(t, 1, len(dnsDetail.Errors))
	assert.Equal(t, 81044, dnsDetail.Errors[0].Code)

	// 验证dcv记录被删除
	dcvDnsDetail, _ = dnsClient.Info(dcvDnsId)
	assert.Equal(t, 1, len(dcvDnsDetail.Errors))
	assert.Equal(t, 81044, dcvDnsDetail.Errors[0].Code)

	// 验证自定义主机记录被删除
	customHostnameDetial, err = customHostnameClient.Info(customHostnameId)
	assert.Equal(t, 1, len(customHostnameDetial.Errors))
	assert.Equal(t, 1436, customHostnameDetial.Errors[0].Code)

	flog.Infof("=======================添加dns记录=======================")
	success1, id1, _, err1 := dnsClient.Create("A", "test1."+testDomain, "127.0.0.1", false, 3600, "测试添加记录"+time.Now().Format("2006-01-02 15:04:05"), true)
	assert.Truef(t, success1, "添加dns记录失败:%+v", err1)

	dns1, _ := dnsClient.Info(id1)
	flog.Infof(fmt.Sprintf("%s [%s] %s %s", dns1.Result.ID, dns1.Result.Type, dns1.Result.Name, dns1.Result.Content))
	assert.Equal(t, "test1."+testDomain, dns1.Result.Name)
	assert.Equal(t, "127.0.0.1", dns1.Result.Content)

	flog.Infof("=======================修改dns记录=======================")
	success2, _, err2 := dnsClient.Update(id1, "CNAME", "test2."+testDomain, "www."+testDomain, true, 1800, "测试修改记录"+time.Now().Format("2006-01-02 15:04:05"))
	assert.Truef(t, success2, "修改dns记录失败:%+v", err2)

	dns2, _ := dnsClient.Info(id1)
	flog.Infof(fmt.Sprintf("%s [%s] %s %s", dns2.Result.ID, dns2.Result.Type, dns2.Result.Name, dns2.Result.Content))
	assert.Equal(t, "test2."+testDomain, dns2.Result.Name)
	assert.Equal(t, "www."+testDomain, dns2.Result.Content)

	flog.Infof("=======================删除dns记录=======================")
	dnsClient.Delete(id1)
	dns3, _ := dnsClient.Info(id1)
	assert.Equal(t, 1, len(dns3.Errors))
	assert.Equal(t, 81044, dns3.Errors[0].Code)
}
