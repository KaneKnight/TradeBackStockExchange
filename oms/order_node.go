package oms

import (
    "github.com/tomdionysus/binarytree"
    "time"
)

type InfoAtLimit struct {
    /* Price of the limit.*/
    Price LimitPrice
    /* The number of shares traded at that price.
     * Updated when match of orders found.*/
    TotalVolume int
    /* A slice of order pointers. Lower indices will be earlier orders.
     * Ordered by event time.*/
    OrderList []*Order
}

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

}

/* Key of the node of the binarytree.*/
type LimitPrice int

func (me LimitPrice) EqualTo(other binarytree.Comparable) bool {
    return me == other.(LimitPrice)
}

func (me LimitPrice) GreaterThan(other binarytree.Comparable) bool {
    return me > other.(LimitPrice)
}

func (me LimitPrice) ValueOf() interface{} {
    return int(me)
}

func (me LimitPrice) LessThan(other binarytree.Comparable) bool {
    return me < other.(LimitPrice)
}


