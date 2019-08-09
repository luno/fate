package fate

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-stack/stack"
)

var (
	// ErrTempt is the error returned when tempting fate and losing.
	ErrTempt = errors.New("tempt fate error")
)

// Fate is an interface that wraps the tempt method.
//
// The main use case of the interface is to support deterministic behaviour
// for tests. See mockfate package.
type Fate interface {
	Tempt() error
}

// Tempt sometimes returns ErrTempt.
func Tempt() error {
	return New().Tempt()
}

// New returns a new fate instance.
func New() Fate {
	return &fate{}
}

type fate struct{}

// Tempt sometimes returns ErrTempt.
func (f *fate) Tempt() error {
	if maybeInOfficeHours(time.Now()) {
		return nil
	}

	temptCount.Inc()

	p := conf.DefaultP
	if f, ok := getPackageP(); ok {
		p = f
	}

	if rand.Float64() < p {
		temptErrors.Inc()
		return ErrTempt
	}

	return nil
}

// maybeInOfficeHours returns true if office hours is configured and
// t falls inside it.
func maybeInOfficeHours(t time.Time) bool {
	if !conf.OfficeHours.Enabled {
		return false
	}

	now := t.In(conf.OfficeHours.Location)

	weekday := now.Weekday()
	if weekday == 6 || weekday == 7 {
		return false
	}

	hour := now.Hour()
	return conf.OfficeHours.HourStart <= hour && hour <= conf.OfficeHours.HourEnd
}

// getPackageP returns the error probability and true if it is explicitly
// configured for the calling package.
func getPackageP() (float64, bool) {
	if len(conf.PackageP) == 0 {
		return 0, false
	}

	c := stack.Caller(2)
	pkg := fmt.Sprintf("%+k", c)

	prob, ok := conf.PackageP[pkg]
	return prob, ok
}
