package session

import (
	"sync"
	"time"
	"io"
	"encoding/base64"
	"crypto/rand"
	"net/http"
)
var SessionManager = &SessionStore{
	cookieName:"SESSION_ID",
	maxlifetime:3600*time.Second,
	cookieMaxAge:3600,
	store:make(map[string]*Session,100),
}
//session------------------------------------------------------------------------------
type Session struct {
	sessionID string
	lifetimeEnd time.Time
	content map[string]string
}
func (s *Session)Set(key,value string) bool{
	s.content[key] = value
	return true
}
func (s *Session)Get(key string) (string) {
	v:=s.content[key]
	return v
}
func (s *Session)Remove(key string) bool{
	delete(s.content,key)
	return true
}
func (s *Session) GetSID()string{
	return s.sessionID
}
func (s *Session)Destroy() {
	SessionManager.destroySession(s.sessionID)
}
//sessionc存储器------------------------------------------------------------------------------------------
type SessionStore struct {
	cookieName string
	maxlifetime time.Duration	//session存活时间 以纳秒为单位
	cookieMaxAge int
	lock sync.RWMutex	//读写锁
	store map[string]*Session
}
func (ss *SessionStore) newSession() *Session{
	//存储session
	var sid string
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err == nil {
		sid = base64.URLEncoding.EncodeToString(b)
	}
	session := &Session{sid,time.Now().Add(ss.maxlifetime),make(map[string]string,20)}
	//建立索引
	ss.lock.Lock()
	ss.store[sid] = session
	ss.lock.Unlock()
	return session
}
func (ss *SessionStore)destroySession(sid string)  {
	ss.lock.Lock()
	delete(ss.store,sid)
	ss.lock.Unlock()
}
func (ss *SessionStore)GetSession(sid string) *Session {
	ss.lock.RLock()
	s,ok := ss.store[sid]
	ss.lock.RUnlock()
	if !ok{
		return nil
	}else{
		if s.lifetimeEnd.Before(time.Now()){
			ss.destroySession(sid)
			return nil
		}else{
			s.lifetimeEnd = time.Now().Add(ss.maxlifetime)
			return s
		}

	}
}
/**
	session回收
 */
func (ss *SessionStore)sessionGC(){
	t:=time.Now()
	for k,v:=range ss.store {
		if v.lifetimeEnd.Before(t){
			ss.destroySession(k)
		}
	}
	time.AfterFunc(ss.maxlifetime, func() {
		ss.sessionGC()
	})
}
/**
	@param s 设置session回话时间
 */
func (ss *SessionStore)SetLifeTime(s int64)  {
	ss.cookieMaxAge = int(s)
	ss.maxlifetime = time.Duration(s*int64(time.Second))
}
/**
	开启session
 */
func (ss *SessionStore)SessionStart(w http.ResponseWriter,req *http.Request) *Session {
	cookie,err := req.Cookie(ss.cookieName)
	var session *Session
	if err!=nil||cookie==nil ||cookie.Value==""{
		session = ss.newSession()
		nCookie := http.Cookie{
			Name:ss.cookieName,
			Value:session.sessionID,
			Path:"/",
			HttpOnly:true,
			MaxAge:ss.cookieMaxAge,
		}
		http.SetCookie(w,&nCookie)
	}else {
		session = ss.GetSession(cookie.Value)
		if session==nil{
			session = ss.newSession()
			cookie.Value=session.sessionID
			cookie.MaxAge =ss.cookieMaxAge
			http.SetCookie(w,cookie)
		}
	}
	return session
}
/**
	包初始化
 */
func init()  {
	go SessionManager.sessionGC()
}