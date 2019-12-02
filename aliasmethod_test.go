package aliasmethod

import (
	"fmt"
	"testing"
)

type Christmas struct {
	name string
	p    float64
	c    int
}

func (this *Christmas) Probability() float64 {
	var v = float64(this.p)
	return v
}

func Test_AliasMethod(t *testing.T) {
	var results = make(map[string]int)

	for i := 0; i < 1000; i++ {
		var am = NewAliasMethod()

		am.AddProbability(&Christmas{name: "1", p: 10})
		am.AddProbability(&Christmas{name: "2", p: 10})
		am.AddProbability(&Christmas{name: "3", p: 10})
		am.AddProbability(&Christmas{name: "4", p: 10})
		am.AddProbability(&Christmas{name: "5", p: 10})
		am.AddProbability(&Christmas{name: "6", p: 10})
		am.AddProbability(&Christmas{name: "7", p: 10})
		am.AddProbability(&Christmas{name: "8", p: 10})
		am.AddProbability(&Christmas{name: "9", p: 10})
		am.AddProbability(&Christmas{name: "10", p: 110})

		if err := am.Prepare(); err != nil {
			t.Fatal(err)
		}
		var p = am.NextValue()

		var c = p.(*Christmas)
		results[c.name] = results[c.name] + 1
	}

	fmt.Println(results)
}
