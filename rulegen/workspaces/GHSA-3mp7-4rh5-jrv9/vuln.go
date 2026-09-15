package main

	"net"
	"net/http"
	"net/http/httputil"

	"github.com/siyuan-note/logging"
	"github.com/siyuan-note/siyuan/kernel/util"
		addr := host + ":" + util.FixedPort

		// 启动一个固定 6806 端口的反向代理服务器，这样浏览器扩展才能直接使用 127.0.0.1:6806，不用配置端口
		proxy := httputil.NewSingleHostReverseProxy(util.ServerURL)
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		if "" != certPath {
			logging.LogInfof("fixed port service [%s] is running (HTTP/HTTPS dual mode)", addr)
		logging.LogInfof("fixed port service [%s] is stopped", addr)
	}
}
	}

	//logging.LogInfof("check auth for [%s]", c.Request.RequestURI)
	localhost := util.IsLocalHost(c.Request.RemoteAddr)

	// 未设置锁屏密码
	if "" == Conf.AccessAuthCode {
	c.Next()
}

func CheckAdminRole(c *gin.Context) {
	if IsAdminRoleContext(c) {
		c.Next()
func Serve(fastMode bool, cookieKey string) {
	gin.SetMode(gin.ReleaseMode)
	ginServer := gin.New()
	ginServer.MaxMultipartMemory = 1024 * 1024 * 32 // 插入较大的资源文件时内存占用较大 https://github.com/siyuan-note/siyuan/issues/5023
	ginServer.Use(
		model.ControlConcurrency, // 请求串行化 Concurrency control when requesting the kernel API https://github.com/siyuan-note/siyuan/issues/9939
