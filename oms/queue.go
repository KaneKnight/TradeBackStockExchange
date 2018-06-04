package oms

type Queue []*Order

func (queue *Queue) Push(order *Order) {
  *queue = append(*queue, order)
}

func (queue *Queue) Pop() (order *Order) {
  order = (*queue)[0]
  *queue = (*queue)[1:]
  return
}

func (queue *Queue) Len() int {
  return len(*queue)
}
