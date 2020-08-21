package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

var dbx *sqlx.DB

func main() {
	var wg sync.WaitGroup
	var err error

	dsn := "test:test@tcp(127.0.0.1:3306)/test_db"
	dbx, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %s.", err.Error())
	}
	defer dbx.Close()

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			heavySelect()
		}
		wg.Done()
	}()
	go busy()
	wg.Wait()
}

func busy() {
	i := 0
	for {
		i++
	}
}

func heavySelect() {
	log.Println("heavy select...")
	_, err := dbx.Exec("SELECT SLEEP(5)")
	if err != nil {
		log.Println(err)
	}
	log.Println("end")
}
