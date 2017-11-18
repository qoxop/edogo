package ftime

import (
	"time"
	"strconv"
)

type Time int64
type Duration int64
func(nt Time)Year()int{
	return time.Unix(int64(nt/1000),0).Year()
}
func(nt Time)Month() (string,int){
	m:=time.Unix(int64(nt/1000),0).Month()
	return m.String(),int(m)
}
func(nt Time)Day()int{
	return time.Unix(int64(nt/1000),0).Day()
}
func(nt Time)Format()string{
	y:=strconv.Itoa(nt.Year())
	_,m:=nt.Month()
	d:=strconv.Itoa(nt.Day())
	return y+"-"+strconv.Itoa(m)+"-"+d
}
func Now()Time{
	ut:=time.Now().Unix()
	return Time(ut/1000)
}
func(nt Time)tTime()time.Time{
	return time.Unix(int64(nt/1000),0)
}
