package main

import (
  "fmt"
  "github.com/louiscarteron/WebApps2018/oms"
)

func main() {
  order := oms.Order{1, true, 2, 3, 4, 5, nil, nil, nil}
  fmt.Printf("%d, %b\n" , order.IdNumber, order.BuyOrSell)
}
