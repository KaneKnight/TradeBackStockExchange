package oms

type Order struct {
    IdNumber       int
    /* Buy is true, sell is false.*/
    BuyOrSell      bool
    NumberOfShares int
    /* For bids this is the maximum price, for asks, lowest price.*/
    LimitPrice int
    /* Time that order was inserted into book.*/
    EntryTime   int
    /* Time order was placed on website.*/
    EventTime   int
    NextOrder   *Order
    PrevOrder   *Order
    ParentLimit *Limit
}

type Limit struct {
    LimitPrice  int
    Size        int
    TotalVolume int
    Parent      *Limit
    LeftChild   *Limit
    RightChild  *Limit
    HeadOrder   *Order
    TailOrder   *Order
}

/* There will be 2 different trees for buy and sell.
Order map which maps IDs to Orders. Limit order which maps prices to limits.*/
type Book struct {
    BuyTree    *Limit
    SellTree   *Limit
    LowestSell *Limit
    HighestBuy *Limit
    OrderMap *map[int]Order
    LimitMap *map[int]Limit
}

func (l Limit) listIsEmpty() bool {
    return (l.HeadOrder == nil && l.TailOrder == nil);
}

func (l Limit) pushOrder(order *Order) {
    if (l.listIsEmpty()) {
        l.HeadOrder = order;
        l.TailOrder = order;
    } else {
        lastOrder := l.TailOrder
        lastOrder.NextOrder = order
        order.PrevOrder = lastOrder
        l.TailOrder = order
    }
}

/* 1 arg, an order to be inserted into the book*/
func (b Book) insertOrder(order *Order) {
    //TODO: implement
}

/* 1 arg order to be removed from book.*/
func (b Book) cancelOrder(order *Order) {
  //TODO: implement
}

func (b Book) execute() {
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
