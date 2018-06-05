package oms

import (
    "github.com/tomdionysus/binarytree"
    "time"
)

/* A mapping of ids to orders.*/
type OrderMap map[int]*Order

/* A mapping of prices to information about the orders at that price.*/
type LimitMap map[LimitPrice]*InfoAtLimit

/* There will be 2 different trees for buy and sell.
 * Order map which maps IDs to Orders.
 * Limit order which maps prices to limits.*/
type Book struct {
    BuyTree    *binarytree.Tree
    SellTree   *binarytree.Tree
    LowestSell *InfoAtLimit
    HighestBuy *InfoAtLimit
    OrderMap   OrderMap
    BuyLimitMap   LimitMap
    SellLimitMap LimitMap
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
        BuyLimitMap:   make(map[LimitPrice]*InfoAtLimit),
        SellLimitMap: make(map[LimitPrice]*InfoAtLimit)}
}

/* Initialises a limit struct with a price and initialises a slice with base
 * length of 10. Fields that are linked to the tree are ignored.*/
func InitLimitInfo (price LimitPrice) *InfoAtLimit {
    return &InfoAtLimit{
        Price:price,
        TotalVolume: 0,
        OrderList: make([]*Order, 10) }
}

/* Initalises order struct with buy or sell, number of shares,
 * limit price and the time the order button was clicked.
 * Fields linked to the tree are ingnored. Also updates the current ID,
 * to allow mapping to orders.*/
func InitOrder(buy bool, numberOfShares int,
    limitPrice LimitPrice, eventTime time.Time) *Order {
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
func (b *Book) InsertOrderIntoBook(order *Order) {
    if (order.Buy) {
        b.insertOrderIntoBuyTree(order)
    } else {
        b.insertOrderIntoSellTree(order)
    }

}


func (b *Book) insertOrderIntoBuyTree(order *Order) {
    buyMap := b.BuyLimitMap
    info := buyMap[order.LimitPrice]
    if (info == nil) {
        b.insertBuyOrderAtNewLimit(order)
    } else {
        b.insertBuyOrderAtLimit(info, order)
    }
}

func (b *Book) insertOrderIntoSellTree(order *Order) {
    info := b.SellLimitMap[order.LimitPrice]
    if (info == nil) {
        b.insertSellOrderAtNewLimit(order)

    } else {
        b.insertSellOrderAtLimit(info, order)
    }
}

func (b *Book) insertBuyOrderAtNewLimit(order *Order) {
    limitPrice := order.LimitPrice
    info := InitLimitInfo(limitPrice)
    b.BuyTree.Set(limitPrice, info)
    b.BuyLimitMap.insertLimitInfoIntoMap(info)
}

func (b *Book) insertSellOrderAtNewLimit(order *Order) {
    limitPrice := order.LimitPrice
    info := InitLimitInfo(limitPrice)
    b.SellTree.Set(limitPrice, info)
    b.SellLimitMap.insertLimitInfoIntoMap(info)
}


func (b *Book) insertBuyOrderAtLimit(limit *InfoAtLimit, order *Order) {
    limit.OrderList = append(limit.OrderList, order)
    b.OrderMap.insertOrderIntoMap(order)
}

func (b *Book) insertSellOrderAtLimit(limit *InfoAtLimit, order *Order) {
    limit.OrderList = append(limit.OrderList, order)
    b.OrderMap.insertOrderIntoMap(order)
}


func (m LimitMap) insertLimitInfoIntoMap(
    limit *InfoAtLimit)  {
    m[limit.Price] = limit
}

func (m OrderMap) insertOrderIntoMap(order *Order) {
    m[order.IdNumber] = order
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
