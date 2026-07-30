package main

import (
	"log"

	env "github.com/dneedsleep/Social/internal"
	"github.com/dneedsleep/Social/internal/db"
	"github.com/dneedsleep/Social/internal/store"
	"github.com/joho/godotenv"
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
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Database connection pool established")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
