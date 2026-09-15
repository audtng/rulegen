package main

// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"

	"github.com/beego/beego/v2/core/logs"

	"github.com/beego/beego/v2/core/utils"
	"github.com/beego/beego/v2/server/web/context/param"
)

var globalRouterTemplate = `package {{.routersDir}}

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"{{.globalimport}}
)

func init() {
{{.globalinfo}}
}
`

var (
	lastupdateFilename = "lastupdate.tmp"
	commentFilename    string
	pkgLastupdate      map[string]int64
	genInfoList        map[string][]ControllerComments

	routerHooks = map[string]int{
		"beego.BeforeStatic": BeforeStatic,
		"beego.BeforeRouter": BeforeRouter,
		"beego.BeforeExec":   BeforeExec,
		"beego.AfterExec":    AfterExec,
		"beego.FinishRouter": FinishRouter,
	}

	routerHooksMapping = map[int]string{
		BeforeStatic: "beego.BeforeStatic",
		BeforeRouter: "beego.BeforeRouter",
		BeforeExec:   "beego.BeforeExec",
		AfterExec:    "beego.AfterExec",
		FinishRouter: "beego.FinishRouter",
	}
)

const commentPrefix = "commentsRouter_"

func init() {
	pkgLastupdate = make(map[string]int64)
}

func parserPkg(pkgRealpath string) error {
	rep := strings.NewReplacer("\\", "_", "/", "_", ".", "_")
	commentFilename, _ = filepath.Rel(AppPath, pkgRealpath)
	commentFilename = commentPrefix + rep.Replace(commentFilename) + ".go"
	if !compareFile(pkgRealpath) {
		logs.Info(pkgRealpath + " no changed")
		return nil
	}
	genInfoList = make(map[string][]ControllerComments)
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedSyntax,
		Dir:  pkgRealpath,
	}, "./...")

	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		for _, fl := range pkg.Syntax {
			for _, d := range fl.Decls {
				switch specDecl := d.(type) {
				case *ast.FuncDecl:
					if specDecl.Recv != nil {
						exp, ok := specDecl.Recv.List[0].Type.(*ast.StarExpr) // Check that the type is correct first beforing throwing to parser
						if ok {
							parserComments(specDecl, fmt.Sprint(exp.X), pkg.PkgPath)
						}
					}
				}
			}
		}
	}
	genRouterCode(pkgRealpath)
	savetoFile(pkgRealpath)
	return nil
}

type parsedComment struct {
	routerPath string
	methods    []string
	params     map[string]parsedParam
	filters    []parsedFilter
	imports    []parsedImport
}

type parsedImport struct {
	importPath  string
	importAlias string
}

type parsedFilter struct {
	pattern string
	pos     int
	filter  string
	params  []bool
}

type parsedParam struct {
	name     string
	datatype string
	location string
	defValue string
	required bool
}

func parserComments(f *ast.FuncDecl, controllerName, pkgpath string) error {
	if f.Doc != nil {
		parsedComments, err := parseComment(f.Doc.List)
		if err != nil {
			return err
		}
		for _, parsedComment := range parsedComments {
			if parsedComment.routerPath != "" {
				key := pkgpath + ":" + controllerName
				cc := ControllerComments{}
				cc.Method = f.Name.String()
				cc.Router = parsedComment.routerPath
				cc.AllowHTTPMethods = parsedComment.methods
				cc.MethodParams = buildMethodParams(f.Type.Params.List, parsedComment)
				cc.FilterComments = buildFilters(parsedComment.filters)
				cc.ImportComments = buildImports(parsedComment.imports)
				genInfoList[key] = append(genInfoList[key], cc)
			}
		}
	}
	return nil
}

func buildImports(pis []parsedImport) []*ControllerImportComments {
	var importComments []*ControllerImportComments

	for _, pi := range pis {
		importComments = append(importComments, &ControllerImportComments{
			ImportPath:  pi.importPath,
			ImportAlias: pi.importAlias,
		})
	}

	return importComments
}

func buildFilters(pfs []parsedFilter) []*ControllerFilterComments {
	var filterComments []*ControllerFilterComments

	for _, pf := range pfs {
		var (
			returnOnOutput bool
			resetParams    bool
		)

		if len(pf.params) >= 1 {
			returnOnOutput = pf.params[0]
		}

		if len(pf.params) >= 2 {
			resetParams = pf.params[1]
		}

		filterComments = append(filterComments, &ControllerFilterComments{
			Filter:         pf.filter,
			Pattern:        pf.pattern,
			Pos:            pf.pos,
			ReturnOnOutput: returnOnOutput,
			ResetParams:    resetParams,
		})
	}

	return filterComments
}

func buildMethodParams(funcParams []*ast.Field, pc *parsedComment) []*param.MethodParam {
	result := make([]*param.MethodParam, 0, len(funcParams))
	for _, fparam := range funcParams {
		for _, pName := range fparam.Names {
			methodParam := buildMethodParam(fparam, pName.Name, pc)
			result = append(result, methodParam)
		}
	}
	return result
}

func buildMethodParam(fparam *ast.Field, name string, pc *parsedComment) *param.MethodParam {
	options := []param.MethodParamOption{}
	if cparam, ok := pc.params[name]; ok {
		// Build param from comment info
		name = cparam.name
		if cparam.required {
			options = append(options, param.IsRequired)
		}
		switch cparam.location {
		case "body":
			options = append(options, param.InBody)
		case "header":
			options = append(options, param.InHeader)
		case "path":
			options = append(options, param.InPath)
		}
		if cparam.defValue != "" {
			options = append(options, param.Default(cparam.defValue))
		}
	} else {
		if paramInPath(name, pc.routerPath) {
			options = append(options, param.InPath)
		}
	}
	return param.New(name, options...)
}

func paramInPath(name, route string) bool {
	return strings.HasSuffix(route, ":"+name) ||
		strings.Contains(route, ":"+name+"/")
}

var routeRegex = regexp.MustCompile(`@router\s+(\S+)(?:\s+\[(\S+)\])?`)

func parseComment(lines []*ast.Comment) (pcs []*parsedComment, err error) {
	pcs = []*parsedComment{}
	params := map[string]parsedParam{}
	filters := []parsedFilter{}
	imports := []parsedImport{}

	for _, c := range lines {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(t, "@Param") {
			pv := getparams(strings.TrimSpace(strings.TrimLeft(t, "@Param")))
			if len(pv) < 4 {
				logs.Error("Invalid @Param format. Needs at least 4 parameters")
			}
			p := parsedParam{}
			names := strings.SplitN(pv[0], "=>", 2)
			p.name = names[0]
			funcParamName := p.name
			if len(names) > 1 {
				funcParamName = names[1]
			}
			p.location = pv[1]
			p.datatype = pv[2]
			switch len(pv) {
			case 5:
				p.required, _ = strconv.ParseBool(pv[3])
			case 6:
				p.defValue = pv[3]
				p.required, _ = strconv.ParseBool(pv[4])
			}
			params[funcParamName] = p
		}
	}

	for _, c := range lines {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(t, "@Import") {
			iv := getparams(strings.TrimSpace(strings.TrimLeft(t, "@Import")))
			if len(iv) == 0 || len(iv) > 2 {
				logs.Error("Invalid @Import format. Only accepts 1 or 2 parameters")
				continue
			}

			p := parsedImport{}
			p.importPath = iv[0]

			if len(iv) == 2 {
				p.importAlias = iv[1]
			}

			imports = append(imports, p)
		}
	}

filterLoop:
	for _, c := range lines {
		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(t, "@Filter") {
			fv := getparams(strings.TrimSpace(strings.TrimLeft(t, "@Filter")))
			if len(fv) < 3 {
				logs.Error("Invalid @Filter format. Needs at least 3 parameters")
				continue filterLoop
			}

			p := parsedFilter{}
			p.pattern = fv[0]
			posName := fv[1]
			if pos, exists := routerHooks[posName]; exists {
				p.pos = pos
			} else {
				logs.Error("Invalid @Filter pos: ", posName)
				continue filterLoop
			}

			p.filter = fv[2]
			fvParams := fv[3:]
			for _, fvParam := range fvParams {
				switch fvParam {
				case "true":
					p.params = append(p.params, true)
				case "false":
					p.params = append(p.params, false)
				default:
					logs.Error("Invalid @Filter param: ", fvParam)
					continue filterLoop
				}
			}

			filters = append(filters, p)
		}
	}

	for _, c := range lines {
		var pc = &parsedComment{}
		pc.params = params
		pc.filters = filters
		pc.imports = imports

		t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(t, "@router") {
			t := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
			matches := routeRegex.FindStringSubmatch(t)
			if len(matches) == 3 {
				pc.routerPath = matches[1]
				methods := matches[2]
				if methods == "" {
					pc.methods = []string{"get"}
					// pc.hasGet = true
				} else {
					pc.methods = strings.Split(methods, ",")
					// pc.hasGet = strings.Contains(methods, "get")
				}
				pcs = append(pcs, pc)
			} else {
				return nil, errors.New("Router information is missing")
			}
		}
	}
	return
}

// direct copy from bee\g_docs.go
// analysis params return []string
// @Param	query		form	 string	true		"The email for login"
// [query form string true "The email for login"]
func getparams(str string) []string {
	var s []rune
	var j int
	var start bool
	var r []string
	var quoted int8
	for _, c := range str {
		if unicode.IsSpace(c) && quoted == 0 {
			if !start {
				continue
			} else {
				start = false
				j++
				r = append(r, string(s))
				s = make([]rune, 0)
				continue
			}
		}

		start = true
		if c == '"' {
			quoted ^= 1
			continue
		}
		s = append(s, c)
	}
	if len(s) > 0 {
		r = append(r, string(s))
	}
	return r
}

func genRouterCode(pkgRealpath string) {
	os.Mkdir(getRouterDir(pkgRealpath), 0755)
	logs.Info("generate router from comments")
	var (
		globalinfo   string
		globalimport string
		sortKey      []string
	)
	for k := range genInfoList {
		sortKey = append(sortKey, k)
	}
	sort.Strings(sortKey)
	for _, k := range sortKey {
		cList := genInfoList[k]
		sort.Sort(ControllerCommentsSlice(cList))
		for _, c := range cList {
			allmethod := "nil"
			if len(c.AllowHTTPMethods) > 0 {
				allmethod = "[]string{"
				for _, m := range c.AllowHTTPMethods {
					allmethod += `"` + m + `",`
				}
				allmethod = strings.TrimRight(allmethod, ",") + "}"
			}

			params := "nil"
			if len(c.Params) > 0 {
				params = "[]map[string]string{"
				for _, p := range c.Params {
					for k, v := range p {
						params = params + `map[string]string{` + k + `:"` + v + `"},`
					}
				}
				params = strings.TrimRight(params, ",") + "}"
			}

			methodParams := "param.Make("
			if len(c.MethodParams) > 0 {
				lines := make([]string, 0, len(c.MethodParams))
				for _, m := range c.MethodParams {
					lines = append(lines, fmt.Sprint(m))
				}
				methodParams += "\n				" +
					strings.Join(lines, ",\n				") +
					",\n			"
			}
			methodParams += ")"

			imports := ""
			if len(c.ImportComments) > 0 {
				for _, i := range c.ImportComments {
					var s string
					if i.ImportAlias != "" {
						s = fmt.Sprintf(`
	%s "%s"`, i.ImportAlias, i.ImportPath)
					} else {
						s = fmt.Sprintf(`
	"%s"`, i.ImportPath)
					}
					if !strings.Contains(globalimport, s) {
						imports += s
					}
				}
			}

			filters := ""
			if len(c.FilterComments) > 0 {
				for _, f := range c.FilterComments {
					filters += fmt.Sprintf(`                &beego.ControllerFilter{
                    Pattern: "%s",
                    Pos: %s,
                    Filter: %s,
                    ReturnOnOutput: %v,
                    ResetParams: %v,
                },`, f.Pattern, routerHooksMapping[f.Pos], f.Filter, f.ReturnOnOutput, f.ResetParams)
				}
			}

			if filters == "" {
				filters = "nil"
			} else {
				filters = fmt.Sprintf(`[]*beego.ControllerFilter{
%s
            }`, filters)
			}

			globalimport += imports

			globalinfo = globalinfo + `
    beego.GlobalControllerRouter["` + k + `"] = append(beego.GlobalControllerRouter["` + k + `"],
        beego.ControllerComments{
            Method: "` + strings.TrimSpace(c.Method) + `",
            ` + `Router: "` + c.Router + `"` + `,
            AllowHTTPMethods: ` + allmethod + `,
            MethodParams: ` + methodParams + `,
            Filters: ` + filters + `,
            Params: ` + params + `})
`
		}
	}

	if globalinfo != "" {
		f, err := os.Create(filepath.Join(getRouterDir(pkgRealpath), commentFilename))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		routersDir := AppConfig.DefaultString("routersdir", "routers")
		content := strings.Replace(globalRouterTemplate, "{{.globalinfo}}", globalinfo, -1)
		content = strings.Replace(content, "{{.routersDir}}", routersDir, -1)
		content = strings.Replace(content, "{{.globalimport}}", globalimport, -1)
		f.WriteString(content)
	}
}

func compareFile(pkgRealpath string) bool {
	if !utils.FileExists(filepath.Join(getRouterDir(pkgRealpath), commentFilename)) {
		return true
	}
	if utils.FileExists(lastupdateFilename) {
		content, err := ioutil.ReadFile(lastupdateFilename)
		if err != nil {
			return true
		}
		json.Unmarshal(content, &pkgLastupdate)
		lastupdate, err := getpathTime(pkgRealpath)
		if err != nil {
			return true
		}
		if v, ok := pkgLastupdate[pkgRealpath]; ok {
			if lastupdate <= v {
				return false
			}
		}
	}
	return true
}

func savetoFile(pkgRealpath string) {
	lastupdate, err := getpathTime(pkgRealpath)
	if err != nil {
		return
	}
	pkgLastupdate[pkgRealpath] = lastupdate
	d, err := json.Marshal(pkgLastupdate)
	if err != nil {
		return
	}
	ioutil.WriteFile(lastupdateFilename, d, os.ModePerm)
}

func getpathTime(pkgRealpath string) (lastupdate int64, err error) {
	fl, err := ioutil.ReadDir(pkgRealpath)
	if err != nil {
		return lastupdate, err
	}
	for _, f := range fl {
		var t int64
		if f.IsDir() {
			t, err = getpathTime(filepath.Join(pkgRealpath, f.Name()))
			if err != nil {
				return lastupdate, err
			}
		} else {
			t = f.ModTime().UnixNano()
		}
		if lastupdate < t {
			lastupdate = t
		}
	}
	return lastupdate, nil
}

func getRouterDir(pkgRealpath string) string {
	dir := filepath.Dir(pkgRealpath)
	for {
		routersDir := AppConfig.DefaultString("routersdir", "routers")
		d := filepath.Join(dir, routersDir)
		if utils.FileExists(d) {
			return d
		}

		if r, _ := filepath.Rel(dir, AppPath); r == "." {
			return d
		}
		// Parent dir.
		dir = filepath.Dir(dir)
	}
}
	"net/http"
	"path"
	"reflect"
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
	// these beego.Controller's methods shouldn't reflect to AutoRouter
	exceptMethod = []string{"Init", "Prepare", "Finish", "Render", "RenderString",
		"RenderBytes", "Redirect", "Abort", "StopRun", "UrlFor", "ServeJSON", "ServeJSONP",
		"ServeYAML", "ServeXML", "Input", "ParseForm", "GetString", "GetStrings", "GetInt", "GetBool",
		"GetFloat", "GetFile", "SaveToFile", "StartSession", "SetSession", "GetSession",
		"DelSession", "SessionRegenerateID", "DestroySession", "IsAjax", "GetSecureCookie",
		"SetSecureCookie", "XsrfToken", "CheckXsrfCookie", "XsrfFormHtml",
		"GetControllerAndAction", "ServeFormatted"}

	urlPlaceholder = "{{placeholder}}"
	// DefaultAccessLogFilter will skip the accesslog if return true
}

// default log filter static file will not show
type logFilter struct {
}

func (l *logFilter) Filter(ctx *beecontext.Context) bool {
	requestPath := path.Clean(ctx.Request.URL.Path)
	exceptMethod = append(exceptMethod, action)
}

// ControllerInfo holds information about the controller.
type ControllerInfo struct {
	pattern        string
	controllerType reflect.Type
	methods        map[string]string
	handler        http.Handler
	runFunction    FilterFunc
	routerType     int
	initialize     func() ControllerInterface
	methodParams   []*param.MethodParam
}

func (c *ControllerInfo) GetPattern() string {
	return c.pattern
}

// ControllerRegister containers registered router rules, controller handlers and filters.
type ControllerRegister struct {
	routers      map[string]*Tree
	enablePolicy bool
	policies     map[string]*Tree
	enableFilter bool
	filters      [FinishRouter + 1][]*FilterRouter
	pool         sync.Pool

	// the filter created by FilterChain
	chainRoot *FilterRouter

	cfg *Config
}

				return beecontext.NewContext()
			},
		},
		cfg: cfg,
	}
	res.chainRoot = newFilterRouter("/*", res.serveHttp, WithCaseSensitive(false))
	return res
}

// Add controller handler and pattern rules to ControllerRegister.
// usage:
//	default methods is the same name as method
//	Add("/api/delete",&RestController{},"delete:DeleteFood")
//	Add("/api",&RestController{},"get,post:ApiFunc"
//	Add("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
func (p *ControllerRegister) Add(pattern string, c ControllerInterface, mappingMethods ...string) {
	p.addWithMethodParams(pattern, c, nil, mappingMethods...)
}

func (p *ControllerRegister) addWithMethodParams(pattern string, c ControllerInterface, methodParams []*param.MethodParam, mappingMethods ...string) {
	reflectVal := reflect.ValueOf(c)
	t := reflect.Indirect(reflectVal).Type()
	methods := make(map[string]string)
	if len(mappingMethods) > 0 {
		semi := strings.Split(mappingMethods[0], ";")
		for _, v := range semi {
			colon := strings.Split(v, ":")
			if len(colon) != 2 {
				panic("method mapping format is invalid")
			}
			comma := strings.Split(colon[0], ",")
			for _, m := range comma {
				if m == "*" || HTTPMETHOD[strings.ToUpper(m)] {
					if val := reflectVal.MethodByName(colon[1]); val.IsValid() {
						methods[strings.ToUpper(m)] = colon[1]
					} else {
						panic("'" + colon[1] + "' method doesn't exist in the controller " + t.Name())
					}
				} else {
					panic(v + " is an invalid method mapping. Method doesn't exist " + m)
				}
			}
		}
	}

	route := &ControllerInfo{}
	route.pattern = pattern
	route.methods = methods
	route.routerType = routerTypeBeego
	route.controllerType = t
	route.initialize = func() ControllerInterface {
		vc := reflect.New(route.controllerType)
		execController, ok := vc.Interface().(ControllerInterface)

		return execController
	}

	route.methodParams = methodParams
	if len(methods) == 0 {
		for m := range HTTPMETHOD {
			p.addToRouter(m, pattern, route)
		}
	} else {
		for k := range methods {
			if k == "*" {
				for m := range HTTPMETHOD {
					p.addToRouter(m, pattern, route)
				}
			} else {
				p.addToRouter(k, pattern, route)
			}
		}
	}
}

func (p *ControllerRegister) addToRouter(method, pattern string, r *ControllerInfo) {
				for _, f := range a.Filters {
					p.InsertFilter(f.Pattern, f.Pos, f.Filter, WithReturnOnOutput(f.ReturnOnOutput), WithResetParams(f.ResetParams))
				}
				p.addWithMethodParams(a.Router, c, a.MethodParams, strings.Join(a.AllowHTTPMethods, ",")+":"+a.Method)
			}
		}
	}
	p.pool.Put(ctx)
}

// Get add get method
// usage:
//    Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Get(pattern string, f FilterFunc) {
	p.AddMethod("get", pattern, f)
}

//    Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Post(pattern string, f FilterFunc) {
	p.AddMethod("post", pattern, f)
}

//    Put("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Put(pattern string, f FilterFunc) {
	p.AddMethod("put", pattern, f)
}

//    Delete("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Delete(pattern string, f FilterFunc) {
	p.AddMethod("delete", pattern, f)
}

//    Head("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Head(pattern string, f FilterFunc) {
	p.AddMethod("head", pattern, f)
}

//    Patch("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Patch(pattern string, f FilterFunc) {
	p.AddMethod("patch", pattern, f)
}

//    Options("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Options(pattern string, f FilterFunc) {
	p.AddMethod("options", pattern, f)
}

//    Any("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Any(pattern string, f FilterFunc) {
	p.AddMethod("*", pattern, f)
}

//    AddMethod("get","/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) AddMethod(method, pattern string, f FilterFunc) {
	method = strings.ToUpper(method)
	if method != "*" && !HTTPMETHOD[method] {
		panic("not support http method: " + method)
	}
	route := &ControllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeRESTFul
	route.runFunction = f
	methods := make(map[string]string)
	if method == "*" {
		for val := range HTTPMETHOD {
			methods[val] = val
		}
	} else {
		methods[method] = method
	}
	route.methods = methods
	for k := range methods {
		if k == "*" {
			for m := range HTTPMETHOD {
				p.addToRouter(m, pattern, route)
			}
		} else {
			p.addToRouter(k, pattern, route)
		}
	}
}

// Handler add user defined Handler
func (p *ControllerRegister) Handler(pattern string, h http.Handler, options ...interface{}) {
	route := &ControllerInfo{}
	route.pattern = pattern
	route.routerType = routerTypeHandler
	route.handler = h
	if len(options) > 0 {
		if _, ok := options[0].(bool); ok {
			pattern = path.Join(pattern, "?:all(.*)")
}

// AddAuto router to ControllerRegister.
// example beego.AddAuto(&MainContorlller{}),
// MainController has method List and Page.
// visit the url /main/list to execute List function
// /main/page to execute Page function.
}

// AddAutoPrefix Add auto router to ControllerRegister with prefix.
// example beego.AddAutoPrefix("/admin",&MainContorlller{}),
// MainController has method List and Page.
// visit the url /admin/main/list to execute List function
// /admin/main/page to execute Page function.
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	for i := 0; i < rt.NumMethod(); i++ {
		if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
			route := &ControllerInfo{}
			route.routerType = routerTypeBeego
			route.methods = map[string]string{"*": rt.Method(i).Name}
			route.controllerType = ct
			pattern := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name), "*")
			patternInit := path.Join(prefix, controllerName, rt.Method(i).Name, "*")
			patternFix := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
			patternFixInit := path.Join(prefix, controllerName, rt.Method(i).Name)
			route.pattern = pattern
			for m := range HTTPMETHOD {
				p.addToRouter(m, pattern, route)
				p.addToRouter(m, patternInit, route)
				p.addToRouter(m, patternFix, route)
				p.addToRouter(m, patternFixInit, route)
			}
		}
	}
}
//     }
// }
func (p *ControllerRegister) InsertFilterChain(pattern string, chain FilterChain, opts ...FilterOpt) {
	root := p.chainRoot
	filterFunc := chain(root.filterFunc)
	opts = append(opts, WithCaseSensitive(p.cfg.RouterCaseSensitive))
	p.chainRoot = newFilterRouter(pattern, filterFunc, opts...)
	p.chainRoot.next = root

}

// add Filter into
	for _, l := range t.leaves {
		if c, ok := l.runObject.(*ControllerInfo); ok {
			if c.routerType == routerTypeBeego &&
				strings.HasSuffix(path.Join(c.controllerType.PkgPath(), c.controllerType.Name()), controllerName) {
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
		runRouter    reflect.Type
		findRouter   bool
		runMethod    string
		methodParams []*param.MethodParam
		routerInfo   *ControllerInfo
		isRunnable   bool
	)

	if p.cfg.RecoverFunc != nil {
	}

	// session init
	if p.cfg.WebConfig.Session.SessionOn {
		ctx.Input.CruSession, err = GlobalSessions.SessionStart(rw, r)
		if err != nil {
			logs.Error(err)

// FindRouter Find Router info for URL
func (p *ControllerRegister) FindRouter(context *beecontext.Context) (routerInfo *ControllerInfo, isFind bool) {
	var urlPath = context.Input.URL()
	if !p.cfg.RouterCaseSensitive {
		urlPath = strings.ToLower(urlPath)
	}
import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"github.com/beego/beego/v2"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/session"

	"github.com/beego/beego/v2/core/utils"
	"github.com/beego/beego/v2/server/web/context"
)

// Config is the main struct for BConfig
// TODO after supporting multiple servers, remove common config to somewhere else
type Config struct {
	AppName             string // Application name
	RunMode             string // Running Mode: dev | prod
	RouterCaseSensitive bool
	ServerName          string
	RecoverPanic        bool
	RecoverFunc         func(*context.Context, *Config)
	CopyRequestBody     bool
	EnableGzip          bool
	// MaxMemory and MaxUploadSize are used to limit the request body
	// if the request is not uploading file, MaxMemory is the max size of request body
	// if the request is uploading file, MaxUploadSize is the max size of request body
	MaxMemory          int64
	MaxUploadSize      int64
	EnableErrorsShow   bool
	EnableErrorsRender bool
	Listen             Listen
	WebConfig          WebConfig
	Log                LogConfig
}

// Listen holds for http and https related config
type Listen struct {
	Graceful          bool // Graceful means use graceful module to start the server
	ServerTimeOut     int64
	ListenTCP4        bool
	EnableHTTP        bool
	HTTPAddr          string
	HTTPPort          int
	AutoTLS           bool
	Domains           []string
	TLSCacheDir       string
	EnableHTTPS       bool
	EnableMutualHTTPS bool
	HTTPSAddr         string
	HTTPSPort         int
	HTTPSCertFile     string
	HTTPSKeyFile      string
	TrustCaFile       string
	EnableAdmin       bool
	AdminAddr         string
	AdminPort         int
	EnableFcgi        bool
	EnableStdIo       bool // EnableStdIo works with EnableFcgi Use FCGI via standard I/O
	ClientAuth        int
}

// WebConfig holds web related config
type WebConfig struct {
	AutoRender             bool
	EnableDocs             bool
	FlashName              string
	FlashSeparator         string
	DirectoryIndex         bool
	StaticDir              map[string]string
	StaticExtensionsToGzip []string
	StaticCacheFileSize    int
	StaticCacheFileNum     int
	TemplateLeft           string
	TemplateRight          string
	ViewsPath              string
	CommentRouterPath      string
	EnableXSRF             bool
	XSRFKey                string
	XSRFExpire             int
	Session                SessionConfig
}

// SessionConfig holds session related config
type SessionConfig struct {
	SessionOn                    bool
	SessionProvider              string
	SessionName                  string
	SessionGCMaxLifetime         int64
	SessionProviderConfig        string
	SessionCookieLifeTime        int
	SessionAutoSetCookie         bool
	SessionDomain                string
	SessionDisableHTTPOnly       bool // used to allow for cross domain cookies/javascript cookies.
	SessionEnableSidInHTTPHeader bool // enable store/get the sessionId into/from http headers
	SessionNameInHTTPHeader      string
	SessionEnableSidInURLQuery   bool // enable get the sessionId from Url Query params
}

// LogConfig holds Log related config
type LogConfig struct {
	AccessLogs       bool
	EnableStaticLogs bool   // log static files requests default: false
	AccessLogsFormat string // access log format: JSON_FORMAT, APACHE_FORMAT or empty string
	FileLineNum      bool
	Outputs          map[string]string // Store Adaptor : config
}

var (
	if err != nil {
		panic(err)
	}
	var filename = "app.conf"
	if os.Getenv("BEEGO_RUNMODE") != "" {
		filename = os.Getenv("BEEGO_RUNMODE") + ".app.conf"
	}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		}
		if cfg.RunMode == DEV && cfg.EnableErrorsRender {
			showErr(err, ctx, stack)
		}
		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}
	}
}

				SessionEnableSidInHTTPHeader: false, // enable store/get the sessionId into/from http headers
				SessionNameInHTTPHeader:      "Beegosessionid",
				SessionEnableSidInURLQuery:   false, // enable get the sessionId from Url Query params
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
