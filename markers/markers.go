package markers

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// SVG paths for aircraft

	// Airbus A380
	AIRBUS_A380_SVGPATH = "m 244.73958,0 c -19.45177,2.9398148 -21.49332,76.729166 -21.49332,76.729166 v 35.718754 c -2.84181,7.02289 -10.27301,13.22916 -10.27301,13.22916 l -45.64879,37.48264 c 0.57163,-5.30799 0.32665,-9.71772 0.32665,-9.71772 0,-10.45268 -2.1232,-14.8624 -2.1232,-14.8624 h -19.35378 c -2.28653,2.36819 -2.04154,16.16899 -2.04154,16.16899 0.24498,13.22916 1.95987,18.4555 1.95987,18.4555 h 2.7765 l 1.03464,3.86133 -49.2966,36.3978 c 0.97994,-4.57305 0.89827,-11.10597 0.89827,-11.10597 -0.0817,-15.18904 -2.123197,-16.90393 -2.123197,-16.90393 H 79.946631 c -1.796554,1.0616 -2.1232,15.59735 -2.1232,15.59735 0.244984,13.55582 2.204861,19.02713 2.204861,19.02713 h 2.69483 l 1.306585,5.38966 -71.698817,52.99833 c -9.1460906,7.0229 -9.0644291,20.66037 -9.0644291,20.66037 l -0.3266461,22.13027 1.714892,-12.41255 80.6815842,-35.03279 c 0.408308,11.10596 1.388246,10.77932 1.388246,10.77932 1.388246,-0.48997 2.368184,-12.65754 2.368184,-12.65754 l 21.721969,-8.65612 c 0.0817,13.63747 1.22492,13.55581 1.22492,13.55581 2.04154,-0.89827 3.42978,-15.51569 3.42978,-15.51569 l 20.98701,-8.32947 c 0.32665,14.61741 1.55157,14.45409 1.55157,14.45409 2.85816,-5.55298 2.93982,-16.49563 2.93982,-16.49563 l 20.33372,-8.08449 c 0.24498,6.12461 1.38824,6.12461 1.38824,6.12461 1.30659,-0.0817 2.20486,-7.67618 2.20486,-7.67618 l 9.96271,-3.91975 10.94264,-2.93982 c 0.40831,6.28794 1.46991,6.12462 1.46991,6.12462 1.38825,0.0817 2.04154,-7.18622 2.04154,-7.18622 l 29.72479,-7.8395 c 0.73495,21.39532 4.16474,35.35943 4.16474,35.35943 v 47.28203 c -0.0817,8.32947 3.34812,32.17464 3.34812,32.17464 2.44985,10.20769 3.91976,16.74061 3.91976,16.74061 0.16332,4.89969 -5.71631,8.41114 -5.71631,8.41114 l -58.71463,43.52559 c -9.5544,7.4312 -11.10597,19.5171 -11.10597,19.5171 l -2.44985,12.65754 86.15291,-31.11304 c 1.63323,7.02289 4.81803,14.29076 4.81803,14.29076 0.24499,12.24923 1.30658,18.0472 1.30658,18.0472 1.0616,-5.79797 1.63323,-18.0472 1.63323,-18.0472 2.53151,-6.04295 4.73637,-14.45409 4.73637,-14.45409 l 86.07125,31.43969 -2.04154,-11.43261 c -2.93981,-15.43403 -11.59594,-21.06868 -11.59594,-21.06868 l -58.79629,-43.44393 c -4.81803,-3.02147 -5.55299,-8.24781 -5.55299,-8.24781 0.89828,-4.73637 3.83809,-16.74061 3.83809,-16.74061 3.1848,-16.49563 3.59311,-37.5643 3.59311,-37.5643 v -42.30067 c 3.26646,-15.10739 4.00142,-35.03279 4.00142,-35.03279 l 29.47981,8.00282 c 1.22491,7.18622 2.20486,7.0229 2.20486,7.0229 1.0616,0 1.55157,-6.20628 1.55157,-6.20628 l 10.86098,3.02148 10.04437,4.00141 c 0.81662,7.92117 2.04153,7.75785 2.04153,7.75785 1.46991,-0.48997 1.55157,-6.3696 1.55157,-6.3696 l 20.25206,8.24781 c 1.71489,16.65895 3.10314,16.41397 3.10314,16.41397 1.30658,-0.24499 1.55157,-14.53575 1.55157,-14.53575 l 20.98701,8.49279 c 1.79655,15.02573 3.26646,15.67902 3.26646,15.67902 1.22492,-1.87822 1.38825,-13.96412 1.38825,-13.96412 l 21.80362,9.14609 c 0.73496,12.49421 2.20486,12.33089 2.20486,12.33089 0.89828,-1.30659 1.71489,-10.94265 1.71489,-10.94265 l 80.51827,35.35944 1.55156,12.24923 -0.57163,-25.07009 c -0.81661,-12.82085 -8.81944,-17.80221 -8.81944,-17.80221 l -71.78048,-52.835 1.46991,-5.55299 h 2.69483 c 2.20486,-5.96129 2.1232,-18.94547 2.1232,-18.94547 -0.0817,-14.53576 -2.1232,-15.59735 -2.1232,-15.59735 h -19.5171 c -2.53151,6.94122 -2.04154,15.43403 -2.04154,15.43403 -0.0817,5.96128 0.81661,12.41255 0.81661,12.41255 l -48.99691,-36.17606 0.89828,-3.91975 h 2.69483 c 2.04154,-5.47132 2.1232,-18.4555 2.1232,-18.4555 0.4083,-12.7392 -2.20487,-16.08732 -2.20487,-16.08732 h -19.5171 c -2.53151,7.67618 -1.95988,16.08732 -1.95988,16.08732 0,5.14467 0.48997,8.49279 0.48997,8.49279 l -43.85223,-36.01273 c -7.10455,-4.73637 -12.00425,-14.78073 -12.00425,-14.78073 V 76.68017 C 262.29681,-4.246399 244.73958,0 244.73958,0 Z"
	AIRBUS_A380_SCALE   = 0.08

	// Fokker F100
	FOKKER_F100_SVGPATH = "M 46.861745,1.2786934 C 43.475576,1.1366164 41.533857,13.402601 41.533857,13.402601 v 36.339066 l -12.694274,6.085416 -25.9291663,9.525 c -1.5875,0.79375 -1.5875,1.5875 -1.5875,1.5875 L 0.26458333,71.172917 19.05,69.585417 l 0.264583,2.645833 0.529167,-2.645833 8.995833,-0.529167 0.79375,3.704167 0.79375,-3.96875 10.054167,-0.79375 c 0,2.116666 1.058333,4.233333 1.058333,4.233333 V 79.375 h -0.79375 v -1.322917 c -0.79375,-1.322916 -3.175,-1.322916 -3.175,-1.322916 -2.116666,0 -2.116666,0.79375 -2.116666,0.79375 v 6.085416 c 0,5.027084 0.79375,9.789584 0.79375,9.789584 0.79375,0.529166 2.116666,0.529166 2.116666,0.529166 1.322917,0 1.5875,-0.529166 1.5875,-0.529166 l 2.645834,2.38125 c 1.5875,8.731253 3.439583,13.493753 3.439583,13.493753 l -14.552083,8.99583 c -1.5875,1.32292 -1.5875,3.43958 -1.5875,3.43958 v 1.85209 L 46.0375,119.85625 c 0,1.85208 0.79375,2.91042 0.79375,2.91042 0.529167,-0.52917 1.058333,-2.91042 1.058333,-2.91042 l 15.875,3.70417 v -2.11667 c 0,-2.11667 -1.5875,-3.175 -1.5875,-3.175 L 47.625,109.27292 C 49.741667,104.775 51.064583,95.779167 51.064583,95.779167 L 53.975,93.397917 c 0.529167,0.529166 1.5875,0.529166 1.5875,0.529166 1.322917,0 1.852083,-0.529166 1.852083,-0.529166 0.529167,-3.439584 0.79375,-10.054167 0.79375,-10.054167 0.264584,-2.645833 0,-6.35 0,-6.35 -1.058333,-0.264583 -2.116666,-0.264583 -2.116666,-0.264583 -1.852084,0 -2.910417,0.79375 -2.910417,0.79375 C 52.916667,78.052083 52.916667,79.375 52.916667,79.375 H 52.3875 v -8.73125 c 0.529167,-0.529167 0.529167,-2.645833 0.529167,-2.645833 L 63.5,68.791667 c 0,1.5875 0.79375,3.96875 0.79375,3.96875 0.264583,-0.79375 0.529167,-3.704167 0.529167,-3.704167 l 8.995833,0.529167 c 0,0.79375 0.529167,2.645833 0.529167,2.645833 C 74.6125,71.172917 74.6125,69.585417 74.6125,69.585417 l 18.785417,1.5875 L 92.86875,68.2625 c -0.529167,-2.645833 -2.116667,-2.910417 -2.116667,-2.910417 l -25.664583,-9.525 -12.7,-6.35 V 13.49375 C 50.270833,0.79375 46.861745,1.2786934 46.861745,1.2786934 Z"
	FOKKER_F100_SCALE   = 0.28

	PILATUS_PC12_SVGPATH = "m 79.106876,0.66218011 c -1.359493,-0.00697 -2.160254,4.17289299 -2.289625,4.39264259 -1.967034,-1.0118725 -4.15542,-0.9005986 -4.15542,-0.9005986 -3.486582,0 -4.796659,1.2825704 -4.796659,1.2825704 l -0.0978,0.9550896 7.305392,-0.1976392 1.530753,0.2014878 -0.220407,0.9077521 c -0.399426,0.3516962 -1.034997,2.7690542 -1.034997,2.7690542 -2.069779,-0.2315048 -3.62981,2.84763 -3.62981,2.84763 0.672464,0.720252 0.795355,1.382282 0.795355,1.382282 0.172968,0.630384 0.122412,1.171252 0.122412,1.171252 0.797976,-1.770997 1.960666,-2.728388 1.960666,-2.728388 -2.616132,10.685321 -2.750904,13.463044 -2.750904,13.463044 L 70.74069,34.949962 70.70507,48.09274 6.8994795,50.822026 5.8881872,51.363393 4.5792198,53.366579 1.1044002,60.472344 0.87835887,62.227938 1.1879514,64.40816 l 2.1898163,-1.61448 2.0466784,-0.766509 1.105854,0.117381 17.3707949,1.202297 7.186496,0.709679 0.286276,1.225773 0.325785,-1.113523 16.996822,1.912617 0.382133,1.132766 0.285628,-1.130842 19.888144,2.381823 0.895033,0.686523 0.640947,2.156937 -0.04767,13.69537 c 1.310302,11.427521 2.365858,15.702858 2.365858,15.702858 l 1.03707,10.53356 1.002613,1.51057 2.556729,10.05558 -21.422117,2.68951 -1.357673,2.48773 -1.228396,3.56571 0.241586,2.55225 23.989079,1.26092 c 0.316795,2.70452 1.114144,4.22562 1.114144,4.22562 0.942653,-1.95158 0.996137,-4.16289 0.996137,-4.16289 l 24.097051,-1.29569 0.24132,-2.45077 -1.22898,-3.69515 -1.35747,-2.48798 -21.419012,-2.68926 c 1.06656,-4.22428 2.451545,-10.04006 2.451545,-10.04006 l 1.093742,-1.53314 1.002225,-10.52651 c 1.462837,-7.701803 2.329457,-16.103111 2.329457,-16.103111 l 0,-13.574203 0.491462,-1.975798 1.061812,-0.588384 19.863726,-2.335191 0.42812,1.084018 0.19106,-1.130842 17.07319,-1.895042 0.29729,1.329685 0.50519,-1.370095 6.86926,-0.845982 17.4947,-1.157076 1.78061,-0.07049 1.95211,0.84925 2.31275,1.508195 0.0622,-3.988413 -3.55321,-7.245596 c -1.07923,-2.001236 -2.53814,-2.388237 -3.39632,-2.404273 L 88.637338,48.278242 87.270403,47.974075 87.102135,35.021354 86.104185,25.485467 C 84.887101,17.89936 83.381326,12.926262 83.381326,12.926262 c 1.322167,1.131917 2.194868,2.546029 2.194868,2.546029 -0.521042,-1.81351 0.700145,-2.509275 0.700145,-2.509275 C 84.648758,9.8327359 82.743358,10.169395 82.743358,10.169395 82.167513,7.2904327 81.510364,7.0605751 81.510364,7.0605751 l 0.261081,-0.6011484 c 0.626602,0.3506788 0.950892,-0.1146079 0.950892,-0.1146079 1.536857,-0.7331791 8.077226,0.047065 8.077226,0.047065 L 90.701113,5.4367945 C 90.096427,4.2298792 86.063942,4.1452816 86.063942,4.1452816 83.526015,4.1170826 81.115277,5.1529616 81.115277,5.1529616 80.367705,0.32429073 79.106876,0.66218011 79.106876,0.66218011 Z"
	PILATUS_PC12_SCALE   = 0.2
)

func InitMarker(svgpath string, scale float32) (vs []ebiten.Vertex, is []uint16, maxX, maxY int, err error) {
	// converts SVG path to image

	// Prepare the path object
	path := vector.Path{}

	// SVG path to ebiten vector.Path
	maxX, maxY, err = PathFromSVG(&path, scale, svgpath)
	if err != nil {
		return vs, is, maxX, maxY, err
	}

	// Get a list of verticies and indicies
	vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)

	// Set colours (not sure what SrcX/SrcY are doing... Copied from the Ebiten example: https://ebiten.org/examples/vector.html)
	for i := range vs {
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = 0xff / float32(0xff)
		vs[i].ColorG = 0xff / float32(0xff)
		vs[i].ColorB = 0xff / float32(0xff)
		vs[i].ColorA = 1
	}

	return vs, is, maxX, maxY, nil
}
