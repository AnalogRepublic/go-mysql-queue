package msq

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ConnectionConfig struct {
	Type     string
	Host     string
	Proto    string
	Port     int
	Username string
	Password string
	Database string
	Charset  string
	Locale   string
	Logging  bool
}

type Connection struct {
	db     *gorm.DB
	Config ConnectionConfig
}

func (c *Connection) Attempt() error {
	if c.Config.Locale == "" {
		c.Config.Locale = "Local"
	}

	if c.Config.Charset == "" {
		c.Config.Charset = "utf8mb4"
	}

	if c.Config.Proto == "" {
		c.Config.Proto = "tcp"
	}

	connectionString := c.getConnectionString()

	if c.Config.Logging {
		fmt.Println("Connecting to " + connectionString)
	}

	db, err := gorm.Open(c.getType(), connectionString)

	if err != nil {
		return err
	}

	db.LogMode(c.Config.Logging)

	c.db = db

	return nil
}

func (c *Connection) Close() error {
	return c.db.Close()
}

func (c *Connection) Database() *gorm.DB {
	return c.db
}

func (c *Connection) SetupDatabase() error {
	dbScope := c.db

	if c.getType() != "sqlite3" {
		tableOptions := fmt.Sprintf("ENGINE=InnoDB DEFAULT CHARSET=%s", c.Config.Charset)
		dbScope = dbScope.Set("gorm:table_options", tableOptions)
	}

	dbScope = dbScope.AutoMigrate(&Event{})

	hasTable := dbScope.HasTable(&Event{})

	if !hasTable {
		return errors.New("Events table was not created")
	}

	return nil
}

func (c *Connection) getType() string {
	if c.Config.Type == "sqlite" {
		return "sqlite3"
	}

	return c.Config.Type
}

func (c *Connection) getConnectionString() string {
	dbType := c.getType()

	port := strconv.Itoa(c.Config.Port)
	hostname := c.Config.Host

	if port != "" && port != "0" {
		hostname = fmt.Sprintf("%s:%s", c.Config.Host, port)
	}

	if dbType == "mysql" {
		return fmt.Sprintf(
			"%s:%s@%s(%s)/%s?charset=%s&parseTime=True&loc=%s",
			c.Config.Username,
			c.Config.Password,
			c.Config.Proto,
			hostname,
			c.Config.Database,
			c.Config.Charset,
			c.Config.Locale,
		)
	} else if dbType == "sqlite3" {
		return c.Config.Database
	}

	panic("Invalid database type provided, must be 'myqsl' or 'sqlite3'/'sqlite'")
}

func Connect(config ConnectionConfig) (*Queue, error) {
	connection := &Connection{
		Config: config,
	}

	err := connection.Attempt()

	if err != nil {
		return &Queue{}, err
	}

	connection.SetupDatabase()

	queue := &Queue{
		Connection: connection,
	}

	return queue, nil
}
