package main

	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/siyuan-note/logging"
	"github.com/siyuan-note/siyuan/kernel/util"
		addr := host + ":" + util.FixedPort

		// 启动一个固定 6806 端口的反向代理服务器，这样浏览器扩展才能直接使用 127.0.0.1:6806，不用配置端口
		proxy := newFixedPortReverseProxy(util.ServerURL)

		if "" != certPath {
			logging.LogInfof("fixed port service [%s] is running (HTTP/HTTPS dual mode)", addr)
		logging.LogInfof("fixed port service [%s] is stopped", addr)
	}
}

func newFixedPortReverseProxy(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(request *httputil.ProxyRequest) {
			request.SetURL(target)
			request.Out.Host = request.In.Host
			request.SetXForwarded()
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
	}

	//logging.LogInfof("check auth for [%s]", c.Request.RequestURI)
	localhost := IsLocalRequest(c)

	// 未设置锁屏密码
	if "" == Conf.AccessAuthCode {
	c.Next()
}

// IsLocalRequest 判断请求是否由本机客户端直接发起或经可信本机代理转发。
func IsLocalRequest(c *gin.Context) bool {
	// 仅当直接连接和可信代理解析出的原始客户端均为环回地址时，才视为本机请求。
	return util.IsLocalHost(c.Request.RemoteAddr) && util.IsLocalHostname(c.ClientIP())
}

func CheckAdminRole(c *gin.Context) {
	if IsAdminRoleContext(c) {
		c.Next()
func Serve(fastMode bool, cookieKey string) {
	gin.SetMode(gin.ReleaseMode)
	ginServer := gin.New()
	if err := ginServer.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		logging.LogFatalf(logging.ExitCodeSecurityRisk, "set trusted proxies failed: %s", err)
	}
	ginServer.ForwardedByClientIP = true
	ginServer.RemoteIPHeaders = []string{"X-Forwarded-For"}
	ginServer.MaxMultipartMemory = 1024 * 1024 * 32 // 插入较大的资源文件时内存占用较大 https://github.com/siyuan-note/siyuan/issues/5023
	ginServer.Use(
		model.ControlConcurrency, // 请求串行化 Concurrency control when requesting the kernel API https://github.com/siyuan-note/siyuan/issues/9939
