package main

import (
	"errors"
	"log"
	"sync"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
    "gorm.io/datatypes"
)

// MySQL连接池
type MySQLPool struct {
	connections chan *gorm.DB
	factory     func() (*gorm.DB, error)
	close       bool
	size        int
	mu          sync.Mutex
}

var GlobalMySQLPool *MySQLPool

func init() {
	var err error
	GlobalMySQLPool, err = NewMySQLPool(10)
	if err != nil {
		log.Printf("MySQL连接池初始化失败: %v", err)
	} else {
		log.Printf("MySQL连接池初始化成功")
	}
	initModel()
}

func NewMySQLPool(size int) (*MySQLPool, error) {
	p := &MySQLPool{
		connections: make(chan *gorm.DB, size),
		factory:     mysqlFactory,
		size:        size,
	}

	for i := 0; i < size; i++ {
		conn, err := p.factory()
		if err != nil {
			p.Close()
			return nil, err
		}
		p.connections <- conn
	}
	return p, nil
}

func (p *MySQLPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.close {
		return
	}
	p.close = true
	close(p.connections)

	for conn := range p.connections {
		if conn != nil {
			sqlDB, _ := conn.DB()
			sqlDB.Close()
		}
	}
}

func (p *MySQLPool) Get() (*gorm.DB, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.close {
		return nil, errors.New("pool is closed")
	}
	select {
	case conn := <-p.connections:
		return conn, nil
	default:
		conn, err := p.factory()
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func (p *MySQLPool) Put(conn *gorm.DB) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.close {
		sqlDB, _ := conn.DB()
		return sqlDB.Close()
	}
	select {
	case p.connections <- conn:
		return nil
	default:
		sqlDB, _ := conn.DB()
		return sqlDB.Close()
	}
}

// MySQL连接工厂
func mysqlFactory() (*gorm.DB, error) {
	dsn := "root:pass@tcp(db:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

type User struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	OpenID   string    `gorm:"type:varchar(100);not null"`
	Comments datatypes.JSON `gorm:"type:json"`
}

func initModel(){
	cli,err:=GlobalMySQLPool.Get()
	if err!=nil {
		log.Println("数据库连接失败")
		return
	}
	defer GlobalMySQLPool.Put(cli)
	err = cli.AutoMigrate(&User{})
    if err != nil {
        log.Fatalf("建表失败: %v", err)
    }
}