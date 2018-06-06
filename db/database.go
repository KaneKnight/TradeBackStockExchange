package db

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

type DBConfig struct {
    Host     string
    User     string
    Password string
    Name     string
    Port     int
}

/*--------------TYPE STRUCTS USED FOR QUERIES---------------*/

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

type Company struct {
    Value string `json:"value", db:"ticker"`
    Label string `json:"label", db:"name"`
}

type CompanyList struct {
    Companies []Company `json:"results"`
}

type CompanyDataRequest struct {
    CompanyName string `json:"companyName"`
    DataNums    int    `json:"dataNums"`
}

type CompanyDataResponse struct {
    CompanyName string        `json:"companyName"`
    CompanyData []CompanyData `json:"data"`
}

type CompanyData struct {
    Time time.Time `json:"time"`
    Price int64    `json:"price"`
}

/* No args, called on the DataBase struct and returns a pointer to
 * sqlx database struct. Opens a connection to the database.*/
func (db DBConfig) OpenDataBase() (*sqlx.DB) {
    connStr := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d`,
                           db.Host,
                           db.User,
                           db.Password,
                           db.Name,
                           db.Port)
    return sqlx.MustConnect("postgres", connStr)
}

/* 6 args, first is the sqlx database struct pointer and the rest are
 * the fields of Transaction struct, returns void. Inserts a transaction
 * into the database.*/
func InsertTransaction(db *sqlx.DB, t Transaction) {
    ax := db.MustBegin()
    ax.MustExec(`insert into transactionTable (buyerId, sellerId,
                                          ticker, amountTraded,
                                          cashTraded, timeOfTrade)
                                          values ($1, $2, $3, $4, $5, $6)`,
                                          t.BuyerId, t.SellerId,  t.Ticker,
                                          t.AmountTraded, t.CashTraded, t.TimeOfTrade)
    ax.Commit()
}

/* 2 args, the sqlx database struct pointer and the userId of the user
 * which we want to retrieve the transactions for. Returns an array of
 * transaction struct.*/
func GetAllTransactionsOfUser(db *sqlx.DB,
                              userId int) []Transaction {
    transactions := []Transaction{}
    err := db.Select(&transactions, `select * from transactionTable
                             where buyerId=$1 or sellerId=$1`, userId)
    if (err != nil) {
      log.Fatalln(err)
    }
    return transactions
}

/* 4 args, sqlx database pointer, userId of user requesting sale, ticker of
 * proposed sale and the amount of shares requested to sell. Will return
 * true if user has the required amount, false if not.*/
func UserCanSellAmountOfShares(db *sqlx.DB,
                              userId int,
                              ticker string,
                              requestedAmount int) bool {
    var numberOfSharesOwned int
    err := db.Get(&numberOfSharesOwned, `select amount from positionTable
                                         where userId=$1
                                         and ticker=$2`,
                                         userId,
                                         ticker)
    if (err != nil) {
      return false;
    }
    return (numberOfSharesOwned >= requestedAmount)
}

/* 3 args, sqlx database pointer, userId and the price of the requested sale,
 * returns true if user has enough cash to buy.*/
func UserCanBuyAmountRequested(db *sqlx.DB,
                               userId int,
                               priceOfSale float64) bool {
    var userCash float64
    err := db.Get(&userCash, `select userCash from userTable
                              where userId=$1`, userId)
    if (err != nil) {
      log.Fatalln(err)
    }

    return (priceOfSale <= userCash)
}

/* 2 args, the sqlx database struct pointer and the transaction that
 * we need to update the positons of the buyer and the seller.*/
func UpdatePositionOfUsersFromTransaction(db *sqlx.DB,
                                          t Transaction) {
     ax := db.MustBegin()
     UpdateBuyerPosition(db, ax, t.BuyerId, t.Ticker,
                         t.AmountTraded, t.CashTraded)
     UpdateSellerPosition(db, ax, t.SellerId, t.Ticker,
                          t.AmountTraded, t.CashTraded)
     err := ax.Commit()
     if (err != nil) {
       log.Fatalln(err)
     }
}

/* 6 args, the sqlx pointers, sellerId, ticker followed by the amount of
 * shares traded and the cash exchanged. Will check if user has had the
 * position before and if not will create a new positon, then updates cash.*/
func UpdateBuyerPosition(db *sqlx.DB,
                         ax *sqlx.Tx,
                         buyerId int,
                         ticker string,
                         amountTraded int,
                         cashTraded float64) {
   var numberOfPositions int
   err := db.Get(&numberOfPositions , `select count(*) from positionTable
                                       where userId=$1 and ticker=$2`,
                                       buyerId, ticker)
   if (err != nil) {
      log.Fatalln(err)
   }
   if (numberOfPositions == 0) {
       CreateNewPosition(ax, buyerId, ticker, amountTraded, cashTraded)
   } else {
       UpdatePosition(ax, buyerId, ticker, amountTraded, cashTraded)
   }
   //Minus may not recognise.
   UpdateUserCash(ax, buyerId, -cashTraded)
}

/* 6 args, the sqlx pointers, sellerId, ticker followed by the amount of
 * shares traded and the cash exchanged. Will update the postion of the
 * seller in the database and also change the cash the user has in the
 * user table.*/
func UpdateSellerPosition(db *sqlx.DB,
                          ax *sqlx.Tx,
                          sellerId int,
                          ticker string,
                          amountTraded int,
                          cashTraded float64) {
    UpdatePosition(ax, sellerId, ticker, -amountTraded, -cashTraded)
    UpdateUserCash(ax, sellerId, cashTraded)
}

/* 5 args, sqlx transaction struct pointer, then the userId and the ticker,
 * followed by the amount the user bought or sold. Positive means bought,
 * negative means sold. Creates a new entry for positions users have never,
 * taken before. IMPORTANT: POSITIONS THAT USERS TAKE ARE NOT DELETED
 * WHEN ALL OF THE SHARES ARE SOLD.*/
func CreateNewPosition(ax *sqlx.Tx,
                       userId int,
                       ticker string,
                       amountTraded int,
                       cashTraded float64) {
    ax.MustExec(`insert into positionTable (userId,
                                            ticker,
                                            amount,
                                            cashSpentOnPosition)
                 values ($1, $2, $3, $4)`, userId,
                                           ticker,
                                           amountTraded,
                                           cashTraded)
}

/* 5 args, sqlx transaction struct pointer, then the userId and the ticker,
 * followed by the amount the user bought or sold. Positive means bought,
 * negative means sold. Updates position that already exists.*/
func UpdatePosition(ax *sqlx.Tx,
                    userId int,
                    ticker string,
                    amountTraded int,
                    cashTraded float64) {
    ax.MustExec(`update positionTable
                 set amount=amount+$1,
                     cashSpentOnPosition=cashSpentOnPosition+$2
                 where userId=$3 and ticker=$4`,
                 amountTraded,
                 cashTraded,
                 userId,
                 ticker)
}

/* 3 args, the sqlx transaction object pointer, the userId of the user
 * which we want to update their cash and the difference in cash which
 * may be negative.*/
func UpdateUserCash(ax *sqlx.Tx,
                    userId int,
                    cashTraded float64) {
    ax.MustExec(`update userTable
                 set userCash=userCash+$1
                 where userId=$2`, cashTraded, userId)
}



/* 3 args, first is the sqlx database struct pointer, the second is
 * the username and the last is the password hash.*/
func CreateUser(db *sqlx.DB,
                userName string,
                userPasswordHash string,
                startingCash float64) {
    ax := db.MustBegin()
    ax.MustExec(`insert into userTable (userName, userPasswordHash, userCash)
                 values ($1, $2, $3)`, userName, userPasswordHash, startingCash)
    ax.Commit()
}

func GetAllCompanies(db *sqlx.DB) CompanyList {
    var companyList CompanyList
    err := db.Select(&companyList.Companies, `select * from companyTabel`, nil)
    if (err != nil) {
      log.Fatalln(err)
    }
    return companyList
}

func QueryCompanyDataPoints(db *sqlx.DB, name string, num int) CompanyDataResponse {
    var resp CompanyDataResponse
    resp.CompanyName = name

    //TODO: Division in the SQL Query might cause errors
    err := db.Select(&resp.CompanyData, `select timeOfTrade, cashTraded/amountTraded 
                                         from transactionTable join
                                              companyTable 
                                              using ticker
                                         limit (num) values ($1)`, num)

    if err != nil {
      log.Fatalln(err)
    }
    return resp
}

/* 2 args, first is the sqlx database struct pointer, the second is
 * the userId of the user you wish to remove. TODO: DELETE POSITIONS FROM
 * DATABASE.*/
func RemoveUser(db *sqlx.DB,
                userId int) {
    ax := db.MustBegin()
    ax.MustExec(`delete from userTable
                 where userId=$1`, userId)
    ax.Commit()
}
/*
func main() {
  database := DataBase{"g1727122_u",
     "PTqnydAPoe",
     "g1727122_u",
     "require",
     5432}
  db := database.openDataBase()

  // t := time.Now()
  // transaction := Transaction{
  //   1,
  //   2,
  //   "AAPL",
  //   3,
  //   300,
  //   t}

  canBuy := UserCanBuyAmountRequested(db, 1, 1001)
  fmt.Println(canBuy)
}*/

