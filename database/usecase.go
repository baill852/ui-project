package database

import (
	"context"
	"fmt"
	"ui-project/logger"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	client *gorm.DB
	log    logger.LogUsecase
}

func NewDatabase(ctx context.Context, log logger.LogUsecase) DatabaseUsecase {
	return &database{
		log: log,
	}
}

func (d *database) GetClient() *gorm.DB {
	return d.client
}

func (d *database) Bootstrap(ctx context.Context) {
	d.log.LogInfo(ctx, fmt.Sprintf("Database Bootstrap on %v:%v\n", viper.GetString("db.host"), viper.GetInt("db.port")))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.name"),
		viper.GetInt("db.port"),
	)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		d.log.LogErr(ctx, "connection failed", err)
		panic(err)
	}
	d.client = client
}

func (d *database) Shutdown(ctx context.Context) {
	if db, err := d.client.DB(); err != nil {
		d.log.LogErr(ctx, "Database Shutdown failed", err)
	} else {
		if err := db.Close(); err != nil {
			d.log.LogErr(ctx, "Database Shutdown failed", err)
		} else {
			d.log.LogInfo(ctx, "Database Shutdown success")
		}
	}
}
