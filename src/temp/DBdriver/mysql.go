package DBdriver
import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	cfg "../config"
	"reflect"
	"strings"
	"strconv"
	"encoding/json"
	//"fmt"
	"sync"
	"fmt"
	"time"
)
var lock sync.Mutex
var map_SQLSTR map[string]string
var map_ORMTYPE map[string][]string
var map_JSON map[string]string
var map_INSERT map[string]string
var map_TABLE map[string][]string
func GE() (map[string]string,map[string][]string,map[string]string,map[string]string,map[string][]string)  {
	return map_SQLSTR,map_ORMTYPE,map_JSON,map_INSERT,map_TABLE
}
func init ()  {
	map_SQLSTR = make(map[string]string,1)
	map_ORMTYPE = make(map[string][]string,1)
	map_JSON = make(map[string]string,1)
	map_INSERT = make(map[string]string,1)
	map_TABLE = make(map[string][]string,1)
}
func MysqlDB() *sql.DB {
	db, err := sql.Open(cfg.MysqlConfig.DriverName, cfg.MysqlConfig.Url)
	if err!=nil{panic(err)}
	return db
}
func MysqlClose(db *sql.DB)  {
	db.Close()
}
func objFieldsAddr(reObj reflect.Value) ([]interface{}) {
	//
	Addrs := make([]interface{},0,0)
	for i:=0;i<reObj.NumField();i++{
		var field = reObj.Field(i)
		if field.Kind().String()=="struct"{
			Addrs = append(Addrs,objFieldsAddr(field)...)
		}else{
			Addrs = append(Addrs,field.Addr().Interface())
		}
	}
	return Addrs
}
func QueryGoObj(db *sql.DB,table string,objPtr interface{},getCondition func()(string,[]interface{})) {
	if db==nil{
		db = MysqlDB()
		defer db.Close()
	}
	sql,ok:=map_SQLSTR[table];if !ok{panic("表映射不存在..")}
	cdt,params:=getCondition()
	stmt,_:=db.Prepare(sql+cdt)
	row:=stmt.QueryRow(params...)
	reObj := reflect.ValueOf(objPtr).Elem()
	addrs :=objFieldsAddr(reObj)
	row.Scan(addrs...)
	stmt.Close()
}
func QueryGoArray(db *sql.DB,table string,objPtr interface{},getEntity func()(interface{}),getCondition func()(string,[]interface{}))([]interface{}) {
	if db==nil{
		db = MysqlDB()
		defer db.Close()
	}
	sql,ok:=map_SQLSTR[table];if !ok{panic("表映射不存在..")}
	cdt,params:=getCondition()
	stmt,_:=db.Prepare(sql+cdt)
	rows,_:=stmt.Query(params...)
	reObj := reflect.ValueOf(objPtr).Elem()
	addrs :=objFieldsAddr(reObj)
	entitys:=make([]interface{},0,0)
	for rows.Next(){
		rows.Scan(addrs)
		entitys = append(entitys,getEntity())
	}
	rows.Close()
	stmt.Close()
	return entitys
}
/*
	@getCondition() 返回的字符串查询字符串必须以and开头，第二个是查询参数的slice
*/
func QueryJsObj(db *sql.DB,table string,getCondition func()(string,[]interface{})) string {
	if db==nil{
		db = MysqlDB()
		defer db.Close()
	}
	sql,ok:=map_SQLSTR[table];if !ok{panic("表映射不存在..")}
	cdt,params:=getCondition()
	stmt,_:=db.Prepare(sql+cdt)
	row:=stmt.QueryRow(params...)
	addrs,fc:=getAddrs(table)
	row.Scan(addrs...)
	stmt.Close()
	return fc()
}
/*
	@getCondition() 返回的字符串查询字符串必须以and开头，第二个是查询参数的slice
*/
func QueryJsArray(db *sql.DB,table string,getCondition func()(condition string,params []interface{})) string {
	if db==nil{
		db = MysqlDB()
		defer db.Close()
	}
	sql,ok:=map_SQLSTR[table];if !ok{panic("表映射不存在..")}

	cdt,params:=getCondition()
	stmt,_:=db.Prepare(sql+cdt)
	fmt.Println(sql+cdt)
	rows,err:=stmt.Query(params...);if err!=nil {panic(err)}
	js:="["
	a:=0
	addrs,fc:=getAddrs(table)
	for rows.Next(){
		//fmt.Println(addrs)
		rows.Scan(addrs...)
		time.Sleep(3000000)
		if a==0{ js = js+fc() }else { js = js+","+fc() }
		a++
	}
	rows.Close()
	stmt.Close()
	return (js+"]")
}
/*
	@getCondition() 返回的字符串查询字符串不必以and开头，第二个是查询参数的slice
	@content 是更新的字典。key代表数据的字段，value是对应字段的值
*/
func Update(db *sql.DB,table string,content map[string]string,getCondition func()(string,[]interface{})) (af int64 ){
	if db==nil{
		db = MysqlDB()
		defer db.Close()
	}
	var set string
	a:=0;
	for k,v:=range content{
		if a == 0{	set+=" SET "+k+" = " + v  }else { set+=" ,SET "+k+" = " + v }
		a++
	}
	cdt,params:=getCondition()
	sql:=" UPDATE " + table + set +" where "+cdt
	stmt,_:=db.Prepare(sql)
	result,err:=stmt.Exec(params...);if err!=nil {panic(err)}
	af,_=result.RowsAffected()
	return
}
/*
	@getCondition() 返回的字符串查询字符串不必以and开头，第二个是查询参数的slice
*/
func Delete(db *sql.DB,table string,getCondition func()(string,[]interface{}))(af int64 )  {
	if db==nil{
		db = MysqlDB()
		defer db.Close()
	}
	cdt,params:=getCondition()
	sql := "DELETE FROM "+table +" WHERE " +cdt
	stmt,_:=db.Prepare(sql)
	result,err:=stmt.Exec(params...);if err!=nil {panic(err)}
	af,_=result.RowsAffected()
	stmt.Close()
	return
}
func Insert(db *sql.DB,table string,params map[string]string)(lid int64) {
	if db==nil{
		db = MysqlDB()
		defer db.Close()
	}
	stmt,_:=db.Prepare(map_INSERT[table])
	pArr,ok:=map_TABLE[table];if !ok{panic("表映射不存在")}
	num:=len(pArr)
	var ps =make([]interface{},num,num)
	for i,v:=range pArr {
		ps[i]=params[v]
	}
	result,_ := stmt.Exec(ps...)
	lid,_ =  result.LastInsertId()
	return
}
func getAddrs(tabel string)([]interface{},func()(string)){
	tystr,ok:=map_ORMTYPE[tabel]
	if !ok {panic("找不到表映射...")}
	n:=len(tystr)
	ptr:=make([]interface{},n,n)
	values:=make([]interface{},n,n)
	for i,v:=range tystr{
		switch v{
		case "bool":
			var vs bool
			values[i]=vs
			ptr[i]=&values[i]
		case "int":
			var vs int
			values[i]=vs
			ptr[i]=&values[i]
		case "int64":
			var vs int64
			values[i]=vs
			ptr[i]=&values[i]
		case "float32":
			var vs float32
			values[i]=vs
			ptr[i]=&values[i]
		case "float64":
			var vs float64
			values[i]=vs
			ptr[i]=&values[i]
		case "string":
			var vs string
			values[i]=vs
			ptr[i]=&values[i]
		}
	}
	return ptr,func()(string){
		tp,ok:=map_JSON[tabel]
		if !ok {panic("找不到表映射...")}
		for i,v:=range tystr{
			var s string
			var vs = values[i]
			switch v{
			case "bool":
				if vs==nil{ s="" }else { s = strconv.FormatBool(vs.(bool)) }
			case "int":
				if vs==nil{ s="0" }else { s = strconv.Itoa(vs.(int)) }
			case "int64":
				if vs==nil{ s="0" }else { s = strconv.FormatInt(vs.(int64), 10)}
			case "float32":
				if vs==nil{ s="0" }else { s = strconv.FormatFloat(float64(vs.(float32)),'f',-1,32) }
			case "float64":
				if vs==nil{ s="0" }else { s = strconv.FormatFloat(vs.(float64),'f',-1,64) }
			case "string":
				if vs==nil{ s="" }else {
					if vby,ok1:=vs.([]byte);ok1{
						s = string(vby)
					}else if vstr,ok2:=vs.(string);ok2{
						s = vstr
					}
				}
			}
			tp = strings.Replace(tp,"#",s,1)
			fmt.Println(v,"=>" ,s)
		}
		return tp
	}
}
func entityInit(rt reflect.Type,table string,join string) ([]string,string,string) {
	sqlStr:=" "
	typeArr:=make([]string,0)
	for i:=0;i<rt.NumField();i++{
		if i>0{
			sqlStr = sqlStr+","
		}
		typeStr:=rt.Field(i).Type.Kind().String()
		var sqlElem string
		if rt.Field(i).Tag.Get("sql")!=""{
			sqlElem =rt.Field(i).Tag.Get("sql")
		}else{
			sqlElem =rt.Field(i).Tag.Get("json")
			if sqlElem==""{	panic("缺少必要的tag标签...")	}
		}
		switch typeStr {
		case "bool":
			sqlStr = sqlStr+table+"."+sqlElem
			typeArr = append(typeArr,"bool")
		case "int":
			sqlStr = sqlStr+table+"."+sqlElem
			typeArr = append(typeArr,"int")
		case "int64":
			sqlStr = sqlStr+table+"."+sqlElem
			typeArr = append(typeArr,"int64")
		case "float32":
			sqlStr = sqlStr+table+"."+sqlElem
			typeArr = append(typeArr,"float32")
		case "float64":
			sqlStr = sqlStr+table+"."+sqlElem
			typeArr = append(typeArr,"float64")
		case "string":
			sqlStr = sqlStr+table+"."+sqlElem
			typeArr = append(typeArr,"string")
		case "struct":
			if rt.Field(i).Tag.Get("table")!=""{
				sqlElem = rt.Field(i).Tag.Get("table")
			}
			pk,fk:=rt.Field(i).Tag.Get("pk"),rt.Field(i).Tag.Get("fk")
			if fk==""{ panic("必须标明外键") }
			if pk==""{	pk="id" }

			arr,s,j:= entityInit(rt.Field(i).Type,sqlElem," ")
			join = (join +" LEFT OUTER JOIN " + sqlElem + " ON " + table + "."+fk+" = " + sqlElem+"."+pk + j)
			sqlStr = sqlStr + s
			typeArr =  append(typeArr,arr...)
		default:
			panic("entity 不支持的类型:"+typeStr)
		}
	}
	return typeArr,sqlStr,join
}
func EntityRegister(ptr interface{},table string){
	lock.Lock()
	//初始化json字符串模版
	tp,_:=json.Marshal(ptr);js := string(tp)
	js = strings.Replace(js, "0,", "#,", -1)
	js = strings.Replace(js, "\"\"", "\"#\"", -1)
	js = strings.Replace(js, "\"\"}", "\"#\"}", -1)
	js = strings.Replace(js, "0}", "#}", -1)
	map_JSON[table]=js

	//初始化插入语句
	var DB= MysqlDB()
	stmt,_:=DB.Prepare("select * from "+table+" where 1=2 ")
	rows,_:=stmt.Query()
	cl,_:=rows.Columns()

	var ips string = " ("
	for i,_:=range cl{	if i==0 { ips = ips+"?" }else{ ips = (ips+",?") } }
	ips=ips+")"
	insert:= "INSERT INTO "+ table + "(" +strings.Join(cl,",") + ") VALUES " +ips
	map_INSERT[table]=insert
	map_TABLE[table] = cl
	rows.Close();stmt.Close();DB.Close()

	//初始化查询语句
	tArr,itemStr,join:=entityInit(reflect.TypeOf(ptr).Elem(),table,table)
	itemStr=strings.TrimSpace(itemStr)
	itemStr=strings.TrimRight(itemStr,",")
	sqlStr := "SELECT "+itemStr+" FROM "+join+" WHERE 1=1 "
	map_ORMTYPE[table] = tArr //类型映射
	map_SQLSTR[table] = sqlStr //查询模版
	lock.Unlock()
}