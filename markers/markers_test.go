package markers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitMarkers(t *testing.T) {
	_, err := InitMarkers()
	require.NoError(t, err)
}
