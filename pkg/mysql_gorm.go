package livelead

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // go mysql driver
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql import driver for gorm
)

// Database is an struct with the needed properties to instantiate
// a gorm database connection and use it properly
type Database struct {
	Host      string
	Port      int64
	User      string
	Pass      string
	Dbname    string
	Charset   string
	ParseTime string
	Loc       string

	*gorm.DB
}

// Storer is an interface used to force client to implement
// the declared methods
type Storer interface {
	Open() error
	Close()
	Instance() *gorm.DB
	Update(element interface{}, wCond string, wFields []string) error
	Insert(element interface{}) error
}

// Open opens a database connection
func (d *Database) Open() error {
	connstr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		d.User, d.Pass, d.Host, d.Port, d.Dbname, d.Charset, d.ParseTime, d.Loc)

	db, err := gorm.Open("mysql", connstr)
	if err != nil {
		log.Fatalf("Error opening database connection %v", err)
		return err
	}

	if err = db.DB().Ping(); err != nil {
		log.Fatalf("Error pinging database %v", err)
		return err
	}

	d.DB = db

	return nil
}

// Close Database.DB instance
func (d *Database) Close() {
	d.DB.Close()
}

// Instance returns a Database.DB instance
func (d *Database) Instance() *gorm.DB {
	return d.DB
}

// Update executes an update sentence using gorm
func (d *Database) Update(element interface{}, wCond string, wFields []string) error {
	wFieldsArr := []interface{}{}
	for _, z := range wFields {
		wFieldsArr = append(wFieldsArr, z)
	}

	d.Model(element).Where(wCond).Update(wFieldsArr...)
	return nil
}

// Insert executes a insert statement using gorm
func (d *Database) Insert(element interface{}) error {
	if result := d.Create(element); result.Error != nil {
		return fmt.Errorf("error inserting element %v", result.Error)
	}
	return nil
}

// AutoMigrate automatically migrate your schema, to keep your schema update to date.
// and create the table if not exists
func (d *Database) AutoMigrate() error {
	if err := d.DB.AutoMigrate(LeadLive{}).Error; err != nil {
		return err
	}
	return nil
}
