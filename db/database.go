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
    buyerId bigint,
    sellerId bigint,
    ticker text,
    amountTraded integer,
    cashTraded integer,
    timeOfTrade timestamp
);

create table positionTable (
    userId bigint,
    ticker text,
    amount integer,
    cashSpentOnPosition integer
);

create table userTable (
    userId bigint,
    userCash integer,
    cashReserved integer
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

type UserRequest struct {
    UserIdString string `json:"userIdString"`
    UserId       int
}


type User struct {
    UserId int              `db:"userid"`
    UserName string         `db:"username"`
    UserCash int        `db:"usercash"`
}

type Transaction struct {
    BuyerId int           `db:"buyerid"`
    SellerId int          `db:"sellerid"`
    Ticker string         `db:"ticker"`
    AmountTraded int      `db:"amounttraded"`
    CashTraded  int       `db:"cashtraded"`
    TimeOfTrade time.Time `db:"timeoftrade"`
}

type Position struct {
    UserId int                    `db:"userid"`
    Ticker string                 `db:"ticker"`
    Amount int                    `db:"amount"`
    CashSpentOnPosition int       `db:"cashspentonposition"`
}

type OrderRequest struct {
  UserIdString string `json:"userIdString"`
  EquityTicker string `json:"equityTicker"`
  Amount int `json:"amount"`
  OrderType string `json:"orderType"`
  LimitPrice float64  `json:"limitPrice"`
  UserId int `json:"userId"`
}

type CancelOrderRequest struct {
    LimitPrice int `json:"limitPrice"`
    UserIdString string `json:"userIdString"`
    Ticker string  `json:"ticker"`
    Bid    bool    `json:"bid"`
    UserId int `json:"userId"`
}

type PriceRequest struct {
    Ticker string `json:"ticker"`
}

type PriceResponse struct {
    LowestAsk  float64 `json:"lowestAsk"`
    HighestBid float64 `json:"highestBid"`
}


type PositionResponse struct {
  Positions []JSONPosition `json:"positions"`
}

type JSONPosition struct {
  Ticker string  `json:"ticker"`
  Amount int     `json:"numberOfSharesOwned"`
  Value  float64 `json:"valueOfPosition"`
  Gain   float64 `json:"percentageGain"`
  Name   string  `json:"name"`
}

type PositionRequest struct {
    UserIdString string `json:"userIdString"`
    UserId       int
}

type Company struct {
    Value string `db:"ticker"`
    Label string `db:"name"`
}

type CompanyList struct {
    Companies []Company `json:"results"`
}

type CompanyDataRequest struct {
    Ticker   string `json:"ticker"`
    DataNums int    `json:"dataNums"`
}

type CompanyDataResponse struct {
    CompanyData []CompanyData `json:"data"`
}

type CompanyData struct {
    Price float64 `json: "price"`
}

type CompanyInfoRequest struct {
  UserIdString string `json:"userIdString"`
  Ticker       string `json:"Ticker"`
  UserId       int    `json:"userId"`
}

type CompanyInfoResponse struct {
  Amount      int    `db:"amount"`
}

type UserTransactionsRequest struct {
  UserIdString string `json:"userIdString"`
  UserId       int
}

type UserTransactionsResponse struct {
  BuyTransactions []UserTransaction  `json:"BuyTransactions"`
  SellTransactions []UserTransaction `json:"SellTransactions"`
}

type UserTransaction struct {
  Ticker       string  `json:"ticker"`
  AmountTraded int     `json:"amountTraded"`
  CashTraded   int     `json:"cashSpent"`
  Price        float64 `json:"price"`
  TimeOfTrade  string  `json:"time"`
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
func InsertTransaction(db *sqlx.DB, t *Transaction) {
    ax := db.MustBegin()
    fmt.Println(*t)
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
    if userId == -1 {
      return true
    }
    var numberOfSharesOwned int
    err := db.Get(&numberOfSharesOwned, `select amount from positionTable
                                         where userId=$1
                                         and ticker=$2`,
                                         userId,
                                         ticker)
    if (err != nil) {
      log.Fatalln(err)
    }
    return (numberOfSharesOwned >= requestedAmount)
}

/* 3 args, sqlx database pointer, userId and the price of the requested sale,
 * returns true if user has enough cash to buy.*/
func UserCanBuyAmountRequested(db *sqlx.DB,
                               userId int,
                               priceOfSale int) bool {
    if userId == -1 {
      return true
    }

    var userCash int
    err := db.Get(&userCash, `select userCash - cashReserved from userTable
                              where userId=$1`, userId)
    if (err != nil) {
      log.Fatalln(err)
    }

    return (priceOfSale <= userCash)
}

/* 2 args, the sqlx database struct pointer and the transaction that
 * we need to update the positons of the buyer and the seller.*/
func UpdatePositionOfUsersFromTransaction(db *sqlx.DB,
                                          t *Transaction) {
     ax := db.MustBegin()
     if t.BuyerId != -1 {
        UpdateBuyerPosition(db, ax, t.BuyerId, t.Ticker, t.AmountTraded, t.CashTraded)
     }
     if (t.SellerId != -1) {
        UpdateSellerPosition(db, ax, t.SellerId, t.Ticker, t.AmountTraded, t.CashTraded)
     }
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
                         cashTraded int) {
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

   ax.MustExec(`update userTable
                       set cashReserved=cashReserved-$1
                       where userId=$2`, cashTraded, buyerId)

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
                          cashTraded int) {
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
                       cashTraded int) {
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
                    cashTraded int) {
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
                    cashTraded int) {
    ax.MustExec(`update userTable
                 set userCash=userCash+$1
                 where userId=$2`, cashTraded, userId)
}



/* 3 args, first is the sqlx database struct pointer, the second is
 * the username and the last is the password hash.*/
func CreateUser(db *sqlx.DB,
                userId int,
                startingCash int) {
    db.MustExec(`insert into userTable (userId, userCash, cashReserved)
                 values ($1, $2, $3)`, userId, startingCash, 0)
}

func UserExists(db *sqlx.DB, userId int) bool {
  var num int
  err := db.Get(&num, `select count(*) from userTable where userId=$1`, userId)
  if err != nil {
    log.Fatalln(err)
  }
  return num != 0
}

func ReserveCash(db *sqlx.DB,
                 userId int,
                 numberOfShares int,
                 limitPrice int) {
    fmt.Println(userId, numberOfShares, limitPrice)
    db.MustExec(`update userTable
                 set cashReserved=cashReserved+$1
                 where userId=$2`, numberOfShares * limitPrice, userId)
}

func ZeroReserveCashOfAllUsers(db *sqlx.DB) {
    db.MustExec(`update userTable
                 set cashReserved=0`)
}

func GetAvailableCash(db *sqlx.DB,
                      userId int) int {
    var available int
    err := db.Select(&available, `select userCash - cashReserved from userTable
                                 where userId=$1`, userId)
    if (err != nil) {
        log.Fatalln(err)
    }
    return available
}

func GetPosition(db *sqlx.DB, ticker string, userId int) Position {
    var position []Position
    err := db.Select(&position, `select * from positionTable
                                 where userid=$1 and ticker=$2`, userId, ticker)
    if (err != nil) {
        log.Fatalln(err)
    }
    return position[0]
}

type Pos struct {
  Ticker string
  Name   string
}
func GetAllUserPositions(db *sqlx.DB, userId int) []Pos {
  var positions []Pos
  err := db.Select(&positions, `select positionTable.ticker, name from positionTable join companyTable on positionTable.ticker=companyTable.ticker where userId=$1`, userId)
  if err != nil {
    log.Fatalln(err)
  }

  return positions
}

func GetAllUserTransactions(db *sqlx.DB, userId int) UserTransactionsResponse {
  var response UserTransactionsResponse
  fmt.Println(userId)
  err1 := db.Select(&response.BuyTransactions, `select ticker, amountTraded, cashTraded,
        cast(cashTraded as float(53))/cast(amountTraded as float(53))/100 as Price,
        timeOfTrade 
        from transactionTable
        where buyerId=$1`, userId)
  if err1 != nil {
    log.Fatalln(err1)
  }
  err2 := db.Select(&response.SellTransactions, `select ticker, amountTraded, cashTraded,
        cast(cashTraded as float(53))/cast(amountTraded as float(53))/100 as Price,
        timeOfTrade 
        from transactionTable
        where sellerId=$1`, userId)
  if err2 != nil {
    log.Fatalln(err2)
  }
  fmt.Println(response)
  return response
}

func GetAllCompanies(db *sqlx.DB) CompanyList {
    var companyList CompanyList
    err := db.Select(&companyList.Companies, `select * from companyTable`)
    if (err != nil) {
      log.Fatalln(err)
    }
    return companyList
}

func QueryCompanyDataPoints(db *sqlx.DB, ticker string, num int) CompanyDataResponse {
    var resp CompanyDataResponse

    err := db.Select(&resp.CompanyData,
        `select cast(cashTraded as float(53))/cast(amountTraded as float(53))/100 as Price
               from transactionTable
               where transactionTable.ticker=$1
               limit $2`, ticker, num)

    if err != nil {
      log.Fatalln(err)
    }
    return resp
}

func QueryCompanyInfo(db *sqlx.DB, userId int, ticker string) CompanyInfoResponse {
  var resp CompanyInfoResponse
  err := db.Get(&resp.Amount, `select amount
                                  from positionTable
                                  where ticker=$1 and userId=$2`,
                                  ticker, userId)

  if err != nil {
    resp.Amount = 0
  }
  return resp
}

/* 2 args, first is the sqlx database struct pointer, the second is
 * the userId of the user you wish to remove. TODO: DELETE POSITIONS FROM
 * DATABASE.*/
func RemoveUser(db *sqlx.DB,
                userId int) {
    db.MustExec(`delete from userTable
                 where userId=$1`, userId)
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

