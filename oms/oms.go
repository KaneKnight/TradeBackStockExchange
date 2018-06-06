package oms

import (
  "fmt"
  //"net/http"
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

func init() {
  database = dbConfig.OpenDataBase()
  orderQueue = queue.New(100)

  //initiate the processor routine
  go processOrder()
}

//orderHandler assume that API is supplied with correct JSON format
func OrderHandler(c *gin.Context) {
  var order Order = Order{101, true, 10, 1001, time.Now(), time.Now(), nil}
  //Binds supplied JSON to Order struct from order_book defs
  //c.BindJSON(&order)
  orderQueue.Put(order)
}

type equityList struct {
  equities []string `json:"equities"`
}

//API handler that returns a list of all equity we serve
func GetEquityList(c *gin.Context) {

}

type equityDataRequest struct {
  equityName string `json:"equityName"`
  dataNums   int    `json:"dataNums"`
}

type equityDataResponse struct {
  equityName string     `json:"equityName"`
  equityData []equityData `json:"data"`
}

type equityData struct {
  time time.Time `json:"time"`
  price int64    `json:"price"`
}

//API handler that returns n number of datapoints for a requested equity
func GetEquityDataPoints(c *gin.Context) {

}

//To be run continuously as a goroutine whilst the platform is functioning
func processOrder() {
  for true {
    var order Order
    i, _ := orderQueue.Poll(1, -1) //blocks if orderQueue empty
    order = i[0].(Order)
    //Process the order, need Kane's stuff...
    time.Sleep(1 * time.Second)
    fmt.Println(order)
  }
}


