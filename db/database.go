package main

import (
    "fmt"
    "log"
    "time"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)


var schema = `
create table transaction (
    transactionId serial,
    buyerId integer,
    sellerId integer,
    ticker text,
    amountTraded integer,
    cashTraded float(64)
);

create table user (
    userId serial,
    userName text,
    userPasswordHash text
);

create table company (
    ticker text,
    name text
);`

type DataBase struct {
    User string
    Password string
    Name string
    Sslmode string
    Port int
}

type Transaction struct {
    TransactionId int
    BuyerId int
    SellerId int
    EquityTicker string
    AmountTraded int
    CashTraded  float64
    TimeOfTrade time.Time
}

/* No args, called on the DataBase struct and returns a pointer to
 * sqlx database object.*/
func (db DataBase) openDataBase() (*sqlx.DB, error) {
    connStr := fmt.Sprintf(`user=%s password=%s dbname=%s sslmode=%s port=%d`,
                           db.User,
                           db.Password,
                           db.Name,
                           db.Sslmode,
                           db.Port)
    fmt.Println(connStr)
    return sqlx.Open("postgres", connStr)
}

/* 5 args, takes fields of Transaction struct */
func (db *sqlx.DB) insertTransaction(buyerId int,
                                     sellerId int,
                                     ticker string,
                                     amountTraded int,
                                     cash float64,
                                     timeOfTrade time.Time) {
    db
}

func main() {
  database := DataBase{"g1727122_u",
     "PTqnydAPoe",
     "g1727122_u",
     "require",
     5432}

  db, err := database.openDataBase()

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
