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
  //"fmt"
  //"github.com/streadway/amqp"
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

//Order book instance
var book *Book

func init() {
  database = dbConfig.OpenDataBase()
  orderQueue = queue.New(100)
  book = InitBook()

  //initiate the processor routine
  go processOrder()
}

//orderHandler assume that API is supplied with correct JSON format
func OrderHandler(c *gin.Context) {
  var orderRequest db.OrderRequest
  c.BindJSON(&orderRequest)

  //TODO:Improve
  var buyOrSell bool
  if orderRequest.OrderType == "marketBid" {
    buyOrSell = true
  } else {
    buyOrSell = false
  }

  order := InitOrder(orderRequest.UserId, buyOrSell, orderRequest.EquityTicker, orderRequest.Amount, 10, time.Now())
  orderQueue.Put(order)
  c.JSON(http.StatusOK, nil)
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
      var order *Order
      i, _ := orderQueue.Poll(1, -1) //blocks if orderQueue empty
      order = i[0].(*Order)
      success, transaction := ExecuteFake(book, order)
      //Process the order, need Kane's stuff...
      //success, transactions := book.Execute(order)
      if success {
          db.InsertTransaction(database, *transaction)
          db.UpdatePositionOfUsersFromTransaction(database, *transaction)
      }
  }
}