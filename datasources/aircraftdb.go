package datasources

import (
	_ "embed"

	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FORGET_AIRCRAFT_AFTER_SECONDS = 60
)

type Aircraft struct {
	Callsign     string
	Lat          float64
	Long         float64
	Track        int
	LastUpdated  int64
	AircraftType string
}

type AircraftDB struct {
	Aircraft    map[int]*Aircraft
	Mutex       sync.Mutex
	idleTimeout int64 // seconds

}

func (adb *AircraftDB) GetAircraft() map[int]Aircraft {
	output := make(map[int]Aircraft)
	adb.Mutex.Lock()
	defer adb.Mutex.Unlock()
	for k, v := range adb.Aircraft {
		output[k] = Aircraft{
			Callsign:     v.Callsign,
			Lat:          v.Lat,
			Long:         v.Long,
			Track:        v.Track,
			AircraftType: v.AircraftType,
		}
	}
	return output
}

func (adb *AircraftDB) newAircraft(icao int) {
	adb.Mutex.Lock()
	_, icaoInDB := adb.Aircraft[icao]
	adb.Mutex.Unlock()
	if !icaoInDB {

		aircraftType := readsbAircraft[icao].aircraftType

		// logmsg := fmt.Sprintf("AircraftDB[%6X]: Now recieving", icao)
		// if aircraftType != "" {
		// 	logmsg = fmt.Sprintf("%s, type: %s", logmsg, aircraftType)
		// } else {
		// 	logmsg = fmt.Sprintf("%s, type unknown", logmsg)
		// }
		// log.Println(logmsg)

		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		adb.Aircraft[icao] = &Aircraft{
			AircraftType: aircraftType,
		}
	}
}

func (adb *AircraftDB) SetCallsign(icao int, callsign string) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Callsign != callsign {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated callsign to: %s", icao, callsign)
		adb.Aircraft[icao].Callsign = callsign
	}
}

func (adb *AircraftDB) SetLat(icao int, lat float64) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Lat != lat {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated lat to: %f", icao, lat)
		adb.Aircraft[icao].Lat = lat
	}
}

func (adb *AircraftDB) SetLong(icao int, long float64) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Long != long {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated long to: %f", icao, long)
		adb.Aircraft[icao].Long = long
	}
}

func (adb *AircraftDB) SetTrack(icao int, track int) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Track != track {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated track to: %d", icao, track)
		adb.Aircraft[icao].Track = track
	}
}

func (adb *AircraftDB) SetLastSeen(icao int) {
	adb.newAircraft(icao)
	adb.Mutex.Lock()
	defer adb.Mutex.Unlock()
	adb.Aircraft[icao].LastUpdated = time.Now().Unix()
}

func (adb *AircraftDB) forgetter() {

	// sleep for 1 second
	time.Sleep(time.Second)

	// "lock" the database
	adb.Mutex.Lock()
	defer adb.Mutex.Unlock()

	// for each aircraft entry
	for k, _ := range adb.Aircraft {

		// forget entries with LastUpdated older than adb.idleTimeout
		if time.Now().Unix() > adb.Aircraft[k].LastUpdated+adb.idleTimeout {
			// log.Printf("AircraftDB[%X]: Forgetting inactive aircraft", k)
			defer delete(adb.Aircraft, k)
		}
	}

	// run again
	go adb.forgetter()
}

func NewAircraftDB(idleTimeout int64) *AircraftDB {
	// Initialises and returns a pointer to an aircraft db
	adb := AircraftDB{
		idleTimeout: idleTimeout,
	}
	adb.Aircraft = make(map[int]*Aircraft)
	go adb.forgetter()
	return &adb
}

func GetReadsbDBVersion() int {
	return readsbDBVersion
}
