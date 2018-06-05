package oms

import (
  "fmt"
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

var database *sqlx.DB

//A dynamically sized queue of orders submitted from API
//OMS will collect from this queue and process the orders
var orderQueue queue.Queue

func InitDB(config db.DBConfig) {
  database = config.OpenDataBase()
}

//orderHandler assume that API is supplied with correct JSON format
func orderHandler(c *gin.Context) {
  var order Order
  //Binds supplied JSON to Order struct from order_book defs
  c.BindJSON(&order)
  orderQueue.Put(order)
}
