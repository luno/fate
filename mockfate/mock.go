// Package mockfate provides a mock implementation for the fate.Fate interface
// useful for testing.
package mockfate

import (
	"sync"
	"testing"

	"github.com/luno/fate"
)

// Compile-time check that MockFate implements fate.Fate
var _ fate.Fate = (*MockFate)(nil)

type option func(*MockFate)

type MockFate struct {
	n          int
	explicit   []error
	mu         sync.Mutex
	defaultRes error
}

func New(_ testing.TB, ol ...option) *MockFate {
	f := &MockFate{}

	for _, o := range ol {
		o(f)
	}

	return f
}

// WithExplicit allows configuring mockfate to return these errors
// sequentially when tempted.
func WithExplicit(errors ...error) option {
	return func(f *MockFate) {
		f.explicit = errors
	}
}

// WithDefault allows configuring mockfate to always return the
// error when tempted.
func WithDefault(error error) option {
	return func(f *MockFate) {
		f.defaultRes = error
	}
}

func (f *MockFate) Tempt() error {
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
