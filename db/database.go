package main

import (
    "fmt"
    "log"
    "time"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
)


var schema = `
create table transactionTable (
    buyerId integer,
    sellerId integer,
    ticker text,
    amountTraded integer,
    cashTraded float(53),
    timeOfTrade timestamp
);

create table userTable (
    userId serial,
    userName text,
    userPasswordHash text
);

create table companyTable (
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

type User struct {
  UserId int              `db:"userId"`
  UserName string         `db:"userName"`
  UserPasswordHash string `db:"userPasswordHash"`
}

type Transaction struct {
    BuyerId int           `db:"buyerId"`
    SellerId int          `db:"sellerId"`
    Ticker string         `db:"ticker"`
    AmountTraded int      `db:"amountTraded"`
    CashTraded  float64   `db:"cashTraded"`
    TimeOfTrade time.Time `db:"timeOfTrade"`
}

/* No args, called on the DataBase struct and returns a pointer to
 * sqlx database struct. Opens a connection to the database.*/
func (db DataBase) openDataBase() (*sqlx.DB, error) {
    connStr := fmt.Sprintf(`user=%s password=%s dbname=%s sslmode=%s port=%d`,
                           db.User,
                           db.Password,
                           db.Name,
                           db.Sslmode,
                           db.Port)
    return sqlx.Open("postgres", connStr)
}

/* 6 args, first is the sqlx database struct pointer and the rest are
 * the fields of Transaction struct, returns void. Inserts a transaction
 * into the database.*/
func insertTransaction(db *sqlx.DB,
                       buyerId int,
                       sellerId int,
                       ticker string,
                       amountTraded int,
                       cash float64,
                       timeOfTrade time.Time) {
    ax := db.MustBegin()
    ax.MustExec(`insert into transactionTable (buyerId, sellerId,
                                          ticker, amountTraded,
                                          cashTraded, timeOfTrade)
                                          values ($1, $2, $3, $4, $5, $6)`,
                                          buyerId, sellerId,  ticker,
                                          amountTraded, cash, timeOfTrade)
    ax.Commit()
}


/* 3 args, first is the sqlx database struct pointer, the second is
 * the username and the last is the password hash.*/
func createUser(db *sqlx.DB,
                userName string,
                userPasswordHash string) {
    ax := db.MustBegin()
    ax.MustExec(`insert into userTable (userName, userPasswordHash)
                 values ($1, $2)`, userName, userPasswordHash)
    ax.Commit()
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

  users := []User{}
  err = db.Select(&users, "select * from userTable")
    if err != nil {
        fmt.Println(err)
        return
    }
  for _, user := range users {
    fmt.Printf("name: %s\n", user.UserName)
  }
}
