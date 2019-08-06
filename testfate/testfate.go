package testfate

import (
	"sync"
	"testing"

	"github.com/luno/fate"
	"github.com/stretchr/testify/assert"
)

// Compile-time check that TestFate implements fate.Fate
var _ fate.Fate = (*TestFate)(nil)

type option func(*TestFate)

type TestFate struct {
	n          int
	explicit   []error
	mu         sync.Mutex
	defaultRes error
}

func New(_ testing.TB, ol ...option) *TestFate {
	f := &TestFate{}

	for _, o := range ol {
		o(f)
	}

	return f
}

func WithExplicit(errors ...error) option {
	return func(f *TestFate) {
		f.explicit = errors
	}
}

func WithDefault(error error) option {
	return func(f *TestFate) {
		f.defaultRes = error
	}
}

func (f *TestFate) Tempt() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.n++
	if len(f.explicit) > 0 {
		err := f.explicit[0]
		f.explicit = f.explicit[1:]
		return err
	}
	return f.defaultRes
}

func AssertCount(t *testing.T, fate *TestFate, n int) bool {
	return assert.Equal(t, n, fate.n)
}
