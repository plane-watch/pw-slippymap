package datasources

import (
	_ "embed"

	"sync"
	"time"

	"pw_slippymap/datasources/readsb_protobuf"

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
	AltBaro      int
	Category     int
	GroundSpeed  int
	AirGround    readsb_protobuf.AircraftMeta_AirGround
	History      []AircraftHistoryLocation
}

type AircraftHistoryLocation struct {
	Lat  float64
	Long float64
	Alt  int
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
			AltBaro:      v.AltBaro,
			Category:     v.Category,
			GroundSpeed:  v.GroundSpeed,
			AirGround:    v.AirGround,
		}
	}
	return output
}

func (adb *AircraftDB) newAircraft(icao int) {
	adb.Mutex.Lock()
	_, icaoInDB := adb.Aircraft[icao]
	adb.Mutex.Unlock()
	if !icaoInDB {

		// lookup aircraft type
		aircraftType := readsbAircraft[icao].aircraftType

		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()

		// create Aircraft object
		adb.Aircraft[icao] = &Aircraft{
			AircraftType: aircraftType,
		}

		// init History slice
		adb.Aircraft[icao].History = make([]AircraftHistoryLocation, 0)
	}
}

func (adb *AircraftDB) AddHistory(icao int, lat, long float64, alt int) {
	adb.newAircraft(icao)
	adb.Mutex.Lock()
	defer adb.Mutex.Unlock()
	ahl := AircraftHistoryLocation{
		Lat:  lat,
		Long: long,
		Alt:  alt,
	}
	adb.Aircraft[icao].History = append(adb.Aircraft[icao].History, ahl)
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

func (adb *AircraftDB) SetPosition(icao int, lat, long float64, altBaro int) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Lat != lat || adb.Aircraft[icao].Long != long || adb.Aircraft[icao].AltBaro != altBaro {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		adb.Aircraft[icao].Lat = lat
		adb.Aircraft[icao].Long = long
		adb.Aircraft[icao].AltBaro = altBaro
		go adb.AddHistory(icao, lat, long, altBaro)
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

func (adb *AircraftDB) SetCategory(icao int, category int) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Category != category {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		adb.Aircraft[icao].Category = category
	}
}

func (adb *AircraftDB) SetGs(icao int, gs int) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].GroundSpeed != gs {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		adb.Aircraft[icao].GroundSpeed = gs
	}
}

func (adb *AircraftDB) SetAirGround(icao int, ag readsb_protobuf.AircraftMeta_AirGround) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].AirGround != ag {
		defer ebiten.ScheduleFrame()
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		adb.Aircraft[icao].AirGround = ag
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
