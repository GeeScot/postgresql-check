package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/GeeScot/go-common/fileio"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

type Config struct {
	Postgres struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"postgres"`
	Port int `json:"port"`
}

func main() {
	configFile := flag.String("c", "/etc/postgresql-check/config.json", "config file")
	flag.Parse()

	var config Config
	fileio.ReadJSON(*configFile, &config)

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/postgres?sslmode=disable&application_name=postgresql-check",
		config.Postgres.Username,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port)

	db, err := sql.Open("pg", connectionString)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	defer db.Close()

	isInRecoveryHandler := func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		rows, err := db.Query("select pg_is_in_recovery()")
		if err != nil || !rows.Next() {
			log.Printf("%s\n", "[503] Service unavailable")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		defer rows.Close()

		result := ""
		err = rows.Scan(&result)

		if err != nil {
			log.Printf("%s\n", "[503] Service unavailable")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		isInRecovery, err := strconv.ParseBool(result)
		if err != nil {
			log.Printf("%s\n", "[503] Service unavailable")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		var statusCode int
		if isInRecovery {
			statusCode = http.StatusPartialContent
		} else {
			statusCode = http.StatusOK
		}

		w.WriteHeader(statusCode)
		log.Printf("[%d] Request time: %s\n", statusCode, time.Since(t).Truncate(time.Microsecond))
	}

	http.HandleFunc("/", isInRecoveryHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
