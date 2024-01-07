package database

import (
	"sync"
	"time"

	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DB defines db connections.
type DB struct {
	*gorm.DB

	config *config.DatabaseConfig

	connLock *sync.Mutex
}

var defaultDB DB

// Initialize inits default db.
func Initialize(cfg *config.DatabaseConfig) {
	defaultDB.initialize(cfg)
}

// Initialize initializes models.
// It only creates the connection instance, doesn't reset or migrate anything.
func (d *DB) initialize(cfg *config.DatabaseConfig) {
	d.config = cfg

	if d.connLock == nil {
		d.connLock = &sync.Mutex{}
	}
	d.connLock.Lock()
	defer d.connLock.Unlock()

	if d.DB == nil {
		d.DB = d.dialDB()
	}
}

// Finalize closes db.
func Finalize() {
	defaultDB.finalize()
}

// Finalize closes the database.
func (d *DB) finalize() {
	d.connLock.Lock()
	defer d.connLock.Unlock()

	if d.DB != nil {
		sql, err := d.DB.DB()
		if err != nil {
			panic(err)
		}
		sql.Close()
		d.DB = nil
	}
}

// GetDB gets db from singleton.
func GetDB() *DB {
	return defaultDB.getDB()
}

// GetDB returns the database handle.
func (d *DB) getDB() *DB {
	d.connLock.Lock()
	conn := d.DB
	d.connLock.Unlock()

	if conn == nil {
		panic("uninitialized database")
	}

	return &defaultDB
}

func (d *DB) dialDB() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	db, err = d.connect()
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (d *DB) connect() (db *gorm.DB, err error) {
	db, err = gorm.Open(d.config.Open(), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		return
	}

	// Set database parameters.
	sql, err := db.DB()
	if err != nil {
		return nil, err
	}

	sql.SetMaxIdleConns(d.config.MaxIdleConns)
	sql.SetMaxOpenConns(d.config.MaxOpenConns)
	sql.SetConnMaxLifetime(time.Duration(d.config.MaxLifetime) * time.Second)

	return
}
