package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/smartwalle/alipay/v3"
)

var (
	// appId
	appId = "2016101800712387"
	// 应用公钥
	//aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArMR10yQY4067dFE4pbJWRcz2v1BYqVDNKb5T63cPMLkIX5s/Z5XvBn1sAvpoL74lJBzYCi6xcUJ1PWKKkkqOckWJ3urYe7MCaMjb1Hgko0PyShmqUuE79pqaoqGWx0p6KZjDfKHLgNPsBFYZ9YoKpz8Ef0oPE+prlwhDBSoNRCVcVNg+zrs9PRFiWCscD/qgo+j4cve4oGW6Vt7A8loKNk0nlRaV5v/haq8CfsSdQSrSr00tK8z7RZYGPd8490ay+QRNeo7bsxtalH0FZ5V+GOrkfED8Hx/9jo32TXZjhO+EMMgnxTSOQkVaefEm3HcHGFMVu+0GsNwr7MKSSA6fNQIDAQAB"
	// 应用私钥
	privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCCbFREyM5aT6FCgmtpjGhiRY6HaVy3HR7IrKe6rb/ZKpUw96v9w1vJLK/MuMzdqCZDj00pcgaeP0Hfz970JUZzDXIxoMhFkzkwELf7SCkmHHm98KmTZmtfKZ2UmwLq7MD3C4zg76Sq9OXWnV2Mlu4lHm2L5YP3UJeN8K/x9Xmt97mVXbLleo5chXJF5Fn9mt7qBF3nzxyR05iDm6m53yvOB1FAw1t1d3XREnXikJ1Jfao1PaRetLFEbTxjhCiKBwMXLEq+ElzaejFwdFj2mhiB3RH4FeMlqb2XdBFkalszCe8HbNgjSUjLEHMP25Kd7rBuLsbybttV2VNGA934eDBLAgMBAAECggEAeAJWSDS+LDu5lwmK7MRWCJcYX2GD/rq86Mp2TEA0mA+m2DgN+qMYgjSsTyR83jkrfF8QEWLq1BJA+0C7Vsw11wg4W5Imtl6R8LhxkrcDph3tp+kbIJbNlfMlHOjF9oWDyc7HAvOAAg72rlR/EC6RU3Z8No4qsBdoXSNrwqiuHnTRO6tkNHxAWiTMVnztNmSGJURXqh6SAS1+H06c4jnJg/2PVIQwUFjo7Z5i2hKvGINBQ5RS57AZkgmbfjxNkglS7GD01rhl5mGCh33Otq4pcMF70FLbJrpVrst8yscTb02fyOJ6d84Xn6K0pAO7zuGza909qBrzDPqIcmFzvuPAAQKBgQDPmCi2T+DRDYXmSfvnKG/+waF+nmZPUTSJQAR4l/taE+LpE8nCzQ8wBvvyBDP59AJXjklPUuClgYw37yrGhGqfXjJFbzuPxr9U7eFhDGexy/AtVJfrfbvBzRlNmIKH1E+N9eMp1zVydvpG6i60ePO23sDLQsgAwPj6N+VJ/kMeAQKBgQCg1Z9alvUBv1H6LKMVfjginMLKtndr66cCcLR7kmjC/1NCZuXg1xlOqNfAsPdTnP6B9xRL0182v0P9smWx1oLlxR5JRA0V5SvqH9A1TDjr1+h+IsjWVVAC0M1ebfPJoUmhs/+3SYcsDmwTh9sOSlkscdAF2hEa8D4O1DUaMNpmSwKBgApaCA2vRgKmrfqhzdHlDlChzy/FLkzeO8RsUMzCp2ICg9ojhngUSaGXd5DF7OGV7Vf4XGd8Nn+KSjev0W48xCRWSiN0PIAa5QeTJR31xGX1SXC5OyofBvHPDGf2JuwnBiCKFl3LwXqHvEs0+kc9kMmZqft4xQhklwXDK8fYyfgBAoGBAIUILYcA1idb2LLVuQ9OF6CJiZWi16Sshre+AYs0zvJ7vqJt+ja/tG8buVnpBqpicSGO/Xq6m0btbY+qv/MZO6xSH3r6jthNdsVxCwcKxQpOzD+JBhZC+qtZioVQ7RUaE41tFVbFusj2JO8CsG5hkODyQt6UQRHHJY2eeU3wmrWBAoGBAKcazwXk/bN9g8IB77aErEE48X4MChwP8zZVwaQs1yvkSQcgFy7FkKMi74W5fTK+vLgylAqW1ZwEtV4KF6MgYuqd12UrW8mh9B0wz13cMoYb1WalqEykevCXhRw85tUqw9AjMbTRxwTToL5Lz36orAo6hoNA5n4K/EnO5/B//iMp"
	Client, _  = alipay.New(appId, privateKey, false)
)

func init() {
	Client.LoadAppPublicCertFromFile("./conf/app_public.crt")
	Client.LoadAliPayPublicCertFromFile("./conf/public.txt")
	Client.LoadAliPayRootCertFromFile("/home/lixu/git/golang/src/mn-hosted/conf/alipay/alipayRootCert.crt")
}

func VerifySignAlipay(req *http.Request) (bool, error) {
	return Client.VerifySign(req.Form)
}

func WebPageAlipay(orderid int64, price int32) (string, error) {
	pay := alipay.TradePagePay{}
	// 支付宝回调地址（需要在支付宝后台配置）
	// 支付成功后，支付宝会发送一个POST消息到该地址
	pay.NotifyURL = "http://pay.vpubchain.cn:8088/alipay"
	// 支付成功之后，浏览器将会重定向到该 URL
	pay.ReturnURL = "http://pay.vpubchain.cn:8088/return"
	//支付标题
	pay.Subject = "支付宝支付测试"

	//订单号，一个订单号只能支付一次
	pay.OutTradeNo = strconv.FormatInt(orderid, 10)
	fmt.Println("tradeNo: ", pay.OutTradeNo)
	//销售产品码，与支付宝签约的产品码名称,目前仅支持FAST_INSTANT_TRADE_PAY
	pay.ProductCode = "FAST_INSTANT_TRADE_PAY"
	//金额
	var amount = float64(price) / 100
	pay.TotalAmount = strconv.FormatFloat(amount, 'g', 1, 64)

	fmt.Println("amount=", pay.TotalAmount)

	url, err := Client.TradePagePay(pay)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	payURL := url.String()
	//这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	fmt.Println(payURL)

	return payURL, nil
	//打开默认浏览器
	//payURL = strings.Replace(payURL, "&", "^&", -1)
	//exec.Command("cmd", "/c", "start", payURL).Start()
}
