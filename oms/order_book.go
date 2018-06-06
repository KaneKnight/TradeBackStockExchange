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
        OrderList: make([]*Order, 0) }
}

/* Initalises order struct with buy or sell, number of shares,
 * limit price and the time the order button was clicked.
 * Fields linked to the tree are ingnored. Also updates the current ID,
 * to allow mapping to orders.*/
func InitOrder(userId int, buy bool, numberOfShares int,
    limitPrice LimitPrice, eventTime time.Time) *Order {
    order := Order{
        IdNumber:currentId,
        UserId:userId,
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
        b.BuyTree.Balance()
    } else {
        b.insertOrderIntoSellTree(order)
        b.SellTree.Balance()
    }
}

/* Auxiliary helper to insert. Checks if limit is in the buy map. If not,
 * a new limit is created and the order is pushed on. */
func (b *Book) insertOrderIntoBuyTree(order *Order) {
    info := b.BuyLimitMap[order.LimitPrice]
    if (info == nil) {
        b.insertBuyOrderAtNewLimit(order)
    } else {
        b.insertOrderAtLimit(info, order)
    }
}

/* Auxiliary helper to insert. Checks if limit is in the sell map. If not,
 * a new limit is created and the order is pushed on. */
func (b *Book) insertOrderIntoSellTree(order *Order) {
    info := b.SellLimitMap[order.LimitPrice]
    if (info == nil) {
        b.insertSellOrderAtNewLimit(order)
    } else {
        b.insertOrderAtLimit(info, order)
    }
}

/* Creates a new limit pushes the order onto its list and inserts it into the
 * buy binary tree. Adds the limit to the buy map.*/
func (b *Book) insertBuyOrderAtNewLimit(order *Order) {
    limitPrice := order.LimitPrice
    info := InitLimitInfo(limitPrice)
    info.pushToList(order)
    b.BuyTree.Set(limitPrice, info)
    b.BuyLimitMap.insertLimitInfoIntoMap(info)
}

/* Creates a new limit pushes the order onto its list and inserts it into the
 * sell binary tree. Adds the limit to the sell map.*/
func (b *Book) insertSellOrderAtNewLimit(order *Order) {
    limitPrice := order.LimitPrice
    info := InitLimitInfo(limitPrice)
    info.pushToList(order)
    b.SellTree.Set(limitPrice, info)
    b.SellLimitMap.insertLimitInfoIntoMap(info)
}

/* Takes a limit and pushes the order onto its list,
 * then inserts the order into its map. */
func (b *Book) insertOrderAtLimit(limit *InfoAtLimit, order *Order) {
    limit.OrderList = append(limit.OrderList, order)
    b.OrderMap.insertOrderIntoMap(order)
}

/* Inserts the limit given into the map.
Mapping limit price to the limit info.*/
func (m LimitMap) insertLimitInfoIntoMap(limit *InfoAtLimit)  {
    m[limit.Price] = limit
}

/* Inserts the order given into the books order map.
 * Mapping order id to order.*/
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
/*
func (b *Book) GetVolumeAtLimit(limit *Limit) int {
    //TODO: implement
    return 0
}

func (b *Book) GetBestBid(limit *Limit) *Limit {
    //TODO: implement
    return nil
} */
