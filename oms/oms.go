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

func InitDB(config db.DBConfig) {
  database = config.OpenDataBase()
}

func AskHandler(c *gin.Context) {
  var transaction db.Transaction
  c.BindJSON(&transaction)

  transaction.TimeOfTrade = time.Now()

  db.InsertTransaction(database, transaction)

  transactions := db.GetAllTransactionsOfUser(database, 101)
  fmt.Println(transactions)
  fmt.Println("")
  fmt.Println("SIZE OF TRANSACTION TABLE:", len(transactions))
  fmt.Println("")

  c.JSON(http.StatusOK, transaction.TimeOfTrade)
}

func BidHandler(c *gin.Context) {

}
