package fate_test

import (
	"testing"

	"github.com/luno/fate"
)

func TestPackageProbability1(t *testing.T) {
	defer fate.SetConfig(
		fate.WithoutOfficeHours(),
		fate.WithPackageP(map[string]float64{"github.com/luno/fate_test": 1}),
	)()

	for i := 0; i < 100; i++ {
		err := fate.New().Tempt()
		if err != fate.ErrTempt {
			t.Fatalf("expected fate.ErrTempt: %v", err)
		}
	}
}

func TestPackageProbability0(t *testing.T) {
	defer fate.SetConfig(
		fate.WithoutOfficeHours(),
		fate.WithPackageP(map[string]float64{"github.com/luno/fate_test": 0}),
	)()

	for i := 0; i < 100; i++ {
		err := fate.New().Tempt()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestDefaultProbability1(t *testing.T) {
	defer fate.SetConfig(
		fate.WithoutOfficeHours(),
		fate.WithDefaultP(1),
	)()

	for i := 0; i < 100; i++ {
		err := fate.New().Tempt()
		if err != fate.ErrTempt {
			t.Fatalf("expected fate.ErrTempt: %v", err)
		}
	}
}

func TestDefaultProbability0(t *testing.T) {
	defer fate.SetConfig(
		fate.WithoutOfficeHours(),
		fate.WithDefaultP(0),
	)()

	for i := 0; i < 100; i++ {
		err := fate.New().Tempt()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}
