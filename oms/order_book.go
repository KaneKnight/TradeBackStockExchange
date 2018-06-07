package oms

import (
    "github.com/tomdionysus/binarytree"
    "time"
    "github.com/louiscarteron/WebApps2018/db"
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
    BuyOrder *Order
    SellOrder *Order
}

func (b Book) InsertBuyOrderFake(order *Order) {
    b.BuyOrder = order
}

func (b Book) InsertSellOrderFake(order *Order)  {
    b.SellOrder = order
}

func (b Book) ExecuteFake() *db.Transaction {
    return &db.Transaction{
        BuyerId:b.BuyOrder.UserId,
        SellerId:b.SellOrder.UserId,
        Ticker:b.BuyOrder.CompanyTicker,
        AmountTraded:b.BuyOrder.NumberOfShares,
        CashTraded:int(b.SellOrder.LimitPrice) * b.BuyOrder.NumberOfShares,
        TimeOfTrade:time.Now() }
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
        Size: 0,
        OrderList: make([]*Order, 0) }
}

/* Initalises order struct with buy or sell, number of shares,
 * limit price and the time the order button was clicked.
 * Fields linked to the tree are ingnored. Also updates the current ID,
 * to allow mapping to orders.*/
func InitOrder(userId int, buy bool, companyTicker string, numberOfShares int,
    limitPrice LimitPrice, eventTime time.Time) *Order {
    order := Order{
        IdNumber:currentId,
        UserId:userId,
        Buy:buy,
        CompanyTicker: companyTicker,
        NumberOfShares:numberOfShares,
        LimitPrice:limitPrice,
        EventTime:eventTime}
    currentId += 1
    return &order
}

/* 1 arg, an order to be inserted into the book.
 * This order will be a partial that is the result of the init order function
 * defined above. Will update the lowest sell and highest buy.*/
func (b *Book) InsertOrderIntoBook(order *Order) {
    if (order.Buy) {
        b.insertOrderIntoBuyTree(order)
        b.BuyTree.Balance()
        if (b.HighestBuy.Price < order.LimitPrice) {
            b.HighestBuy = b.BuyLimitMap[order.LimitPrice]
        }
    } else {
        b.insertOrderIntoSellTree(order)
        b.SellTree.Balance()
        if (b.LowestSell.Price > order.LimitPrice) {
            b.LowestSell = b.SellLimitMap[order.LimitPrice]
        }
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
    info.Size += order.NumberOfShares
    b.BuyTree.Set(limitPrice, info)
    b.BuyLimitMap.insertLimitInfoIntoMap(info)
}

/* Creates a new limit pushes the order onto its list and inserts it into the
 * sell binary tree. Adds the limit to the sell map.*/
func (b *Book) insertSellOrderAtNewLimit(order *Order) {
    limitPrice := order.LimitPrice
    info := InitLimitInfo(limitPrice)
    info.pushToList(order)
    info.Size += order.NumberOfShares
    b.SellTree.Set(limitPrice, info)
    b.SellLimitMap.insertLimitInfoIntoMap(info)
}

/* Takes a limit and pushes the order onto its list,
 * then inserts the order into its map. */
func (b *Book) insertOrderAtLimit(limit *InfoAtLimit, order *Order) {
    limit.OrderList = append(limit.OrderList, order)
    b.OrderMap.insertOrderIntoMap(order)
    limit.Size += order.NumberOfShares
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



func (b *Book) Execute(order *Order) (bool, []*db.Transaction) {
    if (order.Buy) {
        return b.MatchBuy(order)
    } else {
        return b.MatchSell(order)
    }
}

func (b *Book) MatchBuy(order *Order) (bool, []*db.Transaction) {
    if (order.LimitPrice >= b.LowestSell.Price) {
        amountLeftToFill := order.NumberOfShares
        currentPrice := b.LowestSell
        for (amountLeftToFill > 0 && currentPrice.Price <= order.LimitPrice) {
            amountLeftToFill -= currentPrice.Size
            isNextPrice, newPrice , _  := b.SellTree.Next(
                currentPrice.
                Price)
            if (isNextPrice) {
                currentPrice = b.SellLimitMap[newPrice.(LimitPrice)]
            } else {
                return false, nil
            }
        }
        if (amountLeftToFill > 0) {
            return false, nil
        }
        return true, b.CalcualteTransactionsBuy(order)
    } else {
        b.InsertOrderIntoBook(order)
        return false, nil
    }
}
func (b *Book) MatchSell(order *Order) (bool, []*db.Transaction) {
    if (order.LimitPrice <= b.HighestBuy.Price) {
        amountLeftToFill := order.NumberOfShares
        currentPrice := b.HighestBuy
        for (amountLeftToFill > 0 && currentPrice.Price >= order.LimitPrice) {
            amountLeftToFill -= currentPrice.Size
            isNextPrice, newPrice , _  := b.BuyTree.Previous(
                currentPrice.
                    Price)
            if (isNextPrice) {
                currentPrice = b.BuyLimitMap[newPrice.(LimitPrice)]
            } else {
                return false, nil
            }
        }
        if (amountLeftToFill > 0) {
            return false, nil
        }
        return true, b.CalcualteTransactionsSell(order)
    } else {
        b.InsertOrderIntoBook(order)
        return false, nil
    }
}
func (b *Book) CalcualteTransactionsBuy(order *Order) []*db.Transaction {
    return nil
}
func (b *Book) CalcualteTransactionsSell(order *Order) []*db.Transaction {
    return nil
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
