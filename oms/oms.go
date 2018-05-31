package oms

import (
  //"fmt"
  "net/http"
  //"encoding/json"
  //"errors"
  //"strconv"
  //"github.com/gin-gonic/contrib/static"
  "github.com/gin-gonic/gin"
  "github.com/louiscarteron/WebApps2018/db"
  "github.com/jmoiron/sqlx"
)

var database *sqlx.DB

func InitDB(config db.DataBase) {
  database, _ = config.OpenDataBase()
}

type test struct{
  UserId int `json:"userId"`
}
func AskHandler(c *gin.Context) {
  var testJSON test
  c.BindJSON(&testJSON)
  c.JSON(http.StatusOK, testJSON)
}

func BidHandler(c *gin.Context) {

}
