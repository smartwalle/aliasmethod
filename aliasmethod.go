package aliasmethod

import (
	"container/list"
	"math/rand"
	"time"
	"errors"
)

type AliasMethod struct {
	alias 		[]int
	probability []float64
}

func NewAliasMethod(p []float64) (alias *AliasMethod, err error) {
	if p == nil {
		return nil, errors.New("概率不能为空")
	}

	if len(p) == 0 {
		return nil, errors.New("概率不能为空")
	}

	alias = &AliasMethod{}
	alias.preprocess(p)
	return alias, nil
}

func (this *AliasMethod) preprocess(prob []float64) (err error) {
	var p = make([]float64, len(prob))
	copy(p, prob)

	this.alias       = make([]int, len(p))
	this.probability = make([]float64, len(p))

	var average float64 = 1.0 / float64(len(p))

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
		var smallElement *list.Element = small.Back()
		var largeElement *list.Element = large.Back()

		if smallElement == nil || largeElement == nil {
			break
		}

		var less int = 0;
		var more int = 0;

		if v, ok := smallElement.Value.(int); ok {
			less = v
		}
		if v, ok := largeElement.Value.(int); ok {
			more = v
		}

		this.probability[less] = p[less] * float64(len(p))
		this.alias[less] = more

		p[more] = p[more] + p[less] - average

		if (p[more] >= 1.0 / float64(len(p))) {
			large.PushBack(more)
		} else {
			small.PushBack(more)
		}

		large.Remove(largeElement)
		small.Remove(smallElement)
	}

	for {
		var smallElement *list.Element = small.Back()
		if smallElement == nil {
			break
		}
		if v, ok := smallElement.Value.(int); ok {
			this.probability[v] = 1.0
		}
		small.Remove(smallElement)
	}

	for {
		var largeElement *list.Element = large.Back()
		if largeElement == nil {
			break
		}
		if v, ok := largeElement.Value.(int); ok {
			this.probability[v] = 1.0
		}
		large.Remove(largeElement)
	}

	return err
}

func (this *AliasMethod) Next() int {
	rand.Seed(time.Now().UnixNano())

	var column = rand.Intn(len(this.probability))
	var f = rand.Float64()

	var coinToss = f < this.probability[column]

	if coinToss {
		return column
	}
	return this.alias[column]
}
