package ftype

import (
	"net/http"
	"../DBdriver"
	"../session"
	"database/sql"
)

type ConnInfo struct {
	ConnID string
	Session *session.Session
	UrlParam string
	ReqMethod string
	ReqType string
	MysqlDB *sql.DB
	SqlParam map[string]string
}
type MongoCongig struct {
	Url string
	DataBase string
}
type HandleMethod func(http.ResponseWriter,*http.Request,*ConnInfo)

func (ci *ConnInfo)GET_DB() *sql.DB{

	if ci.MysqlDB==nil {
		ci.MysqlDB = DBdriver.MysqlDB()
	}
	return ci.MysqlDB
}
func (ci *ConnInfo) SQLGetP(key string) string {
	return ci.SqlParam[key]
}
func (ci *ConnInfo) SQLSetP(key,value string) bool {
	ci.SqlParam[key] = value
	return true
}
func (ci *ConnInfo)Destroy() {
	if ci.MysqlDB!=nil{
		ci.MysqlDB.Close()
	}
	ci.SqlParam =nil
	ci = nil
}