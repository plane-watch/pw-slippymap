package markers

import (
	"log"
	"pw_slippymap/slippymap"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type aircraftMarker struct {
	model   string
	svgPath string
	scale   float64
	iata    string
}

// map key (and the order of the map) is ICAO
var Aircraft = map[string]aircraftMarker{
	"A388": {
		model:   "Airbus A380-800",
		svgPath: "m 244.73958,0 c -19.45177,2.9398148 -21.49332,76.729166 -21.49332,76.729166 v 35.718754 c -2.84181,7.02289 -10.27301,13.22916 -10.27301,13.22916 l -45.64879,37.48264 c 0.57163,-5.30799 0.32665,-9.71772 0.32665,-9.71772 0,-10.45268 -2.1232,-14.8624 -2.1232,-14.8624 h -19.35378 c -2.28653,2.36819 -2.04154,16.16899 -2.04154,16.16899 0.24498,13.22916 1.95987,18.4555 1.95987,18.4555 h 2.7765 l 1.03464,3.86133 -49.2966,36.3978 c 0.97994,-4.57305 0.89827,-11.10597 0.89827,-11.10597 -0.0817,-15.18904 -2.123197,-16.90393 -2.123197,-16.90393 H 79.946631 c -1.796554,1.0616 -2.1232,15.59735 -2.1232,15.59735 0.244984,13.55582 2.204861,19.02713 2.204861,19.02713 h 2.69483 l 1.306585,5.38966 -71.698817,52.99833 c -9.1460906,7.0229 -9.0644291,20.66037 -9.0644291,20.66037 l -0.3266461,22.13027 1.714892,-12.41255 80.6815842,-35.03279 c 0.408308,11.10596 1.388246,10.77932 1.388246,10.77932 1.388246,-0.48997 2.368184,-12.65754 2.368184,-12.65754 l 21.721969,-8.65612 c 0.0817,13.63747 1.22492,13.55581 1.22492,13.55581 2.04154,-0.89827 3.42978,-15.51569 3.42978,-15.51569 l 20.98701,-8.32947 c 0.32665,14.61741 1.55157,14.45409 1.55157,14.45409 2.85816,-5.55298 2.93982,-16.49563 2.93982,-16.49563 l 20.33372,-8.08449 c 0.24498,6.12461 1.38824,6.12461 1.38824,6.12461 1.30659,-0.0817 2.20486,-7.67618 2.20486,-7.67618 l 9.96271,-3.91975 10.94264,-2.93982 c 0.40831,6.28794 1.46991,6.12462 1.46991,6.12462 1.38825,0.0817 2.04154,-7.18622 2.04154,-7.18622 l 29.72479,-7.8395 c 0.73495,21.39532 4.16474,35.35943 4.16474,35.35943 v 47.28203 c -0.0817,8.32947 3.34812,32.17464 3.34812,32.17464 2.44985,10.20769 3.91976,16.74061 3.91976,16.74061 0.16332,4.89969 -5.71631,8.41114 -5.71631,8.41114 l -58.71463,43.52559 c -9.5544,7.4312 -11.10597,19.5171 -11.10597,19.5171 l -2.44985,12.65754 86.15291,-31.11304 c 1.63323,7.02289 4.81803,14.29076 4.81803,14.29076 0.24499,12.24923 1.30658,18.0472 1.30658,18.0472 1.0616,-5.79797 1.63323,-18.0472 1.63323,-18.0472 2.53151,-6.04295 4.73637,-14.45409 4.73637,-14.45409 l 86.07125,31.43969 -2.04154,-11.43261 c -2.93981,-15.43403 -11.59594,-21.06868 -11.59594,-21.06868 l -58.79629,-43.44393 c -4.81803,-3.02147 -5.55299,-8.24781 -5.55299,-8.24781 0.89828,-4.73637 3.83809,-16.74061 3.83809,-16.74061 3.1848,-16.49563 3.59311,-37.5643 3.59311,-37.5643 v -42.30067 c 3.26646,-15.10739 4.00142,-35.03279 4.00142,-35.03279 l 29.47981,8.00282 c 1.22491,7.18622 2.20486,7.0229 2.20486,7.0229 1.0616,0 1.55157,-6.20628 1.55157,-6.20628 l 10.86098,3.02148 10.04437,4.00141 c 0.81662,7.92117 2.04153,7.75785 2.04153,7.75785 1.46991,-0.48997 1.55157,-6.3696 1.55157,-6.3696 l 20.25206,8.24781 c 1.71489,16.65895 3.10314,16.41397 3.10314,16.41397 1.30658,-0.24499 1.55157,-14.53575 1.55157,-14.53575 l 20.98701,8.49279 c 1.79655,15.02573 3.26646,15.67902 3.26646,15.67902 1.22492,-1.87822 1.38825,-13.96412 1.38825,-13.96412 l 21.80362,9.14609 c 0.73496,12.49421 2.20486,12.33089 2.20486,12.33089 0.89828,-1.30659 1.71489,-10.94265 1.71489,-10.94265 l 80.51827,35.35944 1.55156,12.24923 -0.57163,-25.07009 c -0.81661,-12.82085 -8.81944,-17.80221 -8.81944,-17.80221 l -71.78048,-52.835 1.46991,-5.55299 h 2.69483 c 2.20486,-5.96129 2.1232,-18.94547 2.1232,-18.94547 -0.0817,-14.53576 -2.1232,-15.59735 -2.1232,-15.59735 h -19.5171 c -2.53151,6.94122 -2.04154,15.43403 -2.04154,15.43403 -0.0817,5.96128 0.81661,12.41255 0.81661,12.41255 l -48.99691,-36.17606 0.89828,-3.91975 h 2.69483 c 2.04154,-5.47132 2.1232,-18.4555 2.1232,-18.4555 0.4083,-12.7392 -2.20487,-16.08732 -2.20487,-16.08732 h -19.5171 c -2.53151,7.67618 -1.95988,16.08732 -1.95988,16.08732 0,5.14467 0.48997,8.49279 0.48997,8.49279 l -43.85223,-36.01273 c -7.10455,-4.73637 -12.00425,-14.78073 -12.00425,-14.78073 V 76.68017 C 262.29681,-4.246399 244.73958,0 244.73958,0 Z",
		scale:   0.07,
		iata:    "388",
	},
	"DH8D": {
		model:   "De Havilland Canada DHC-8-400 Dash 8Q",
		svgPath: "m 47.148749,2.06375 c -1.641604,0 -4.28625,7.1889238 -4.28625,11.985625 V 46.381458 L 42.2275,46.646042 c -3.201567,0 -7.593542,0.423333 -7.593542,0.423333 v -7.328958 c 0,-1.481667 -0.661458,-3.677709 -1.296458,-3.677709 -0.555625,0 -1.402292,2.143125 -1.402292,3.704167 v 7.46125 L 4.0216667,49.424167 c -0.873125,0.05292 -0.9789584,0.3175 -1.0054167,0.608541 l -0.2645833,2.936875 c -0.026458,0.423334 0,0.47625 0.238125,0.502709 l 12.9910413,0.952499 0.238125,0.873126 0.3175,-0.820208 8.016875,0.502708 0.211667,1.296458 0.47625,-1.243542 6.6675,0.3175 v 4.180417 c 0,0.926042 0.714375,3.889375 1.243542,3.889375 0.532141,0 1.613958,-2.883917 1.613958,-3.862917 v -3.254375 h 7.752292 c 0.291041,0.02646 0.343958,0.291042 0.343958,0.502709 v 18.626666 c 0,3.915834 2.910417,24.659162 2.910417,24.659162 -2.116667,0.10585 -11.059584,1.48167 -11.059584,1.48167 v 4.68313 h 12.250209 c 0.07938,1.34937 0.370416,2.7252 0.370416,2.7252 0,0 0.370417,-1.48166 0.47625,-2.7252 h 12.012084 v -4.70959 c 0,0 -8.916459,-1.29645 -10.95375,-1.561038 0,0 2.619375,-20.611042 2.619375,-24.579792 V 56.885417 c 0.02646,-0.343959 0.264583,-0.555625 0.555625,-0.555625 H 60.0075 v 3.095625 c 0,1.190625 0.687917,3.96875 1.190625,3.96875 0.555625,0 1.613958,-2.751667 1.613958,-3.836459 v -3.915833 l 6.6675,-0.343958 0.47625,0.978958 0.291042,-1.058333 8.202083,-0.555625 0.264584,0.714375 0.264583,-0.79375 12.85875,-0.978958 c 0.343958,-0.05292 0.449792,-0.291042 0.449792,-0.635001 l -0.07938,-2.301875 C 92.154377,49.900416 91.678123,49.503542 91.09604,49.424166 L 62.785625,47.333958 v -7.567083 c 0,-1.561042 -0.846667,-3.730625 -1.42875,-3.730625 -0.555625,0 -1.481667,2.169583 -1.481667,3.730625 V 46.99 c -2.804583,-0.185208 -8.043333,-0.370417 -8.043333,-0.370417 l -0.529167,-0.47625 V 14.102292 c 0,-4.8418753 -2.513544,-12.038542 -4.153959,-12.038542 z",
		scale:   0.3,
		iata:    "DH4",
	},
	"E190": {
		model:   "Embraer 190 / Lineage 1000",
		svgPath: "m 57.784999,1.5875 c -5.667507,2.9104167 -5.767916,14.975416 -5.767916,14.975416 v 36.565417 l -2.434166,4.1275 L 43.18,60.589583 C 44.185417,59.107917 44.238333,54.2925 44.238333,54.2925 44.502917,50.64125 43.65625,49.53 43.65625,49.53 h -7.249583 c -0.9525,0.899584 -0.899584,5.185833 -0.899584,5.185833 -0.05292,4.497917 1.11125,5.715 1.11125,5.715 l 0.635,0.05292 L 37.7825,63.1825 4.4979167,80.062917 c -1.1112501,0.740833 -1.27,1.693333 -1.27,1.693333 l -2.06375,7.46125 v 1.42875 h 0.635 l 2.1166666,-3.545417 16.2983337,-4.28625 v 2.010834 L 20.584583,85.725 21.007917,84.719583 V 82.602916 L 31.90875,79.6925 v 1.693333 l 0.370417,0.899584 0.370416,-0.582084 v -2.2225 l 4.868334,-1.322916 H 42.06875 V 80.01 l 0.423333,0.582083 0.423334,-0.370416 V 78.105 h 9.048749 v 29.79208 c 0.250906,10.4753 2.381251,18.46792 2.381251,18.46792 l -14.605,10.21292 c -1.534584,1.11124 -1.42875,3.12208 -1.42875,3.12208 l -0.05292,3.38667 17.727086,-4.70959 c 0.9525,6.0325 1.74625,5.97959 1.74625,5.97959 0.740834,0.0529 1.957917,-6.08542 1.957917,-6.08542 l 17.674167,4.86833 -0.05292,-2.96333 c -0.158747,-2.80459 -0.846664,-3.175 -0.846664,-3.175 l -15.24,-10.68916 c 2.698751,-12.64709 2.487084,-18.30917 2.487084,-18.30917 V 77.999167 h 8.73125 l 0.370416,2.487083 0.740834,-2.487083 H 78.105 l 4.60375,1.534583 0.47625,2.54 0.687916,-2.169584 10.742084,3.069168 0.211667,2.592916 0.687916,-2.275417 16.033747,4.28625 1.74625,2.275417 0.47625,0.846666 0.58209,0 v -2.010834 l -2.01084,-7.196667 c -0.37041,-1.058329 -1.64041,-1.746245 -1.64041,-1.746245 L 77.681667,63.1825 l 0.47625,-2.645833 h 0.635 C 79.957084,60.325 80.221667,54.9275 80.221667,54.9275 80.115833,49.318334 79.163333,49.318333 79.163333,49.318333 H 71.91375 c -1.11125,1.957917 -0.899583,4.92125 -0.899583,4.92125 0.105833,4.233334 1.27,6.244167 1.27,6.244167 l -5.87375,-2.963333 -2.69875,-4.392083 V 16.51 C 63.552917,4.0745833 57.784999,1.5875 57.784999,1.5875 Z",
		scale:   0.22,
		iata:    "E90",
	},
	"F100": {
		model:   "Fokker 100",
		svgPath: "M 46.861745,1.2786934 C 43.475576,1.1366164 41.533857,13.402601 41.533857,13.402601 v 36.339066 l -12.694274,6.085416 -25.9291663,9.525 c -1.5875,0.79375 -1.5875,1.5875 -1.5875,1.5875 L 0.26458333,71.172917 19.05,69.585417 l 0.264583,2.645833 0.529167,-2.645833 8.995833,-0.529167 0.79375,3.704167 0.79375,-3.96875 10.054167,-0.79375 c 0,2.116666 1.058333,4.233333 1.058333,4.233333 V 79.375 h -0.79375 v -1.322917 c -0.79375,-1.322916 -3.175,-1.322916 -3.175,-1.322916 -2.116666,0 -2.116666,0.79375 -2.116666,0.79375 v 6.085416 c 0,5.027084 0.79375,9.789584 0.79375,9.789584 0.79375,0.529166 2.116666,0.529166 2.116666,0.529166 1.322917,0 1.5875,-0.529166 1.5875,-0.529166 l 2.645834,2.38125 c 1.5875,8.731253 3.439583,13.493753 3.439583,13.493753 l -14.552083,8.99583 c -1.5875,1.32292 -1.5875,3.43958 -1.5875,3.43958 v 1.85209 L 46.0375,119.85625 c 0,1.85208 0.79375,2.91042 0.79375,2.91042 0.529167,-0.52917 1.058333,-2.91042 1.058333,-2.91042 l 15.875,3.70417 v -2.11667 c 0,-2.11667 -1.5875,-3.175 -1.5875,-3.175 L 47.625,109.27292 C 49.741667,104.775 51.064583,95.779167 51.064583,95.779167 L 53.975,93.397917 c 0.529167,0.529166 1.5875,0.529166 1.5875,0.529166 1.322917,0 1.852083,-0.529166 1.852083,-0.529166 0.529167,-3.439584 0.79375,-10.054167 0.79375,-10.054167 0.264584,-2.645833 0,-6.35 0,-6.35 -1.058333,-0.264583 -2.116666,-0.264583 -2.116666,-0.264583 -1.852084,0 -2.910417,0.79375 -2.910417,0.79375 C 52.916667,78.052083 52.916667,79.375 52.916667,79.375 H 52.3875 v -8.73125 c 0.529167,-0.529167 0.529167,-2.645833 0.529167,-2.645833 L 63.5,68.791667 c 0,1.5875 0.79375,3.96875 0.79375,3.96875 0.264583,-0.79375 0.529167,-3.704167 0.529167,-3.704167 l 8.995833,0.529167 c 0,0.79375 0.529167,2.645833 0.529167,2.645833 C 74.6125,71.172917 74.6125,69.585417 74.6125,69.585417 l 18.785417,1.5875 L 92.86875,68.2625 c -0.529167,-2.645833 -2.116667,-2.910417 -2.116667,-2.910417 l -25.664583,-9.525 -12.7,-6.35 V 13.49375 C 50.270833,0.79375 46.861745,1.2786934 46.861745,1.2786934 Z",
		scale:   0.28,
		iata:    "100",
	},
	"PC12": {
		model:   "Pilatus PC-12",
		svgPath: "m 79.106876,0.66218011 c -1.359493,-0.00697 -2.160254,4.17289299 -2.289625,4.39264259 -1.967034,-1.0118725 -4.15542,-0.9005986 -4.15542,-0.9005986 -3.486582,0 -4.796659,1.2825704 -4.796659,1.2825704 l -0.0978,0.9550896 7.305392,-0.1976392 1.530753,0.2014878 -0.220407,0.9077521 c -0.399426,0.3516962 -1.034997,2.7690542 -1.034997,2.7690542 -2.069779,-0.2315048 -3.62981,2.84763 -3.62981,2.84763 0.672464,0.720252 0.795355,1.382282 0.795355,1.382282 0.172968,0.630384 0.122412,1.171252 0.122412,1.171252 0.797976,-1.770997 1.960666,-2.728388 1.960666,-2.728388 -2.616132,10.685321 -2.750904,13.463044 -2.750904,13.463044 L 70.74069,34.949962 70.70507,48.09274 6.8994795,50.822026 5.8881872,51.363393 4.5792198,53.366579 1.1044002,60.472344 0.87835887,62.227938 1.1879514,64.40816 l 2.1898163,-1.61448 2.0466784,-0.766509 1.105854,0.117381 17.3707949,1.202297 7.186496,0.709679 0.286276,1.225773 0.325785,-1.113523 16.996822,1.912617 0.382133,1.132766 0.285628,-1.130842 19.888144,2.381823 0.895033,0.686523 0.640947,2.156937 -0.04767,13.69537 c 1.310302,11.427521 2.365858,15.702858 2.365858,15.702858 l 1.03707,10.53356 1.002613,1.51057 2.556729,10.05558 -21.422117,2.68951 -1.357673,2.48773 -1.228396,3.56571 0.241586,2.55225 23.989079,1.26092 c 0.316795,2.70452 1.114144,4.22562 1.114144,4.22562 0.942653,-1.95158 0.996137,-4.16289 0.996137,-4.16289 l 24.097051,-1.29569 0.24132,-2.45077 -1.22898,-3.69515 -1.35747,-2.48798 -21.419012,-2.68926 c 1.06656,-4.22428 2.451545,-10.04006 2.451545,-10.04006 l 1.093742,-1.53314 1.002225,-10.52651 c 1.462837,-7.701803 2.329457,-16.103111 2.329457,-16.103111 l 0,-13.574203 0.491462,-1.975798 1.061812,-0.588384 19.863726,-2.335191 0.42812,1.084018 0.19106,-1.130842 17.07319,-1.895042 0.29729,1.329685 0.50519,-1.370095 6.86926,-0.845982 17.4947,-1.157076 1.78061,-0.07049 1.95211,0.84925 2.31275,1.508195 0.0622,-3.988413 -3.55321,-7.245596 c -1.07923,-2.001236 -2.53814,-2.388237 -3.39632,-2.404273 L 88.637338,48.278242 87.270403,47.974075 87.102135,35.021354 86.104185,25.485467 C 84.887101,17.89936 83.381326,12.926262 83.381326,12.926262 c 1.322167,1.131917 2.194868,2.546029 2.194868,2.546029 -0.521042,-1.81351 0.700145,-2.509275 0.700145,-2.509275 C 84.648758,9.8327359 82.743358,10.169395 82.743358,10.169395 82.167513,7.2904327 81.510364,7.0605751 81.510364,7.0605751 l 0.261081,-0.6011484 c 0.626602,0.3506788 0.950892,-0.1146079 0.950892,-0.1146079 1.536857,-0.7331791 8.077226,0.047065 8.077226,0.047065 L 90.701113,5.4367945 C 90.096427,4.2298792 86.063942,4.1452816 86.063942,4.1452816 83.526015,4.1170826 81.115277,5.1529616 81.115277,5.1529616 80.367705,0.32429073 79.106876,0.66218011 79.106876,0.66218011 Z",
		scale:   0.2,
		iata:    "PL2",
	},
	"SF34": {
		model:   "Saab SF340A/B",
		svgPath: "M 94.720832,2.06375 C 84.084583,6.4558334 84.666666,31.485416 84.666666,31.485416 V 68.58 L 69.267917,70.379167 V 66.83375 64.187917 l 0.79375,-2.06375 -1.058334,-1.005417 c -0.370417,-7.037917 -0.846666,-7.9375 -0.846666,-7.9375 1.957916,0.211666 3.280833,0.158754 3.280833,0.158754 7.9375,-0.158755 7.884584,-0.846671 7.884584,-0.846671 C 77.893333,51.59375 71.4375,51.964167 71.4375,51.964167 L 67.839167,52.07 C 66.622083,48.365834 65.7225,48.365833 65.7225,48.365833 c -1.27,0.3175 -2.275417,3.757084 -2.275417,3.757084 L 60.0075,52.07 c -7.9375,0.05292 -7.884583,0.635 -7.884583,0.635 0.105833,0.846667 7.831666,0.687917 7.831666,0.687917 1.27,0.211666 3.227917,-0.05292 3.227917,-0.05292 -1.5875,5.60917 -1.27,13.440836 -1.27,13.440836 V 71.27875 L 4.7095833,78.052083 C 2.2225,78.84583 2.2754167,82.126667 2.2754167,82.126667 L 1.905,82.920417 2.2754167,83.449583 v 3.968751 l 60.8541663,3.598333 c 0.05292,2.38125 2.38125,2.328333 2.38125,2.328333 2.116667,0.211666 2.38125,-2.116667 2.38125,-2.116667 l 16.615834,1.005417 v 40.11083 c 0,8.6785 1.481666,13.07042 1.481666,13.07042 l -26.564166,4.445 c -2.328333,0.58209 -2.434167,2.8575 -2.434167,2.8575 l 0.05292,5.18583 c 0,0.635 1.05833,0.68792 1.05833,0.68792 l 30.215417,4.1275 2.487083,2.16958 c 1.5875,8.36083 3.386667,8.36084 3.386667,8.36084 l 0.211666,1.16416 0.211667,-1.11125 c 2.010833,-0.15875 3.4925,-8.57249 3.4925,-8.57249 l 2.43416,-2.27542 30.63876,-4.445 c 0.89958,-0.21168 0.84666,-0.635 0.84666,-0.635 v -4.18042 c -0.10583,-3.28084 -2.27541,-3.65125 -2.27541,-3.65125 l -26.77584,-4.33917 c 1.64042,-8.14917 1.48167,-12.91166 1.48167,-12.91166 V 91.969167 l 16.29833,-0.9525 c 0.26457,2.2225 2.43417,2.116666 2.43417,2.116666 2.38125,0.05292 2.43417,-2.487083 2.43417,-2.487083 l 61.11875,-3.704167 v -3.65125 l 0.68791,-0.687916 -0.635,-0.740834 C 186.47833,77.522916 184.52042,77.47 184.52042,77.47 l -57.41459,-6.561667 v -4.28625 -2.963333 c 0.42334,-0.47625 0.42334,-1.164167 0.42334,-1.164167 0.0529,-1.375833 -0.68792,-1.74625 -0.68792,-1.74625 0,-4.339166 -1.00542,-7.725833 -1.00542,-7.725833 1.32292,0.370417 3.12209,0.05292 3.12209,0.05292 8.20208,0.634997 7.99041,-0.740837 7.99041,-0.740837 -0.0529,-0.793749 -8.04333,-0.423333 -8.04333,-0.423333 -1.32292,-0.05292 -3.38667,0.05292 -3.38667,0.05292 -1.27,-3.810003 -2.27541,-3.810003 -2.27541,-3.810003 -0.95251,0.05292 -2.16959,3.862916 -2.16959,3.862916 -11.64167,-0.423333 -11.27125,0.582084 -11.27125,0.582084 0,0.740833 7.72584,0.740833 7.72584,0.740833 1.69333,0.105833 3.175,-0.15875 3.175,-0.15875 C 119.22125,58.89625 119.38,66.675 119.38,66.675 v 3.439583 L 104.61625,68.42125 V 31.538333 c 0,-26.0905456 -9.895418,-29.474583 -9.895418,-29.474583 z",
		scale:   0.17,
		iata:    "SF3",
	},
}

// for lookups from IATA to ICAO if needed, likely faster than iterating over Aircraft to find? maybe?
// map key-type is IATA
var IATAtoICAO = map[string]string{
	"388": "A388",
	"100": "F100",
	"E90": "E190",
	"PL2": "PC12",
	"SF3": "SF34",
}

// Marker object
type Marker struct {
	Img     *ebiten.Image
	CentreX float64
	CentreY float64
	icao    string
}

// Colour object
type RGBA struct {
	r, g, b, a float64
}

func InitMarkers() (imgs map[string]Marker, err error) {
	// Initialise all markers.
	// Renders all SVGs to images.
	// Produces a map of marker images.

	var strokeColour, fillColour, bgColour RGBA
	imgs = make(map[string]Marker)

	// Set default colours
	bgColour = RGBA{ // temp background colour.
		r: 0,
		g: 0.5,
		b: 0,
		a: 0.3,
	}
	strokeColour = RGBA{ // black
		r: 0,
		g: 0,
		b: 0,
		a: 1,
	}
	fillColour = RGBA{ // white
		r: 1,
		g: 1,
		b: 1,
		a: 1,
	}

	var wg sync.WaitGroup
	c := make(chan Marker, len(Aircraft))

	// Pre-render aircraft concurrently
	for k, v := range Aircraft {

		wg.Add(1)

		go func(k string, v aircraftMarker) {
			defer wg.Done()
			log.Printf("Pre-rendering sprites: %s (ICAO: %s)", v.model, k)

			r := renderSVG{
				scale:        v.scale,
				d:            v.svgPath,
				pathStroked:  true,
				pathFilled:   true,
				bgFilled:     false,
				strokeWidth:  2,
				strokeColour: strokeColour,
				fillColour:   fillColour,
				bgColour:     bgColour,
				offsetX:      1,
				offsetY:      1,
			}

			img, err := imgFromSVG(r)
			if err != nil {
				log.Fatal(err)
			}
			c <- Marker{
				Img:     img,
				CentreX: float64(img.Bounds().Dx()) / 2,
				CentreY: float64(img.Bounds().Dy()) / 2,
				icao:    k,
			}
		}(k, v)
	}

	wg.Wait()
	log.Println("Pre-rendering finished, building Marker map")

	// Read markers out of channel, into object to be returned
	for elem := range c {
		imgs[elem.icao] = elem
		if len(c) == 0 {
			break
		}
	}

	return imgs, nil
}

func (m *Marker) MarkerDrawOpts(angleDegrees, xPos, yPos float64) (drawOpts ebiten.DrawImageOptions) {
	drawOpts.GeoM.Translate(-m.CentreX, -m.CentreY)
	drawOpts.GeoM.Rotate(slippymap.DegreesToRadians(angleDegrees))
	drawOpts.GeoM.Translate(m.CentreX, m.CentreY)
	drawOpts.GeoM.Translate(xPos, yPos)
	return drawOpts
}
