package oms

import (
  //"fmt"
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
  var order *Order = InitOrder(101, true, 1, 1001, time.Now())
  //Binds supplied JSON to Order struct from order_book defs
  //c.BindJSON(&order)
  orderQueue.Put(*order)
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

  response := db.QueryCompanyInfo(database, data.UserId, data.CompanyName)
  c.JSON(http.StatusOK, response)
}

//To be run continuously as a goroutine whilst the platform is functioning
func processOrder() {
  for true {
    var order Order
    i, _ := orderQueue.Poll(1, -1) //blocks if orderQueue empty
    order = i[0].(Order)
    //Process the order, need Kane's stuff...
    success, transactions := book.Execute(order)
    if success {
      //put into db
    }


    /*time.Sleep(1 * time.Second)
    fmt.Println(order)
    fmt.Println(orderQueue.Len())*/
  }
}


