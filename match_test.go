package match_test

import (
	"testing"

	"github.com/Instantan/match"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	_, err := match.Compile("namespace. real | virtual ].[ root* ].value")
	assert.Error(t, err)

	m, err := match.Compile("namespace.[ real | virtual ].[ root* ].value")
	assert.NoError(t, err)

	assert.Equal(t, true, m.Matches("namespace.real.root.path.value"))
	assert.Equal(t, true, m.Matches("namespace.real.root.bla.value"))
	assert.Equal(t, true, m.Matches("namespace.virtual.root.bla.value"))
	assert.Equal(t, true, m.Matches("namespace.virtual.root.dwakjdnajkwd.value"))
	assert.Equal(t, true, m.Matches("namespace.virtual.root.value"))

	assert.Equal(t, false, m.Matches("nadmespace.rea.root.path.value"))
	assert.Equal(t, false, m.Matches("namespace.re.root.bla.value"))
	assert.Equal(t, false, m.Matches("namespace.virtual.root.lue"))
	assert.Equal(t, false, m.Matches("namespace.virtual.oot.dwakjdnajkwd.value"))
	assert.Equal(t, false, m.Matches("namespace.virtual.roo.value"))
}

func BenchmarkMatch(b *testing.B) {
	m, _ := match.Compile("namespace.[ real | virtual ].[ root* ].value")

	items := []string{
		"namespace.real.root.path.value",
		"namespace.real.root.bla.value",
		"namespace.virtual.root.bla.value",
		"namespace.virtual.root.dwakjdnajkwd.value",
		"namespace.virtual.root.value",
		"nadmespace.rea.root.path.value",
		"namespace.virtual.root.lue",
		"namespace.virtual.oot.dwakjdnajkwd.value",
		"namespace.virtual.roo.value",
	}

	for n := 0; n < b.N; n++ {
		for i := range items {
			m.Matches(items[i])
		}
	}
}
