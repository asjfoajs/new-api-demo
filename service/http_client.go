package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"new-api-demo/constant"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

var (
	httpClient      *http.Client
	proxyClientLock sync.Mutex
	proxyClients    = make(map[string]*http.Client)
)

// checkRedirect SSRF防护校验，目前先不实现
func checkRedirect(req *http.Request, via []*http.Request) error {
	//fetchSetting := system_setting.GetFetchSetting()
	//urlStr := req.URL.String()
	//if err := common.ValidateURLWithFetchSetting(urlStr, fetchSetting.EnableSSRFProtection, fetchSetting.AllowPrivateIp, fetchSetting.DomainFilterMode, fetchSetting.IpFilterMode, fetchSetting.DomainList, fetchSetting.IpList, fetchSetting.AllowedPorts, fetchSetting.ApplyIPFilterForDomain); err != nil {
	//	return fmt.Errorf("redirect to %s blocked: %v", urlStr, err)
	//}
	if len(via) >= 10 {
		return fmt.Errorf("stopped after 10 redirects")
	}
	return nil
}

func InitHttpClient() {
	//if common.RelayTimeout == 0 {
	//	httpClient = &http.Client{
	//		CheckRedirect: checkRedirect,
	//	}
	//} else {
	//	httpClient = &http.Client{
	//		Timeout:       time.Duration(common.RelayTimeout) * time.Second,
	//		CheckRedirect: checkRedirect,
	//	}
	//}

	httpClient = &http.Client{
		Timeout:       time.Duration(constant.RelayTimeout) * time.Second,
		CheckRedirect: checkRedirect,
		Transport: &http.Transport{
			// MaxIdleConns 定义了所有主机（Host）共享的最大空闲连接数。
			// 即使你访问了多个不同的网站，总共留在内存中等待复用的连接不会超过这个数。
			MaxIdleConns: 100,

			// MaxIdleConnsPerHost 定义了“每个主机”允许保留的最大空闲连接数。
			// 这是最重要的参数：在高并发访问同一个 API 服务时，默认值通常很小（默认 2），
			// 将其调大（如 100）可以显著减少频繁创建 TCP 连接带来的延迟。
			MaxIdleConnsPerHost: 100,

			// IdleConnTimeout 定义了一个空闲连接在被自动关闭之前，可以在连接池中存放的最长时间。
			// 90 秒表示如果一个连接在 90 秒内没有被再次使用，它将被强制关闭以释放系统资源。
			IdleConnTimeout: 90 * time.Second,

			// DisableKeepAlives 控制是否禁用 HTTP 长连接。
			// 设置为 false 表示“开启 Keep-Alive”，即允许连接复用。
			// 在绝大多数生产场景下，为了性能都应该保持为 false。
			DisableKeepAlives: false,
		},
	}

}

func GetHttpClient() *http.Client {
	return httpClient
}

// ResetProxyClientCache 清空代理客户端缓存，确保下次使用时重新初始化
func ResetProxyClientCache() {
	proxyClientLock.Lock()
	defer proxyClientLock.Unlock()
	for _, client := range proxyClients {
		if transport, ok := client.Transport.(*http.Transport); ok && transport != nil {
			transport.CloseIdleConnections()
		}
	}
	proxyClients = make(map[string]*http.Client)
}

// NewProxyHttpClient 创建支持代理的 HTTP 客户端
func NewProxyHttpClient(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		return http.DefaultClient, nil
	}

	proxyClientLock.Lock()
	if client, ok := proxyClients[proxyURL]; ok {
		proxyClientLock.Unlock()
		return client, nil
	}
	proxyClientLock.Unlock()

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	switch parsedURL.Scheme {
	case "http", "https":
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(parsedURL),
			},
			CheckRedirect: checkRedirect,
		}
		client.Timeout = time.Duration(constant.RelayTimeout) * time.Second
		proxyClientLock.Lock()
		proxyClients[proxyURL] = client
		proxyClientLock.Unlock()
		return client, nil

	case "socks5", "socks5h":
		// 获取认证信息
		var auth *proxy.Auth
		if parsedURL.User != nil {
			auth = &proxy.Auth{
				User:     parsedURL.User.Username(),
				Password: "",
			}
			if password, ok := parsedURL.User.Password(); ok {
				auth.Password = password
			}
		}

		// 创建 SOCKS5 代理拨号器
		// proxy.SOCKS5 使用 tcp 参数，所有 TCP 连接包括 DNS 查询都将通过代理进行。行为与 socks5h 相同
		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
		if err != nil {
			return nil, err
		}

		client := &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dialer.Dial(network, addr)
				},
			},
			CheckRedirect: checkRedirect,
		}
		client.Timeout = time.Duration(constant.RelayTimeout) * time.Second
		proxyClientLock.Lock()
		proxyClients[proxyURL] = client
		proxyClientLock.Unlock()
		return client, nil

	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s, must be http, https, socks5 or socks5h", parsedURL.Scheme)
	}
}
