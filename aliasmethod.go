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

type Probability interface {
	Probability() float64
}

type AliasMethod struct {
	alias          []int
	probability    []float64
	rawProbability []Probability
}

func NewAliasMethod() *AliasMethod {
	var alias = &AliasMethod{}
	return alias
}

func (this *AliasMethod) AddProbability(p Probability) {
	if p == nil {
		return
	}
	this.rawProbability = append(this.rawProbability, p)
}

func (this *AliasMethod) Prepare() error {
	if len(this.rawProbability) == 0 {
		return errors.New("概率不能为空")
	}

	var total float64 = 0
	for _, p := range this.rawProbability {
		total += p.Probability()
	}

	var scale = total / kProbability

	var values = make([]float64, 0, len(this.rawProbability))
	for _, p := range this.rawProbability {
		values = append(values, p.Probability()/scale)
	}

	this.preprocess(values)

	return nil
}

func (this *AliasMethod) preprocess(prob []float64) (err error) {
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

	return err
}

func (this *AliasMethod) Next() int {
	rand.Seed(time.Now().UnixNano())

	var proLen = len(this.probability)
	if proLen == 0 {
		return -1
	}

	var column = rand.Intn(proLen)
	var f = rand.Float64()

	var coinToss = f < this.probability[column]

	if coinToss {
		return column
	}
	return this.alias[column]
}

func (this *AliasMethod) NextValue() interface{} {
	var index = this.Next()
	if index < 0 {
		return nil
	}
	return this.rawProbability[index]
}
