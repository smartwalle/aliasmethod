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
	return this.p
}

func Test_AliasMethod(t *testing.T) {
	var am = NewAliasMethod()

	am.AddProbability(&Christmas{name: "圣诞老人", p: 0.05})
	am.AddProbability(&Christmas{name: "圣诞树", p: 0.15})
	am.AddProbability(&Christmas{name: "圣诞袜", p: 0.15})
	am.AddProbability(&Christmas{name: "圣诞小鹿", p: 0.15})
	am.AddProbability(&Christmas{name: "谢谢参与", p: 0.5})

	if err := am.Prepare(); err != nil {
		t.Fatal(err)
	}

	var results = make(map[string]int)

	for i := 0; i < 100; i++ {
		var p = am.NextValue()

		var c = p.(*Christmas)
		results[c.name] = results[c.name] + 1
	}

	fmt.Println(results)
}
