package datasources

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAircraftDB(t *testing.T) {

	var wg sync.WaitGroup
	BuildReadsbAircraftsJSON(&wg)
	wg.Wait()

	assert.IsType(t, int(1), GetReadsbDBVersion())

	adb := NewAircraftDB(2)
	assert.IsType(t, &AircraftDB{}, adb)

	adb.SetCallsign(0xAAAAAA, "TEST1")
	adb.SetLat(0xAAAAAA, -31.9523)
	adb.SetLong(0xAAAAAA, 115.8613)
	adb.SetTrack(0xAAAAAA, 123)
	adb.SetLastSeen(0xAAAAAA)

	adb.SetCallsign(0x7C1465, "TEST2")
	adb.SetLat(0x7C1465, -31.9523)
	adb.SetLong(0x7C1465, 115.8613)
	adb.SetTrack(0x7C1465, 123)
	adb.SetLastSeen(0x7C1465)

	output := adb.GetAircraft()

	assert.Equal(t, "TEST1", output[0xAAAAAA].Callsign)
	assert.Equal(t, -31.9523, output[0xAAAAAA].Lat)
	assert.Equal(t, 115.8613, output[0xAAAAAA].Long)
	assert.Equal(t, 123, output[0xAAAAAA].Track)
	assert.Equal(t, "", output[0xAAAAAA].AircraftType)

	log.Println("Waiting for forgetter...")
	time.Sleep(time.Second * 5)

	output = adb.GetAircraft()
	_, ok := output[0xAAAAAA]
	assert.Equal(t, false, ok)

}
