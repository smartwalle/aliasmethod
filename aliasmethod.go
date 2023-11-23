package aliasmethod

import (
	"container/list"
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
	r     *rand.Rand
	alias []int     // 別名表
	probs []float64 // 概率表
	items []T
}

func New[T Item]() *AliasMethod[T] {
	var m = &AliasMethod[T]{}
	m.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return m
}

func (m *AliasMethod[T]) Add(item T) {
	m.items = append(m.items, item)
}

func (m *AliasMethod[T]) Prepare() bool {
	if len(m.items) == 0 {
		return false
	}

	var total = int32(0)
	for _, item := range m.items {
		total += item.GetWeight()
	}

	var scale = float64(total) / kProbability

	var probs = make([]float64, 0, len(m.items))
	for _, item := range m.items {
		probs = append(probs, float64(item.GetWeight())/scale)
	}

	return m.process(probs)
}

func (m *AliasMethod[T]) process(probs []float64) bool {
	m.alias = make([]int, len(probs))
	m.probs = make([]float64, len(probs))

	var avg = kProbability / float64(len(probs))

	var smallList = list.New()
	var largeList = list.New()

	for index, value := range probs {
		if value >= avg {
			largeList.PushBack(index)
		} else {
			smallList.PushBack(index)
		}
	}

	for {
		var smallElement = smallList.Back()
		var largeElement = largeList.Back()

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

		m.probs[less] = probs[less] * float64(len(probs))
		m.alias[less] = more

		probs[more] = probs[more] + probs[less] - avg

		if probs[more] >= kProbability/float64(len(probs)) {
			largeList.PushBack(more)
		} else {
			smallList.PushBack(more)
		}

		largeList.Remove(largeElement)
		smallList.Remove(smallElement)
	}

	for {
		var smallElement = smallList.Back()
		if smallElement == nil {
			break
		}
		if v, ok := smallElement.Value.(int); ok {
			m.probs[v] = kProbability
		}
		smallList.Remove(smallElement)
	}

	for {
		var largeElement = largeList.Back()
		if largeElement == nil {
			break
		}
		if v, ok := largeElement.Value.(int); ok {
			m.probs[v] = kProbability
		}
		largeList.Remove(largeElement)
	}

	return true
}

func (m *AliasMethod[T]) Next() int {
	var pLen = len(m.probs)
	if pLen == 0 {
		return -1
	}

	var index = m.r.Intn(pLen)
	var value = m.r.Float64()

	var coin = value < m.probs[index]
	if coin {
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
