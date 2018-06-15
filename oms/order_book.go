package oms

import (
    "math"
    "github.com/tomdionysus/binarytree"
    "time"
    "github.com/louiscarteron/WebApps2018/db"
    "fmt"
    "container/list"
)

/* A mapping of ids to orders.*/
type OrderMap map[int]*Order

/* A mapping of prices to information about the orders at that price.*/
type LimitMap map[LimitPrice]*InfoAtLimit

/* There will be 2 different trees for buy and sell.
 * Order map which maps IDs to Orders.
 * Limit order which maps prices to limits.*/
type Book struct {
    BuyTree      *binarytree.Tree
    SellTree     *binarytree.Tree
    LowestSell   *InfoAtLimit
    HighestBuy   *InfoAtLimit
    OrderMap     OrderMap
    BuyLimitMap  LimitMap
    SellLimitMap LimitMap
    TickerOFBook string
}

func InitTransaction(buyerId int, sellerId int,
    ticker string, amount int, cashTraded int,
    timeOfTrade time.Time) *db.Transaction {
    return &db.Transaction{
        BuyerId:      buyerId,
        SellerId:     sellerId,
        Ticker:       ticker,
        AmountTraded: amount,
        CashTraded:   cashTraded,
        TimeOfTrade:  timeOfTrade}
}

/* Initialises the book struct,
 * maps are created and all other fields are set to nil*/
func InitBook(ticker string) *Book {
    return &Book{
        BuyTree:      binarytree.NewTree(),
        SellTree:     binarytree.NewTree(),
        LowestSell:   nil,
        HighestBuy:   nil,
        OrderMap:     make(map[int]*Order),
        BuyLimitMap:  make(map[LimitPrice]*InfoAtLimit),
        SellLimitMap: make(map[LimitPrice]*InfoAtLimit),
        TickerOFBook: ticker}
}

/* Initialises a limit struct with a price and initialises a slice with base
 * length of 10. Fields that are linked to the tree are ignored.*/
func InitLimitInfo(price LimitPrice) *InfoAtLimit {
  return &InfoAtLimit{
        Price:       price,
        TotalVolume: 0,
        Size:        0,
        OrderList:   list.New(),
        UserOrderMap make(map[int]*[]**Order)
}

/* Initalises order struct with buy or sell, number of shares,
 * limit price and the time the order button was clicked.
 * Fields linked to the tree are ingnored. Also updates the current ID,
 * to allow mapping to orders.*/
func InitOrder(userId int, buy bool, marketOrder bool, companyTicker string,
    numberOfShares int,
    limitPrice LimitPrice, eventTime time.Time) *Order {
    order := Order{
        IdNumber:       currentId,
        UserId:         userId,
        Buy:            buy,
        MarketOrder:    marketOrder,
        CompanyTicker:  companyTicker,
        NumberOfShares: numberOfShares,
        LimitPrice:     limitPrice,
        EventTime:      eventTime}
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
        if (b.HighestBuy == nil || b.HighestBuy.Price < order.LimitPrice) {
            b.HighestBuy = b.BuyLimitMap[order.LimitPrice]
        }
    } else {
        b.insertOrderIntoSellTree(order)
        b.SellTree.Balance()
        if (b.LowestSell == nil || b.LowestSell.Price > order.LimitPrice) {
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
    info.OrderList.PushBack(order)
    info.Size += order.NumberOfShares
    b.BuyTree.Set(limitPrice, info)
    b.BuyLimitMap.insertLimitInfoIntoMap(info)
    if (info.UserOrderMap[order.UserId] == nil) {
      list := &make([]**Order)
      info.UserOrderMap[order.UserId] = list
    }
    info.UserOrderMap[order.UserId].PushToList(order)
}

/* Creates a new limit pushes the order onto its list and inserts it into the
 * sell binary tree. Adds the limit to the sell map.*/
func (b *Book) insertSellOrderAtNewLimit(order *Order) {
    limitPrice := order.LimitPrice
    info := InitLimitInfo(limitPrice)
    info.OrderList.PushBack(order)
    info.Size += order.NumberOfShares
    b.SellTree.Set(limitPrice, info)
    b.SellLimitMap.insertLimitInfoIntoMap(info)
    if (info.UserOrderMap[order.UserId] == nil) {
      list := &make([]**Order)
      info.UserOrderMap[order.User] = list
    }
    info.UserOrderMap[order.UserId].PushToList(order)
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
func (m LimitMap) insertLimitInfoIntoMap(limit *InfoAtLimit) {
    m[limit.Price] = limit
}

/* Inserts the order given into the books order map.
 * Mapping order id to order.*/
func (m OrderMap) insertOrderIntoMap(order *Order) {
    m[order.IdNumber] = order
}

func (b *Book) Execute(order *Order) (bool,
    *[]*db.Transaction) {
    if (order.Buy) {
        return b.MatchBuy(order)
    } else {
        return b.MatchSell(order)
    }
}

func (b *Book) MatchBuy(order *Order) (bool,
    *[]*db.Transaction) {
    //TODO:Remove
    fmt.Println(b.CanFillBuyOrder(order))
    if (b.CanFillBuyOrder(order)) {
        return true, b.CalculateTransactionsBuy(order)
    } else {
        b.InsertOrderIntoBook(order)
        return false, nil
    }
}

func (b *Book) MatchSell(order *Order) (bool,
    *[]*db.Transaction) {
    if (b.canFillSellOrder(order)) {
        return true, b.CalculateTransactionsSell(order)
    } else {
        b.InsertOrderIntoBook(order)
        return false, nil
    }
}

func (b *Book) CanFillBuyOrder(order *Order) bool {
    amountLeftToFill := order.NumberOfShares
    currentPrice := b.LowestSell
    if (currentPrice == nil) {
        return false
    }
    for (amountLeftToFill > 0 && (currentPrice.Price <= order.
        LimitPrice || order.MarketOrder)) {
        amountLeftToFill -= currentPrice.Size
        isNextPrice, newPrice, _ := b.SellTree.Next(
            currentPrice.
                Price)
        if (isNextPrice) {
            currentPrice = b.SellLimitMap[newPrice.(LimitPrice)]
        } else {
            break
        }
    }
    return amountLeftToFill <= 0
}

func (b *Book) canFillSellOrder(order *Order) bool {
    amountLeftToFill := order.NumberOfShares
    currentPrice := b.HighestBuy
    if (currentPrice == nil) {
        return false
    }
    for (amountLeftToFill > 0 && (currentPrice.Price >= order.
        LimitPrice || order.MarketOrder)) {
        amountLeftToFill -= currentPrice.Size
        isNextPrice, newPrice, _ := b.SellTree.Previous(
            currentPrice.
                Price)
        if (isNextPrice) {
            currentPrice = b.BuyLimitMap[newPrice.(LimitPrice)]
        } else {
            break
        }
    }
    return amountLeftToFill <= 0
}

func (b *Book) CalculateTransactionsBuy(order *Order) *[]*db.Transaction {
    amountLeftToFill := order.NumberOfShares
    currentPrice := b.LowestSell
    transactions := make([]*db.Transaction, 0)
    for (amountLeftToFill > 0) {
        exists, sellOrder := currentPrice.popFromList()
        if exists {
            amountLeftToFill -= sellOrder.NumberOfShares
            cashTraded := order.NumberOfShares * int(currentPrice.Price)
            transaction := InitTransaction(order.UserId, sellOrder.UserId,
                order.CompanyTicker, sellOrder.NumberOfShares, cashTraded,
                time.Now())
            transactions = append(transactions, transaction)
            currentPrice.TotalVolume += sellOrder.NumberOfShares
            if (amountLeftToFill < 0) {
                sellOrder.NumberOfShares = -1 * amountLeftToFill
                currentPrice.OrderList = append(currentPrice.OrderList,
                    sellOrder)
            }
        } else {
            isNextPrice, newPrice, _ := b.SellTree.Next(currentPrice.Price)
            if (isNextPrice) {
                currentPrice = b.SellLimitMap[newPrice.(LimitPrice)]
            } else {
                return &transactions
            }
        }
    }
    return &transactions
}

func (b *Book) CalculateTransactionsSell(order *Order) *[]*db.Transaction {
    amountLeftToFill := order.NumberOfShares
    currentPrice := b.HighestBuy
    transactions := make([]*db.Transaction, 0)
    for (amountLeftToFill > 0) {
        exists, buyOrder := currentPrice.popFromList()
        if exists {
            amountLeftToFill -= buyOrder.NumberOfShares
            cashTraded := order.NumberOfShares * int(currentPrice.Price)
            transaction := InitTransaction(buyOrder.UserId, order.UserId,
                order.CompanyTicker, buyOrder.NumberOfShares, cashTraded,
                time.Now())
            transactions = append(transactions, transaction)
            currentPrice.TotalVolume += buyOrder.NumberOfShares
            if (amountLeftToFill < 0) {
                buyOrder.NumberOfShares = -1 * amountLeftToFill
                currentPrice.OrderList = append(currentPrice.OrderList,
                    buyOrder)
            }
        } else {
            isNextPrice, newPrice, _ := b.BuyTree.Next(currentPrice.Price)
            if (isNextPrice) {
                currentPrice = b.BuyLimitMap[newPrice.(LimitPrice)]
            } else {
                return &transactions
            }
        }
    }
    return &transactions
}

func GetHighestBidOfStock(ticker string) int {
    if bookMap[ticker] == nil {
        return 0
    }
    return int(bookMap[ticker].HighestBuy.Price)
}

func GetLowestAskOfStock(ticker string) int {
    if bookMap[ticker] == nil {
        return 0
    }
    return int(bookMap[ticker].LowestSell.Price)
}

func getValueAndGain(ticker string, userId int) (int, int) {
    position := db.GetPosition(database, ticker, userId)
    cashSpent := position.CashSpentOnPosition
    currentPriceOfStock := GetHighestBidOfStock(ticker)
    number := position.Amount
    value := currentPriceOfStock * number
    return value, (int(value / cashSpent) - 1) * 100
}

func getPositionResponse(ticker string, userId int) db.PositionResponse {
    position := db.GetPosition(database, ticker, userId)
    currentPriceOfStock := float64(GetHighestBidOfStock(ticker)) / 100
    var value float64
    var gain float64
    if (&position != nil && currentPriceOfStock != 0) {
        value = float64(position.Amount) * currentPriceOfStock
        cashSpent := float64(position.CashSpentOnPosition)
        gain = ((value / cashSpent) - 1) * 100
        gain = Round(gain, 0.01)
    } else {
        value = 0
        gain = 0
    }
    return db.PositionResponse{
        ticker,
        position.Amount,
        value,
        gain}
}

func CancelOrder(cancelRequest *db.CancelOrderRequest) {
    //TODO:Implement
}

func Round(x float64, unit float64) float64 {
    return math.Round(x/unit) * unit
}
