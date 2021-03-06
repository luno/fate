package fate

import (
	"time"
)

var (
	globalConf = defaultConfig()
)

type config struct {
	// DefaultP is the default error probability when tempting fate.
	DefaultP float64

	// PackageP is the per package error probability when tempting fate"
	PackageP map[string]float64

	// OfficeHours allows limiting errors to office hours.
	OfficeHours officeHours
}

type officeHours struct {
	Enabled   bool
	Location  *time.Location
	HourStart int
	HourEnd   int
}

// SetConfig sets the global config defined by the functional options.
// It returns an undo function mostly used for testing.
// Note configuring and using Fate instances is advised over global config.
func SetConfig(opts ...option) func() {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}

	undo := globalConf
	globalConf = conf
	return func() {
		globalConf = undo
	}
}

func cloneConfig(conf *config) *config {
	c := *conf
	c.PackageP = make(map[string]float64)
	for key, value := range conf.PackageP {
		c.PackageP[key] = value
	}
	return &c
}

// defaultConfig defines the default config.
func defaultConfig() *config {
	return &config{
		DefaultP: 1.0 / 1e6,
		PackageP: make(map[string]float64),
		OfficeHours: officeHours{
			Enabled:   true,
			Location:  time.UTC,
			HourStart: 9,
			HourEnd:   17,
		},
	}
}

type option func(*config)

// WithDefaultP allows defining the default error probability when tempting fate.
func WithDefaultP(p float64) option {
	return func(c *config) {
		c.DefaultP = p
	}
}

// WithPackageP allows defining the per package error probability when tempting fate.
func WithPackageP(p map[string]float64) option {
	return func(c *config) {
		c.PackageP = p
	}
}

// WithOfficeHours allows limiting errors to office hours.
func WithOfficeHours(loc *time.Location, hourStart, hourEnd int) option {
	return func(c *config) {
		c.OfficeHours = officeHours{
			Enabled:   true,
			Location:  loc,
			HourStart: hourStart,
			HourEnd:   hourEnd,
		}
	}
}

// WithoutOfficeHours allows disabling limiting of errors to office hours.
func WithoutOfficeHours() option {
	return func(c *config) {
		c.OfficeHours = officeHours{
			Enabled: false,
		}
	}
}
