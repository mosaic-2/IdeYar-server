package dbutil

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mosaic-2/IdeYar-server/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	gormConnectionPool *gorm.DB
	gormDebugLogger    logger.Interface
	gormLog            bool
)

func init() {

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password= %s sslmode=disable",
		config.GetDBHost(),
		config.GetDBPort(),
		config.GetDBName(),
		config.GetDBUser(),
		config.GetDBPass(),
	)

	gormLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Minute,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	log.Printf("connecting to database ...")
	cp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger, NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(fmt.Sprintf("database connection failed: %s", err.Error()))
	}

	gormLog, err = strconv.ParseBool(config.GetDBLog())
	if err != nil {
		panic(fmt.Sprintf("could not parse %s env value", config.GetDBLog()))
	}

	gormDebugLogger = logger.New(log.Default(), logger.Config{
		SlowThreshold:             time.Minute,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      false,
		Colorful:                  false,
	})

	gormConnectionPool = cp
	log.Printf("connected to database: dbname=%s user=%s host=%s port=%s",
		config.GetDBName(), config.GetDBUser(), config.GetDBHost(), config.GetDBPort())

}

func GormDB(ctx context.Context) *gorm.DB {

	db := gormConnectionPool.WithContext(ctx)
	if gormLog {
		db = db.Session(&gorm.Session{Logger: gormDebugLogger})
	}
	return db
}

func WithGormDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, "gormKey", db)
}
