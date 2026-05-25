package db

import (
	"context"
	"log/slog"
	"os"
	"product-mall/conf"
	"product-mall/internal/model"
	"product-mall/pkg/pkg_logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

type ormLog struct{}

//orm 日志记录
func (l ormLog) Printf(format string, args ...any) {
	if conf.ENV == "dev" {
		pkg_logger.Logger.Info(format, "args", args)
	}
}

func Database(connRead, connWrite string) *gorm.DB {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据版本自动配置
	}), &gorm.Config{
		Logger: pkg_logger.NewGORMLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) //打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	DB = db
	_ = DB.Use(dbresolver.
		Register(dbresolver.Config{
			// `db2` 作为 sources，`db3`、`db4` 作为 replicas
			Sources:  []gorm.Dialector{mysql.Open(connRead)},                         // 写操作
			Replicas: []gorm.Dialector{mysql.Open(connWrite), mysql.Open(connWrite)}, // 读操作
			Policy:   dbresolver.RandomPolicy{},                                      // sources/replicas 负载均衡策略
		}))
	//迁移 schema
	if conf.ENV == "dev" {
		Migration()
	}
	return DB
}

//自动迁移模式
func Migration() {
	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{},
			&model.Notice{},
			&model.Product{},
			&model.ProductImg{},
			&model.Address{},
			&model.Cart{})
	if err != nil {
		slog.Info("table migration failed")
		os.Exit(0)
	}
	slog.Info("table migration success")
}

// 获取db
func GetDBCtx(ctx context.Context) *gorm.DB {
	return DB.WithContext(ctx)
}

// 获取db
func GetDB() *gorm.DB {
	return DB
}
