package main

import (
	"log"

	env "github.com/dneedsleep/Social/internal"
	"github.com/dneedsleep/Social/internal/db"
	"github.com/dneedsleep/Social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://root:newpassword@localhost/social? sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	store := store.NewPostgresStorage(conn)

	db.Seed(store)

}
