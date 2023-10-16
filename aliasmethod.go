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
	r           *rand.Rand
	alias       []int
	probability []float64
	items       []T
}

func New[T Item]() *AliasMethod[T] {
	var m = &AliasMethod[T]{}
	m.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return m
}

func (m *AliasMethod[T]) Add(item T) {
	//if item == nil {
	//	return
	//}
	m.items = append(m.items, item)
}

func (m *AliasMethod[T]) Prepare() error {
	if len(m.items) == 0 {
		return errors.New("概率不能为空")
	}

	var total = int32(0)
	for _, item := range m.items {
		total += item.GetWeight()
	}

	var scale = float64(total) / kProbability

	var values = make([]float64, 0, len(m.items))
	for _, item := range m.items {
		values = append(values, float64(item.GetWeight())/scale)
	}

	return m.process(values)
}

func (m *AliasMethod[T]) process(prob []float64) error {
	var p = make([]float64, len(prob))
	copy(p, prob)

	m.alias = make([]int, len(p))
	m.probability = make([]float64, len(p))

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

		m.probability[less] = p[less] * float64(len(p))
		m.alias[less] = more

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
			m.probability[v] = kProbability
		}
		small.Remove(smallElement)
	}

	for {
		var largeElement = large.Back()
		if largeElement == nil {
			break
		}
		if v, ok := largeElement.Value.(int); ok {
			m.probability[v] = kProbability
		}
		large.Remove(largeElement)
	}

	return nil
}

func (m *AliasMethod[T]) Next() int {
	var pLen = len(m.probability)
	if pLen == 0 {
		return -1
	}

	var index = m.r.Intn(pLen)
	var value = m.r.Float64()

	var coinToss = value < m.probability[index]

	if coinToss {
		return index
	}
	return m.alias[index]
}

func (m *AliasMethod[T]) NextItem() T {
	var index = m.Next()
	var item T
	if index >= 0 {
		item = m.items[index]
	}
	return item
}
