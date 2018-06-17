package main

import (
  "github.com/auth0/go-jwt-middleware"
  "github.com/dgrijalva/jwt-go"
  "os"
  "errors"
  //"log"
  "runtime"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/static"
  "net/http"
  "github.com/louiscarteron/WebApps2018/oms"
  "fmt"
  "encoding/json"
)

//Jwks stores a slice of JSON Web Keys
type Jwks struct {
  Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
  Kty string   `json:"kty"`
  Kid string   `json:"kid"`
  Use string   `json:"use"`
  N   string   `json:"n"`
  E   string   `json:"e"`
  X5c []string `json:"x5c"`
}

//Use by passing to route definitions, along with the handler
var jwtMiddleWare *jwtmiddleware.JWTMiddleware


func main() {

  /*jwtMiddleWare_temp := jwtmiddleware.New(jwtmiddleware.Options{
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      aud := os.Getenv("AUTH0_API_AUDIENCE")
      checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
      if !checkAudience {
        return token, errors.New("Invalid audience.")
      }

      //verify ISS claim
      iss := os.Getenv("AUTH0_DOMAIN")
      checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
      if !checkIss {
        //return token, errors.New("Invalid issuer.")
      }

      cert, err := getPemCert(token)
      if err != nil {
        log.Fatalf("could not get certificate: %+v", err)
      }

      result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
      return result, nil
    },

    SigningMethod: jwt.SigningMethodRS256,
  })

  //Assign global jwtMiddleWare
  jwtMiddleWare = jwtMiddleWare_temp*/

  //Allow go runtime to utilise 2 CPU cores
  runtime.GOMAXPROCS(2)

  //Set default router
  router := gin.Default()

  router.LoadHTMLGlob("web/*.html")
  router.Use(static.Serve("/", static.LocalFile("./web", false)))

  router.GET("/", func(c *gin.Context){
    c.HTML(http.StatusOK, "index.html", nil)
  })

  //Set up API routing
  api := router.Group("/api")

  api.POST("/bid", oms.OrderHandler)
  api.POST("/ask", oms.OrderHandler)
  api.POST("/cancel-orders", oms.CancelHandler)
  api.GET("/get-company-list", oms.GetCompanyList)
  api.POST("/highest-bid-lowest-ask", oms.HighestBidLowestAsk)
  api.GET("/get-datapoints", oms.GetCompanyDataPoints)
  api.POST("/get-company-info", oms.GetCompanyInfo)
  api.POST("/check-user-exists", oms.CreateUser)
  api.POST("/get-all-user-positions", oms.GetPositionData)
  api.POST("/order", oms.OrderHandler)
  api.POST("/get-transaction-history", oms.UserTransactionsHandler)


  //run on default port 8080
  router.Run()
}

//authMiddleWare intercepts the requests, and checks for a valid token
func authMiddleWare() gin.HandlerFunc {
  return func(c *gin.Context) {
    //Get the client secret key
    err := jwtMiddleWare.CheckJWT(c.Writer, c.Request)
    if err != nil {
      //Token not found
      fmt.Println(err)
      c.Abort()
      c.Writer.WriteHeader(http.StatusUnauthorized)
      c.Writer.Write([]byte("Unauthorized"))
      return
    }
  }
}

func getPemCert(token *jwt.Token) (string, error) {
  cert := ""
  resp, err := http.Get(os.Getenv("AUTH0_DOMAIN") + ".well-known/jwks.json")
  if err != nil {
    return cert, err
  }
  defer resp.Body.Close()

  var jwks = Jwks{}
  err = json.NewDecoder(resp.Body).Decode(&jwks)

  if err != nil {
    return cert, err
  }

  x5c := jwks.Keys[0].X5c
  for k, v := range x5c {
    if token.Header["kid"] == jwks.Keys[k].Kid {
      cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
    }
  }

  if cert == "" {
    return cert, errors.New("Unable to find appropriate key.")
  }

  return cert, nil

}
