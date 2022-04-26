package datasources

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	FORGET_AIRCRAFT_AFTER_SECONDS = 60
)

//go:embed readsb_json/aircrafts.json
var readsbAircraftsJSONBlob []byte // format: {"icao":["registration","type","flags"],...}

var readsbAircraftsJSON map[int]readsbAircraft

type readsbAircraft struct {
	registration string
	aircraftType string
	// flags        int
}

//go:embed readsb_json/dbversion.json
var readsbDBVersionJSONBlob []byte

//go:embed readsb_json/operators.json
var readsbDBOperatorsJSONBlob []byte // format: {"id":["name","country","radio"],...}

//go:embed readsb_json/types.json
var readsbDBTypesJSONBlob []byte // format: {"type":["model","species","wtc"],

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

func (adb *AircraftDB) GetAircraft() map[int]Aircraft {
	output := make(map[int]Aircraft)
	adb.Mutex.Lock()
	defer adb.Mutex.Unlock()
	for k, v := range adb.Aircraft {
		output[k] = Aircraft{
			Callsign: v.Callsign,
			Lat:      v.Lat,
			Long:     v.Long,
			Track:    v.Track,
		}
	}
	return output
}

func BuildReadsbAircraftsJSON(wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Processing aircraft.json")

	readsbAircraftsJSON = make(map[int]readsbAircraft)

	data := make(map[string]interface{})
	err := json.Unmarshal(readsbAircraftsJSONBlob, &data)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range data {

		// get icao
		icao64, err := strconv.ParseInt(k, 16, 64)
		if err != nil {
			log.Fatal(err)
		}
		icao := int(icao64)

		// get data from interface{}

		// sanity checks
		if reflect.TypeOf(v).Kind() != reflect.Slice {
			log.Fatal("aircraft.json: JSON data not type of slice")
		}
		vr := reflect.ValueOf(v)

		readsbAircraftsJSON[icao] = readsbAircraft{
			registration: vr.Index(0).Elem().String(),
			aircraftType: vr.Index(1).Elem().String(),
		}
	}
	log.Println("Finished processing aircraft.json")
	wg.Done()
}

func (adb *AircraftDB) newAircraft(icao int) {
	_, icaoInDB := adb.Aircraft[icao]
	if !icaoInDB {

		aircraftType := readsbAircraftsJSON[icao].aircraftType

		logmsg := fmt.Sprintf("AircraftDB[%6X]: Now recieving", icao)
		if aircraftType != "" {
			logmsg = fmt.Sprintf("%s, type: %s", logmsg, aircraftType)
		} else {
			logmsg = fmt.Sprintf("%s, type unknown", logmsg)
		}

		log.Println(logmsg)

		adb.Mutex.Lock()
		defer adb.Mutex.Unlock()
		adb.Aircraft[icao] = &Aircraft{}
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
