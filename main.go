package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/CyCoreSystems/agi"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connect_to_db() bool {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "1234",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "idg_23",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		slog.Error("Error Occured : err", slog.String("err", err.Error()))
		return false
	}

	pingErr := db.Ping()
	if pingErr != nil {
		slog.Error("Error Occured : err", slog.String("pingErr", pingErr.Error()))
		return false
	}
	slog.Info("DB Connected")
	return true
}

func main() {
	db_status := connect_to_db()
	if !db_status {
		os.Exit(1)
	}
	err := agi.Listen("0.0.0.0:4574", handler)
	if err != nil {
		println("Error Occured")
	}

	slog.Info("FastAgi Server Running on Port 8080")
}

func handler(a *agi.AGI) {
	defer func(a *agi.AGI) {
		err := a.Close()
		if err != nil {
			slog.Error("Error Occured : err", slog.String("err", err.Error()))
		}
	}(a)

	err := a.Answer()
	if err != nil {
		_ = a.Hangup()
		return
	}
	err = a.Set("MYVAR", "foo")
	if err != nil {
		println("failed to set variable MYVAR")
		return
	}

	digit, err := a.StreamFile("/usr/TECKINFO/audiofile", "#", -1)
	slog.Info("Test : " + digit)
	if err != nil {
		slog.Error(" Error Playing File ", slog.String("Res", digit))
		os.Exit(1)
	}
	_ = a.Hangup()
}
