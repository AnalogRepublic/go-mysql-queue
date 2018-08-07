package msq

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ConnectionConfig struct {
	Type     string
	Host     string
	Username string
	Password string
	Database string
	Charset  string
	Locale   string
}

type Connection struct {
	db     *gorm.DB
	Config ConnectionConfig
}

func (c *Connection) Attempt() error {
	db, err := gorm.Open(c.getType(), c.getConnectionString())

	if err != nil {
		return err
	}

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
	c.db.AutoMigrate(&Event{})

	hasTable := c.db.HasTable(&Event{})

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

	if c.Config.Locale == "" {
		c.Config.Locale = "Local"
	}

	if dbType == "mysql" {
		return fmt.Sprintf(
			"%s:%s@%s/%s?charset=%s&parseTime=True&loc=%s",
			c.Config.Username,
			c.Config.Password,
			c.Config.Host,
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
