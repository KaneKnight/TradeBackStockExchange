package oms

import (
    "github.com/tomdionysus/binarytree"
    "time"
)

/* There will be 2 different trees for buy and sell.
 * Order map which maps IDs to Orders.
 * Limit order which maps prices to limits.*/
type Book struct {
    BuyTree    *binarytree.Tree
    SellTree   *binarytree.Tree
    LowestSell *InfoAtLimit
    HighestBuy *InfoAtLimit
    OrderMap   map[int]*Order
    LimitMap   map[LimitPrice]*InfoAtLimit
}

/* Initialises the book struct,
 * maps are created and all other fields are set to nil*/
func InitBook() *Book {
    return &Book{
        BuyTree:    binarytree.NewTree(),
        SellTree:   binarytree.NewTree(),
        LowestSell: nil,
        HighestBuy: nil,
        OrderMap:   make(map[int]*Order),
        LimitMap:   make(map[LimitPrice]*InfoAtLimit) }
}

/* Initialises a limit struct with a price and initialises a slice with base
 * length of 10. Fields that are linked to the tree are ignored.*/
func InitLimitInfo (price int) *InfoAtLimit {
    return &InfoAtLimit{
        TotalVolume: 0,
        OrderList: make([]*Order, 10) }
}

/* Initalises order struct with buy or sell, number of shares,
 * limit price and the time the order button was clicked.
 * Fields linked to the tree are ingnored. Also updates the current ID,
 * to allow mapping to orders.*/
func InitOrder(buy bool, numberOfShares int,
    limitPrice int, eventTime time.Time) *Order {
    order := Order{
        IdNumber:currentId,
        Buy:buy,
        NumberOfShares:numberOfShares,
        LimitPrice:limitPrice,
        EventTime:eventTime}
    currentId += 1
    return &order
}

/* 1 arg, an order to be inserted into the book.
 * This order will be a partial that is the result of the init order function
 * defined above.*/
func (b *Book) InsertOrder(order *Order) {
    if (order.Buy) {
        //insert into buy tree.
    } else {
        //insert into sell tree.
    }
}

func (b *Book) insertOrderIntoTree(order Order) {
    limitPrice := order.LimitPrice

}

func (b *Book) getLimitFromMap(price LimitPrice) (bool, *InfoAtLimit) {
    info := b.LimitMap[price]
    if (info == nil) {
        return false, info
    } else {
        return true, info
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
