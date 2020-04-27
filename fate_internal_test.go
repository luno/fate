package fate

import (
	"testing"
)

func TestFate(t *testing.T) {
	cases := []struct {
		name string
		fate Fate
		err  error
	}{
		{
			name: "package probability 1",
			fate: New(WithPackageP(map[string]float64{"github.com/luno/fate": 1}), WithoutOfficeHours()),
			err:  ErrTempt,
		}, {
			name: "package probability 0",
			fate: New(WithPackageP(map[string]float64{"github.com/luno/fate": 0}), WithoutOfficeHours()),
			err:  nil,
		}, {
			name: "default probability 1",
			fate: New(WithDefaultP(1), WithoutOfficeHours()),
			err:  ErrTempt,
		}, {
			name: "default probability 0",
			fate: New(WithDefaultP(0), WithoutOfficeHours()),
			err:  nil,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				err := test.fate.Tempt()
				if err != test.err {
					t.Fatalf("expected: %v\ngot: %v", test.err, err)
				}
			}
		})
	}
}
