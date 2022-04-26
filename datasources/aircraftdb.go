package datasources

import (
	_ "embed"
	"encoding/json"

	"log"
	"sync"
	"time"
)

const (
	FORGET_AIRCRAFT_AFTER_SECONDS = 60
)

//go:embed readsb_json/aircrafts.json
var readsbAircraftsJSONBlob []byte

//go:embed readsb_json/dbversion.json
var readsbDBVersionJSONBlob []byte

//go:embed readsb_json/operators.json
var readsbDBOperatorsJSONBlob []byte

//go:embed readsb_json/types.json
var readsbDBTypesJSONBlob []byte

func init() {

	// ensure the embedded JSON is valid
	if !json.Valid(readsbAircraftsJSONBlob) {
		log.Fatal("Embedded aircrafts.json is invalid")
	}
	if !json.Valid(readsbDBVersionJSONBlob) {
		log.Fatal("Embedded dbversion.json is invalid")
	}
	if !json.Valid(readsbDBOperatorsJSONBlob) {
		log.Fatal("Embedded operators.json is invalid")
	}
	if !json.Valid(readsbDBTypesJSONBlob) {
		log.Fatal("Embedded types.json is invalid")
	}
}

type Aircraft struct {
	Callsign    string
	Lat         float64
	Long        float64
	Track       int
	LastUpdated int64
}

type AircraftDB struct {
	Aircraft map[int]*Aircraft
	Mutex    sync.Mutex
}

func (adb *AircraftDB) newAircraft(icao int) {
	_, icaoInDB := adb.Aircraft[icao]
	if !icaoInDB {
		log.Printf("AircraftDB[%6X]: Added to DB", icao)
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		adb.Aircraft[icao] = &Aircraft{}
	}
}

func (adb *AircraftDB) UpdateCallsign(icao int, callsign string) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Callsign != callsign {
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated callsign to: %s", icao, callsign)
		adb.Aircraft[icao].Callsign = callsign
	}
}

func (adb *AircraftDB) UpdateLat(icao int, lat float64) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Lat != lat {
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated lat to: %f", icao, lat)
		adb.Aircraft[icao].Lat = lat
	}
}

func (adb *AircraftDB) UpdateLong(icao int, long float64) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Long != long {
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated long to: %f", icao, long)
		adb.Aircraft[icao].Long = long
	}
}

func (adb *AircraftDB) UpdateTrack(icao int, track int) {
	adb.newAircraft(icao)
	if adb.Aircraft[icao].Track != track {
		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		// log.Printf("AircraftDB[%X]: Updated track to: %d", icao, track)
		adb.Aircraft[icao].Track = track
	}
}

func (adb *AircraftDB) UpdateLastSeen(icao int) {
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

		// forget entries with LastUpdated older than FORGET_AIRCRAFT_AFTER_SECONDS
		if time.Now().Unix() > adb.Aircraft[k].LastUpdated+FORGET_AIRCRAFT_AFTER_SECONDS {
			log.Printf("AircraftDB[%X]: Forgetting inactive aircraft", k)
			defer delete(adb.Aircraft, k)
		}
	}

	// run again
	go adb.forgetter()
}

func NewAircraftDB() *AircraftDB {
	// Initialises and returns a pointer to an aircraft db
	adb := AircraftDB{}
	adb.Aircraft = make(map[int]*Aircraft)
	go adb.forgetter()
	return &adb
}

func GetReadsbDBVersion() int {
	data := make(map[string]interface{})
	err := json.Unmarshal(readsbDBVersionJSONBlob, &data)
	if err != nil {
		log.Fatal(err)
	}
	return int(data["version"].(float64))
}
