package main

import (
	"log"

	env "github.com/dneedsleep/Social/internal"
	"github.com/dneedsleep/Social/internal/db"
	"github.com/dneedsleep/Social/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://root:newpassword@localhost/social? sslmode=disable"),
			maxOpenConns: env.GetIng("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetIng("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// logger

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// database

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
