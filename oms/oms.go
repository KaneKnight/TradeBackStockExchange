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

    //initiate the processor routine
    go processOrder()
}

//orderHandler assume that API is supplied with correct JSON format
func OrderHandler(c *gin.Context) {
    var orderRequest db.OrderRequest
    c.BindJSON(&orderRequest)

    var buy bool
    var market bool
    switch orderRequest.OrderType {
    case "marketBid":
        market = true
        buy = true

        /*TODO: update thinking for market orders.*/
        db.ReserveCash(database, orderRequest.UserId,
            orderRequest.Amount, int(orderRequest.Price * 100))
    case "marketAsk":
        market = true
        buy = false
    case "limitBid":
        market = false
        buy = true
        db.ReserveCash(database, orderRequest.UserId,
            orderRequest.Amount, int(orderRequest.Price * 100))
    case "limitAsk":
        market = false
        buy = false
    }

    price := LimitPrice(orderRequest.Price * 100)
    order := InitOrder(orderRequest.UserId, buy, market,
        orderRequest.EquityTicker, orderRequest.Amount, price, time.Now())
    orderQueue.Put(order)
    c.JSON(http.StatusOK, nil)
}

func CreateUser(c *gin.Context) {
    var userData db.UserRequest
    c.BindJSON(&userData)

    db.CreateUser(database, userData.UserId, userData.UserName,
        userData.UserCash * 100)
    c.JSON(http.StatusOK, nil)
}

func GetPositionData(c *gin.Context) {
    var positionRequest db.PositionRequest
    c.BindJSON(&positionRequest)

    response := getPositionResponse(positionRequest.EquityTicker,
        positionRequest.UserId)

    c.JSON(http.StatusOK, response)
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

    response := db.QueryCompanyDataPoints(database, data.CompanyName, data.DataNums)
    c.JSON(http.StatusOK, response)
}

//API handler that returns the amount of stock a user has for a given company
func GetCompanyInfo(c *gin.Context) {
    var data db.CompanyInfoRequest
    c.BindJSON(&data)

    response := db.QueryCompanyInfo(database, data.UserId, data.Ticker)
    c.JSON(http.StatusOK, response)
}

//To be run continuously as a goroutine whilst the platform is functioning
func processOrder() {
    for true {
        order1 := InitOrder(5, true, false, "AAPL",1, 5000, time.Now())
        orderQueue.Put(order1)
        var order *Order
        i, _ := orderQueue.Poll(1, -1) //blocks if orderQueue empty
        order = i[0].(*Order)
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


