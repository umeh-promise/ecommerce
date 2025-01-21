package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/umeh-promise/ecommerce/cmd/api"
	"github.com/umeh-promise/ecommerce/internal/db"
	"github.com/umeh-promise/ecommerce/utils"
)

func main() {
	config := utils.Config{
		Addr:        utils.GetString("DB_ADDR", "postgres://user:password@localhost:5432/ecommerce?sslmode=disable"),
		MaxOpenConn: utils.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConn: utils.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime: utils.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	db, err := db.NewDBConnection(config.Addr, config.MaxOpenConn, config.MaxIdleConn, config.MaxIdleTime)
	if err != nil {
		utils.Logger.Fatal("failed to open database connection %w", err)
	}

	defer db.Close()
	utils.Logger.Info("DB connected successfully")

	sever := api.NewAPIServer(":8080", db)
	if err := sever.Run(); err != nil {
		log.Fatal(err)
	}
}
