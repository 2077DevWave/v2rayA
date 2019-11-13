package tools

import (
	"V2RayA/global"
	"V2RayA/models/touch"
	"V2RayA/models/v2ray"
	"V2RayA/proxyWithHttp"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func GetHttpClientWithProxy(proxyURL string) (client *http.Client, err error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return
	}
	dialer, err := proxyWithHttp.FromURL(u, proxyWithHttp.Direct)
	if err != nil {
		return
	}
	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial
	client = &http.Client{Transport: httpTransport}
	return
}

func GetHttpClientWithV2RayAProxy() (client *http.Client, err error) {
	host := "localhost"
	//是否在docker环境
	if global.ServiceControlMode == v2ray.Docker {
		//连接网关，即宿主机的端口，失败则用同网络下v2ray容器的
		out, err := exec.Command("sh", "-c", "ip route|grep default|awk '{print $3}'").Output()
		if err == nil {
			host = strings.TrimSpace(string(out))
		} else {
			host = "v2ray"
		}
	}
	return GetHttpClientWithProxy("socks5://" + host + ":20170")
}

func GetHttpClientWithV2RayAPac() (client *http.Client, err error) {
	host := "localhost"
	//是否在docker环境
	if global.ServiceControlMode == v2ray.Docker {
		//连接网关，即宿主机的端口，失败则用同网络下v2ray容器的
		out, err := exec.Command("sh", "-c", "ip route|grep default|awk '{print $3}'").Output()
		if err == nil {
			host = strings.TrimSpace(string(out))
		} else {
			host = "v2ray"
		}
	}
	return GetHttpClientWithProxy("http://" + host + ":20172")
}

func GetHttpClientAutomatically() (c *http.Client, err error) {
	tr := global.GetTouchRaw()
	if tr.ConnectedServer == nil {
		return http.DefaultClient, nil
	}
	switch tr.Setting.ProxyModeWhenSubscribe {
	case touch.ProxyModePac:
		c, err = GetHttpClientWithV2RayAPac()
	case touch.ProxyModeProxy:
		c, err = GetHttpClientWithV2RayAProxy()
	default:
		c = http.DefaultClient
	}
	return
}
