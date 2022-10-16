package main

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

const path string = "ddns"

func main() {
	run()
}

func readOldIp() (f *os.File, oldIP string) {
	cache, err1 := ioutil.ReadAll(f)

	if err1 != nil {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("不能创建临时文件！err:", err)
			os.Exit(-1)
		}
		// 文件不存在
		return f, "0.0.0.0"
	} else {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("不能创建临时文件！err:", err)
			os.Exit(-1)
		}
		oldIP = string(cache)
		return f, oldIP
	}
}

func run() {
	// 获取配置文件
	cf, err := ioutil.ReadFile("config.json")
	if err != nil {
		curPath, _ := os.Getwd()
		fmt.Println("不能读取 config.json 文件,目录：", curPath)
	}

	con := &Config{}
	err = json.Unmarshal(cf, con)
	if err != nil {
		fmt.Println(err)
	}
	var MX uint64 = uint64(con.MX)
	var TTL uint64 = uint64(con.TTL)
	var RecordId uint64 = uint64(con.RecordId)

	// 识别当前 ip
	res, _ := http.Get("https://ipinfo.io/ip")
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	ipv4 := string(body)
	address := net.ParseIP(ipv4)
	if address == nil {
		fmt.Println("ip地址格式不正确")
	} else {
		fmt.Println("当前获取到的 IP：", address.String())
	}

	// 读取旧 ip
	f, oldIP := readOldIp()
	_, _ = f.WriteString(ipv4)
	defer f.Close()
	if oldIP == ipv4 {
		fmt.Println("不需要更新 IP！")
		os.Exit(0)
	}

	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(
		con.SecretId,
		con.SecretKey,
	)

	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := dnspod.NewClient(credential, "", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewModifyRecordRequest()

	request.Domain = common.StringPtr(con.Domain)
	request.SubDomain = common.StringPtr(con.SubDomain)
	request.RecordType = common.StringPtr(con.RecordType)
	request.RecordLine = common.StringPtr(con.RecordLine)
	request.Value = common.StringPtr(address.String())
	request.MX = common.Uint64Ptr(MX)
	request.TTL = common.Uint64Ptr(TTL)
	request.RecordId = common.Uint64Ptr(RecordId)

	// 返回的resp是一个ModifyRecordResponse的实例，与请求对象对应
	response, err := client.ModifyRecord(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}
