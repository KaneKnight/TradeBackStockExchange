package oms

import (
    "net/http"
    //"encoding/json"
    //"errors"
    //"strconv"
    //"github.com/gin-gonic/contrib/static"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/louiscarteron/WebApps2018/db"
    "github.com/jmoiron/sqlx"
    "github.com/Workiva/go-datastructures/queue"
    "fmt"
    "hash/fnv"
)

var dbConfig = db.DBConfig{
    "db.doc.ic.ac.uk",
    "g1727122_u",
    "PTqnydAPoe",
    "g1727122_u",
    5432}

//Represent the connection to the Postgres database provided in the config
var database *sqlx.DB

//A dynamically sized queue of orders submitted from API
//OMS will collect from this queue and process the orders
//Better than a simple channel because it can grow indefinitely without
//the need for buffers
var orderQueue *queue.Queue

var bookMap map[string]*Book

func init() {
    database = dbConfig.OpenDataBase()
    orderQueue = queue.New(100)
    bookMap = make(map[string]*Book)
    db.ZeroReserveCashOfAllUsers(database)

    //initiate the processor routine
    go processOrder()
}

func hash(s string) int {
    h := fnv.New32a()
    h.Write([]byte(s))
    return int(h.Sum32())
}

//orderHandler assume that API is supplied with correct JSON format
func OrderHandler(c *gin.Context) {
    var orderRequest db.OrderRequest
    c.BindJSON(&orderRequest)
    orderRequest.UserId = hash(orderRequest.UserIdString)

    var buy bool
    var market bool
    switch orderRequest.OrderType {
    case "marketBid":
        market = true
        buy = true
    case "marketAsk":
        market = true
        buy = false
    case "limitBid":
        market = false
        buy = true
    case "limitAsk":
        market = false
        buy = false
    }

    price := LimitPrice(orderRequest.LimitPrice * 100)
    order := InitOrder(orderRequest.UserId, buy, market,
        orderRequest.EquityTicker, orderRequest.Amount, price, time.Now())
    orderQueue.Put(order)
    c.JSON(http.StatusOK, nil)
}

func UserTransactionsHandler(c *gin.Context) {
  var request db.UserTransactionsRequest
  var response db.UserTransactionsResponse
  c.BindJSON(&response)
  request.UserId = hash(request.UserIdString)

  response = db.GetAllUserTransactions(database, request.UserId)
  c.JSON(http.StatusOK, response)
}

func CancelHandler(c *gin.Context) {
    var cancelOrder db.CancelOrderRequest
    c.BindJSON(&cancelOrder)
    //TODO: check with bloomberg script

    CancelOrder(&cancelOrder)
    c.JSON(http.StatusOK, nil)
}

func HighestBidLowestAsk(c *gin.Context) {
    var priceRequest db.PriceRequest
    c.BindJSON(&priceRequest)
    ticker := priceRequest.Ticker
    response := db.PriceResponse{
        GetLowestAskOfStock(ticker),
        GetHighestBidOfStock(ticker)}
    c.JSON(http.StatusOK, response)
}

func CreateUser(c *gin.Context) {
    var userData db.UserRequest
    c.BindJSON(&userData)
    userData.UserId = hash(userData.UserIdString)
    if !db.UserExists(database, userData.UserId) {
      db.CreateUser(database, userData.UserId, 10000000 * 100)
    }
    c.JSON(http.StatusOK, nil)
}

func GetPositionData(c *gin.Context) {
    var positionRequest db.PositionRequest
    c.BindJSON(&positionRequest)
    positionRequest.UserId = hash(positionRequest.UserIdString)

    var positionResponse db.PositionResponse
    positionResponse = GetUserPositionsResponse(positionRequest.UserId)

    c.JSON(http.StatusOK, positionResponse)
}

//API handler that returns a list of all equity we serve
func GetCompanyList(c *gin.Context) {
    companyList := db.GetAllCompanies(database)
    c.JSON(http.StatusOK, companyList)
}

//API handler that returns n number of datapoints for a requested equity
func GetCompanyDataPoints(c *gin.Context) {
    var data db.CompanyDataRequest
    c.BindJSON(&data)

    response := db.QueryCompanyDataPoints(database, data.Ticker, data.DataNums)
    for i := 0; i < len(response.CompanyData); i++ {
      response.CompanyData[i].Price = Round(response.CompanyData[i].Price, 0.01)
    }
    c.JSON(http.StatusOK, response)
}

//API handler that returns the amount of stock a user has for a given company
func GetCompanyInfo(c *gin.Context) {
    var data db.CompanyInfoRequest
    c.BindJSON(&data)
    data.UserId = hash(data.UserIdString)

    response := db.QueryCompanyInfo(database, data.UserId, data.Ticker)
    c.JSON(http.StatusOK, response)
}

//To be run continuously as a goroutine whilst the platform is functioning
func processOrder() {
    for true {
        var order *Order
        i, _ := orderQueue.Poll(1, -1) //blocks if orderQueue empty
        order = i[0].(*Order)
        priceOfSale := int(order.LimitPrice) * order.NumberOfShares
        /* Checks if buyer can afford and that the seller can sell.*/
        fmt.Println("Going into big if statement...")
        if ((order.Buy && db.UserCanBuyAmountRequested(database, order.UserId,
            priceOfSale)) ||
            !order.Buy && db.UserCanSellAmountOfShares(database,
                order.UserId, order.CompanyTicker, order.NumberOfShares)) {

            fmt.Println("Going into order.Buy if check...")
            if (order.Buy) {
                fmt.Println("Going into ReserveCash()...")
                db.ReserveCash(database, order.UserId,
                    order.NumberOfShares, int(order.LimitPrice))
            }
            book := bookMap[order.CompanyTicker]
            if book == nil {
                book = InitBook(order.CompanyTicker)
                bookMap[order.CompanyTicker] = book
            }
            fmt.Println("Going into Execute()...")
            success, transactions := book.Execute(order)
            if success {
                length := len(*transactions)
                for j := 0; j < length; j++ {
                    db.InsertTransaction(database, (*transactions)[j])
                    db.UpdatePositionOfUsersFromTransaction(database,
                        (*transactions)[j])
                }
            }
        }
    }
}
