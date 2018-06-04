package oms

import "time"

var currentId int = 0

type Order struct {
    IdNumber int
    /* Buy is true, sell is false.*/
    BuyOrSell bool
    NumberOfShares int
    /* For bids this is the maximum price, for asks, lowest price.*/
    LimitPrice int
    /* Time that order was inserted into book.*/
    EntryTime time.Time
    /* Time order was placed on website.*/
    EventTime time.Time
    /* Only initialised when order is in map.*/
    ParentLimit *Limit
}

type Limit struct {
    /* Unique identifier that is key of map.*/
    LimitPrice  int
    /* Sum of number of shares in each order.*/
    TotalVolume int
    /* Parent price in tree.*/
    Parent      *Limit
    /* Left child price in tree.*/
    LeftChild   *Limit
    /* Right child price in tree.*/
    RightChild  *Limit
    /* A slice of order pointers. Lower indicies will be earlier orders.
     * Ordered by event time.*/
    OrderList *[]*Order
}

/* There will be 2 different trees for buy and sell.
 * Order map which maps IDs to Orders.
 * Limit order which maps prices to limits.*/
type Book struct {
    BuyTree    *Limit
    SellTree   *Limit
    LowestSell *Limit
    HighestBuy *Limit
    OrderMap *map[int]Order
    LimitMap *map[int]Limit
}

/* Initialises the book struct,
 * maps are created and all other fields are set to nil*/
func InitBook(book *Book) {
    *book.OrderMap = make(map[int]Order)
    *book.LimitMap = make(map[int]Limit)
}

/* Initialises a limit struct with a price and initialises a slice with base
 * length of 10. Fields that are linked to the tree are ignored.*/
func InitLimit (l *Limit, price int) {
    l.LimitPrice = price
    *l.OrderList = make([]*Order, 10)
}

/* Initalises order struct with buy or sell, number of shares,
 * limit price and the time the order button was clicked.
 * Fields linked to the tree are ingnored.*/
func InitOrder(o *Order, buyOrSell bool, numberOfShares int,
    limitPrice int, eventTime time.Time) {
    o.IdNumber = currentId
    currentId += 1
    o.BuyOrSell = buyOrSell
    o.NumberOfShares = numberOfShares
    o.LimitPrice = limitPrice
    o.EventTime = eventTime
}

/* 1 arg, an order to be inserted into the book*/
func (b Book) InsertOrder(order *Order) {

}

/* 1 arg order to be removed from book.*/
func (b Book) CancelOrder(order *Order) {
  //TODO: implement
}

func (b Book) Execute() {
    //TODO: implement
}

func (b Book) GetVolumeAtLimit(limit *Limit) int {
    //TODO: implement
    return 0
}

func (b Book) GetBestBid(limit *Limit) *Limit {
    //TODO: implement
    return nil
}
