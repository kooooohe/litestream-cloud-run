package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func main() {
	fmt.Println("test")
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	_, err = db.Exec("INSERT INTO users(name) VALUES(?)",randString(16))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	tx.Commit()
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer rows.Close()

	us := []user{}
	for rows.Next() {
		u := user{}
		err := rows.Scan(&u.Id, &u.Name)
		if err != nil {
			fmt.Printf("%v", err)
		}
		us = append(us, u)
	}
	j, err := json.Marshal(us)
	if err != nil {
		fmt.Printf("%v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

var db *sql.DB

type user struct {
	Id   int
	Name string
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	dataPath := flag.String("dp", "", "the path for data")

	flag.Parse()
	fmt.Println(*dataPath)
	if *dataPath == "" {
		flag.Usage()
		return errors.New("data path option error")
	}
	_, err := os.Stat(*dataPath)
	// TOOD do it as shell
	if os.IsNotExist(err) {
		// restore
	} else {
	}

	db, err = sql.Open("sqlite3", *dataPath)
	if err != nil {
		return err
	}
	defer db.Close()
	cSQLstmt := "CREATE TABLE IF NOT EXISTS users(id integer not null primary key, name text);"
	_, err = db.Exec(cSQLstmt)

	if err != nil {
		return err
	}

	http.HandleFunc("/", helloHandler)
	go http.ListenAndServe(":8080", nil)
	<-ctx.Done()
	return nil
}
