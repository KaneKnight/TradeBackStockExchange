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

/* There will be 2 different trees for buy and sell.*/
type Book struct {
    BuyTree    *Limit
    SellTree   *Limit
    LowestSell *Limit
    HighestBuy *Limit
}

/* 1 arg, an order to be inserted into the book*/
func (b Book) InsertOrder(order Order) {
    //TODO: implement
}

/* 1 arg order to be removed from book.*/
func (b Book) CancelOrder(order Order) {
  //TODO: implement
}

func (b Book) Execute() {
    //TODO: implement
}

func (b Book) GetVolumeAtLimit(limit *Limit) int {

}

func (b Book) GetBestBid(limit *Limit) *Limit {

}
