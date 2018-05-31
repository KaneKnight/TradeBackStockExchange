package oms

import (
  //"fmt"
  "net/http"
  //"encoding/json"
  //"errors"
  //"strconv"
  //"github.com/gin-gonic/contrib/static"
  "github.com/gin-gonic/gin"
)

func AskHandler(c *gin.Context) {
  //params := c.Params
  //userId, _ := strconv.Atoi(params.ByName("userId"))
  //c.Header("Content-Type", "application/json")
  c.JSON(http.StatusOK, "Hi Louis")
}

func BidHandler(c *gin.Context) {

}
