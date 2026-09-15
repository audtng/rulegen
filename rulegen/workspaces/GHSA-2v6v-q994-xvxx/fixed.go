package main

	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/context/param"
		"LOCK":      true,
		"UNLOCK":    true,
	}
	// these web.Controller's methods shouldn't reflect to AutoRouter
	// see registerControllerExceptMethods
	exceptMethod = initExceptMethod()

	urlPlaceholder = "{{placeholder}}"
	// DefaultAccessLogFilter will skip the accesslog if return true
}

// default log filter static file will not show
type logFilter struct{}

func (l *logFilter) Filter(ctx *beecontext.Context) bool {
	requestPath := path.Clean(ctx.Request.URL.Path)
	exceptMethod = append(exceptMethod, action)
}

func initExceptMethod() []string {
	res := make([]string, 0, 32)
	c := &Controller{}
	t := reflect.TypeOf(c)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		res = append(res, m.Name)
	}
	return res
}

// ControllerInfo holds information about the controller.
type ControllerInfo struct {
	pattern        string
	controllerType reflect.Type
	methods        map[string]string
	handler        http.Handler
	runFunction    HandleFunc
	routerType     int
	initialize     func() ControllerInterface
	methodParams   []*param.MethodParam
	sessionOn      bool
}

type ControllerOption func(*ControllerInfo)

func (c *ControllerInfo) GetPattern() string {
	return c.pattern
}

func WithRouterMethods(ctrlInterface ControllerInterface, mappingMethod ...string) ControllerOption {
	return func(c *ControllerInfo) {
		c.methods = parseMappingMethods(ctrlInterface, mappingMethod)
	}
}

func WithRouterSessionOn(sessionOn bool) ControllerOption {
	return func(c *ControllerInfo) {
		c.sessionOn = sessionOn
	}
}

type filterChainConfig struct {
	pattern string
	chain   FilterChain
	opts    []FilterOpt
}

// ControllerRegister containers registered router rules, controller handlers and filters.
type ControllerRegister struct {
	routers      map[string]*Tree
	enablePolicy bool
	enableFilter bool
	policies     map[string]*Tree
	filters      [FinishRouter + 1][]*FilterRouter
	pool         sync.Pool

	// the filter created by FilterChain
	chainRoot *FilterRouter

	// keep registered chain and build it when serve http
	filterChains []filterChainConfig

	cfg *Config
}

				return beecontext.NewContext()
			},
		},
		cfg:          cfg,
		filterChains: make([]filterChainConfig, 0, 4),
	}
	res.chainRoot = newFilterRouter("/*", res.serveHttp, WithCaseSensitive(false))
	return res
}

// Init will be executed when HttpServer start running
func (p *ControllerRegister) Init() {
	for i := len(p.filterChains) - 1; i >= 0; i-- {
		fc := p.filterChains[i]
		root := p.chainRoot
		filterFunc := fc.chain(func(ctx *beecontext.Context) {
			var preFilterParams map[string]string
			root.filter(ctx, p.getUrlPath(ctx), preFilterParams)
		})
		p.chainRoot = newFilterRouter(fc.pattern, filterFunc, fc.opts...)
		p.chainRoot.next = root
	}
}

// Add controller handler and pattern rules to ControllerRegister.
// usage:
//	default methods is the same name as method
//	Add("/api/delete",&RestController{},"delete:DeleteFood")
//	Add("/api",&RestController{},"get,post:ApiFunc"
//	Add("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
func (p *ControllerRegister) Add(pattern string, c ControllerInterface, opts ...ControllerOption) {
	p.addWithMethodParams(pattern, c, nil, opts...)
}

func parseMappingMethods(c ControllerInterface, mappingMethods []string) map[string]string {
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	methods := make(map[string]string)

	if len(mappingMethods) == 0 {
		return methods
	}

	semi := strings.Split(mappingMethods[0], ";")
	for _, v := range semi {
		colon := strings.Split(v, ":")
		if len(colon) != 2 {
			panic("method mapping format is invalid")
		}
		comma := strings.Split(colon[0], ",")
		for _, m := range comma {
			if m != "*" && !HTTPMETHOD[strings.ToUpper(m)] {
				panic(v + " is an invalid method mapping. Method doesn't exist " + m)
			}
			if val := reflectVal.MethodByName(colon[1]); val.IsValid() {
				methods[strings.ToUpper(m)] = colon[1]
				continue
			}
			panic("'" + colon[1] + "' method doesn't exist in the controller " + t.Name())
		}
	}

	return methods
}

func (p *ControllerRegister) addRouterForMethod(route *ControllerInfo) {
	if len(route.methods) == 0 {
		for m := range HTTPMETHOD {
			p.addToRouter(m, route.pattern, route)
		}
		return
	}
	for k := range route.methods {
		if k != "*" {
			p.addToRouter(k, route.pattern, route)
			continue
		}
		for m := range HTTPMETHOD {
			p.addToRouter(m, route.pattern, route)
		}
	}
}

func (p *ControllerRegister) addWithMethodParams(pattern string, c ControllerInterface, methodParams []*param.MethodParam, opts ...ControllerOption) {
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()

	route := p.createBeegoRouter(t, pattern)
	route.initialize = func() ControllerInterface {
		vc := reflect.New(route.controllerType)
		execController, ok := vc.Interface().(ControllerInterface)

		return execController
	}
	route.methodParams = methodParams
	for i := range opts {
		opts[i](route)
	}

	globalSessionOn := p.cfg.WebConfig.Session.SessionOn
	if !globalSessionOn && route.sessionOn {
		logs.Warn("global sessionOn is false, sessionOn of router [%s] can't be set to true", route.pattern)
		route.sessionOn = globalSessionOn
	}

	p.addRouterForMethod(route)
}

func (p *ControllerRegister) addToRouter(method, pattern string, r *ControllerInfo) {
				for _, f := range a.Filters {
					p.InsertFilter(f.Pattern, f.Pos, f.Filter, WithReturnOnOutput(f.ReturnOnOutput), WithResetParams(f.ResetParams))
				}
				p.addWithMethodParams(a.Router, c, a.MethodParams, WithRouterMethods(c, strings.Join(a.AllowHTTPMethods, ",")+":"+a.Method))
			}
		}
	}
	p.pool.Put(ctx)
}

// CtrlGet add get method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlGet("/api/:id", MyController.Ping)
// If the receiver of function Ping is pointer, you should use CtrlGet("/api/:id", (*MyController).Ping)
func (p *ControllerRegister) CtrlGet(pattern string, f interface{}) {
	p.AddRouterMethod(http.MethodGet, pattern, f)
}

// CtrlPost add post method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlPost("/api/:id", MyController.Ping)
// If the receiver of function Ping is pointer, you should use CtrlPost("/api/:id", (*MyController).Ping)
func (p *ControllerRegister) CtrlPost(pattern string, f interface{}) {
	p.AddRouterMethod(http.MethodPost, pattern, f)
}

// CtrlHead add head method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlHead("/api/:id", MyController.Ping)
// If the receiver of function Ping is pointer, you should use CtrlHead("/api/:id", (*MyController).Ping)
func (p *ControllerRegister) CtrlHead(pattern string, f interface{}) {
	p.AddRouterMethod(http.MethodHead, pattern, f)
}

// CtrlPut add put method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlPut("/api/:id", MyController.Ping)

func (p *ControllerRegister) CtrlPut(pattern string, f interface{}) {
	p.AddRouterMethod(http.MethodPut, pattern, f)
}

// CtrlPatch add patch method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlPatch("/api/:id", MyController.Ping)
func (p *ControllerRegister) CtrlPatch(pattern string, f interface{}) {
	p.AddRouterMethod(http.MethodPatch, pattern, f)
}

// CtrlDelete add delete method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlDelete("/api/:id", MyController.Ping)
func (p *ControllerRegister) CtrlDelete(pattern string, f interface{}) {
	p.AddRouterMethod(http.MethodDelete, pattern, f)
}

// CtrlOptions add options method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlOptions("/api/:id", MyController.Ping)
func (p *ControllerRegister) CtrlOptions(pattern string, f interface{}) {
	p.AddRouterMethod(http.MethodOptions, pattern, f)
}

// CtrlAny add all method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlAny("/api/:id", MyController.Ping)
func (p *ControllerRegister) CtrlAny(pattern string, f interface{}) {
	p.AddRouterMethod("*", pattern, f)
}

// AddRouterMethod add http method router
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    AddRouterMethod("get","/api/:id", MyController.Ping)
func (p *ControllerRegister) AddRouterMethod(httpMethod, pattern string, f interface{}) {
	httpMethod = p.getUpperMethodString(httpMethod)
	ct, methodName := getReflectTypeAndMethod(f)

	p.addBeegoTypeRouter(ct, methodName, httpMethod, pattern)
}

// addBeegoTypeRouter add beego type router
func (p *ControllerRegister) addBeegoTypeRouter(ct reflect.Type, ctMethod, httpMethod, pattern string) {
	route := p.createBeegoRouter(ct, pattern)
	methods := p.getHttpMethodMapMethod(httpMethod, ctMethod)
	route.methods = methods

	p.addRouterForMethod(route)
}

// createBeegoRouter create beego router base on reflect type and pattern
func (p *ControllerRegister) createBeegoRouter(ct reflect.Type, pattern string) *ControllerInfo {
	route := &ControllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeBeego
	route.sessionOn = p.cfg.WebConfig.Session.SessionOn
	route.controllerType = ct
	return route
}

// createRestfulRouter create restful router with filter function and pattern
func (p *ControllerRegister) createRestfulRouter(f HandleFunc, pattern string) *ControllerInfo {
	route := &ControllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeRESTFul
	route.sessionOn = p.cfg.WebConfig.Session.SessionOn
	route.runFunction = f
	return route
}

// createHandlerRouter create handler router with handler and pattern
func (p *ControllerRegister) createHandlerRouter(h http.Handler, pattern string) *ControllerInfo {
	route := &ControllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeHandler
	route.sessionOn = p.cfg.WebConfig.Session.SessionOn
	route.handler = h
	return route
}

// getHttpMethodMapMethod based on http method and controller method, if ctMethod is empty, then it will
// use http method as the controller method
func (p *ControllerRegister) getHttpMethodMapMethod(httpMethod, ctMethod string) map[string]string {
	methods := make(map[string]string)
	// not match-all sign, only add for the http method
	if httpMethod != "*" {

		if ctMethod == "" {
			ctMethod = httpMethod
		}
		methods[httpMethod] = ctMethod
		return methods
	}

	// add all http method
	for val := range HTTPMETHOD {
		if ctMethod == "" {
			methods[val] = val
		} else {
			methods[val] = ctMethod
		}
	}
	return methods
}

// getUpperMethodString get upper string of method, and panic if the method
// is not valid
func (p *ControllerRegister) getUpperMethodString(method string) string {
	method = strings.ToUpper(method)
	if method != "*" && !HTTPMETHOD[method] {
		panic("not support http method: " + method)
	}
	return method
}

// get reflect controller type and method by controller method expression
func getReflectTypeAndMethod(f interface{}) (controllerType reflect.Type, method string) {
	// check f is a function
	funcType := reflect.TypeOf(f)
	if funcType.Kind() != reflect.Func {
		panic("not a method")
	}

	// get function name
	funcObj := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	if funcObj == nil {
		panic("cannot find the method")
	}
	funcNameSli := strings.Split(funcObj.Name(), ".")
	lFuncSli := len(funcNameSli)
	if lFuncSli == 0 {
		panic("invalid method full name: " + funcObj.Name())
	}

	method = funcNameSli[lFuncSli-1]
	if len(method) == 0 {
		panic("method name is empty")
	} else if method[0] > 96 || method[0] < 65 {
		panic(fmt.Sprintf("%s is not a public method", method))
	}

	// check only one param which is the method receiver
	if numIn := funcType.NumIn(); numIn != 1 {
		panic("invalid number of param in")
	}

	controllerType = funcType.In(0)

	// check controller has the method
	_, exists := controllerType.MethodByName(method)
	if !exists {
		panic(controllerType.String() + " has no method " + method)
	}

	// check the receiver implement ControllerInterface
	if controllerType.Kind() == reflect.Ptr {
		controllerType = controllerType.Elem()
	}
	controller := reflect.New(controllerType)
	_, ok := controller.Interface().(ControllerInterface)
	if !ok {
		panic(controllerType.String() + " is not implemented ControllerInterface")
	}

	return
}

// HandleFunc define how to process the request
type HandleFunc func(ctx *beecontext.Context)

// Get add get method
// usage:
//    Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Get(pattern string, f HandleFunc) {
	p.AddMethod("get", pattern, f)
}

//    Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Post(pattern string, f HandleFunc) {
	p.AddMethod("post", pattern, f)
}

//    Put("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Put(pattern string, f HandleFunc) {
	p.AddMethod("put", pattern, f)
}

//    Delete("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Delete(pattern string, f HandleFunc) {
	p.AddMethod("delete", pattern, f)
}

//    Head("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Head(pattern string, f HandleFunc) {
	p.AddMethod("head", pattern, f)
}

//    Patch("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Patch(pattern string, f HandleFunc) {
	p.AddMethod("patch", pattern, f)
}

//    Options("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Options(pattern string, f HandleFunc) {
	p.AddMethod("options", pattern, f)
}

//    Any("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Any(pattern string, f HandleFunc) {
	p.AddMethod("*", pattern, f)
}

//    AddMethod("get","/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) AddMethod(method, pattern string, f HandleFunc) {
	method = p.getUpperMethodString(method)

	route := p.createRestfulRouter(f, pattern)
	methods := p.getHttpMethodMapMethod(method, "")
	route.methods = methods

	p.addRouterForMethod(route)
}

// Handler add user defined Handler
func (p *ControllerRegister) Handler(pattern string, h http.Handler, options ...interface{}) {
	route := p.createHandlerRouter(h, pattern)
	if len(options) > 0 {
		if _, ok := options[0].(bool); ok {
			pattern = path.Join(pattern, "?:all(.*)")
}

// AddAuto router to ControllerRegister.
// example beego.AddAuto(&MainController{}),
// MainController has method List and Page.
// visit the url /main/list to execute List function
// /main/page to execute Page function.
}

// AddAutoPrefix Add auto router to ControllerRegister with prefix.
// example beego.AddAutoPrefix("/admin",&MainController{}),
// MainController has method List and Page.
// visit the url /admin/main/list to execute List function
// /admin/main/page to execute Page function.
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	for i := 0; i < rt.NumMethod(); i++ {
		methodName := rt.Method(i).Name
		if !utils.InSlice(methodName, exceptMethod) {
			p.addAutoPrefixMethod(prefix, controllerName, methodName, ct)
		}
	}
}

func (p *ControllerRegister) addAutoPrefixMethod(prefix, controllerName, methodName string, ctrl reflect.Type) {
	pattern := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(methodName), "*")
	patternInit := path.Join(prefix, controllerName, methodName, "*")
	patternFix := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(methodName))
	patternFixInit := path.Join(prefix, controllerName, methodName)

	route := p.createBeegoRouter(ctrl, pattern)
	route.methods = map[string]string{"*": methodName}
	for m := range HTTPMETHOD {

		p.addToRouter(m, pattern, route)

		// only case sensitive, we add three more routes
		if p.cfg.RouterCaseSensitive {
			p.addToRouter(m, patternInit, route)
			p.addToRouter(m, patternFix, route)
			p.addToRouter(m, patternFixInit, route)
		}
	}
}
//     }
// }
func (p *ControllerRegister) InsertFilterChain(pattern string, chain FilterChain, opts ...FilterOpt) {
	opts = append([]FilterOpt{WithCaseSensitive(p.cfg.RouterCaseSensitive)}, opts...)
	p.filterChains = append(p.filterChains, filterChainConfig{
		pattern: pattern,
		chain:   chain,
		opts:    opts,
	})
}

// add Filter into
	for _, l := range t.leaves {
		if c, ok := l.runObject.(*ControllerInfo); ok {
			if c.routerType == routerTypeBeego &&
				strings.HasSuffix(path.Join(c.controllerType.PkgPath(), c.controllerType.Name()), `/`+controllerName) {
				find := false
				if HTTPMETHOD[strings.ToUpper(methodName)] {
					if len(c.methods) == 0 {

// Implement http.Handler interface.
func (p *ControllerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := p.GetContext()

	ctx.Reset(rw, r)
	r := ctx.Request
	rw := ctx.ResponseWriter.ResponseWriter
	var (
		runRouter        reflect.Type
		findRouter       bool
		runMethod        string
		methodParams     []*param.MethodParam
		routerInfo       *ControllerInfo
		isRunnable       bool
		currentSessionOn bool
		originRouterInfo *ControllerInfo
		originFindRouter bool
	)

	if p.cfg.RecoverFunc != nil {
	}

	// session init
	currentSessionOn = p.cfg.WebConfig.Session.SessionOn
	originRouterInfo, originFindRouter = p.FindRouter(ctx)
	if originFindRouter {
		currentSessionOn = originRouterInfo.sessionOn
	}
	if currentSessionOn {
		ctx.Input.CruSession, err = GlobalSessions.SessionStart(rw, r)
		if err != nil {
			logs.Error(err)

// FindRouter Find Router info for URL
func (p *ControllerRegister) FindRouter(context *beecontext.Context) (routerInfo *ControllerInfo, isFind bool) {
	urlPath := context.Input.URL()
	if !p.cfg.RouterCaseSensitive {
		urlPath = strings.ToLower(urlPath)
	}
import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"github.com/beego/beego/v2"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/utils"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

// Config is the main struct for BConfig
// TODO after supporting multiple servers, remove common config to somewhere else
type Config struct {
	// AppName
	// @Description Application's name. You'd better set it because we use it to do some logging and tracing
	// @Default beego
	AppName string // Application name
	// RunMode
	// @Description it's the same as environment. In general, we have different run modes.
	// For example, the most common case is using dev, test, prod three environments
	// when you are developing the application, you should set it as dev
	// when you completed coding and want QA to test your code, you should deploy your application to test environment
	// and the RunMode should be set as test
	// when you completed all tests, you want to deploy it to prod, you should set it to prod
	// You should never set RunMode="dev" when you deploy the application to prod
	// because Beego will do more things which need Go SDK and other tools when it found out the RunMode="dev"
	// @Default dev
	RunMode string // Running Mode: dev | prod

	// RouterCaseSensitive
	// @Description If it was true, it means that the router is case sensitive.
	// For example, when you register a router with pattern "/hello",
	// 1. If this is true, and the request URL is "/Hello", it won't match this pattern
	// 2. If this is false and the request URL is "/Hello", it will match this pattern
	// @Default true
	RouterCaseSensitive bool
	// RecoverPanic
	// @Description if it was true, Beego will try to recover from panic when it serves your http request
	// So you should notice that it doesn't mean that Beego will recover all panic cases.
	// @Default true
	RecoverPanic bool
	// CopyRequestBody
	// @Description if it's true, Beego will copy the request body. But if the request body's size > MaxMemory,
	// Beego will return 413 as http status
	// If you are building RESTful API, please set it to true.
	// And if you want to read data from request Body multiple times, please set it to true
	// In general, if you don't meet any performance issue, you could set it to true
	// @Default false
	CopyRequestBody bool
	// EnableGzip
	// @Description If it was true, Beego will try to compress data by using zip algorithm.
	// But there are two points:
	// 1. Only static resources will be compressed
	// 2. Only those static resource which has the extension specified by StaticExtensionsToGzip will be compressed
	// @Default false
	EnableGzip bool
	// EnableErrorsShow
	// @Description If it's true, Beego will show error message to page
	// it will work with ErrorMaps which allows you register some error handler
	// You may want to set it to false when application was deploy to prod environment
	// because you may not want to expose your internal error msg to your users
	// it's a little bit unsafe
	// @Default true
	EnableErrorsShow bool
	// EnableErrorsRender
	// @Description If it's true, it will output the error msg as a page. It's similar to EnableErrorsShow
	// And this configure item only work in dev run mode (see RunMode)
	// @Default true
	EnableErrorsRender bool
	// ServerName
	// @Description server name. For example, in large scale system,
	// you may want to deploy your application to several machines, so that each of them has a server name
	// we suggest you'd better set value because Beego use this to output some DEBUG msg,
	// or integrated with other tools such as tracing, metrics
	// @Default
	ServerName string

	// RecoverFunc
	// @Description when Beego want to recover from panic, it will use this func as callback
	// see RecoverPanic
	// @Default defaultRecoverPanic
	RecoverFunc func(*context.Context, *Config)
	// @Description MaxMemory and MaxUploadSize are used to limit the request body
	// if the request is not uploading file, MaxMemory is the max size of request body
	// if the request is uploading file, MaxUploadSize is the max size of request body
	// if CopyRequestBody is true, this value will be used as the threshold of request body
	// see CopyRequestBody
	// the default value is 1 << 26 (64MB)
	// @Default 67108864
	MaxMemory int64
	// MaxUploadSize
	// @Description  MaxMemory and MaxUploadSize are used to limit the request body
	// if the request is not uploading file, MaxMemory is the max size of request body
	// if the request is uploading file, MaxUploadSize is the max size of request body
	// the default value is 1 << 30 (1GB)
	// @Default 1073741824
	MaxUploadSize int64
	// Listen
	// @Description the configuration about socket or http protocol
	Listen Listen
	// WebConfig
	// @Description the configuration about Web
	WebConfig WebConfig
	// LogConfig
	// @Description log configuration
	Log LogConfig
}

// Listen holds for http and https related config
type Listen struct {
	// Graceful
	// @Description means use graceful module to start the server
	// @Default false
	Graceful bool
	// ListenTCP4
	// @Description if it's true, means that Beego only work for TCP4
	// please check net.Listen function
	// In general, you should not set it to true
	// @Default false
	ListenTCP4 bool
	// EnableHTTP
	// @Description if it's true, Beego will accept HTTP request.
	// But if you want to use HTTPS only, please set it to false
	// see EnableHTTPS
	// @Default true
	EnableHTTP bool
	// AutoTLS
	// @Description If it's true, Beego will use default value to initialize the TLS configure
	// But those values could be override if you have custom value.
	// see Domains, TLSCacheDir
	// @Default false
	AutoTLS bool
	// EnableHTTPS
	// @Description If it's true, Beego will accept HTTPS request.
	// Now, you'd better use HTTPS protocol on prod environment to get better security
	// In prod, the best option is EnableHTTPS=true and EnableHTTP=false
	// see EnableHTTP
	// @Default false
	EnableHTTPS bool
	// EnableMutualHTTPS
	// @Description if it's true, Beego will handle requests on incoming mutual TLS connections
	// see Server.ListenAndServeMutualTLS
	// @Default false
	EnableMutualHTTPS bool
	// EnableAdmin
	// @Description if it's true, Beego will provide admin service.
	// You can visit the admin service via browser.
	// The default port is 8088
	// see AdminPort
	// @Default false
	EnableAdmin bool
	// EnableFcgi
	// @Description
	// @Default false
	EnableFcgi bool
	// EnableStdIo
	// @Description EnableStdIo works with EnableFcgi Use FCGI via standard I/O
	// @Default false
	EnableStdIo bool
	// ServerTimeOut
	// @Description Beego use this as ReadTimeout and WriteTimeout
	// The unit is second.
	// see http.Server.ReadTimeout, WriteTimeout
	// @Default 0
	ServerTimeOut int64
	// HTTPAddr
	// @Description Beego listen to this address when the application start up.
	// @Default ""
	HTTPAddr string
	// HTTPPort
	// @Description Beego listen to this port
	// you'd better change this value when you deploy to prod environment
	// @Default 8080
	HTTPPort int
	// Domains
	// @Description Beego use this to configure TLS. Those domains are "white list" domain
	// @Default []
	Domains []string
	// TLSCacheDir
	// @Description Beego use this as cache dir to store TLS cert data
	// @Default ""
	TLSCacheDir string
	// HTTPSAddr
	// @Description Beego will listen to this address to accept HTTPS request
	// see EnableHTTPS
	// @Default ""
	HTTPSAddr string
	// HTTPSPort
	// @Description  Beego will listen to this port to accept HTTPS request
	// @Default 10443
	HTTPSPort int
	// HTTPSCertFile
	// @Description Beego read this file as cert file
	// When you are using HTTPS protocol, please configure it
	// see HTTPSKeyFile
	// @Default ""
	HTTPSCertFile string
	// HTTPSKeyFile
	// @Description Beego read this file as key file
	// When you are using HTTPS protocol, please configure it
	// see HTTPSCertFile
	// @Default ""
	HTTPSKeyFile string
	// TrustCaFile
	// @Description Beego read this file as CA file
	// @Default ""
	TrustCaFile string
	// AdminAddr
	// @Description Beego will listen to this address to provide admin service
	// In general, it should be the same with your application address, HTTPAddr or HTTPSAddr
	// @Default ""
	AdminAddr string
	// AdminPort
	// @Description  Beego will listen to this port to provide admin service
	// @Default 8088
	AdminPort int
	// @Description Beego use this tls.ClientAuthType to initialize TLS connection
	// The default value is tls.RequireAndVerifyClientCert
	// @Default 4
	ClientAuth int
}

// WebConfig holds web related config
type WebConfig struct {
	// AutoRender
	// @Description If it's true, Beego will render the page based on your template and data
	// In general, keep it as true.
	// But if you are building RESTFul API and you don't have any page,
	// you can set it to false
	// @Default true
	AutoRender bool
	// Deprecated: Beego didn't use it anymore
	EnableDocs bool
	// EnableXSRF
	// @Description If it's true, Beego will help to provide XSRF support
	// But you should notice that, now Beego only work for HTTPS protocol with XSRF
	// because it's not safe if using HTTP protocol
	// And, the cookie storing XSRF token has two more flags HttpOnly and Secure
	// It means that you must use HTTPS protocol and you can not read the token from JS script
	// This is completed different from Beego 1.x because we got many security reports
	// And if you are in dev environment, you could set it to false
	// @Default false
	EnableXSRF bool
	// DirectoryIndex
	// @Description When Beego serves static resources request, it will look up the file.
	// If the file is directory, Beego will try to find the index.html as the response
	// But if the index.html is not exist or it's a directory,
	// Beego will return 403 response if DirectoryIndex is **false**
	// @Default false
	DirectoryIndex bool
	// FlashName
	// @Description the cookie's name when Beego try to store the flash data into cookie
	// @Default BEEGO_FLASH
	FlashName string
	// FlashSeparator
	// @Description When Beego read flash data from request, it uses this as the separator
	// @Default BEEGOFLASH
	FlashSeparator string
	// StaticDir
	// @Description Beego uses this as static resources' root directory.
	// It means that Beego will try to search static resource from this start point
	// It's a map, the key is the path and the value is the directory
	// For example, the default value is /static => static,
	// which means that when Beego got a request with path /static/xxx
	// Beego will try to find the resource from static directory
	// @Default /static => static
	StaticDir map[string]string
	// StaticExtensionsToGzip
	// @Description The static resources with those extension will be compressed if EnableGzip is true
	// @Default [".css", ".js" ]
	StaticExtensionsToGzip []string
	// StaticCacheFileSize
	// @Description If the size of static resource < StaticCacheFileSize, Beego will try to handle it by itself,
	// it means that Beego will compressed the file data (if enable) and cache this file.
	// But if the file size > StaticCacheFileSize, Beego just simply delegate the request to http.ServeFile
	// the default value is 100KB.
	// the max memory size of caching static files is StaticCacheFileSize * StaticCacheFileNum
	// see StaticCacheFileNum
	// @Default 102400
	StaticCacheFileSize int
	// StaticCacheFileNum
	// @Description Beego use it to control the memory usage of caching static resource file
	// If the caching files > StaticCacheFileNum, Beego use LRU algorithm to remove caching file
	// the max memory size of caching static files is StaticCacheFileSize * StaticCacheFileNum
	// see StaticCacheFileSize
	// @Default 1000
	StaticCacheFileNum int
	// TemplateLeft
	// @Description Beego use this to render page
	// see TemplateRight
	// @Default {{
	TemplateLeft string
	// TemplateRight
	// @Description Beego use this to render page
	// see TemplateLeft
	// @Default }}
	TemplateRight string
	// ViewsPath
	// @Description The directory of Beego application storing template
	// @Default views
	ViewsPath string
	// CommentRouterPath
	// @Description Beego scans this directory and its sub directory to generate router
	// Beego only scans this directory when it's in dev environment
	// @Default controllers
	CommentRouterPath string
	// XSRFKey
	// @Description the name of cookie storing XSRF token
	// see EnableXSRF
	// @Default beegoxsrf
	XSRFKey string
	// XSRFExpire
	// @Description the expiration time of XSRF token cookie
	// second
	// @Default 0
	XSRFExpire int
	// @Description session related config
	Session SessionConfig
}

// SessionConfig holds session related config
type SessionConfig struct {
	// SessionOn
	// @Description if it's true, Beego will auto manage session
	// @Default false
	SessionOn bool
	// SessionAutoSetCookie
	// @Description if it's true, Beego will put the session token into cookie too
	// @Default true
	SessionAutoSetCookie bool
	// SessionDisableHTTPOnly
	// @Description used to allow for cross domain cookies/javascript cookies
	// In general, you should not set it to true unless you understand the risk
	// @Default false
	SessionDisableHTTPOnly bool
	// SessionEnableSidInHTTPHeader
	// @Description enable store/get the sessionId into/from http headers
	// @Default false
	SessionEnableSidInHTTPHeader bool
	// SessionEnableSidInURLQuery
	// @Description enable get the sessionId from Url Query params
	// @Default false
	SessionEnableSidInURLQuery bool
	// SessionProvider
	// @Description session provider's name.
	// You should confirm that this provider has been register via session.Register method
	// the default value is memory. This is not suitable for distributed system
	// @Default memory
	SessionProvider string
	// SessionName
	// @Description If SessionAutoSetCookie is true, we use this value as the cookie's name
	// @Default beegosessionID
	SessionName string
	// SessionGCMaxLifetime
	// @Description Beego will GC session to clean useless session.
	// unit: second
	// @Default 3600
	SessionGCMaxLifetime int64
	// SessionProviderConfig
	// @Description the config of session provider
	// see SessionProvider
	// you should read the document of session provider to learn how to set this value
	// @Default ""
	SessionProviderConfig string
	// SessionCookieLifeTime
	// @Description If SessionAutoSetCookie is true,
	// we use this value as the expiration time and max age of the cookie
	// unit second
	// @Default 0
	SessionCookieLifeTime int
	// SessionDomain
	// @Description If SessionAutoSetCookie is true, we use this value as the cookie's domain
	// @Default ""
	SessionDomain string
	// SessionNameInHTTPHeader
	// @Description if SessionEnableSidInHTTPHeader is true, this value will be used as the http header
	// @Default Beegosessionid
	SessionNameInHTTPHeader string
	// SessionCookieSameSite
	// @Description If SessionAutoSetCookie is true, we use this value as the cookie's same site policy
	// the default value is http.SameSiteDefaultMode
	// @Default 1
	SessionCookieSameSite http.SameSite

	// SessionIDPrefix
	// @Description session id's prefix
	// @Default ""
	SessionIDPrefix string
}

// LogConfig holds Log related config
type LogConfig struct {
	// AccessLogs
	// @Description If it's true, Beego will log the HTTP request info
	// @Default false
	AccessLogs bool
	// EnableStaticLogs
	// @Description log static files requests
	// @Default false
	EnableStaticLogs bool
	// FileLineNum
	// @Description if it's true, it will log the line number
	// @Default true
	FileLineNum bool
	// AccessLogsFormat
	// @Description access log format: JSON_FORMAT, APACHE_FORMAT or empty string
	// @Default APACHE_FORMAT
	AccessLogsFormat string
	// Outputs
	// @Description the destination of access log
	// the key is log adapter and the value is adapter's configure
	// @Default "console" => ""
	Outputs map[string]string // Store Adaptor : config
}

var (
	if err != nil {
		panic(err)
	}
	filename := "app.conf"
	if os.Getenv("BEEGO_RUNMODE") != "" {
		filename = os.Getenv("BEEGO_RUNMODE") + ".app.conf"
	}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		}

		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}

		if cfg.RunMode == DEV && cfg.EnableErrorsRender {
			showErr(err, ctx, stack)
		}
	}
}

				SessionEnableSidInHTTPHeader: false, // enable store/get the sessionId into/from http headers
				SessionNameInHTTPHeader:      "Beegosessionid",
				SessionEnableSidInURLQuery:   false, // enable get the sessionId from Url Query params
				SessionCookieSameSite:        http.SameSiteDefaultMode,
			},
		},
		Log: LogConfig{
// For 1.x, it use assignSingleConfig to parse the file
// but for 2.x, we use Unmarshaler method
func assignConfig(ac config.Configer) error {
	parseConfigForV1(ac)

	err := ac.Unmarshaler("", BConfig)
			// do nothing here
		}
	}
}

// LoadAppConfig allow developer to apply a config file
