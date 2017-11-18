package router
import (
	"fmt"
	"regexp"
	"strings"
	"io"
	"encoding/base64"
	"crypto/rand"
	"../ftype"
)

var sRegx = regexp.MustCompile("^/[a-zA-Z0-9]+")
var npR = regexp.MustCompile("(/+)$")
var pR =  regexp.MustCompile("([^/]+)$")


type  Router struct{
	pMap map[string]ftype.HandleMethod
	params map[string]string
	npMap map[string]ftype.HandleMethod
	resource [2]string
	fileServeMethod ftype.HandleMethod
	notFoundMethod ftype.HandleMethod
}

func (rs Router)InitRouter(pMap int,npMap int) (r Router ) {
	r.pMap = make(map[string]ftype.HandleMethod,pMap)
	r.params = make(map[string]string,pMap)
	r.npMap = make(map[string]ftype.HandleMethod,npMap)
	r.resource[0]="/file"

	return
}
func(r * Router) AddRouterRules(pattern string,handler ftype.HandleMethod) {
	var pRegx = regexp.MustCompile("^(/?[a-zA-Z0-9]+)(/[a-zA-Z0-9]+)*(/:[a-zA-Z0-9]+)$")
	var pR = regexp.MustCompile("(:[a-zA-Z0-9]+)$")
	var npRegx = regexp.MustCompile("^(/?[a-zA-Z0-9]*)(/[a-zA-Z0-9]+)*(/?)$")
	ok,_:=regexp.MatchString("^/",pattern)
	if !ok {
		pattern = "/"+pattern
	}
	pattern = strings.TrimSpace(pattern)
	if pattern=="/"{
		r.npMap["/"] = handler
		return
	}
	if pRegx.MatchString(pattern){
		path := pR.ReplaceAllString(pattern,"")
		if path==r.resource[0]{
			panic(path+" -> 此url与路径已经被注册为静态资源的范围路径")
		}
		r.pMap[path] = handler
		r.params[path] = pR.FindString(pattern)
	}else if npRegx.MatchString(pattern){
		path := npR.ReplaceAllString(pattern,"")
		if path==r.resource[0]{
			panic(path+" -> 此url与路径已经被注册为静态资源的范围路径")
		}
		r.npMap[path] = handler
	}else{
		fmt.Println(pattern+"-->该路由器规则拼写有误，所允许字符集[a-zA-Z0-9]")
	}
}
func(r * Router) SetStaticPath(url string,path string){
	//   /file -> c:file

	sRegx = regexp.MustCompile("^/[a-zA-Z0-9]+")
	url = sRegx.FindString(url)
	if r.npMap[url]!=nil || r.pMap[url]!=nil {
		panic(url+" -> "+path+"动态路由器已经添加了该规则,此静态路由设置不会生效")
	}else{
		r.resource[0] = url
	}
	var spRegx = regexp.MustCompile("(/+)$")
	r.resource[1] = spRegx.ReplaceAllString(path,"")
}


func (r * Router) GetHandler(urlPath string) (f ftype.HandleMethod,cInfo *ftype.ConnInfo) {
	//生成唯一的ConnID
	cInfo = &ftype.ConnInfo{}
	if urlPath=="/" {
		f = r.npMap["/"]
		b := make([]byte, 48)
		if _, err := io.ReadFull(rand.Reader, b); err == nil {
			cInfo.ConnID = base64.URLEncoding.EncodeToString(b)
		}
		cInfo.ReqType = "DYNAMIC_RESOURCE"
		cInfo.SqlParam = make(map[string]string,20)
	}else {
		s :=sRegx.FindString(urlPath)
		if s ==r.resource[0] {//静态资源请求
			cInfo.ReqType = "STATIC_RESOURCE"
			f = r.fileServeMethod
		}else {//动态资源请求
			b := make([]byte, 48)
			if _, err := io.ReadFull(rand.Reader, b); err == nil {
				cInfo.ConnID = base64.URLEncoding.EncodeToString(b)
			}
			cInfo.ReqType = "DYNAMIC_RESOURCE"
			cInfo.SqlParam = make(map[string]string,20)
			nps := npR.ReplaceAllString(urlPath,"")
			var npf = r.npMap[nps]
			if npf!=nil{
				f = npf
			}else{
				ps := pR.ReplaceAllString(nps,"")
				cInfo.UrlParam = pR.FindString(nps)
				var pf = r.pMap[ps]
				if pf!=nil {
					f = pf
				}else {
					cInfo.ReqType = "NOTFOUND"
					f =  r.notFoundMethod
				}
			}
		}
	}
	return
}
func (r * Router)SetFileServeMethod(hm ftype.HandleMethod) {
	r.fileServeMethod = hm
}
func (r * Router) SetNotFoundMethod(nfm ftype.HandleMethod) {
	r.notFoundMethod = nfm
}
func (r *Router)LocalFilePath()string {
	return r.resource[1]
}