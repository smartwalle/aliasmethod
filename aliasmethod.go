package aliasmethod

import (
	"container/list"
	"errors"
	"math/rand"
	"time"
)

const (
	kProbability = 1.0
)

type Item interface {
	GetWeight() int32
}

type AliasMethod[T Item] struct {
	alias       []int
	probability []float64
	items       []T
	r           *rand.Rand
}

func New[T Item]() *AliasMethod[T] {
	var m = &AliasMethod[T]{}
	m.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return m
}

func (this *AliasMethod[T]) Add(item T) {
	//if item == nil {
	//	return
	//}
	this.items = append(this.items, item)
}

func (this *AliasMethod[T]) Prepare() error {
	if len(this.items) == 0 {
		return errors.New("概率不能为空")
	}

	var total = int32(0)
	for _, item := range this.items {
		total += item.GetWeight()
	}

	var scale = float64(total) / kProbability

	var values = make([]float64, 0, len(this.items))
	for _, item := range this.items {
		values = append(values, float64(item.GetWeight())/scale)
	}

	return this.process(values)
}

func (this *AliasMethod[T]) process(prob []float64) error {
	var p = make([]float64, len(prob))
	copy(p, prob)

	this.alias = make([]int, len(p))
	this.probability = make([]float64, len(p))

	var average = kProbability / float64(len(p))

	var small = list.New()
	var large = list.New()

	for index, value := range p {
		if value >= average {
			large.PushBack(index)
		} else {
			small.PushBack(index)
		}
	}

	for {
		var smallElement = small.Back()
		var largeElement = large.Back()

		if smallElement == nil || largeElement == nil {
			break
		}

		var less = 0
		var more = 0

		if v, ok := smallElement.Value.(int); ok {
			less = v
		}
		if v, ok := largeElement.Value.(int); ok {
			more = v
		}

		this.probability[less] = p[less] * float64(len(p))
		this.alias[less] = more

		p[more] = p[more] + p[less] - average

		if p[more] >= kProbability/float64(len(p)) {
			large.PushBack(more)
		} else {
			small.PushBack(more)
		}

		large.Remove(largeElement)
		small.Remove(smallElement)
	}

	for {
		var smallElement = small.Back()
		if smallElement == nil {
			break
		}
		if v, ok := smallElement.Value.(int); ok {
			this.probability[v] = kProbability
		}
		small.Remove(smallElement)
	}

	for {
		var largeElement = large.Back()
		if largeElement == nil {
			break
		}
		if v, ok := largeElement.Value.(int); ok {
			this.probability[v] = kProbability
		}
		large.Remove(largeElement)
	}

	return nil
}

func (this *AliasMethod[T]) Next() int {
	var pLen = len(this.probability)
	if pLen == 0 {
		return -1
	}

	var index = this.r.Intn(pLen)
	var value = this.r.Float64()

	var coinToss = value < this.probability[index]

	if coinToss {
		return index
	}
	return this.alias[index]
}

func (this *AliasMethod[T]) NextItem() T {
	var index = this.Next()
	var item T
	if index >= 0 {
		item = this.items[index]
	}
	return item
}
