package fate

import (
	"testing"
)

func TestFate(t *testing.T) {
	cases := []struct {
		name string
		opt  option
		err  bool
	}{
		{
			name: "package probability 1",
			opt:  WithPackageP(map[string]float64{"github.com/luno/fate": 1}),
			err:  true,
		}, {
			name: "package probability 0",
			opt:  WithPackageP(map[string]float64{"github.com/luno/fate": 0}),
			err:  false,
		}, {
			name: "default probability 1",
			opt:  WithDefaultP(1),
			err:  true,
		}, {
			name: "default probability 0",
			opt:  WithDefaultP(0),
			err:  false,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fate := New(test.opt, WithoutOfficeHours())
			for i := 0; i < 100; i++ {
				err := fate.Tempt()
				if test.err {
					if err != ErrTempt {
						t.Fatalf("expected fate.ErrTempt: %v", err)
					}
				} else {
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
				}
			}
		})
	}
}
