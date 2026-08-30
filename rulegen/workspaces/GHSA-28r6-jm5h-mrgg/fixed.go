package main

	if runObject == nil && len(t.fixrouters) > 0 {
		// Filter the .json .xml .html extension
		for _, str := range allowSuffixExt {
			if strings.HasSuffix(seg, str) && strings.HasSuffix(treePattern, seg){
				for _, subTree := range t.fixrouters {
					// strings.HasSuffix(treePattern, seg) avoid cases: /aaa.html/bbb could access /aaa/bbb
					if subTree.prefix == seg[:len(seg)-len(str)] {
						runObject = subTree.match(treePattern, pattern, wildcardValues, ctx)
						if runObject != nil {
import (
	"strings"
	"testing"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)
}

func init() {
	routers = make([]testInfo, 0, 128)
	// match example
	routers = append(routers, matchTestInfo("/topic/?:auth:int", "/topic", nil))
	routers = append(routers, matchTestInfo("/topic/?:auth:int", "/topic/123", map[string]string{":auth": "123"}))
	routers = append(routers, notMatchTestInfo("/read_:id:int\\.htm", "/read_222_htm"))
	routers = append(routers, notMatchTestInfo("/read_:id:int\\.htm", " /read_262shtm"))

	// test .html, .json not suffix
	const abcHtml = "/suffix/abc.html"
	routers = append(routers, notMatchTestInfo(abcHtml, "/suffix.html/abc"))
	routers = append(routers, matchTestInfo("/suffix/abc", abcHtml, nil))
	routers = append(routers, matchTestInfo("/suffix/*", abcHtml, nil))
	routers = append(routers, notMatchTestInfo("/suffix/*", "/suffix.html/a"))
	const abcSuffix = "/abc/suffix/*"
	routers = append(routers, notMatchTestInfo(abcSuffix, "/abc/suffix.html/a"))
	routers = append(routers, matchTestInfo(abcSuffix, "/abc/suffix/a", nil))
	routers = append(routers, notMatchTestInfo(abcSuffix, "/abc.j/suffix/a"))

}

func TestTreeRouters(t *testing.T) {
	for _, r := range routers {

		shouldMatch := r.shouldMatchOrNot
		tr := NewTree()
		tr.AddRouter(r.pattern, "astaxie")
		ctx := context.NewContext()
			if obj != nil {
				t.Fatal("pattern:", r.pattern, ", should not match", r.requestUrl)
			} else {
				continue
			}
		}
		if obj == nil || obj.(string) != "astaxie" {
			}
		}
	}
	time.Sleep(time.Second)
}

func TestStaticPath(t *testing.T) {
