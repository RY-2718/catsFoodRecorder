package main

import (
	"database/sql"
	"github.com/RY-2718/catsFoodRecorder/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("バインドするアドレスを指定してください")
	}

	db := initDB("./storage.db")
	migrate(db)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "./frontend/index.html")
	e.GET("/foods", handler.GetFoods(db))
	e.POST("/foods", handler.CreateFood(db))
	e.DELETE("/foods/:id", handler.DeleteFood(db))

	e.Start(os.Args[1])
}

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db
}

func migrate(db *sql.DB) {
	sql := `
		CREATE TABLE IF NOT EXISTS foods(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name VARCHAR NOT NULL,
		created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))
	);
	`

	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		panic(err)
	}
}
