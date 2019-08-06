package fate_test

import (
	"testing"

	"github.com/luno/fate"
	"github.com/stretchr/testify/require"
)

func TestPackageProbability1(t *testing.T) {
	defer fate.DisableOfficeHoursForTesting(t)()
	defer fate.SetPackagePForTesting(t, "github.com/luno/fate_test:1")()

	for i := 0; i < 1000; i++ {
		err := fate.New().Tempt()
		require.Error(t, err)
		require.Equal(t, err, fate.ErrTempt)
	}
}

func TestPackageProbability0(t *testing.T) {
	defer fate.DisableOfficeHoursForTesting(t)()
	defer fate.SetPackagePForTesting(t, "github.com/luno/fate_test:0")()

	for i := 0; i < 1000; i++ {
		err := fate.New().Tempt()
		require.NoError(t, err)
	}
}
