package aliasmethod_test

import (
	"github.com/smartwalle/aliasmethod"
	"testing"
)

type Christmas struct {
	name   string
	weight int32
}

func (this *Christmas) GetWeight() int32 {
	return this.weight
}

func Test_AliasMethod(t *testing.T) {
	var results = make(map[string]int)

	var m = aliasmethod.New[*Christmas]()

	m.Add(&Christmas{name: "1", weight: 10})
	m.Add(&Christmas{name: "2", weight: 10})
	m.Add(&Christmas{name: "3", weight: 10})
	m.Add(&Christmas{name: "4", weight: 10})
	m.Add(&Christmas{name: "5", weight: 10})
	m.Add(&Christmas{name: "6", weight: 10})
	m.Add(&Christmas{name: "7", weight: 10})
	m.Add(&Christmas{name: "8", weight: 10})
	m.Add(&Christmas{name: "9", weight: 10})
	m.Add(&Christmas{name: "10", weight: 110})

	if err := m.Prepare(); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		var p = m.NextItem()
		results[p.name] = results[p.name] + 1
	}

	t.Log(results)
}

func BenchmarkAliasMethod_NextItem(b *testing.B) {
	var results = make(map[string]int)

	var m = aliasmethod.New[*Christmas]()

	m.Add(&Christmas{name: "1", weight: 10})
	m.Add(&Christmas{name: "2", weight: 10})
	m.Add(&Christmas{name: "3", weight: 10})
	m.Add(&Christmas{name: "4", weight: 10})
	m.Add(&Christmas{name: "5", weight: 10})
	m.Add(&Christmas{name: "6", weight: 10})
	m.Add(&Christmas{name: "7", weight: 10})
	m.Add(&Christmas{name: "8", weight: 10})
	m.Add(&Christmas{name: "9", weight: 10})
	m.Add(&Christmas{name: "10", weight: 110})

	if err := m.Prepare(); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		var p = m.NextItem()
		results[p.name] = results[p.name] + 1
	}

	b.Log(results)
}
