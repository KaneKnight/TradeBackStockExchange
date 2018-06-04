package oms

type Order struct {
  IdNumber int
  BuyOrSell bool
  Shares int
  Limit int
  EntryTime int
  EventTime int
  NextOrder *Order
  PrevOrder *Order
  ParentLimit *Limit
}

type Limit struct {
  LimitPrice int
  Size int
  TotalVolume int
  Parent *Limit
  LeftChild *Limit
  RightChild *Limit
  HeadOrder *Order
  TailOrder *Order
}

type Book struct {
  BuyTree *Limit
  SellTree *Limit
  LowestSell *Limit
  HighestBuy *Limit
}

func (b Book) GetVolumeAtLimit(limit *Limit) int {
  
}

func (b Book) GetBestBid(limit *Limit) *Limit {

}
