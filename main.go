package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/GeeScot/go-common/env"
	"github.com/gin-gonic/gin"
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
		fmt.Printf("%s\n", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	defer db.Close()

	isInRecovery := func(c *gin.Context) {
		rows, err := db.Query("select pg_is_in_recovery()")
		if err != nil || !rows.Next() {
			fmt.Printf("%s\n", err.Error())
			c.Status(http.StatusServiceUnavailable)
			return
		}

		defer rows.Close()

		result := ""
		err = rows.Scan(&result)

		if err != nil {
			fmt.Printf("%s\n", err.Error())
			c.Status(http.StatusServiceUnavailable)
			return
		}

		isInRecovery, err := strconv.ParseBool(result)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			c.Status(http.StatusServiceUnavailable)
			return
		}

		if isInRecovery {
			c.String(http.StatusPartialContent, "Standby")
		} else {
			c.String(http.StatusOK, "Primary")
		}
	}

	r := gin.Default()
	r.OPTIONS("/", isInRecovery)
	// r.GET("/", isInRecovery)

	r.Run("0.0.0.0:26726")
}
