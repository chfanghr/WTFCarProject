package bfs

type Bag struct {
	bag   []int
	index int
}

func NewBag() *Bag {
	b := &Bag{}
	b.bag = make([]int, 0)
	return b
}

func (b *Bag) Add(i int) {
	b.bag = append(b.bag, i)
	b.index++
}

func (b *Bag) IsEmpty(i int) bool { return b.index == 0 }

func (b *Bag) Size() int { return b.index }

func (b *Bag) Iterate() <-chan int {
	ch := make(chan int)
	go func() {
		for _, v := range b.bag {
			ch <- v
		}
		close(ch)
	}()
	return ch
}
