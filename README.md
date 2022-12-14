# DDNS

这是一个根据当前公网 IP 动态修改腾讯云域名记录的工具，如果宽带不是专线，那么 IP 会隔几天变化一次，这样给域名绑定带来麻烦。

此工具通过动态识别当前网络的公网 IP 地址，将变化的 IP 推送到腾讯云已有的域名解析中，实现域名动态绑定 IP。



### 下载

打开：https://github.com/whuanle/ddns-tencent/releases

找到对应操作系统的二进制程序文件。



### 配置文件

需要在**运行程序的目录(执行命令的目录，而不是程序所在目录)**下，创建 `config.json` 文件，内容格式如下：

```json
{
  "SecretId": "id",
  "SecretKey": "密钥",
  "Domain":"666.cn",
  "SubDomain": "域名前缀",
  "RecordType": "记录类型，如 A、NX",
  "RecordLine": "线路名称",
  "Value": "123.123.123.123",
  "MX": 5,
  "TTL": 600,
  "RecordId": 1220273909
}
```



注意，程序使用 `https://ipinfo.io/ip` 检测当前的公网 IP 地址，可能有些慢，建议替换一下国内的其他工具地址。



### 获取腾讯云 API 密钥

打开：https://console.cloud.tencent.com/cam/capi

创建或获取访问腾讯云 API 的密钥，**建议使用子账号**！

![image-20221016163142281](images/image-20221016163142281.png)



给子账号开通 `QcloudCollApiKeyManageAccess`、`QcloudDNSPodFullAccess ` 两个权限。



![image-20221016163312790](images/image-20221016163312790.png)



然后使用子账号登录控制台，打开 https://console.cloud.tencent.com/cam/capi

新建密钥：

![image-20221016163459244](images/image-20221016163459244.png)

复制 `SecretId` 和 `SecretKey`，存储到 `config.json` 中。



### 获取域名信息

打开 https://console.cloud.tencent.com/cns

找到自己的域名，先添加一个域名解析。

![image-20221016163729948](images/image-20221016163729948.png)



接下来要获取此解析记录的 id，即 `RecordId`。



方法①：

然后点击修改域名记录，但是先不保存，按下 F12，在点击保存。

在浏览器控制台中，可以看到一条请求：

```
https://wss.cloud.tencent.com/dns/api/record/update?g_tk={此解析记录的id}
```



根据这个请求，复制后面的记录 id。



方法②：

打开：https://console.cloud.tencent.com/api/explorer?Product=dnspod&Version=2021-03-23&Action=DescribeRecordList

![image-20221016164843962](images/image-20221016164843962.png)



将 `RecordId` 放到 `config.json` 中。



### 配置说明

经过以上步骤，目前还有以下配置需要修改：

```
  "Domain":"666.cn",
  "SubDomain": "域名前缀",
  "RecordType": "记录类型，如 A、NX",
  "RecordLine": "线路名称",
  "Value": "123.123.123.123",
  "MX": 5,
  "TTL": 600,
```



其实就是对应域名解析记录的，此程序只会动态修改 IP，其他修改内容会按照 `config.json` 的配置去区配解析记录并修改。

也可以看这里了解每个参数的说明:

https://docs.dnspod.cn/api/modify-records/





`RecordLine` 这个参数，指的是线路名称，可以填 `默认` 。

![image-20221016165309168](images/image-20221016165309168.png)

![image-20221016165319388](images/image-20221016165319388.png)



### Linux 定时任务

先将程序复制放到 Linux 的目录下，程序可以放在任意目录，但是 `config.json` 需要放在 root 目录下，这是因为 Linux 的 Cron  运行目录是 root 或其他用户目录。

打开 `/etc/cron.d` 目录，创建一个新文件 `ddns`，文件内容如下：

```
* * * * * root /root/ddns
* * * * * root sleep 10; /root/ddns
* * * * * root sleep 20; /root/ddns
* * * * * root sleep 30; /root/ddns
* * * * * root sleep 40; /root/ddns
* * * * * root sleep 50; /root/ddns
```



`* * * * * ` 表示每分钟执行一次；

`root` 表示以什么用户启动程序；

`/root/ddns` 程序目录位置；

`sleep 20;` 休眠时间。



因为 Linux Cron 的粒度是每分钟，因此如果需要每 10s 执行一次脚本的话，需要设置多条记录，使用 `sleep` 延迟执行。

然后执行 `service cron reload` ，刷新定时任务。



因为定时任务看不到程序日志，因此可以改成：

```
* * * * * root /root/ddns >> /tmp/ddns.log
* * * * * root sleep 10; /root/ddns >> /tmp/ddns.log
* * * * * root sleep 20; /root/ddns >> /tmp/ddns.log
* * * * * root sleep 30; /root/ddns >> /tmp/ddns.log
* * * * * root sleep 40; /root/ddns >> /tmp/ddns.log
* * * * * root sleep 50; /root/ddns >> /tmp/ddns.log
```

