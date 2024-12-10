package database

import (
	"fmt"
	"log"

	config "xdemo/internal/config"

	global "xdemo/internal/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DefaultMySQLServerAddr string = "localhost"
var DefaultMySQLServerPort string = "3306"
var DefaultMySQLServerProtocol string = "tcp"
var DefaultMySQLServerTimeOut int = 30

// 数据库连接的接口（支持多种数据库类型）
// type DBer interface {
// 	New() error
// 	Close() error
// 	Connect() error
// }

type DBOption func(*MySQLCServer)

type MySQLCServer struct {
	ClientName string
	Addr       string
	Port       string
	DBName     string
	Protocol   string
	MaxConn    int
	MinConn    int
	Timeout    int
	MaxRetries int
	Username   string
	Password   string
}

// 非必选的选项使用函数选项式进行扩展
func WithAddr(addr string) DBOption {
	return func(s *MySQLCServer) {
		s.Addr = addr
	}
}

func WithPort(port string) DBOption {
	return func(s *MySQLCServer) {
		s.Port = port
	}
}

func WithUsername(username string) DBOption {
	return func(s *MySQLCServer) {
		s.Username = username
	}
}

func WithPassword(password string) DBOption {
	return func(s *MySQLCServer) {
		s.Password = password
	}
}

func WithMaxConn(maxConn int) DBOption {
	return func(s *MySQLCServer) {
		s.MaxConn = maxConn
	}
}

func WithMinConn(minConn int) DBOption {
	return func(s *MySQLCServer) {
		s.MinConn = minConn
	}
}

func WithTimeout(timeout int) DBOption {
	return func(s *MySQLCServer) {
		s.Timeout = timeout
	}
}

func WithMaxRetries(MaxRetries int) DBOption {
	return func(s *MySQLCServer) {
		s.MaxRetries = MaxRetries
	}
}

func WithProtocol(protocol string) DBOption {
	return func(s *MySQLCServer) {
		s.Protocol = protocol
	}
}

// 初始化mysql实例
func NewMySQLServer(name string, dbname string, opts ...DBOption) *gorm.DB {
	// return &MySQLCServer{ClientName: name, Addr: addr, Port: port, Timeout: timeout}
	server := &MySQLCServer{
		ClientName: name,
		DBName:     dbname,
		Addr:       DefaultMySQLServerAddr,
		Port:       DefaultMySQLServerPort,
		Timeout:    DefaultMySQLServerTimeOut,
		Protocol:   DefaultMySQLServerProtocol,
	}
	for _, opt := range opts {
		opt(server)
	}
	// 扩充连接
	dsn := GetDSN(server)
	db, err := server.ConnectWithGORM(dsn)
	if err != nil {
		panic(err.Error())
	}
	log.Println("MySQL Connect Success!!!")
	// 初始化数据table
	err = InitDBTable(db)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// 获取连接mysql的dsn信息
func GetDSN(c *MySQLCServer) string {
	// 目前字符集，解析时间等是固定的
	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Protocol, c.Addr, c.Port, c.DBName)
}

// 通过gorm连接mysql
func (client *MySQLCServer) ConnectWithGORM(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "xdemo_",
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("DB Connect Error: %s", err.Error())
	}
	return db, nil
}

func CloseDB() {
	sqlDB, err := global.GDB.DB()
	if err != nil {
		log.Fatal("Get DB Object Failed - ", err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Fatal("Close DB Object Failed - ", err)
		return
	}
	log.Println("Close DB Success!!!")
}

// 加载DB
func LoadDB(yamlConfig *config.YAMLConfig) {
	// 新建mysql连接用作gorm，后期需要抽象出来
	global.GDB = NewMySQLServer("xdemo", yamlConfig.DataBase.DBName, WithAddr(yamlConfig.DataBase.Host), WithUsername(yamlConfig.DataBase.DBUser), WithPassword(yamlConfig.DataBase.DBPassword), WithPort(yamlConfig.DataBase.Port), WithMaxConn(100), WithMaxRetries(3))
	// 延迟关闭db连接
	// defer func() {
	// 	sqlDB, err := global.GDB.DB()
	// 	if err != nil {
	// 		fmt.Println("Get DB Object Failed,The Error is ", err)
	// 		return
	// 	}
	// 	sqlDB.Close()
	// }()
}
