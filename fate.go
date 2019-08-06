package fate

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-stack/stack"
)

var (
	//TODO(corver): Remove after this date
	enableDate = time.Date(2018, 12, 18, 0, 0, 0, 0, time.UTC)
	enable     = flag.Bool("fate_enable", true, "enable errors when tempting fate")

	officeHoursEnable = flag.Bool("fate_office_hours", true, "limit errors to office hours")
	officeHoursStart  = flag.Int("fate_office_hours_start", 8, "office hours start hour (0-23)")
	officeHoursEnd    = flag.Int("fate_office_hours_end", 17, "office hours end hour (0-23)")
	officeHoursZone   = flag.String("fate_office_hours_zone", "UTC", "office hours timezone (see time.LoadLocation)")

	officeHoursLocation *time.Location
	officeHoursOnce     sync.Once

	defaultP = flag.Float64("fate_default_p", 1.0/1e6, "default error probability when tempting fate")
	packageP = flag.String("fate_package_p", "", "per package error probability when tempting fate;"+
		"fully/qualified/package:0.1,other/package:0")

	// ErrTempt is the error returned when tempting fate and losing.
	ErrTempt = errors.New("tempt fate error")

	packagePMap        map[string]float64
	loadPackageOncePer sync.Once
)

type Fate interface {
	Tempt() error
}

func New() Fate {
	return &fate{}
}

type fate struct{}

func (f *fate) Tempt() error {
	if !isEnabled() {
		return nil
	}

	temptCount.Inc()

	if ok, err := maybeTemptPackage(); ok {
		return err
	}

	return temptDefault()
}

func isEnabled() bool {
	if !*enable {
		return false
	}

	if time.Now().Before(enableDate) {
		return false
	}

	if *officeHoursEnable {
		return isOfficeHours()
	}

	return true
}

func getLocation() *time.Location {
	officeHoursOnce.Do(func() {
		var err error
		officeHoursLocation, err = time.LoadLocation(*officeHoursZone)
		if err != nil {
			log.Printf("fate: Error parsing officeHoursZone flag: %s, %v", *officeHoursZone, err)
			*officeHoursEnable = false
		}
	})
	return officeHoursLocation
}

// isOfficeHours returns true if now is in office hours.
func isOfficeHours() bool {
	now := time.Now().In(getLocation())
	weekday := now.Weekday()
	if weekday == 6 || weekday == 7 {
		return false
	}
	hour := now.Hour()
	return *officeHoursStart <= hour && hour <= *officeHoursEnd
}

func maybeTemptPackage() (bool, error) {
	loadPackageOncePer.Do(func() {
		packagePMap = make(map[string]float64)
		if *packageP == "" {
			return
		}

		var err error
		packagePMap, err = parsePackageP(*packageP)
		if err != nil {
			log.Printf("fate: Error parsing packagep flag: %s, %v", *packageP, err)
		}
	})

	c := stack.Caller(2)
	pkg := fmt.Sprintf("%+k", c)
	prob, ok := packagePMap[pkg]
	if !ok {
		return false, nil
	}

	if rand.Float64() < prob {
		temptErrors.Inc()
		return true, ErrTempt
	}

	return false, nil
}

func temptDefault() error {
	if rand.Float64() < *defaultP {
		temptErrors.Inc()
		return ErrTempt
	}

	return nil
}

func parsePackageP(input string) (map[string]float64, error) {
	result := make(map[string]float64)
	split := strings.Split(input, ",")
	for _, s := range split {
		ss := strings.Split(s, ":")
		if len(ss) != 2 {
			return nil, errors.New("invalid format")
		}
		f, err := strconv.ParseFloat(ss[1], 64)
		if err != nil {
			return nil, errors.New("invalid float probability")
		}
		result[ss[0]] = f
	}
	return result, nil
}

func DisableOfficeHoursForTesting(t *testing.T) func() {
	cache := *officeHoursEnable
	*officeHoursEnable = false
	return func() {
		*officeHoursEnable = cache
	}
}

func SetPackagePForTesting(t *testing.T, s string) func() {
	cache := *packageP
	*packageP = s
	loadPackageOncePer = sync.Once{}
	return func() {
		*packageP = cache
	}
}
