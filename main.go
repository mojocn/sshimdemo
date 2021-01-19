package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var postManVar *PostMam

func main() {
	db := loadDatabase("db.sqlite3")
	err := db.Migrator().AutoMigrate(User{}, Group{}, UserUser{}, Msg{}, GroupUser{})
	if err != nil {
		log.Fatal(err)
	}
	postManVar = NewPostMam(db)
	startSshSvrListen(":2222", db)
}

func loadDatabase(dbPath string) *gorm.DB {
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy:         nil,
		FullSaveAssociations:   false,
		Logger: logger.New(
			log.New(os.Stdout, "", log.Ltime), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // Disable color
			},
		),
		NowFunc: func() time.Time {
			return time.Now().Truncate(time.Second)
		},
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})
	if err != nil {
		log.Fatalf("master fail to open its sqlite db in %s. please install master first. %v", dbPath, err)
		return nil
	}
	return database
}
