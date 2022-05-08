package markers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkers(t *testing.T) {
	t.Run("Aircraft", func(t *testing.T) {
		_, err := InitMarkers(Aircraft)
		require.NoError(t, err)
	})
	t.Run("GroundVehicles", func(t *testing.T) {
		_, err := InitMarkers(GroundVehicles)
		require.NoError(t, err)
	})
}
