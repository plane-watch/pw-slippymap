package datasources

import (
	"io/ioutil"
	"log"
	"net/http"
	"pw_slippymap/datasources/readsb_protobuf"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	ERROR_SECONDS_BACKOFF   = 5
	READSB_UPDATE_FREQUENCY = 250 // milliseconds
)

func ReadsbProtobuf(url string, adb *AircraftDB) {
	// Updates the AircraftDB adb from aircraft.pb located at url

	for {

		// Get the data
		resp, err := http.Get(url)
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
			adb.SetLat(icao, a.GetLat())
			adb.SetLong(icao, a.GetLon())
			adb.SetTrack(icao, int(a.GetTrack()))
			adb.SetAltBaro(icao, int(a.GetAltBaro()))
			adb.SetLastSeen(icao)
		}

		// Get aircraft history (for trails)
		// TODO - aircraft.pb doesn't seem to have history, probably need history.pb...
		// for _, h := range aircraftUpdate.GetHistory() {
		// 	icao := int(h.GetAddr())
		// 	lat := h.GetLat()
		// 	long := h.GetLon()
		// 	altBaro := h.GetAltBaro()
		// 	fmt.Println(icao, lat, long, altBaro)

		// }

		// Wait until next update
		time.Sleep(time.Millisecond * READSB_UPDATE_FREQUENCY)
	}

}
