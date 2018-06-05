package oms

import "time"

var currentId int = 0

type Order struct {
    IdNumber int
    /* Buy is true, sell is false.*/
    Buy bool
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
    /* The number of shares traded at that price.
     * Updated when match of orders found.*/
    TotalVolume int
    /* Parent price in tree.*/
    Parent      *Limit
    /* Left child price in tree.*/
    LeftChild   *Limit
    /* Right child price in tree.*/
    RightChild  *Limit
    /* A slice of order pointers. Lower indices will be earlier orders.
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
 * Fields linked to the tree are ingnored. Also updates the current ID,
 * to allow mapping to orders.*/
func InitOrder(o *Order, buyOrSell bool, numberOfShares int,
    limitPrice int, eventTime time.Time) {
    o.IdNumber = currentId
    currentId += 1
    o.Buy = buyOrSell
    o.NumberOfShares = numberOfShares
    o.LimitPrice = limitPrice
    o.EventTime = eventTime
}

/* 1 arg, an order to be inserted into the book.
 * This order will be a partial that is the result of the init order function
 * defined above.*/
func (b *Book) InsertOrder(order *Order) {
    if(order.Buy) {
        b.BuyTree.insertOrderIntoTree(order)
    } else {
        b.SellTree.insertOrderIntoTree(order)
    }
}

func (tree *Limit) insertOrderIntoTree(order *Order) {
    //TODO: Complete this maybe change insert order above in order to use maps.
    //if root nil make root of tree and add limit to map.
    if (tree == nil) {
        var limit Limit
        InitLimit(&limit, order.LimitPrice)

    } else if (order.LimitPrice < tree.LimitPrice) {
        //Insert into left of tree.
    } else if (order.LimitPrice == tree.LimitPrice){
        //Insert into list of current limit.
    } else {
        //Insert into right of tree.
    }
}

/* 1 arg order to be removed from book.*/
func (b *Book) CancelOrder(order *Order) {
  //TODO: implement
}

func (b *Book) Execute() {
    //TODO: implement
}

func (b *Book) GetVolumeAtLimit(limit *Limit) int {
    //TODO: implement
    return 0
}

func (b *Book) GetBestBid(limit *Limit) *Limit {
    //TODO: implement
    return nil
}
