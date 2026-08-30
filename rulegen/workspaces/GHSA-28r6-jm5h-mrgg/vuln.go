package main

	if runObject == nil && len(t.fixrouters) > 0 {
		// Filter the .json .xml .html extension
		for _, str := range allowSuffixExt {
			if strings.HasSuffix(seg, str) {
				for _, subTree := range t.fixrouters {
					if subTree.prefix == seg[:len(seg)-len(str)] {
						runObject = subTree.match(treePattern, pattern, wildcardValues, ctx)
						if runObject != nil {
import (
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)
}

func init() {
	routers = make([]testInfo, 0)
	// match example
	routers = append(routers, matchTestInfo("/topic/?:auth:int", "/topic", nil))
	routers = append(routers, matchTestInfo("/topic/?:auth:int", "/topic/123", map[string]string{":auth": "123"}))
	routers = append(routers, notMatchTestInfo("/read_:id:int\\.htm", "/read_222_htm"))
	routers = append(routers, notMatchTestInfo("/read_:id:int\\.htm", " /read_262shtm"))

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
				return
			}
		}
		if obj == nil || obj.(string) != "astaxie" {
			}
		}
	}
}

func TestStaticPath(t *testing.T) {
