package main

import (
  //"fmt"
  //"net/http"
  "github.com/gin-gonic/contrib/static"
  "github.com/gin-gonic/gin"
)

func main() {

  //Set default router
  router := gin.Default()

  //Serve frontend static files
  router.Use(static.Serve("/", static.LocalFile("./web", true)))

  //run on default port 8080
  router.Run()
}
