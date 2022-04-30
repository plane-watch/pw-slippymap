package markers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkers(t *testing.T) {
	markers, err := InitMarkers()
	require.NoError(t, err)

	for _, m := range markers {
		testName := fmt.Sprintf("PointInsideMarker: %s", m.icao)
		t.Run(testName, func(t *testing.T) {
			p := []float64{m.CentreX, m.CentreY}
			isInside := m.PointInsideMarker(p[0], p[1])
			assert.True(t, isInside)

			p = []float64{0, 0}
			isInside = m.PointInsideMarker(p[0], p[1])
			assert.False(t, isInside)

		})

	}
}
