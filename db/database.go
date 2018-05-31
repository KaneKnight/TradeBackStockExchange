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
    cashSpentOnPosition float(53)
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
  UserId int                    `db:"userid"`
  Ticker string                 `db:"ticker"`
  Amount int                    `db:"amount"`
  CashSpentOnPosition float64   `db:"cashspentonposition"`
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
                                          cashTraded, timeOfTrade)delete
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

/* 2 args, the sqlx database struct pointer and the transaction that
 * we need to update the positons of the buyer and the seller.*/
func updatePositionOfUsersFromTransaction(db *sqlx.Db,
                                          t Transaction) {
     ax := db.MustBegin()
     updateBuyerPosition(db, ax, t.BuyerId, t.Ticker,
                         t.AmountTraded, t.CashTraded)
     updateSellerPosition(db, ax, t.SellerId, t.Ticker,
                          t.AmountTraded, t.CashTraded)
     err = ax.Commit()
     if (err != nil) {
       log.Fatalln(err)
     }
}

func updateBuyerPosition(db sqlx.DB,
                         ax sqlx.Tx
                         buyerId int,
                         ticker string,
                         amountTraded int,
                         cashTraded float64) {
   var numberOfPositions int
   err := db.Get(&numberOfPositions , `select (count *) from positionTable
                                       where userId=$1 and ticker=$2`,
                                       t.BuyerId, t.Ticker)
   if (numberOfPositions == 0) {
        createNewPosition(ax, buyerId, ticker, amountTraded, cashTraded)
        //Minus may not recognise.
        updateUserCash(ax, buyerId, -cashTraded)
   } else {

   }
}

func updateSellerPosition(db sqlx.DB,
                          ax sqlx.Tx,
                          sellerId int,
                          ticker string,
                          amountTraded int,
                          cashTraded float64) {
                            //TODO:
}

func createNewPosition(ax sqlx.Tx,
                       buyerId int,
                       ticker string,
                       amountTraded int,
                       cashTraded float64) {
    ax.MustExec(`insert into positionTable (userId,
                                            ticker,
                                            amount,
                                            cashSpentOnPosition)
                 values ($1, $2, $3, $4)`, buyerId,
                                           ticker,
                                           amountTraded,
                                           cashTraded)
}

/* 3 args, the sqlx transaction object pointer, the userId of the user
 * which we want to update their cash and the difference in cash which
 * may be negative.*/
func updateUserCash(ax sqlx.Tx,
                    userId int,
                    cashTraded float64) {
    ax.MustExec(`update userTable
                 set userCash=userCash+$1`, cashTraded)
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
