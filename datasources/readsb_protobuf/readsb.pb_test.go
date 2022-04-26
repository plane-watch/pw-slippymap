package readsb_protobuf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// this is currently the data that we care about
// could eventually expand to cover everything
type TestAircraftMeta struct {
	Addr           uint32
	Flight         string
	Squawk         uint32
	Category       uint32
	AltBaro        int32
	MagHeading     int32
	Ias            uint32
	Lat            float64
	Lon            float64
	Messages       uint64
	Seen           uint64
	Rssi           float32
	Distance       uint32
	AirGround      AircraftMeta_AirGround
	AltGeom        int32
	BaroRate       int32
	Gs             uint32
	Tas            uint32
	Mach           float32
	TrueHeading    int32
	Track          int32
	Roll           float32
	NavQnh         float32
	NavAltitudeMcp int32
	NavAltitudeFms int32
	Nic            uint32
	Rc             uint32
	Version        int32
	NicBaro        uint32
	NacP           uint32
	NacV           uint32
	Sil            uint32
	SeenPos        uint32
	Alert          bool
	Spi            bool
	Gva            uint32
	Sda            uint32
	Declination    float64
	WindSpeed      uint32
	WindDirection  uint32
	AddrType       AircraftMeta_AddrType
	Emergency      AircraftMeta_Emergency
	SilType        AircraftMeta_SilType
	NavModes       *AircraftMeta_NavModes
}

func TestReadsbProtobuf(t *testing.T) {

	var testName string

	// define expected results
	// these were obtained from manually running protoc over the aircraft.pb test file
	expected := make(map[uint32]TestAircraftMeta)
	expected[0x7C79CA] = TestAircraftMeta{
		Addr:           0x7C79CA,
		Flight:         "HARR89  ",
		Squawk:         12288,
		Category:       161,
		AltBaro:        1300,
		Lat:            -32.1251220703125,
		Lon:            115.87999877929688,
		Messages:       234,
		Seen:           1650938678764,
		Rssi:           -24.69765,
		Distance:       25692,
		AirGround:      AircraftMeta_AirGround(AircraftMeta_AirGround_value["AG_AIRBORNE"]),
		AltGeom:        975,
		Gs:             104,
		Track:          111,
		NavAltitudeMcp: 24096,
		Nic:            8,
		Rc:             186,
		Version:        2,
		NicBaro:        1,
		NacP:           9,
		NacV:           2,
		Sil:            3,
		SeenPos:        1,
		Gva:            2,
		Sda:            2,
		Declination:    -1.7332517322141492,
		SilType:        3,
	}
	expected[0x7C79D8] = TestAircraftMeta{
		Addr:           0x7C79D8,
		Flight:         "HARR81  ",
		Squawk:         0x1200,
		Category:       0xa1,
		AltBaro:        4750,
		Lat:            -32.34114074707031,
		Lon:            115.931689453125,
		Messages:       0x9c,
		Seen:           0x180639d0e15,
		Rssi:           -25.495672,
		Distance:       0xc096,
		AirGround:      2,
		AltGeom:        4500,
		Gs:             0x65,
		Track:          169,
		NavAltitudeMcp: 33536,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x2,
		Sil:            0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.7937696229174418,
		SilType:        3,
	}
	expected[0x7CF9DD] = TestAircraftMeta{
		Addr:           0x7CF9DD,
		Flight:         "VODO    ",
		Squawk:         0x2050,
		Category:       0xa1,
		AltBaro:        9875,
		MagHeading:     7,
		Lat:            -31.307510375976562,
		Lon:            116.06716380399816,
		Messages:       0x15d,
		Seen:           0x180639d0e76,
		Rssi:           -19.07037,
		Distance:       0x10575,
		AirGround:      2,
		AltGeom:        9725,
		BaroRate:       -2240,
		Gs:             0x124,
		Tas:            0x132,
		Mach:           0.468,
		Track:          8,
		Roll:           5.9765625,
		NavQnh:         1005,
		NavAltitudeMcp: 12000,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x1,
		Sil:            0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.3558498916806274,
		SilType:        3,
		Ias:            0x103,
	}
	expected[0x7CF9DE] = TestAircraftMeta{
		Addr:           0x7CF9DE,
		Flight:         "VIPR04  ",
		Squawk:         0x2001,
		Category:       0xa1,
		AltBaro:        3025,
		MagHeading:     138,
		Ias:            0xc9,
		Lat:            -31.691665649414062,
		Lon:            115.79552145565258,
		Messages:       0x1af,
		Seen:           0x180639d0d93,
		Rssi:           -14.32674,
		Distance:       0x6619,
		AirGround:      2,
		AltGeom:        2775,
		BaroRate:       -64,
		Gs:             0xdc,
		Tas:            0xd2,
		Mach:           0.32,
		Track:          136,
		Roll:           -15.292969,
		NavQnh:         1005,
		NavAltitudeMcp: 2800,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x1,
		Sil:            0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.6082767215688187,
		SilType:        3,
	}
	expected[0x7CF9DF] = TestAircraftMeta{
		Addr:           0x7CF9DF,
		Flight:         "VIPR10  ",
		Squawk:         0x2070,
		Category:       0xa1,
		AltBaro:        3050,
		MagHeading:     328,
		Ias:            0xb9,
		Lat:            -31.630279541015625,
		Lon:            116.00286147173713,
		Messages:       0xf2,
		Seen:           0x180639d0e51,
		Rssi:           -7.988089,
		Distance:       0x7761,
		AirGround:      2,
		AltGeom:        2775,
		BaroRate:       3840,
		Gs:             0xb8,
		Tas:            0xc4,
		Mach:           0.3,
		Track:          331,
		Roll:           0.87890625,
		NavQnh:         1005,
		NavAltitudeMcp: 12000,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacV:           0x1,
		NacP:           0x9,
		Sil:            0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.4956151326506066,
		SilType:        3,
	}
	expected[0x7CF9E2] = TestAircraftMeta{
		Addr:           0x7CF9E2,
		Flight:         "VORT2   ",
		Squawk:         0x2015,
		Category:       0xa1,
		AltBaro:        7200,
		MagHeading:     318,
		Ias:            0xef,
		Lat:            -31.35283906581037,
		Lon:            115.66164550781251,
		Messages:       0xc7,
		Seen:           0x180639d025e,
		Rssi:           -26.951849,
		Distance:       0x10072,
		AirGround:      2,
		AltGeom:        7000,
		BaroRate:       -128,
		Gs:             0xfb,
		Tas:            0x106,
		Mach:           0.412,
		Track:          320,
		Roll:           -5.2734375,
		NavQnh:         1005,
		NavAltitudeMcp: 12000,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x1,
		Sil:            0x3,
		SeenPos:        0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.5432304901119622,
		SilType:        3,
	}
	expected[0x7CF9E8] = TestAircraftMeta{
		Addr:           0x7CF9E8,
		Flight:         "SBOT    ",
		Squawk:         0x2022,
		Category:       0xa1,
		AltBaro:        8700,
		MagHeading:     314,
		Ias:            0xf0,
		Lat:            -31.197509765625,
		Lon:            115.84587545955883,
		Messages:       0x112,
		Seen:           0x180639d0dc6,
		Rssi:           -20.246077,
		Distance:       0x131b1,
		AirGround:      2,
		AltGeom:        8525,
		Gs:             0xfb,
		Tas:            0x112,
		Mach:           0.424,
		Track:          319,
		Roll:           -29.882812,
		NavQnh:         1005,
		NavAltitudeMcp: 12000,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x1,
		Sil:            0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.4102694513764835,
		SilType:        3,
	}
	expected[0x7C4A00] = TestAircraftMeta{
		Addr:           0x7C4A00,
		Flight:         "FD607   ",
		Squawk:         0x4070,
		Category:       0xa1,
		AltBaro:        24000,
		MagHeading:     337,
		Ias:            0xb2,
		Lat:            -30.763720496226163,
		Lon:            115.25399780273438,
		Messages:       0x3a,
		Seen:           0x180639d0837,
		Rssi:           -28.34106,
		Distance:       0x22880,
		AirGround:      2,
		AltGeom:        24375,
		Gs:             0xce,
		Mach:           0.428,
		Track:          345,
		NavQnh:         1013.6,
		NavAltitudeMcp: 24000,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x2,
		Sil:            0x3,
		SeenPos:        0x1,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.5073057927923468,
		SilType:        3,
		NavModes:       &AircraftMeta_NavModes{Autopilot: true, Vnav: true, Althold: false, Approach: false, Lnav: true, Tcas: false},
	}
	expected[0x7C7A4A] = TestAircraftMeta{
		Addr:        0x7C7A4A,
		Flight:      "VOZ1852 ",
		Squawk:      0x4056,
		Category:    0xa0,
		AltBaro:     275,
		MagHeading:  11,
		Ias:         0x93,
		Lat:         -31.942989349365234,
		Lon:         115.9643051147461,
		Messages:    0x139,
		Seen:        0x180639d0a57,
		Rssi:        -23.775469,
		Distance:    0x17a9,
		AirGround:   1,
		AltGeom:     -25,
		BaroRate:    -992,
		Gs:          0x31,
		Tas:         0x96,
		Mach:        0.224,
		Track:       14,
		Roll:        0.17578125,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		NacV:        0x2,
		Sil:         0x2,
		SeenPos:     0x1,
		Declination: -1.6264424232041375,
		SilType:     1,
	}
	expected[0x7C7A6E] = TestAircraftMeta{
		Addr:        0x7C7A6E,
		Flight:      "YGW     ",
		Squawk:      0x3000,
		Category:    0xa1,
		AltBaro:     675,
		Lat:         -32.08644104003906,
		Lon:         115.90114746093751,
		Messages:    0x172,
		Seen:        0x180639d0b65,
		Rssi:        -23.40051,
		Distance:    0x5292,
		AirGround:   2,
		AltGeom:     375,
		BaroRate:    -64,
		Gs:          0x38,
		Track:       237,
		Nic:         0x8,
		Rc:          0xba,
		Version:     1,
		NicBaro:     0x1,
		NacP:        0x9,
		NacV:        0x2,
		Sil:         0x3,
		Declination: -1.7088143056593028,
		SilType:     1,
	}
	expected[0x7C7A90] = TestAircraftMeta{
		Addr:           0x7C7A90,
		Flight:         "YHU     ",
		Squawk:         0x1200,
		Category:       0xa1,
		AltBaro:        2800,
		MagHeading:     198,
		Ias:            0x60,
		Lat:            -32.182891845703125,
		Lon:            116.14246215820313,
		Messages:       0x1c8,
		Seen:           0x180639d0e52,
		Rssi:           -19.354033,
		Distance:       0x92c6,
		AirGround:      2,
		AltGeom:        2475,
		BaroRate:       -64,
		Gs:             0x63,
		Tas:            0x60,
		Mach:           0.152,
		Track:          191,
		Roll:           3.3398438,
		NavQnh:         1004,
		NavAltitudeMcp: 2496,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x2,
		Sil:            0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.6368968320553683,
		SilType:        3,
	}
	expected[0x7C1293] = TestAircraftMeta{
		Addr:        0x7C1293,
		Flight:      "DYD     ",
		Squawk:      0x3603,
		Category:    0xa2,
		AltBaro:     3600,
		Lat:         -31.913681030273438,
		Lon:         115.86769409179688,
		Messages:    0x1e8,
		Seen:        0x180639d0e83,
		Rssi:        -2.4767873,
		Distance:    0x1779,
		AirGround:   3,
		AltGeom:     3350,
		Gs:          0xa4,
		Track:       228,
		Nic:         0x8,
		Rc:          0xba,
		Version:     2,
		NacP:        0xa,
		NacV:        0x1,
		Sil:         0x3,
		Gva:         0x2,
		Sda:         0x2,
		Declination: -1.6592538073691407,
		SilType:     3,
	}
	expected[0x7C7AA9] = TestAircraftMeta{
		Addr:      0x7C7AA9,
		Flight:    "VOZ1481 ",
		Category:  0xa0,
		Lat:       -31.942840576171875,
		Lon:       115.96784820556641,
		Messages:  0x18,
		Seen:      0x180639d0c2b,
		Rssi:      -23.818779,
		Distance:  0x185c,
		AirGround: 1,
		Gs:        0xb,
		Track:     194,
		Nic:       0x8,
		Rc:        0xba,
		NacP:      0x8,
		Sil:       0x2,
		SilType:   1,
	}
	expected[0x7C1AB8] = TestAircraftMeta{
		Addr:        0x7C1AB8,
		Flight:      "UTY3790 ",
		Squawk:      0x4271,
		Category:    0xa0,
		AltBaro:     19350,
		Lat:         -31.51018562963452,
		Lon:         116.85151977539063,
		Messages:    0x1d8,
		Seen:        0x180639d0e78,
		Rssi:        -19.262897,
		Distance:    0x17c5c,
		AirGround:   2,
		AltGeom:     19550,
		Gs:          0x173,
		Track:       37,
		NavQnh:      1004,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		Sil:         0x2,
		Declination: -1.0963232490481147,
		SilType:     1,
	}
	expected[0x7C7AB8] = TestAircraftMeta{
		Addr:           0x7C7AB8,
		Flight:         "VOZ551  ",
		Squawk:         0x4313,
		Category:       0xa0,
		AltBaro:        20150,
		MagHeading:     289,
		Ias:            0xda,
		Lat:            -32.202484130859375,
		Lon:            116.98939819335938,
		Messages:       0x7e,
		Seen:           0x180639d0985,
		Rssi:           -26.000883,
		Distance:       0x19c60,
		AirGround:      2,
		AltGeom:        20350,
		BaroRate:       -1504,
		Gs:             0xef,
		Tas:            0x12c,
		Mach:           0.48,
		Track:          285,
		Roll:           0.3515625,
		NavQnh:         1013.2,
		NavAltitudeMcp: 9008,
		NavAltitudeFms: 9008,
		Nic:            0x8,
		Rc:             0xba,
		NacP:           0x8,
		NacV:           0x2,
		Sil:            0x2,
		SeenPos:        0x2,
		Declination:    -1.269151882338225,
		WindSpeed:      0x3d,
		WindDirection:  0x12a,
		SilType:        1,
	}
	expected[0x7C42CE] = TestAircraftMeta{
		Addr:        0x7C42CE,
		Flight:      "NWK2905 ",
		Squawk:      0x4033,
		Category:    0xa0,
		AltBaro:     28000,
		Lat:         -31.079452514648438,
		Lon:         116.29900764016544,
		Messages:    0x19a,
		Seen:        0x180639d0e81,
		Rssi:        -17.369724,
		Distance:    0x17d01,
		AirGround:   2,
		AltGeom:     28500,
		Gs:          0x1b1,
		Track:       198,
		NavQnh:      1011,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		Sil:         0x2,
		Declination: -1.1869017355068825,
		SilType:     1,
	}
	expected[0x7C42D8] = TestAircraftMeta{
		Addr:        0x7C42D8,
		Flight:      "NWK1617 ",
		Squawk:      0x3654,
		Category:    0xa0,
		AltBaro:     21550,
		Lat:         -31.199752807617188,
		Lon:         116.37898164636948,
		Messages:    0x9c,
		Seen:        0x180639d0d39,
		Rssi:        -17.812311,
		Distance:    0x15a01,
		AirGround:   2,
		AltGeom:     21850,
		Gs:          0x17e,
		Track:       198,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		Sil:         0x2,
		Declination: -1.1920689869381709,
		SilType:     1,
	}
	expected[0x7C42E4] = TestAircraftMeta{
		Addr:      0x7C42E4,
		Flight:    "NWK2909 ",
		Squawk:    0x7242,
		Category:  0xa0,
		AltBaro:   300,
		Lat:       -31.93151092529297,
		Lon:       115.9602264404297,
		Messages:  0x26,
		Seen:      0x180639c538a,
		Rssi:      -18.270395,
		Distance:  0x12bf,
		AirGround: 1,
		Track:     244,
		NavQnh:    1004,
		Nic:       0x8,
		Rc:        0xba,
		NacP:      0x8,
		Sil:       0x2,
		SeenPos:   0x2f,
		SilType:   1,
	}
	expected[0x7C632F] = TestAircraftMeta{
		Addr:      0x7C632F,
		Flight:    "TVL     ",
		Squawk:    0x1200,
		AltBaro:   1800,
		Lat:       -32.05075021517479,
		Lon:       115.7626031369579,
		Messages:  0xf9,
		Seen:      0x180639d0d0c,
		Rssi:      -10.652098,
		AirGround: 3,
		Gs:        0x63,
		Track:     288,
		SeenPos:   0x4,
		SilType:   1,
	}
	expected[0x7C534A] = TestAircraftMeta{
		Addr:        0x7C534A,
		Flight:      "QQK     ",
		Squawk:      0x3771,
		Category:    0xa0,
		AltBaro:     17300,
		Lat:         -31.644012451171875,
		Lon:         116.2130557789522,
		Messages:    0x160,
		Seen:        0x180639d0e51,
		Rssi:        -16.353012,
		Distance:    0x9845,
		AirGround:   2,
		AltGeom:     17525,
		Gs:          0xf6,
		Track:       227,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		Sil:         0x2,
		Declination: -1.4132645494658407,
		SilType:     1,
	}
	expected[0x7C7B87] = TestAircraftMeta{
		Addr:           0x7C7B87,
		Flight:         "YOP     ",
		Squawk:         0x3000,
		Category:       0xa1,
		AltBaro:        1250,
		Lat:            -32.1038818359375,
		Lon:            115.91510009765625,
		Messages:       0x175,
		Seen:           0x180639d0cf9,
		Rssi:           -14.958848,
		Distance:       0x59aa,
		AirGround:      2,
		AltGeom:        900,
		Gs:             0x65,
		Track:          42,
		NavQnh:         1004,
		NavAltitudeMcp: 992,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x2,
		Sil:            0x3,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -1.7092811964451133,
		SilType:        3,
	}
	expected[0x7C2BC5] = TestAircraftMeta{
		Addr:           0x7C2BC5,
		Flight:         "VOZ1876 ",
		Squawk:         0x4232,
		Category:       0xa0,
		AltBaro:        36000,
		MagHeading:     220,
		Ias:            0xdc,
		Lat:            -30.262985229492188,
		Lon:            116.77045036764706,
		Messages:       0xf8,
		Seen:           0x180639d0e11,
		Rssi:           -20.032932,
		Distance:       0x30811,
		AirGround:      2,
		AltGeom:        37050,
		BaroRate:       -32,
		Gs:             0x169,
		Tas:            0x190,
		Mach:           0.672,
		Track:          202,
		Roll:           0.52734375,
		NavQnh:         1013.2,
		NavAltitudeMcp: 36000,
		Nic:            0x8,
		Rc:             0xba,
		NacP:           0x8,
		Sil:            0x2,
		Declination:    -0.7496800749027261,
		WindSpeed:      0x79,
		WindDirection:  0x11a,
		SilType:        1,
	}
	expected[0x7C149C] = TestAircraftMeta{
		Addr:           0x7C149C,
		Flight:         "ECU     ",
		Squawk:         0x3000,
		AltBaro:        875,
		MagHeading:     265,
		Ias:            0x46,
		Lat:            -32.08003157276221,
		Lon:            115.91377725406569,
		Messages:       0x9f,
		Seen:           0x180639d0e36,
		Rssi:           -16.803148,
		AirGround:      2,
		BaroRate:       -64,
		Gs:             0x39,
		Mach:           0.108,
		Track:          287,
		NavQnh:         1004,
		NavAltitudeMcp: 1008,
		SeenPos:        0x1,
		SilType:        1,
	}
	expected[0x7C753B] = TestAircraftMeta{
		Addr:        0x7C753B,
		Flight:      "PY4452  ",
		Squawk:      0x4252,
		Category:    0xa2,
		AltBaro:     4250,
		Lat:         -32.06703186035156,
		Lon:         116.10219726562501,
		Messages:    0x1e8,
		Seen:        0x180639d0e76,
		Rssi:        -8.479529,
		Distance:    0x6177,
		AirGround:   2,
		AltGeom:     4000,
		Gs:          0xb3,
		Track:       267,
		Nic:         0x8,
		Rc:          0xba,
		Version:     1,
		NicBaro:     0x1,
		NacP:        0x9,
		NacV:        0x2,
		Sil:         0x2,
		Declination: -1.6120747963497226,
		SilType:     1,
	}
	expected[0x7C6D8D] = TestAircraftMeta{
		Addr:           0x7C6D8D,
		Flight:         "QFA1203 ",
		Squawk:         0x3242,
		Category:       0xa0,
		AltBaro:        25000,
		MagHeading:     232,
		Ias:            0x121,
		Lat:            -31.08306884765625,
		Lon:            116.64459228515625,
		Messages:       0x13f,
		Seen:           0x180639d0e15,
		Rssi:           -20.23061,
		Distance:       0x1ba38,
		AirGround:      2,
		AltGeom:        25500,
		BaroRate:       -928,
		Gs:             0x197,
		Tas:            0x1a8,
		Mach:           0.692,
		Track:          221,
		Roll:           0.17578125,
		NavQnh:         1013.2,
		NavAltitudeMcp: 12000,
		NavAltitudeFms: 3008,
		Nic:            0x8,
		Rc:             0xba,
		NacP:           0x8,
		NacV:           0x2,
		Sil:            0x2,
		Declination:    -1.045706671803757,
		WindSpeed:      0x4b,
		WindDirection:  0x12d,
		SilType:        1,
	}
	expected[0x7C6DB9] = TestAircraftMeta{
		Addr:          0x7C6DB9,
		Flight:        "QFA1357 ",
		Squawk:        0x3606,
		Category:      0xa0,
		AltBaro:       1375,
		MagHeading:    346,
		Ias:           0x97,
		Lat:           -32.01048252946242,
		Lon:           115.94864203005422,
		Messages:      0x1d7,
		Seen:          0x180639d0e7f,
		Rssi:          -5.4253144,
		Distance:      0x3184,
		AirGround:     2,
		AltGeom:       1075,
		BaroRate:      -576,
		Gs:            0x9b,
		Tas:           0x9c,
		Mach:          0.236,
		Track:         349,
		Roll:          8.964844,
		Nic:           0x8,
		Rc:            0xba,
		NacP:          0x8,
		NacV:          0x2,
		Sil:           0x2,
		Declination:   -1.658948216502071,
		WindSpeed:     0xb,
		WindDirection: 0x110,
		SilType:       1,
	}
	expected[0x7C6DE0] = TestAircraftMeta{
		Addr:           0x7C6DE0,
		Squawk:         0x4377,
		AltBaro:        275,
		Lat:            -31.931355363231596,
		Lon:            115.96119783362565,
		Messages:       0x33,
		Seen:           0x180639d0c81,
		Rssi:           -23.28562,
		Distance:       0x12eb,
		AirGround:      1,
		Gs:             0x7,
		Track:          295,
		NavQnh:         1004,
		NavAltitudeMcp: 3008,
		Nic:            0x8,
		Rc:             0xba,
		NacP:           0x8,
		Sil:            0x2,
		SilType:        1,
	}
	expected[0x7C6DE5] = TestAircraftMeta{
		Addr:        0x7C6DE5,
		Flight:      "QFA933  ",
		Squawk:      0x1116,
		Category:    0xa0,
		AltBaro:     7950,
		MagHeading:  281,
		Ias:         0xf8,
		Lat:         -32.11098881091101,
		Lon:         116.37055066167092,
		Messages:    0x6e,
		Seen:        0x180639d0bb5,
		Rssi:        -25.261557,
		Distance:    0xbb4f,
		AirGround:   2,
		AltGeom:     7825,
		BaroRate:    -1568,
		Gs:          0xff,
		Tas:         0x116,
		Mach:        0.432,
		Track:       279,
		Roll:        -0.3515625,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		NacV:        0x2,
		Sil:         0x2,
		SeenPos:     0x1,
		Declination: -1.5091838930080475,
		SilType:     1,
	}
	expected[0x7C66D3] = TestAircraftMeta{
		Addr:           0x7C66D3,
		Squawk:         0x4312,
		AltBaro:        36000,
		Lat:            -31.195037841796875,
		Lon:            118.17406149471508,
		Messages:       0x18,
		Seen:           0x180639cf990,
		Rssi:           -30.378042,
		Distance:       0x37576,
		AirGround:      2,
		AltGeom:        36950,
		BaroRate:       64,
		Gs:             0x13d,
		Tas:            0x18c,
		Track:          231,
		Roll:           -0.17578125,
		NavQnh:         1013.6,
		NavAltitudeMcp: 25024,
		Nic:            0x8,
		Rc:             0xba,
		Version:        2,
		NicBaro:        0x1,
		NacP:           0x9,
		NacV:           0x2,
		Sil:            0x3,
		SeenPos:        0x31,
		Gva:            0x2,
		Sda:            0x2,
		Declination:    -0.4645265036861029,
		SilType:        3,
		NavModes:       &AircraftMeta_NavModes{Autopilot: true, Vnav: true, Althold: false, Approach: false, Lnav: false, Tcas: true},
	}
	expected[0x76E721] = TestAircraftMeta{
		Addr:      0x76E721,
		Flight:    "EGLE841 ",
		Squawk:    0x2012,
		Category:  0xa1,
		AltBaro:   14200,
		Messages:  0x11,
		Seen:      0x180639d0b24,
		Rssi:      -25.25525,
		AirGround: 2,
		AltGeom:   14250,
		Gs:        0x127,
		Track:     355,
		NacP:      0x8,
		Sil:       0x2,
		SilType:   1,
	}
	expected[0x76E722] = TestAircraftMeta{
		Addr:        0x76E722,
		Squawk:      0x2030,
		AltBaro:     3750,
		Lat:         -31.311500678628192,
		Lon:         115.70147094726563,
		Messages:    0xc,
		Seen:        0x180639d0d16,
		Rssi:        -27.346859,
		Distance:    0x10c23,
		AirGround:   3,
		AltGeom:     4025,
		Gs:          0xcb,
		Track:       236,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		Sil:         0x2,
		SeenPos:     0x25,
		Declination: -1.5105853884307412,
		SilType:     1,
	}
	expected[0x76E725] = TestAircraftMeta{
		Addr:      0x76E725,
		AltBaro:   2525,
		Messages:  0x5,
		Seen:      0x180639d0d87,
		Rssi:      -31.293661,
		AirGround: 2,
		Gs:        0x113,
		Track:     5,
	}
	expected[0x76E72A] = TestAircraftMeta{
		Addr:      0x76E72A,
		AltBaro:   1550,
		Messages:  0x8,
		Seen:      0x180639d0e36,
		Rssi:      -27.4397,
		AirGround: 2,
		AltGeom:   1300,
		Gs:        0x9f,
		Track:     255,
	}
	expected[0x76E72D] = TestAircraftMeta{
		Addr:        0x76E72D,
		Flight:      "TANG836 ",
		Squawk:      0x2033,
		Category:    0xa1,
		AltBaro:     2525,
		Lat:         -31.48397671974311,
		Lon:         115.85335693359376,
		Messages:    0x17,
		Seen:        0x180639cd8a2,
		Rssi:        -25.95451,
		Distance:    0xb5dd,
		AirGround:   2,
		AltGeom:     2275,
		Gs:          0xa6,
		Track:       73,
		Nic:         0x8,
		Rc:          0xba,
		NacP:        0x8,
		Sil:         0x2,
		SeenPos:     0x11,
		Declination: -1.5072859601450626,
		SilType:     1,
	}
	expected[0x7CF7C4] = TestAircraftMeta{
		Addr:     0x7CF7C4,
		Flight:   "PHRX1A  ",
		Category: 0xc0,
		Messages: 0x39,
		Seen:     0x180639d0aaf,
		Rssi:     -18.499548,
		AddrType: 1,
	}
	expected[0x7CF7C5] = TestAircraftMeta{
		Addr:     0x7CF7C5,
		Flight:   "PHRX1B  ",
		Category: 0xc0,
		Messages: 0x25,
		Seen:     0x180639d0c1f,
		Rssi:     -19.185398,
		AddrType: 1,
	}
	expected[0x7CF7C6] = TestAircraftMeta{
		Addr:     0x7CF7C6,
		Flight:   "PHRX2A  ",
		Category: 0xc0,
		Messages: 0x39,
		Seen:     0x180639d0d5c,
		Rssi:     -17.280907,
		AddrType: 1,
	}
	expected[0x7CF7C7] = TestAircraftMeta{
		Addr:     0x7CF7C7,
		Flight:   "PHRX2B  ",
		Category: 0xc0,
		Messages: 0x15,
		Seen:     0x180639d0d7c,
		Rssi:     -21.576744,
		AddrType: 1,
	}

	// open test data file
	testDataFile := filepath.Join("testdata", "aircraft.pb")
	pbData, err := os.ReadFile(testDataFile)
	require.NoError(t, err)

	// url := "http://192.168.69.35:8079/data/aircraft.pb"

	// // Get data
	// resp, err := http.Get(url)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer resp.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(pbData)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	m := &AircraftsUpdate{}
	err = proto.Unmarshal(pbData, m)

	for _, a := range m.GetAircraft() {

		testName = fmt.Sprintf("%s: %X", testDataFile, a.GetAddr())
		t.Run(testName, func(t *testing.T) {
			e, inMap := expected[a.GetAddr()]
			assert.Equal(t, true, inMap)
			if !inMap {
				t.FailNow()
			}

			testName = fmt.Sprintf("Test GetAddr (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Addr, a.GetAddr())
			})

			testName = fmt.Sprintf("Test GetFlight (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Flight, a.GetFlight())
			})

			testName = fmt.Sprintf("Test GetSquawk (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Squawk, a.GetSquawk())
			})

			testName = fmt.Sprintf("Test GetCategory (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Category, a.GetCategory())
			})

			testName = fmt.Sprintf("Test GetAltBaro (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.AltBaro, a.GetAltBaro())
			})

			testName = fmt.Sprintf("Test GetMagHeading (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.MagHeading, a.GetMagHeading())
			})

			testName = fmt.Sprintf("Test GetIas (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Ias, a.GetIas())
			})

			testName = fmt.Sprintf("Test GetLat (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Lat, a.GetLat())
			})

			testName = fmt.Sprintf("Test GetLon (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Lon, a.GetLon())
			})

			testName = fmt.Sprintf("Test GetMessages (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Messages, a.GetMessages())
			})

			testName = fmt.Sprintf("Test GetSeen (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Seen, a.GetSeen())
			})

			testName = fmt.Sprintf("Test GetRssi (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Rssi, a.GetRssi())
			})

			testName = fmt.Sprintf("Test GetDistance (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Distance, a.GetDistance())
			})

			testName = fmt.Sprintf("Test GetAirGround (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.AirGround, a.GetAirGround())
			})

			testName = fmt.Sprintf("Test GetAltGeom (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.AltGeom, a.GetAltGeom())
			})

			testName = fmt.Sprintf("Test GetBaroRate (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.BaroRate, a.GetBaroRate())
			})

			testName = fmt.Sprintf("Test GetGs (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Gs, a.GetGs())
			})

			testName = fmt.Sprintf("Test GetTas (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Tas, a.GetTas())
			})

			testName = fmt.Sprintf("Test GetMach (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Mach, a.GetMach())
			})

			testName = fmt.Sprintf("Test GetTrueHeading (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.TrueHeading, a.GetTrueHeading())
			})

			testName = fmt.Sprintf("Test GetTrack (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Track, a.GetTrack())
			})

			testName = fmt.Sprintf("Test GetRoll (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Roll, a.GetRoll())
			})

			testName = fmt.Sprintf("Test GetNavQnh (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.NavQnh, a.GetNavQnh())
			})

			testName = fmt.Sprintf("Test GetNavAltitudeMcp (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.NavAltitudeMcp, a.GetNavAltitudeMcp())
			})

			testName = fmt.Sprintf("Test GetNavAltitudeFms (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.NavAltitudeFms, a.GetNavAltitudeFms())
			})

			testName = fmt.Sprintf("Test GetNic (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Nic, a.GetNic())
			})

			testName = fmt.Sprintf("Test GetRc (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Rc, a.GetRc())
			})

			testName = fmt.Sprintf("Test GetVersion (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Version, a.GetVersion())
			})

			testName = fmt.Sprintf("Test GetNicBaro (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.NicBaro, a.GetNicBaro())
			})

			testName = fmt.Sprintf("Test GetNacP (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.NacP, a.GetNacP())
			})

			testName = fmt.Sprintf("Test GetNacV (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.NacV, a.GetNacV())
			})

			testName = fmt.Sprintf("Test GetSil (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Sil, a.GetSil())
			})

			testName = fmt.Sprintf("Test GetSeenPos (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.SeenPos, a.GetSeenPos())
			})

			testName = fmt.Sprintf("Test GetAlert (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Alert, a.GetAlert())
			})

			testName = fmt.Sprintf("Test GetSpi (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Spi, a.GetSpi())
			})

			testName = fmt.Sprintf("Test GetGva (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Gva, a.GetGva())
			})

			testName = fmt.Sprintf("Test GetSda (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Sda, a.GetSda())
			})

			testName = fmt.Sprintf("Test GetDeclination (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Declination, a.GetDeclination())
			})

			testName = fmt.Sprintf("Test GetWindSpeed (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.WindSpeed, a.GetWindSpeed())
			})

			testName = fmt.Sprintf("Test GetWindDirection (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.WindDirection, a.GetWindDirection())
			})

			testName = fmt.Sprintf("Test GetAddrType (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.AddrType, a.GetAddrType())
			})

			testName = fmt.Sprintf("Test GetEmergency (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.Emergency, a.GetEmergency())
			})

			testName = fmt.Sprintf("Test GetSilType (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.SilType, a.GetSilType())
			})

			testName = fmt.Sprintf("Test GetNavModes (%X)", a.GetAddr())
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, e.NavModes, a.GetNavModes())
			})
		})
	}
}
