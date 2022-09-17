package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/GeeScot/go-common/env"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

type Config struct {
	Host string `env:"POSTGRES_HOST"`
	Port string `env:"POSTGRES_PORT"`
	User string `env:"POSTGRES_USERNAME"`
	Pass string `env:"POSTGRES_PASSWORD"`
}

func main() {
	var config Config
	env.Read(&config)

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=disable&application_name=postgresql-check",
		config.User,
		config.Pass,
		config.Host,
		config.Port)

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

	stmt, err := db.Prepare("select pg_is_in_recovery()")
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	defer stmt.Close()
	defer db.Close()

	isInRecoveryHandler := func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		writeServiceUnavailable := func() {
			log.Printf("%s\n", "[503] Service unavailable")
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		rows, err := stmt.Query()
		if err != nil || !rows.Next() {
			writeServiceUnavailable()
			return
		}

		defer rows.Close()

		result := ""
		err = rows.Scan(&result)

		if err != nil {
			writeServiceUnavailable()
			return
		}

		isInRecovery, err := strconv.ParseBool(result)
		if err != nil {
			writeServiceUnavailable()
			return
		}

		var statusCode int
		if isInRecovery {
			statusCode = http.StatusServiceUnavailable
		} else {
			statusCode = http.StatusOK
		}

		w.WriteHeader(statusCode)
		log.Printf("[%d] Request time: %s\n", statusCode, time.Since(t).Truncate(time.Microsecond))
	}

	http.HandleFunc("/", isInRecoveryHandler)
	log.Fatal(http.ListenAndServe(":26726", nil))
}
