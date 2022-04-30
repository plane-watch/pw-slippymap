package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/akamensky/argparse"
)

//go:embed dbversion.json
var readsbDBVersionJSONBlob []byte

//go:embed aircrafts.json
var readsbAircraftsJSONBlob []byte // format: {"icao":["registration","type","flags"],...}

type runtimeConfiguration struct {
	outputFile *string
	goPackage  *string
}

func processCommandLine() runtimeConfiguration {
	// process the command line

	output := runtimeConfiguration{}

	// create new parser object
	parser := argparse.NewParser("json2gostruct", "converts readsb json to go struct")

	output.outputFile = parser.String("o", "outputfile", &argparse.Options{Required: true, Help: "Output file. Eg: './data.go'"})
	output.goPackage = parser.String("p", "package", &argparse.Options{Required: true, Help: "Output file. Eg: './data.go'"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return output
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	conf := processCommandLine()

	log.Println("Started.")

	// open output file
	f, err := os.Create(*conf.outputFile)
	check(err)
	defer f.Close()

	// write head
	packageTxt := fmt.Sprintf("package %s\n\n", *conf.goPackage)
	_, err = f.WriteString(packageTxt)
	check(err)

	// read dbversion.json
	log.Println("Processing dbversion.json")
	dbVersion := make(map[string]interface{})
	err = json.Unmarshal(readsbDBVersionJSONBlob, &dbVersion)
	check(err)

	// write dbversion.json as golang
	dbVersionTxt := fmt.Sprintf("var readsbDBVersion int = %d\n\n", int(dbVersion["version"].(float64)))
	_, err = f.WriteString(dbVersionTxt)
	check(err)

	// write readsbAircraft struct
	_, err = f.WriteString(`type readsbAircraftEntry struct {
		registration string
		aircraftType string
		// flags        int
	}`)
	check(err)
	_, err = f.WriteString("\n\n")
	check(err)

	// write var
	_, err = f.WriteString("var readsbAircraft = map[int]readsbAircraftEntry{\n")
	check(err)

	// read aircraft.json
	log.Println("Processing aircraft.json")
	aircraftData := make(map[string]interface{})
	err = json.Unmarshal(readsbAircraftsJSONBlob, &aircraftData)
	check(err)

	for k, v := range aircraftData {
		// get icao
		icao64, err := strconv.ParseInt(k, 16, 64)
		check(err)

		icao := int(icao64)

		// sanity checks
		if reflect.TypeOf(v).Kind() != reflect.Slice {
			log.Fatal("aircraft.json: JSON data not type of slice")
		}

		// get registration & aircraft type
		vr := reflect.ValueOf(v)
		registration := vr.Index(0).Elem().String()
		aircraftType := vr.Index(1).Elem().String()

		// write contents of var
		icaoTxt := fmt.Sprintf("\t0x%X: {\n", icao)
		_, err = f.WriteString(icaoTxt)
		check(err)

		regTxt := fmt.Sprintf("\t\tregistration: \"%s\",\n", registration)
		_, err = f.WriteString(regTxt)
		check(err)

		actTxt := fmt.Sprintf("\t\taircraftType: \"%s\",\n", aircraftType)
		_, err = f.WriteString(actTxt)
		check(err)

		_, err = f.WriteString("\t},\n")
		check(err)
	}

	_, err = f.WriteString("}\n")
	check(err)

	log.Println("Done.")
}
