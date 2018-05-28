package oms

type Order struct {
  int idNumber
  bool buyOrSell
  int shares
  int limit
  int entryTime
  int eventTime
  Order *nextOrder
  Order *prevOrder
  Limit *parentLimit
}

type Limit struct {
  int limitPrice
  int size
  int totalVolume
  Limit *parent
  Limit *leftChild
  Limit *rightChild
  Order *headOrder
  Order *tailOrder
}

type Book struct {
  Limit *buyTree
  Limit *sellTree
  Limit *lowestSell
  Limit *highestBuy
}
