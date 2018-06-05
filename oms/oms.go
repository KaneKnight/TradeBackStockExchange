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
)

var database *sqlx.DB
var orderQueue Queue

func InitDB(config db.DBConfig) {
  database = config.OpenDataBase()
}

func AskHandler(c *gin.Context) {
  var order Order
  c.BindJSON(&order)
}

func BidHandler(c *gin.Context) {

}

func generateOrder(order Order) {

}
