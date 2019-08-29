package fate

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-stack/stack"
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

var (
	// ErrTempt is the error returned when tempting fate and losing.
	ErrTempt = errors.New("tempt fate error", j.C("ERR_9f3adc780288ce11"))
)

// Fate is an interface that wraps the tempt method.
//
// The main use case of the interface is to support deterministic behaviour
// for tests. See mockfate package.
type Fate interface {
	Tempt() error
}

// Tempt sometimes returns ErrTempt.
// Note this uses the global config, configuring and using Fate instances is advised.
func Tempt() error {
	return New().Tempt()
}

// New returns a new fate instance configured with the options provided.
func New(opts ...option) Fate {
	c := cloneConfig(globalConf)
	for _, opt := range opts {
		opt(c)
	}
	return &fate{conf: c}
}

type fate struct {
	conf *config
}

// Tempt sometimes returns ErrTempt.
func (f *fate) Tempt() error {
	if maybeInOfficeHours(f.conf, time.Now()) {
		return nil
	}

	temptCount.Inc()

	p := f.conf.DefaultP
	if f, ok := getPackageP(f.conf); ok {
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
func maybeInOfficeHours(conf *config, t time.Time) bool {
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
func getPackageP(conf *config) (float64, bool) {
	if len(conf.PackageP) == 0 {
		return 0, false
	}

	c := stack.Caller(2)
	pkg := fmt.Sprintf("%+k", c)

	prob, ok := conf.PackageP[pkg]
	return prob, ok
}
