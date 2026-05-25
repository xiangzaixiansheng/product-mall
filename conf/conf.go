package conf

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"product-mall/cache"
	"runtime"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	RunMode  string
	HttpPort string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	AccessKey   string
	SerectKey   string
	Bucket      string
	QiniuServer string

	ValidEmail string
	SmtpHost   string
	SmtpEmail  string
	SmtpPass   string
	ENV        string

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string

	MysqlpathRead  string
	MysqlpathWrite string
)

func Init() {
	//env
	if _env := os.Getenv("ENV"); _env != "" {
		ENV = _env
	} else {
		ENV = "dev"
	}
	slog.Info("environment", "ENV", ENV)
	//获取当前目录
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	configFilePath := filepath.Join(dir, fmt.Sprintf("app.%s.ini", ENV))
	localFilePath := filepath.Join(dir, "/locales/zh-cn.yaml")

	slog.Info("config file", "path", configFilePath)

	file, err := ini.Load(configFilePath)
	if err != nil {
		slog.Error("failed to load config file", "path", configFilePath, "error", err)
		panic(err)
	}
	LoadServer(file)
	LoadMysqlData(file)
	LoadRedisData(file)
	LoadEmail(file)
	LoadQinNiu(file)
	if err := LoadLocales(localFilePath); err != nil {
		slog.Error("failed to load locales", "path", localFilePath, "error", err)
		panic(err)
	}
	slog.Info("redis config", "addr", RedisAddr, "db", RedisDbName)
	//redis
	cache.NewRedis(RedisAddr, RedisDbName, "")

	//MySQL
	MysqlpathRead = strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	MysqlpathWrite = strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
}

func LoadServer(file *ini.File) {
	RunMode = file.Section("server").Key("RunMode").String()
	HttpPort = file.Section("server").Key("HttpPort").String()
}

func LoadMysqlData(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadEmail(file *ini.File) {
	ValidEmail = file.Section("email").Key("ValidEmail").String()
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpEmail = file.Section("email").Key("SmtpEmail").String()
	SmtpPass = file.Section("email").Key("SmtpPass").String()
}

func LoadQinNiu(file *ini.File) {
	AccessKey = file.Section("qiniu").Key("AccessKey").String()
	SerectKey = file.Section("qiniu").Key("SerectKey").String()
	Bucket = file.Section("qiniu").Key("Bucket").String()
	QiniuServer = file.Section("qiniu").Key("QiniuServer").String()
}

func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
