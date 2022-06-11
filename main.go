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

func main() {
	pguser := env.Optional("PGUSER", "postgres")
	pgpass := env.Optional("PGPASS", "")
	pghost := env.Optional("PGHOST", "localhost")
	pgport := env.Optional("PGPORT", "5432")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable&application_name=postgresql-check", pguser, pgpass, pghost, pgport)
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
	log.Fatal(http.ListenAndServe(":26726", nil))
}
