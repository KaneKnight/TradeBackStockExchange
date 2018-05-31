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

create table positionTable (
    userId integer,
    ticker text,
    amount integer,
    initialCashInvestment float(53)
);

create table userTable (
    userId serial,
    userName text,
    userPasswordHash text,
    userCash float(53)
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
  UserId int              `db:"userid"`
  UserName string         `db:"username"`
  UserPasswordHash string `db:"userpasswordhash"`
  UserCash float64        `db:"usercash"`
}

type Transaction struct {
    BuyerId int           `db:"buyerid"`
    SellerId int          `db:"sellerid"`
    Ticker string         `db:"ticker"`
    AmountTraded int      `db:"amounttraded"`
    CashTraded  float64   `db:"cashtraded"`
    TimeOfTrade time.Time `db:"timeoftrade"`
}

type Position struct {
  UserId int                    `db:userid`
  Ticker string                 `db:ticker`
  Amount int                    `db:amount`
  InitialCashInvestment float64 `db:initialcashinvestment`
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

/* 2 args, the sqlx database struct pointer and the userId of the user
 * which we want to retrieve the transactions for. Returns an array of
 * transaction struct.*/
func getAllTransactionsOfUser(db *sqlx.DB,
                              userId int) []Transaction {
    transactions := []Transaction{}
    err := db.Select(&transactions, `select * from transactionTable
                             where buyerId=$1 or sellerId=$1`, userId)
    if (err != nil) {
      log.Fatalln(err)
    }
    return transactions
}

func updatePositionOfUsers(db *sqlx.Db,
                           t Transaction) {

}


/* 3 args, first is the sqlx database struct pointer, the second is
 * the username and the last is the password hash.*/
func createUser(db *sqlx.DB,
                userName string,
                userPasswordHash string,
                startingCash float64) {
    ax := db.MustBegin()
    ax.MustExec(`insert into userTable (userName, userPasswordHash, userCash)
                 values ($1, $2, $3)`, userName, userPasswordHash, startingCash)
    ax.Commit()
}
/* 2 args, first is the sqlx database struct pointer, the second is
 * the userId of the user you wish to remove.*/
func removeUser(db *sqlx.DB,
                userId int) {
    ax := db.MustBegin()
    ax.MustExec(`delete from userTable
                 where userId=$1`, userId)
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


  transactions := getAllTransactionsOfUser(db, 1)
  for _, transaction := range transactions {
    fmt.Printf("ticker: %s\n", transaction.Ticker)
  }
}
