package oms

import (
    "github.com/tomdionysus/binarytree"
    "time"
    "container/list"

)

type InfoAtLimit struct {
    /* Price of the limit.*/
    Price LimitPrice

    /* The number of shares traded at that price.
     * Updated when match of orders found.*/
    TotalVolume int

    /* The number of shares within this price limit*/
    Size int

    /* A slice of order pointers. Lower indices will be earlier orders.
     * Ordered by event time.*/
    OrderList *list.List

    /* Map of userIds to list of orders*/
    UserOrderMap map[int]*([]**Order)
}

/* Pushes order to list.*/
func (list *[]**Order) PushToList(order *Order)  {
    list = append(list, &order)
}

/* Pops head of list, ie oldest order, returns (true,
 * order) if list is non empty and (false, nil) if empty*/
func (list *[]**Order) PopFromList() (bool, **Order){
    if len(list) > 0 {
        order := list[0]
        list = info.list[1:]
        return true, order
    }
    return false, nil
}

var currentId int = 0

type Order struct {
    IdNumber int
    /* Buy is true, sell is false.*/
    UserId int
    Buy bool
    MarketOrder bool
    CompanyTicker string
    NumberOfShares int
    /* For bids this is the maximum price, for asks, lowest price.*/
    LimitPrice LimitPrice
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


