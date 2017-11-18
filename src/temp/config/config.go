package config
type SqlConfig struct {
	DriverName string
	Url string
}
var MysqlConfig = SqlConfig{}