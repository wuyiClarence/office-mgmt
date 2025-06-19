package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"platform-backend/config"
	mylog "platform-backend/utils/log"
)

var MysqlDB *orm

func InitDB() {
	MysqlDB = new(orm)
	MysqlDB.loadDBConfig()
}

type orm struct {
	engine *gorm.DB
}

func (db *orm) loadDBConfig() {
	if db.engine != nil {
		return
	}

	var err error

	retryCount := 10                 // 最大重试次数
	retryInterval := 5 * time.Second // 每次重试的间隔时间

	fmt.Printf("Dsn: %s\n", config.MyConfig.MysqlConfig.Dsn)
	for i := 0; i < retryCount; i++ {
		db.engine, err = gorm.Open(mysql.Open(config.MyConfig.MysqlConfig.Dsn), &gorm.Config{
			Logger: &CustomLogger{},
		})
		if err == nil {
			break // 连接成功，跳出循环
		}
		msg := fmt.Sprintf("open database attempt %d failed: %s\n", i+1, err.Error())
		_, _ = os.Stdout.WriteString(msg)
		log.Println(msg)
		time.Sleep(retryInterval) // 等待后重试
	}
	if err != nil {
		msg := "open database error: " + err.Error()
		_, _ = os.Stdout.WriteString(msg)
		log.Panic(msg)
	}

	sqlDB, _ := db.engine.DB()
	sqlDB.SetMaxIdleConns(config.MyConfig.MysqlConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MyConfig.MysqlConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(110 * time.Second)
	err = sqlDB.Ping()
	if err != nil {
		log.Panic(err)
	}
}

func (db *orm) DB() *gorm.DB {
	return db.engine
}

func (db *orm) SqlDb() *sql.DB {
	e, _ := db.engine.DB()
	return e
}

type CustomLogger struct {
	logger.Interface
}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	_, _ = fmt.Fprintf(mylog.SQLLogger, "INFO: "+msg+"\n", data...)
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	_, _ = fmt.Fprintf(mylog.SQLLogger, "WARN: "+msg+"\n", data...)
}

func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	_, _ = fmt.Fprintf(mylog.SQLLogger, "ERROR: "+msg+"\n", data...)
}

// Trace 在每次查询执行时调用
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sqlDetail, rows := fc()

	if elapsed > time.Millisecond*200 {
		_, _ = fmt.Fprintf(mylog.SQLLogger, "SLOW SQL (>200ms): %s, Rows affected: %d\n", sqlDetail, rows)
	}

	if err != nil {
		_, _ = fmt.Fprintf(mylog.SQLLogger, "SQL ERROR: %s\nError: %v\n", sqlDetail, err)
	} else {
		_, _ = fmt.Fprintf(mylog.SQLLogger, "SQL Executed: %s\nRows affected: %d\n", sqlDetail, rows)
	}
}
