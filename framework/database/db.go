package database

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/sidartaoss/fullcycle/encoder/domain"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "github.com/lib/pq"
)

type Database struct {
	*gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDB() *Database {
	return &Database{}
}

func NewDBTest() *gorm.DB {
	dbTest := NewDB()
	dbTest.Env = "test"
	dbTest.DbTypeTest = "sqlite3"
	dbTest.DsnTest = ":memory:"
	dbTest.AutoMigrateDb = true
	dbTest.Debug = true

	conn, err := dbTest.Connect()
	if err != nil {
		log.Fatalf("database test error: %v", err)
	}

	return conn
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	if d.Env != "test" {
		db, err = gorm.Open(d.DbType, d.Dsn)
		if err != nil {
			return nil, err
		}
	} else {
		db, err = gorm.Open(d.DbTypeTest, d.DsnTest)
		if err != nil {
			return nil, err
		}
	}

	d.DB = db

	if d.Debug {
		d.DB.LogMode(true)
	}

	if d.AutoMigrateDb {
		d.DB.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.DB.Model(domain.Job{}).AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return d.DB, nil

}
