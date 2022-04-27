package datasources

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadsbProtobuf(t *testing.T) {

	var wg sync.WaitGroup
	BuildReadsbAircraftsJSON(&wg)
	wg.Wait()

	assert.IsType(t, int(1), GetReadsbDBVersion())

	adb := NewAircraftDB(2)
	assert.IsType(t, &AircraftDB{}, adb)

	// TODO: serve test data & test ReadsbProtobuf
	// ReadsbProtobuf("file://./readsb_protobuf/testdata/aircraft.pb", adb)

}
