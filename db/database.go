package main

import (
  //"database/sql"
    "fmt"
    "log"
    "time"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)

type DataBase struct {
  User string
  Password string
  Name string
  Sslmode string
  Port int
}

type Transaction struct {
  BuyerId int
  SellerId int
  EquityTicker string
  AmountTraded int
  CashTraded  int
  TimeOfTrade time.Time
}

func main() {
  connStr := `user=g1727122_u password=PTqnydAPoe
    dbname=g1727122_u sslmode=require port=5432`
  db, err := sqlx.Open("postgres", connStr)

  if err != nil {
		log.Fatal(err)
	}

  trades := []Transaction{}
  err = db.Select(&trades, "SELECT * FROM transactions")
    if err != nil {
        fmt.Println(err)
        return
    }
  for _, trade := range trades {
    fmt.Printf("time: %s\n", trade.TimeOfTrade.Format("2006-01-02T15:04:05Z07:00.000"))
  }
}
