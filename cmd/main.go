package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	ctx := context.Background()
	dataStoreName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		"localhost:3306",
		"alex",
		"rocket2288",
		"axxon",
	)

	db, err := sql.Open("mysql", dataStoreName)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	c, err := db.Conn(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
}
