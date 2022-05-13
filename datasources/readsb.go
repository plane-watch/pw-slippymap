package datasources

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"pw_slippymap/datasources/readsb_protobuf"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	ERROR_SECONDS_BACKOFF   = 5
	READSB_UPDATE_FREQUENCY = 500 // milliseconds
)

func readsbProtobufHistory(readsburl string, adb *AircraftDB) {
	// Updates the AircraftDB adb from history.pb located at readsburl/data/history.pb

	i := 0
	finished := false
	for finished == false {

		// build the URL to data/history_i.pb
		urlHistoryPb, err := url.Parse(readsburl)
		if err != nil {
			log.Fatal(err)
		}
		historyPath := fmt.Sprintf("data/history_%d.pb", i)
		urlHistoryPb.Path = path.Join(urlHistoryPb.Path, historyPath)

		// Get the data
		resp, err := http.Get(urlHistoryPb.String())

		// If we have any error downloading the history, then bail
		if err != nil {
			break
		}

		// Close body at end
		defer resp.Body.Close()

		// Check server response
		switch resp.StatusCode {
		case http.StatusOK:
			// pass
		case http.StatusNotFound:
			finished = true
		default:
			log.Fatal(fmt.Sprintf("datasources.ReadsbProtobufHistory: HTTP status was %d, expected %d", resp.StatusCode, http.StatusOK))
		}

		// bail if finished
		if finished {
			break
		}

		// Read response body
		pbData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("datasources.ReadsbProtobufHistory: Error reading HTTP response body: ", err)
			time.Sleep(time.Second * ERROR_SECONDS_BACKOFF)
			continue
		}

		// Attempt to unmarshall
		aircraftUpdate := &readsb_protobuf.AircraftsUpdate{}
		err = proto.Unmarshal(pbData, aircraftUpdate)
		if err != nil {
			log.Println("datasources.ReadsbProtobufHistory: Error unmarshalling protobuf data: ", err)
			time.Sleep(time.Second & ERROR_SECONDS_BACKOFF)
			continue
		}

		// Add history
		for _, v := range aircraftUpdate.GetHistory() {
			adb.AddHistory(int(v.Addr), v.Lat, v.Lon, int(v.AltBaro))
		}

		// Increment index
		i += 1
	}
}

func ReadsbProtobufAircraft(readsburl string, adb *AircraftDB) {
	// Updates the AircraftDB adb from aircraft.pb located at readsburl/data/aircraft.pb

	// readsb history
	log.Println("Reading readsb-protobuf history")
	readsbProtobufHistory(readsburl, adb)

	// build the URL to data/aircraft.pb
	urlAircraftPb, err := url.Parse(readsburl)
	if err != nil {
		log.Fatal(err)
	}
	urlAircraftPb.Path = path.Join(urlAircraftPb.Path, "data/aircraft.pb")

	log.Println("Reading readsb-protobuf live data")
	// Infinite loop pulling data from readsb-protobuf
	for {

		// Get the data
		resp, err := http.Get(urlAircraftPb.String())
		if err != nil {
			log.Println("datasources.ReadsbProtobuf: HTTP error:", err)
			time.Sleep(time.Second * ERROR_SECONDS_BACKOFF)
			continue
		}
		defer resp.Body.Close()

		// Check server response
		if resp.StatusCode != http.StatusOK {
			log.Printf("datasources.ReadsbProtobuf: HTTP status was %d, expected %d", resp.StatusCode, http.StatusOK)
			time.Sleep(time.Second * ERROR_SECONDS_BACKOFF)
			continue
		}

		// Read response body
		pbData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("datasources.ReadsbProtobuf: Error reading HTTP response body: ", err)
			time.Sleep(time.Second * ERROR_SECONDS_BACKOFF)
			continue
		}

		// Attempt to unmarshall
		aircraftUpdate := &readsb_protobuf.AircraftsUpdate{}
		err = proto.Unmarshal(pbData, aircraftUpdate)
		if err != nil {
			log.Println("datasources.ReadsbProtobuf: Error unmarshalling protobuf data: ", err)
			time.Sleep(time.Second & ERROR_SECONDS_BACKOFF)
			continue
		}

		// Update aircraft DB
		for _, a := range aircraftUpdate.GetAircraft() {
			icao := int(a.GetAddr())
			adb.SetCallsign(icao, a.GetFlight())
			adb.SetPosition(icao, a.GetLat(), a.GetLon(), int(a.GetAltBaro()))
			adb.SetTrack(icao, int(a.GetTrack()))
			adb.SetCategory(icao, int(a.GetCategory()))
			adb.SetGs(icao, int(a.GetGs()))
			adb.SetAirGround(icao, a.GetAirGround())
			adb.SetLastSeen(icao)
		}

		// Wait until next update
		time.Sleep(time.Millisecond * READSB_UPDATE_FREQUENCY)
	}

}
