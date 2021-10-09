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

type AliasMethod struct {
	alias       []int
	probability []float64
	items       []Item
	r           *rand.Rand
}

func New() *AliasMethod {
	var m = &AliasMethod{}
	m.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return m
}

func (this *AliasMethod) Add(item Item) {
	if item == nil {
		return
	}
	this.items = append(this.items, item)
}

func (this *AliasMethod) Prepare() error {
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

func (this *AliasMethod) process(prob []float64) error {
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

func (this *AliasMethod) Next() int {
	var proLen = len(this.probability)
	if proLen == 0 {
		return -1
	}

	var c = this.r.Intn(proLen)
	var f = this.r.Float64()

	var coinToss = f < this.probability[c]

	if coinToss {
		return c
	}
	return this.alias[c]
}

func (this *AliasMethod) NextItem() interface{} {
	var index = this.Next()
	if index < 0 {
		return nil
	}
	return this.items[index]
}
