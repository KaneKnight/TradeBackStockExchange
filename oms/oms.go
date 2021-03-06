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
    if orderRequest.UserId != -1 {
      orderRequest.UserId = hash(orderRequest.UserIdString)
    }

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
  c.BindJSON(&request)
  fmt.Println(request.UserIdString)
  if request.UserId != -1 {
    request.UserId = hash(request.UserIdString)
  }

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
        float64(GetLowestAskOfStock(ticker))/100.0,
        float64(GetHighestBidOfStock(ticker))/100.0}
    c.JSON(http.StatusOK, response)
}

func CreateUser(c *gin.Context) {
    var userData db.UserRequest
    c.BindJSON(&userData)
    if userData.UserId != -1 {
      userData.UserId = hash(userData.UserIdString)
    }
    if !db.UserExists(database, userData.UserId) {
      db.CreateUser(database, userData.UserId, 10000000 * 100)
    }
    c.JSON(http.StatusOK, nil)
}

func GetPositionData(c *gin.Context) {
    var positionRequest db.PositionRequest
    c.BindJSON(&positionRequest)
    if positionRequest.UserId != -1 {
      positionRequest.UserId = hash(positionRequest.UserIdString)
    }

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

    response := db.CompanyDataResponse{make([]db.CompanyData, 1)}
    response.CompanyData[0] = db.CompanyData{float64(GetLowestAskOfStock(data.Ticker))/100.0}
    /*response := db.QueryCompanyDataPoints(database, data.Ticker, data.DataNums)
    for i := 0; i < len(response.CompanyData); i++ {
      response.CompanyData[i].Price = Round(response.CompanyData[i].Price, 0.01)
    }*/
    c.JSON(http.StatusOK, response)
}

//API handler that returns the amount of stock a user has for a given company
func GetCompanyInfo(c *gin.Context) {
    var data db.CompanyInfoRequest
    c.BindJSON(&data)
    if data.UserId != -1 {
      data.UserId = hash(data.UserIdString)
    }

    response := db.QueryCompanyInfo(database, data.UserId, data.Ticker)
    c.JSON(http.StatusOK, response)
}

//To be run continuously as a goroutine whilst the platform is functioning
func processOrder() {
    for true {
        fmt.Println(orderQueue.Len())
        var order *Order
        i, _ := orderQueue.Poll(1, -1) //blocks if orderQueue empty
        order = i[0].(*Order)
        priceOfSale := int(order.LimitPrice) * order.NumberOfShares
        fmt.Println("casted: " , int(order.LimitPrice))
        /* Checks if buyer can afford and that the seller can sell.*/
        if ((order.Buy && db.UserCanBuyAmountRequested(database, order.UserId,
            priceOfSale)) ||
            !order.Buy && db.UserCanSellAmountOfShares(database,
                order.UserId, order.CompanyTicker, order.NumberOfShares)) {

            if (order.Buy && order.UserId != -1) {
              fmt.Println("Here:  ", order.UserId)
                db.ReserveCash(database, order.UserId,
                    order.NumberOfShares, int(order.LimitPrice))
            }
            book := bookMap[order.CompanyTicker]
            if book == nil {
                book = InitBook(order.CompanyTicker)
                bookMap[order.CompanyTicker] = book
            }
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
