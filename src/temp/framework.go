package fgo

import "./router"
import (
	"./session"
	"net/http"
	"./ftype"
	"regexp"
	"encoding/json"
)
const (
	DYNAMIC_RESOURCE = "DYNAMIC_RESOURCE"
	STATIC_RESOURCE = "STATIC_RESOURCE"
	NOTFOUND = "NOTFOUND"
)
var sRegx = regexp.MustCompile("^/[a-zA-Z0-9]+")
var MUX = router.Router{}
var SessionManager = session.SessionManager

func init() {
	MUX = MUX.InitRouter(0,0)
	MUX.SetFileServeMethod(func(w http.ResponseWriter, req *http.Request, cInfo *ftype.ConnInfo) {
		http.ServeFile(w,req,sRegx.ReplaceAllString(req.URL.Path,MUX.LocalFilePath()))
	});
	MUX.SetNotFoundMethod(func(w http.ResponseWriter, req *http.Request, cInfo *ftype.ConnInfo) {
		http.NotFound(w,req)
	})
}
func JsonString(obj interface{}) string {
	b,_:=json.Marshal(obj)
	return string(b)
}