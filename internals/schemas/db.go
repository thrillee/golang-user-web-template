package schemas

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db             *gorm.DB
	migratableApps = []interface{}{}
)

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second,   // Slow SQL threshold
		LogLevel:                  logger.Silent, // Log level
		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
		ParameterizedQueries:      true,          // Don't include params in the SQL log
		Colorful:                  false,         // Disable color
	},
)

func AddToMigratables(models *Model) {
	migratableApps = append(migratableApps, models)
}

func startDB() (*gorm.DB, error) {
	err := godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: newLogger,
	})
	return db, err
}

func MigrateApps() {
	migratablesCount := len(migratableApps)
	log.Printf("Running Migration for %v apps...\n", migratablesCount)

	if migratablesCount > 0 {
		defer log.Println("Migration completed")
		log.Fatal(getDB().AutoMigrate(migratableApps...))
	} else {
		log.Println("0 Migrations Found!")
	}
}

func getDB() *gorm.DB {
	var err error
	if db == nil {
		db, err = startDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	return db
}
