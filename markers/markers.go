package markers

import (
	"log"
	"pw_slippymap/slippymap"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iwpnd/piper"
)

type marker struct {
	name    string
	svgPath string
	scale   float64
}

// Missing aircraft markers
var missingMarkers []string

// map key (and the order of the map) is ICAO

// ref: https://www.icao.int/publications/doc8643/pages/search.aspx

var Aircraft = map[string]marker{
	"A320": {
		name:    "AIRBUS A-320",
		svgPath: "m 249.61134,1.38261 c -6.03825,0.2326163 -10.63624,8.2208763 -13.38791,13.72421 l -4.24866,9.080895 c -7.30251,16.404167 -7.8955,23.075004 -8.78417,28.204584 -1.39052,7.877393 -2.77611,19.501522 -2.56646,29.51427 v 85.195841 c -1.12448,4.63021 -1.25677,11.8401 -1.25677,11.8401 -0.26458,4.8948 0.0662,9.85573 -4.29948,12.10469 l -33.73437,17.39636 c 0.79375,-5.15938 1.32291,-7.54063 1.32291,-19.44688 h 1.7198 c -1e-5,-3.9026 0.33072,-6.15156 -2.11667,-8.73125 -0.19844,-4.89479 -0.33073,-10.12031 -2.51354,-12.63385 -2.3151,-1.12448 -23.21719,-1.38906 -25.59844,0 -2.3151,1.85208 -2.51354,8.06979 -2.84427,11.50937 -0.66146,6.81302 -1.05834,13.82448 -0.39688,19.97604 0.39688,4.16719 1.71979,13.82448 2.51354,16.7349 0.21165,0.78989 0.39688,0.97895 0.74084,1.03187 l 3.09562,0.18522 0.52917,3.38666 -145.732508,75.03583 C 1.6814209,301.7621 1.8401709,316.26126 1.8401709,322.42606 l -0.37041,5.23875 c -0.0529,0.635 0.55562,0.68791 0.60854,0.0529 l 0.15875,-5.31811 70.6172911,-20.58459 v 0.7673 l 0.635,-0.18521 c 0.74083,11.72104 1.87854,11.72104 1.87854,11.72104 0,0 1.27,0.0794 1.82563,-12.80583 l 0.8202,-0.26459 v -0.66146 l 48.630428,-14.23458 c 0,1.16417 1.08479,9.86896 2.03729,9.92188 0.68792,0.0529 2.38125,-8.30792 2.06375,-11.00667 l 26.19375,-7.72583 7.62,-0.10584 c -0.0529,1.16417 1.48167,6.50875 1.95792,6.56167 0.635,0 1.905,-5.13292 1.905,-6.45583 h 8.04333 c 0.10585,3.28083 1.27,9.8425 2.01084,9.8425 0.74083,0.0529 2.16958,-7.14376 2.06375,-9.78959 h 39.89916 v 124.83042 c -0.15931,7.1466 3.28523,37.9872 9.13293,63.62097 0.40054,2.32306 -1.60211,7.28959 -5.68749,10.3336 l -61.12048,37.40926 c -3.68485,2.32306 -4.00527,3.68485 -4.80632,9.37234 v 12.73677 l 79.94526,-17.70331 c 2.8838,17.38289 7.12939,30.35998 8.6514,30.35998 h 5.60738 c 1.68222,0 6.32833,-13.45772 8.89171,-30.44008 l 79.94526,17.78341 v -11.21476 c -0.56074,-6.24823 -0.96126,-8.49119 -5.04664,-10.97445 l -60.88016,-37.16895 c -3.84507,-2.56337 -6.16813,-8.49118 -5.7676,-10.2535 5.8477,-26.11438 9.05192,-58.95763 9.05192,-65.12575 V 277.27706 h 39.89253 c 0,3.52464 1.28169,9.77287 2.00264,9.77287 1.04136,0 2.08274,-6.16813 2.08274,-9.77287 h 7.93044 c 0.24032,2.00264 1.4419,6.48854 2.00264,6.48854 0.80105,0 1.77095,-4.36631 2.09137,-6.44905 h 7.60139 l 25.95418,7.57053 c 0,3.68485 1.522,11.13466 2.16284,11.13466 0.80106,0 2.16285,-6.88907 2.16285,-9.93308 l 48.54392,14.17867 v 0.7459 l 0.73863,0.19324 v 0 c 0.0529,3.59834 1.00541,12.80584 1.85208,12.80584 0.79375,0 2.01083,-7.88459 2.01083,-11.64167 l 0.52917,0.15875 v -0.79375 l 70.53792,20.53167 0.21166,5.45041 c -0.008,0.52568 0.85192,0.52569 0.85192,-0.0342 l -0.24954,-5.61429 v -4.77672 c 0,-5.19209 -2.7345,-17.27235 -10.4188,-21.8414 l -145.65204,-75.22901 0.50516,-3.10135 c 3.08787,-0.27147 3.59686,0.10179 3.83439,-1.05192 2.51102,-12.7926 2.85847,-19.91959 2.74854,-25.21194 -0.23753,-11.4353 -1.26009,-23.47163 -4.80468,-23.82071 -9.60293,-0.74652 -12.46497,-0.67865 -23.91568,0.47408 -1.41084,0.13378 -2.80403,5.13357 -3.04156,13.20955 -1.1965,1.17101 -1.96447,2.0437 -1.99938,3.54473 v 4.67763 h 1.57084 c 0,7.62125 0.25964,10.66087 1.29159,19.34796 l -34.24437,-17.58048 c -2.0325,-1.14018 -3.22225,-3.86671 -3.27183,-6.59323 l -0.49573,-10.26166 c -0.1983,-1.18976 -0.64445,-6.69239 -1.04104,-7.18812 V 81.67482 c 0,-10.535709 -1.06201,-21.546001 -2.58067,-29.329138 -0.92408,-4.843666 -1.42927,-11.821506 -9.3018,-29.139305 l -3.70174,-7.97297 C 259.68334,9.5082523 255.84901,1.54136 249.61134,1.38261 Z",
		scale:   0.06,
	},
	"A388": {
		name:    "AIRBUS A-380-800",
		svgPath: "m 244.73958,0 c -19.45177,2.9398148 -21.49332,76.729166 -21.49332,76.729166 v 35.718754 c -2.84181,7.02289 -10.27301,13.22916 -10.27301,13.22916 l -45.64879,37.48264 c 0.57163,-5.30799 0.32665,-9.71772 0.32665,-9.71772 0,-10.45268 -2.1232,-14.8624 -2.1232,-14.8624 h -19.35378 c -2.28653,2.36819 -2.04154,16.16899 -2.04154,16.16899 0.24498,13.22916 1.95987,18.4555 1.95987,18.4555 h 2.7765 l 1.03464,3.86133 -49.2966,36.3978 c 0.97994,-4.57305 0.89827,-11.10597 0.89827,-11.10597 -0.0817,-15.18904 -2.123197,-16.90393 -2.123197,-16.90393 H 79.946631 c -1.796554,1.0616 -2.1232,15.59735 -2.1232,15.59735 0.244984,13.55582 2.204861,19.02713 2.204861,19.02713 h 2.69483 l 1.306585,5.38966 -71.698817,52.99833 c -9.1460906,7.0229 -9.0644291,20.66037 -9.0644291,20.66037 l -0.3266461,22.13027 1.714892,-12.41255 80.6815842,-35.03279 c 0.408308,11.10596 1.388246,10.77932 1.388246,10.77932 1.388246,-0.48997 2.368184,-12.65754 2.368184,-12.65754 l 21.721969,-8.65612 c 0.0817,13.63747 1.22492,13.55581 1.22492,13.55581 2.04154,-0.89827 3.42978,-15.51569 3.42978,-15.51569 l 20.98701,-8.32947 c 0.32665,14.61741 1.55157,14.45409 1.55157,14.45409 2.85816,-5.55298 2.93982,-16.49563 2.93982,-16.49563 l 20.33372,-8.08449 c 0.24498,6.12461 1.38824,6.12461 1.38824,6.12461 1.30659,-0.0817 2.20486,-7.67618 2.20486,-7.67618 l 9.96271,-3.91975 10.94264,-2.93982 c 0.40831,6.28794 1.46991,6.12462 1.46991,6.12462 1.38825,0.0817 2.04154,-7.18622 2.04154,-7.18622 l 29.72479,-7.8395 c 0.73495,21.39532 4.16474,35.35943 4.16474,35.35943 v 47.28203 c -0.0817,8.32947 3.34812,32.17464 3.34812,32.17464 2.44985,10.20769 3.91976,16.74061 3.91976,16.74061 0.16332,4.89969 -5.71631,8.41114 -5.71631,8.41114 l -58.71463,43.52559 c -9.5544,7.4312 -11.10597,19.5171 -11.10597,19.5171 l -2.44985,12.65754 86.15291,-31.11304 c 1.63323,7.02289 4.81803,14.29076 4.81803,14.29076 0.24499,12.24923 1.30658,18.0472 1.30658,18.0472 1.0616,-5.79797 1.63323,-18.0472 1.63323,-18.0472 2.53151,-6.04295 4.73637,-14.45409 4.73637,-14.45409 l 86.07125,31.43969 -2.04154,-11.43261 c -2.93981,-15.43403 -11.59594,-21.06868 -11.59594,-21.06868 l -58.79629,-43.44393 c -4.81803,-3.02147 -5.55299,-8.24781 -5.55299,-8.24781 0.89828,-4.73637 3.83809,-16.74061 3.83809,-16.74061 3.1848,-16.49563 3.59311,-37.5643 3.59311,-37.5643 v -42.30067 c 3.26646,-15.10739 4.00142,-35.03279 4.00142,-35.03279 l 29.47981,8.00282 c 1.22491,7.18622 2.20486,7.0229 2.20486,7.0229 1.0616,4e-5 1.55157,-6.20628 1.55157,-6.20628 l 10.86098,3.02148 10.04437,4.00141 c 0.81662,7.92117 2.04153,7.75785 2.04153,7.75785 1.46991,-0.48997 1.55157,-6.3696 1.55157,-6.3696 l 20.25206,8.24781 c 1.71489,16.65895 3.10314,16.41397 3.10314,16.41397 1.30658,-0.24499 1.55157,-14.53575 1.55157,-14.53575 l 20.98701,8.49279 c 1.79655,15.02573 3.26646,15.67902 3.26646,15.67902 1.22492,-1.87822 1.38825,-13.96412 1.38825,-13.96412 l 21.80362,9.14609 c 0.73496,12.49421 2.20486,12.33089 2.20486,12.33089 0.89828,-1.30659 1.71489,-10.94265 1.71489,-10.94265 l 80.51827,35.35944 1.55156,12.24923 -0.57163,-25.07009 c -0.81661,-12.82085 -8.81944,-17.80221 -8.81944,-17.80221 l -71.78048,-52.835 1.46991,-5.55299 h 2.69483 c 2.20486,-5.96129 2.1232,-18.94547 2.1232,-18.94547 -0.0817,-14.53576 -2.1232,-15.59735 -2.1232,-15.59735 h -19.5171 c -2.53151,6.94122 -2.04154,15.43403 -2.04154,15.43403 -0.0817,5.96128 0.81661,12.41255 0.81661,12.41255 l -48.99691,-36.17606 0.89828,-3.91975 h 2.69483 c 2.04154,-5.47132 2.1232,-18.4555 2.1232,-18.4555 0.4083,-12.7392 -2.20487,-16.08732 -2.20487,-16.08732 h -19.5171 c -2.53151,7.67618 -1.95988,16.08732 -1.95988,16.08732 1e-5,5.14467 0.48997,8.49279 0.48997,8.49279 l -43.85223,-36.01273 c -7.10455,-4.73637 -12.00425,-14.78073 -12.00425,-14.78073 V 76.68017 C 262.29681,-4.246399 244.73958,0 244.73958,0 Z",
		scale:   0.07,
	},
	"B738": {
		name:    "BOEING 737-800",
		svgPath: "m 136.89306,1.5932501 c -2.70599,0.02529 -5.33612,6.4235796 -7.05582,12.1137589 -2.63013,7.384588 -7.42104,24.291882 -7.78922,50.200023 v 50.657548 l -14.49994,12.12555 c 0.60448,-5.08776 0.28888,-16.08574 -0.15112,-18.63831 -0.10075,-0.75561 -0.60415,-1.71989 -1.51122,-1.76309 H 93.896778 c -0.957102,0.0504 -1.360093,0.95711 -1.46084,1.76309 -0.503739,2.72018 -1.05785,15.46476 0.554112,25.28765 h 2.317195 c 0,0 0.05037,0.65486 0.352616,1.1586 0,0 -2.166073,1.41046 -4.483268,2.72018 -16.824852,9.01691 -32.289611,15.16252 -87.2978098,41.70952 -1.6623357,0.7556 -2.367569,1.25934 -2.367569,2.51869 v 13.48771 c -0.028082,0.40581 0.3598147,0.49708 0.4054495,0.0407 0.2053569,-2.07636 0.022817,-4.31246 2.6468214,-4.90572 L 63.330929,176.91969 c 0,2.68388 0.80677,7.22882 1.644959,7.4456 0.779192,-0.12985 1.688248,-5.67078 1.558383,-8.22479 l 21.990507,-4.93488 c 0.129865,2.77046 0.736859,7.54759 1.558382,7.48889 0.606038,-0.0433 1.688248,-5.41105 1.601671,-8.09493 l 6.536548,-1.55838 h 2.380861 c -0.12985,3.33321 1.08221,8.09493 1.60167,8.09493 0.51946,0 1.68825,-4.80501 1.55838,-8.09493 h 18.22442 v 59.65142 c 0.34628,16.4063 2.80756,28.20753 7.96507,50.43098 l -47.70382,26.27606 c -0.779191,0.43288 -1.255364,0.69261 -1.255364,0.99563 v 9.17714 L 135.4494,304.44732 c 0.0866,0.90905 0.69262,2.9869 0.86577,3.11676 v 1.73154 h 1.08221 v -1.75516 c 0.17757,-0.16578 0.88581,-2.16656 1.09828,-3.03415 l 54.42798,11.15473 v -8.83358 c -6.2e-4,-0.5324 -0.1127,-0.75657 -0.70114,-1.14885 l -48.27979,-26.61974 c 5.04562,-22.39677 7.71879,-33.51508 7.73373,-49.93301 v -60.08238 h 18.34822 c 0,4.50589 1.19144,8.05912 1.61194,8.0597 0.55605,0.002 1.55591,-3.389 1.51244,-8.04051 h 2.43443 l 6.60775,1.47805 c 0,2.82568 1.16433,8.09634 1.56499,8.0858 0.65208,-0.13042 1.565,-4.47762 1.52152,-7.39025 l 22.04032,4.91234 c -0.13039,3.56471 1.04332,7.95538 1.52152,8.0858 0.6086,-0.30431 1.73888,-3.69513 1.60846,-7.39025 l 59.51319,13.43286 c 1.54241,0.46066 1.91277,1.91276 2.08666,4.82539 -0.013,0.29466 0.23768,0.25287 0.23768,-0.0396 v -13.07944 c 0.0209,-1.25362 -0.29251,-1.964 -1.17004,-2.50724 -40.64321,-19.62594 -76.9079,-35.8886 -88.41943,-42.28128 l -4.439,-2.65439 c 0.27552,-0.38512 0.23897,-0.78722 0.23897,-1.18932 h 2.37604 c 1.74083,-12.72774 1.28711,-21.03332 0.54832,-25.33233 -0.16801,-1.00696 -0.6073,-1.62197 -1.3541,-1.6659 h -12.12449 c -0.96537,-0.004 -1.38751,0.81341 -1.45321,1.67633 -0.71982,8.06809 -0.72587,12.62239 -0.24133,18.48539 L 151.71593,114.24739 V 64.289337 c 0,-27.031413 -5.46776,-43.108965 -7.71242,-50.648717 -2.01444,-6.5037554 -4.37894,-12.0656698 -7.11045,-12.0473699 z",
		scale:   0.11,
	},
	"B77L": {
		name:    "BOEING 777-200ER",
		svgPath: "m 32.17932,0.23475373 c -0.505623,0 -1.038333,0.90289897 -1.282117,1.48075437 L 30.30129,3.0879145 C 29.7066,4.3622486 28.937913,7.3365901 28.937913,13.299702 v 9.503549 c 0,0.751642 -0.603981,1.313881 -0.922091,1.543216 l -4.061443,2.641048 c 0.113593,-0.556289 0.170544,-1.179861 0.170151,-2.448703 0,-1.879065 -0.214539,-2.818597 -0.347701,-2.818597 h -3.706344 c -0.134643,0 -0.384691,1.169405 -0.384691,2.670639 0,1.561763 0.133002,3.003545 0.303314,3.003545 h 0.480863 c 0,0.273722 0.38469,1.398202 0.473465,1.553558 L 1.3020292,41.960851 C 0.42907781,42.574876 0.38469046,42.574876 0.25892627,44.757255 L 11.503724,41.028717 c 0,0.310711 0.09517,0.963842 0.251236,1.00566 0.126625,-0.03393 0.185651,-0.797484 0.185651,-1.142004 l 5.318674,-1.772313 c 0,0.314418 0.07498,0.9878 0.188069,1.018102 0.117846,-0.03158 0.243351,-0.833239 0.243351,-1.150461 l 2.997185,-1.016087 h 2.258209 c 0,0.471864 0.128812,1.144779 0.235932,1.206625 0.121336,-0.06741 0.235788,-0.685232 0.235788,-1.188602 h 5.520959 v 12.985215 c 0,1.308408 0.453022,4.412699 1.09846,6.821497 0.05389,0.201108 -0.05335,0.343389 -0.140855,0.421429 l -9.161414,7.383049 c -0.130066,0.130066 -0.170268,0.841885 -0.170268,1.220261 0,0.44932 0.03784,1.366881 0.170268,1.338503 l 10.490457,-3.68916 c 0.06149,-0.02365 -0.0095,-0.794588 -0.06149,-1.0831 l 0.898646,3.722268 h 0.260133 l 0.898641,-3.726997 c -0.07095,0.297971 -0.04257,1.078369 -0.04257,1.078369 l 10.481,3.726998 c 0.05203,0.03311 0.108783,-0.879723 0.108783,-1.314855 0,-0.510806 -0.01892,-1.154044 -0.08513,-1.22972 l -9.227633,-7.430346 c -0.103,-0.0984 -0.149038,-0.236806 -0.09921,-0.379884 0.584727,-2.433682 1.107884,-5.544012 1.107884,-6.857563 v -12.9963 h 5.506907 c 0,0.388227 0.122652,1.132346 0.222329,1.177568 0.09099,-0.0445 0.204778,-0.778955 0.204778,-1.172044 h 2.280951 l 3.046441,1.0241 c 0,0.424123 0.128198,1.042922 0.212061,1.09134 0.104274,-0.02794 0.223081,-0.541609 0.223081,-0.963883 l 5.249134,1.744889 c 0,0.429295 0.143808,1.108614 0.248267,1.168923 0.09998,-0.05772 0.200845,-0.594159 0.200845,-1.04188 l 11.219429,3.752126 c -0.0339,-2.152744 -0.08025,-2.204888 -1.030436,-2.819032 L 43.444749,28.980372 C 43.56864,28.71194 43.95064,27.669184 43.95064,27.411076 h 0.464594 c 0.09292,0 0.278757,-1.321513 0.278757,-2.973405 0,-1.641568 -0.216811,-2.725621 -0.402648,-2.725621 h -3.65481 c -0.165511,0 -0.392325,0.867243 -0.392325,2.777242 0,1.208159 0.0855,2.229276 0.154865,2.488162 l -4.047134,-2.653351 c -0.423298,-0.268432 -0.867243,-0.59881 -0.867243,-1.300864 v -9.725512 c 0,-5.9571336 -0.817452,-8.931913 -1.416263,-10.2011175 L 33.492101,1.7241618 C 33.254642,1.1356755 32.735241,0.23475373 32.17932,0.23475373 Z",
		scale:   0.45,
	},
	"B77W": {
		name:    "BOEING 777-300ER",
		svgPath: "m 34.769885,0.5081602 c -0.909111,0 -1.805794,2.8665215 -2.013096,3.322586 -0.524982,1.4431263 -1.192222,2.7167027 -1.192222,8.8928038 v 15.510988 c 0,0.681929 -0.621283,1.354373 -0.966786,1.603135 l -4.02166,2.667287 c 0.165842,-0.691007 0.165842,-1.796618 0.165842,-2.501445 0,-1.06415 -0.138201,-2.750207 -0.400784,-2.750207 h -3.662336 c -0.221122,0 -0.414604,1.520215 -0.414604,2.736387 0,1.326733 0.193482,2.902229 0.304043,2.902229 h 0.469885 c 0.09674,0.538985 0.414604,1.368193 0.525165,1.575495 L 2.5567248,48.370475 c -0.2081497,0.143094 -0.3732872,0.311589 -0.4146037,0.359324 l -1.48249005,1.993386 c -0.0293169,0.04397 -0.0781785,0.09039 -0.0781785,0.180787 v 0.615656 L 2.6727273,50.395812 c 0.06352,-0.03909 0.3046604,-0.147076 0.3860063,-0.175901 L 14.137999,46.535852 c 0,0.224578 0.08983,0.984685 0.193482,0.984685 0.117471,0 0.228032,-0.860303 0.228032,-1.109065 l 5.282747,-1.772433 c 0,0.279858 0.100195,1.012325 0.217667,1.012325 0.110562,0 0.176207,-0.832663 0.176207,-1.119431 l 3.043885,-1.074515 h 2.284349 c 0,0.293222 0.117268,1.182503 0.217434,1.182503 0.100166,0 0.207661,-0.957687 0.207661,-1.177564 h 5.536015 v 18.435466 c 0,1.368123 0.239746,2.547612 0.566794,4.202094 l 0.459299,2.286721 c 0.07818,0.395779 0.141698,0.657188 -0.05863,0.830647 l -9.215707,7.398305 c -0.156724,0.204209 -0.09039,2.604783 0.05807,2.525418 l 10.358187,-3.683066 c 0.147278,-0.0883 0.08983,-1.008869 0.05528,-1.181621 l 0.891399,3.800537 h 0.300731 l 0.894167,-3.802402 c -0.07329,0.38112 -0.02932,1.143361 0.05863,1.18245 l 10.319565,3.679275 c 0.163154,0.02528 0.273624,-2.242745 0.09284,-2.51637 l -9.171398,-7.38093 C 36.879033,69.02467 36.941223,68.803547 37.031054,68.388943 l 0.476795,-2.287233 c 0.29073,-1.53916 0.525165,-2.860768 0.525165,-4.228961 V 43.478148 h 5.555694 c 0,0.330313 0.125977,1.168412 0.188548,1.168412 0.07133,0 0.24073,-0.833641 0.24073,-1.176904 h 2.28248 l 3.054339,1.07518 c 0,0.337048 0.07406,1.099199 0.171506,1.11479 0.09355,-0.01949 0.253361,-0.647046 0.253361,-0.986161 l 5.250428,1.738449 c 0,0.343013 0.109234,1.102776 0.210485,1.106995 0.09355,0.0039 0.218281,-0.615863 0.218281,-0.95108 l 11.260941,3.738055 c 0.204067,0.06975 0.288202,0.122336 0.509055,0.259054 l 1.766511,0.96971 v -0.605429 c -0.0014,-0.13493 -0.01863,-0.155627 -0.125561,-0.290159 l -1.401964,-1.906957 c -0.112749,-0.148189 -0.11228,-0.152745 -0.287505,-0.287505 L 46.034094,34.487368 c 0.133709,-0.324722 0.467982,-1.34664 0.467982,-1.556754 h 0.458431 c 0.105057,0 0.353373,-0.987956 0.353373,-3.046655 0,-1.050571 -0.238766,-2.616876 -0.439329,-2.616876 h -3.629244 c -0.172177,0 -0.420228,1.126976 -0.420228,2.741034 0,0.964615 0.09551,2.062939 0.181462,2.464066 L 38.937967,29.845756 C 38.240771,29.31092 38.049758,28.804736 38.049758,28.250799 V 12.721455 c 0,-6.1678211 -0.730835,-7.4397363 -1.200819,-8.8608767 C 36.591568,3.3905948 35.788182,0.5081602 34.769885,0.5081602 Z",
		scale:   0.45,
	},
	"B788": {
		name:    "BOEING 787-8 Dreamliner",
		svgPath: "m 623.49023,3.1855469 c -14.93714,0 -59.90625,69.4791411 -59.90625,188.1562531 V 392.94336 L 449.1582,472.98242 c 7.88138,-27.087 9.38282,-53.2334 9.38282,-68.8789 0,-19.82472 -1.60718,-32.25566 -3.10743,-40.18555 -0.96442,-4.07211 -3.01112,-12.22266 -7.93554,-12.22266 h -60.54102 c -7.60841,0 -8.14441,9.0081 -8.89453,12.33008 -0.96445,7.71558 -3.10742,17.99751 -3.10742,40.29102 0,23.21797 5.00416,55.40404 10.28711,73.72851 h 13.43359 c 1.01741,2.99985 7.40112,20.46758 9.13672,23.54688 L 53.355469,751.23828 c -5.486763,3.75114 -12.121345,8.13408 -15.732422,14.38867 L 3.359375,820.21484 v 21.22071 L 41.599609,808.84961 C 46.223211,805.00307 55.371313,798.49341 61.082031,795.75 L 223.89453,725.42969 c 0,12.76513 4.93429,19.625 6.33008,19.625 1.56375,0 6.09766,-9.2081 6.09766,-25.05664 L 314.64844,686.125 c 0,10.91754 4.19994,19.5957 6.10351,19.5957 2.07154,0 6.38282,-9.23827 6.38282,-25.13867 l 65.61718,-28.60937 h 43.33399 c 0,13.04507 4.87182,22.95507 6.43945,22.95507 2.07152,0 6.12305,-9.87198 6.12305,-22.85156 H 563.51367 V 883.3125 c 0,39.85395 11.2729,98.40462 14.61328,114.49414 3.25694,14.33646 -1.13607,20.12176 -4.45703,23.22856 l -142.1582,126.4121 c -2.40484,2.2534 -5.70952,5.8914 -6.53516,9.748 l -7.49414,30.2695 187.0586,-68.1601 c 6.10272,25.1581 13.14638,56.1621 13.70507,56.1621 0.0302,0 0.0236,0.043 0.0274,0.045 h 4.15039 1.06055 5.21289 c 0.004,-0 -0.003,-0.045 0.0273,-0.045 0.55869,0 7.60235,-31.004 13.70508,-56.1621 l 187.06054,68.1601 -7.49414,-30.2695 c -0.82563,-3.8566 -4.13031,-7.4946 -6.53515,-9.748 L 673.30273,1021.0352 c -3.32095,-3.1068 -7.71396,-8.8921 -4.45703,-23.22856 3.34039,-16.08952 14.61328,-74.64019 14.61328,-114.49414 V 652.07617 h 114.86524 c 0,12.97958 4.05152,22.85156 6.12305,22.85156 1.56763,0 6.4375,-9.91 6.4375,-22.95507 h 43.33593 l 65.61719,28.60937 c 0,15.9004 4.3093,25.13867 6.38086,25.13867 1.90356,0 6.10352,-8.67816 6.10352,-19.5957 l 78.32613,33.87305 c 0,15.84854 4.5359,25.05664 6.0996,25.05664 1.3958,0 6.3301,-6.85987 6.3301,-19.625 L 1185.8887,795.75 c 5.7107,2.74341 14.8607,9.25307 19.4843,13.09961 l 38.2403,32.58594 v -21.22071 l -34.2656,-54.58789 c -3.6111,-6.25459 -10.2457,-10.63753 -15.7325,-14.38867 L 839.16016,501.5918 c 1.73559,-3.0793 8.11735,-20.54703 9.13476,-23.54688 h 13.4336 c 5.28294,-18.32447 10.2871,-50.51054 10.2871,-73.72851 0,-22.29351 -2.14296,-32.57542 -3.10742,-40.29102 -0.75012,-3.32198 -1.28611,-12.33008 -8.89453,-12.33008 h -60.54101 c -4.92443,0 -6.96918,8.15056 -7.9336,12.22266 -1.50024,7.9299 -3.10742,20.36083 -3.10742,40.18555 0,15.6455 1.49949,41.7919 9.38086,68.8789 L 683.38867,392.94336 V 191.3418 c 0,-118.677112 -44.96129,-188.1562531 -59.89844,-188.1562531 z",
		scale:   0.026,
	},
	"DH8D": {
		name:    "DE HAVILLAND CANADA DHC-8-400 Dash 8",
		svgPath: "m 47.148749,2.06375 c -1.641604,0 -4.28625,7.1889238 -4.28625,11.985625 V 46.381458 L 42.2275,46.646042 c -3.201567,0 -7.593542,0.423333 -7.593542,0.423333 v -7.328958 c 0,-1.481667 -0.661458,-3.677709 -1.296458,-3.677709 -0.555625,0 -1.402292,2.143125 -1.402292,3.704167 v 7.46125 L 4.0216667,49.424167 c -0.873125,0.05292 -0.9789584,0.3175 -1.0054167,0.608541 l -0.2645833,2.936875 c -0.026458,0.423334 0,0.47625 0.238125,0.502709 l 12.9910413,0.952499 0.238125,0.873126 0.3175,-0.820208 8.016875,0.502708 0.211667,1.296458 0.47625,-1.243542 6.6675,0.3175 v 4.180417 c 0,0.926042 0.714375,3.889375 1.243542,3.889375 0.532141,0 1.613958,-2.883917 1.613958,-3.862917 v -3.254375 h 7.752292 c 0.291041,0.02646 0.343958,0.291042 0.343958,0.502709 v 18.626666 c 0,3.915834 2.910417,24.659162 2.910417,24.659162 -2.116667,0.10585 -11.059584,1.48167 -11.059584,1.48167 v 4.68313 h 12.250209 c 0.07938,1.34937 0.370416,2.7252 0.370416,2.7252 0,0 0.370417,-1.48166 0.47625,-2.7252 h 12.012084 v -4.70959 c 0,0 -8.916459,-1.29645 -10.95375,-1.561038 0,0 2.619375,-20.611042 2.619375,-24.579792 V 56.885417 c 0.02646,-0.343959 0.264583,-0.555625 0.555625,-0.555625 H 60.0075 v 3.095625 c 0,1.190625 0.687917,3.96875 1.190625,3.96875 0.555625,0 1.613958,-2.751667 1.613958,-3.836459 v -3.915833 l 6.6675,-0.343958 0.47625,0.978958 0.291042,-1.058333 8.202083,-0.555625 0.264584,0.714375 0.264583,-0.79375 12.85875,-0.978958 c 0.343958,-0.05292 0.449792,-0.291042 0.449792,-0.635001 l -0.07938,-2.301875 C 92.154377,49.900416 91.678123,49.503542 91.09604,49.424166 L 62.785625,47.333958 v -7.567083 c 0,-1.561042 -0.846667,-3.730625 -1.42875,-3.730625 -0.555625,0 -1.481667,2.169583 -1.481667,3.730625 V 46.99 c -2.804583,-0.185208 -8.043333,-0.370417 -8.043333,-0.370417 l -0.529167,-0.47625 V 14.102292 c 0,-4.8418753 -2.513544,-12.038542 -4.153959,-12.038542 z",
		scale:   0.3,
	},
	"E190": {
		name:    "EMBRAER ERJ-190-100",
		svgPath: "m 57.784999,1.5875 c -5.667507,2.9104167 -5.767916,14.975416 -5.767916,14.975416 v 36.565417 l -2.434166,4.1275 L 43.18,60.589583 C 44.185417,59.107917 44.238333,54.2925 44.238333,54.2925 44.502917,50.64125 43.65625,49.53 43.65625,49.53 h -7.249583 c -0.9525,0.899584 -0.899584,5.185833 -0.899584,5.185833 -0.05292,4.497917 1.11125,5.715 1.11125,5.715 l 0.635,0.05292 L 37.7825,63.1825 4.4979167,80.062917 c -1.1112501,0.740833 -1.27,1.693333 -1.27,1.693333 l -2.06375,7.46125 v 1.42875 h 0.635 l 2.1166666,-3.545417 16.2983337,-4.28625 v 2.010834 L 20.584583,85.725 21.007917,84.719583 V 82.602916 L 31.90875,79.6925 v 1.693333 l 0.370417,0.899584 0.370416,-0.582084 v -2.2225 l 4.868334,-1.322916 H 42.06875 V 80.01 l 0.423333,0.582083 0.423334,-0.370416 V 78.105 h 9.048749 v 29.79208 c 0.250906,10.4753 2.381251,18.46792 2.381251,18.46792 l -14.605,10.21292 c -1.534584,1.11124 -1.42875,3.12208 -1.42875,3.12208 l -0.05292,3.38667 17.727086,-4.70959 c 0.9525,6.0325 1.74625,5.97959 1.74625,5.97959 0.740834,0.0529 1.957917,-6.08542 1.957917,-6.08542 l 17.674167,4.86833 -0.05292,-2.96333 c -0.158747,-2.80459 -0.846664,-3.175 -0.846664,-3.175 l -15.24,-10.68916 c 2.698751,-12.64709 2.487084,-18.30917 2.487084,-18.30917 V 77.999167 h 8.73125 l 0.370416,2.487083 0.740834,-2.487083 H 78.105 l 4.60375,1.534583 0.47625,2.54 0.687916,-2.169584 10.742084,3.069168 0.211667,2.592916 0.687916,-2.275417 16.033747,4.28625 1.74625,2.275417 0.47625,0.846666 0.58209,-3e-6 v -2.010834 l -2.01084,-7.196667 c -0.37041,-1.058329 -1.64041,-1.746245 -1.64041,-1.746245 L 77.681667,63.1825 l 0.47625,-2.645833 h 0.635 C 79.957084,60.325 80.221667,54.9275 80.221667,54.9275 80.115833,49.318334 79.163333,49.318333 79.163333,49.318333 H 71.91375 c -1.11125,1.957917 -0.899583,4.92125 -0.899583,4.92125 0.105833,4.233334 1.27,6.244167 1.27,6.244167 l -5.87375,-2.963333 -2.69875,-4.392083 V 16.51 C 63.552917,4.0745833 57.784999,1.5875 57.784999,1.5875 Z",
		scale:   0.22,
	},
	"F100": {
		name:    "FOKKER 100",
		svgPath: "M 46.861745,1.2786934 C 43.475576,1.1366164 41.533857,13.402601 41.533857,13.402601 v 36.339066 l -12.694274,6.085416 -25.9291663,9.525 c -1.5875,0.79375 -1.5875,1.5875 -1.5875,1.5875 L 0.26458333,71.172917 19.05,69.585417 l 0.264583,2.645833 0.529167,-2.645833 8.995833,-0.529167 0.79375,3.704167 0.79375,-3.96875 10.054167,-0.79375 c 0,2.116666 1.058333,4.233333 1.058333,4.233333 V 79.375 h -0.79375 v -1.322917 c -0.79375,-1.322916 -3.175,-1.322916 -3.175,-1.322916 -2.116666,0 -2.116666,0.79375 -2.116666,0.79375 v 6.085416 c 0,5.027084 0.79375,9.789584 0.79375,9.789584 0.79375,0.529166 2.116666,0.529166 2.116666,0.529166 1.322917,0 1.5875,-0.529166 1.5875,-0.529166 l 2.645834,2.38125 c 1.5875,8.731253 3.439583,13.493753 3.439583,13.493753 l -14.552083,8.99583 c -1.5875,1.32292 -1.5875,3.43958 -1.5875,3.43958 v 1.85209 L 46.0375,119.85625 c 0,1.85208 0.79375,2.91042 0.79375,2.91042 0.529167,-0.52917 1.058333,-2.91042 1.058333,-2.91042 l 15.875,3.70417 v -2.11667 c 0,-2.11667 -1.5875,-3.175 -1.5875,-3.175 L 47.625,109.27292 C 49.741667,104.775 51.064583,95.779167 51.064583,95.779167 L 53.975,93.397917 c 0.529167,0.529166 1.5875,0.529166 1.5875,0.529166 1.322917,0 1.852083,-0.529166 1.852083,-0.529166 0.529167,-3.439584 0.79375,-10.054167 0.79375,-10.054167 0.264584,-2.645833 0,-6.35 0,-6.35 -1.058333,-0.264583 -2.116666,-0.264583 -2.116666,-0.264583 -1.852084,0 -2.910417,0.79375 -2.910417,0.79375 C 52.916667,78.052083 52.916667,79.375 52.916667,79.375 H 52.3875 v -8.73125 c 0.529167,-0.529167 0.529167,-2.645833 0.529167,-2.645833 L 63.5,68.791667 c 0,1.5875 0.79375,3.96875 0.79375,3.96875 0.264583,-0.79375 0.529167,-3.704167 0.529167,-3.704167 l 8.995833,0.529167 c 0,0.79375 0.529167,2.645833 0.529167,2.645833 C 74.6125,71.172917 74.6125,69.585417 74.6125,69.585417 l 18.785417,1.5875 L 92.86875,68.2625 c -0.529167,-2.645833 -2.116667,-2.910417 -2.116667,-2.910417 l -25.664583,-9.525 -12.7,-6.35 V 13.49375 C 50.270833,0.79375 46.861745,1.2786934 46.861745,1.2786934 Z",
		scale:   0.28,
	},
	"HAWK": {
		name:    "BAE SYSTEMS T-45 Goshawk",
		svgPath: "m 350.27734,2.4414062 v 1.9472657 h -1.06836 v 4.2714843 h 1.00586 V 25.496094 h -1.00586 v 23.68164 c -8.55231,0 -33.79492,55.175946 -33.79492,122.810546 v 169.35938 h -2.07812 v -5.72657 h -17.68555 c -3.0327,1.21798 -15.2793,5.97238 -15.2793,15.91407 0,27.07357 4.05274,49.5 4.05274,49.5 l -10.48242,8.3164 -27.64649,18.33985 L 32.443359,536.02539 C 8.0976331,550.08142 3.4324084,581.38941 3.3339844,592.49219 v 22.27148 L 305.1875,580.00781 c 1.24481,3.1883 3.97904,9.61345 5.75586,10.16016 0,31.70937 5.33008,124.24023 5.33008,124.24023 -14.76125,16.26467 -18.45092,24.87486 -24.19141,63.96485 l -88.43164,63.41797 c -4.86103,2.94909 -14.21484,15.85538 -14.21484,41.82421 0,1.37893 0.93869,2.50083 2.73437,2.1875 l 12.22461,-1.19531 131.39649,-34.15039 c 1.00687,3.77575 1.34102,5.03516 2.5996,5.03516 h 9.73438 v 2.09765 h -1.0918 v 3.10352 h 3.69922 3.94336 v -3.10352 h -1.08984 v -2.09765 h 9.73242 c 1.25858,0 1.5947,-1.25941 2.60156,-5.03516 l 131.39649,34.15039 12.22461,1.19531 c 1.79565,0.31333 2.73242,-0.80857 2.73242,-2.1875 0,-25.96883 -9.35381,-38.87512 -14.21485,-41.82421 l -88.42968,-63.41797 c -5.7405,-39.08999 -9.43212,-47.70018 -24.19336,-63.96485 0,0 5.33203,-92.53086 5.33203,-124.24023 1.77679,-0.54671 4.50912,-6.97186 5.7539,-10.16016 L 698.375,614.76367 v -22.27148 c -0.0983,-11.10278 -4.76366,-42.41077 -29.10938,-56.4668 l -213.8496,-108.33594 -27.64844,-18.33789 -10.48047,-8.3164 c 0,0 4.05078,-22.42659 4.05078,-49.5 0,-9.94171 -12.24465,-14.69608 -15.27734,-15.91407 h -17.6875 v 5.72657 h -2.07617 V 171.98828 c 0,-67.63461 -25.24457,-122.810546 -33.79688,-122.810546 v -23.68164 h -1.00586 V 8.6601562 H 352.5 V 4.3886719 h -1.06836 V 2.4414062 Z",
		scale:   0.035,
	},
	"PC12": {
		name:    "PILATUS PC-12",
		svgPath: "m 79.106876,0.66218011 c -1.359493,-0.00697 -2.160254,4.17289299 -2.289625,4.39264259 -1.967034,-1.0118725 -4.15542,-0.9005986 -4.15542,-0.9005986 -3.486582,0 -4.796659,1.2825704 -4.796659,1.2825704 l -0.0978,0.9550896 7.305392,-0.1976392 1.530753,0.2014878 -0.220407,0.9077521 c -0.399426,0.3516962 -1.034997,2.7690542 -1.034997,2.7690542 -2.069779,-0.2315048 -3.62981,2.84763 -3.62981,2.84763 0.672464,0.720252 0.795355,1.382282 0.795355,1.382282 0.172968,0.630384 0.122412,1.171252 0.122412,1.171252 0.797976,-1.770997 1.960666,-2.728388 1.960666,-2.728388 -2.616132,10.685321 -2.750904,13.463044 -2.750904,13.463044 L 70.74069,34.949962 70.70507,48.09274 6.8994795,50.822026 5.8881872,51.363393 4.5792198,53.366579 1.1044002,60.472344 0.87835887,62.227938 1.1879514,64.40816 l 2.1898163,-1.61448 2.0466784,-0.766509 1.105854,0.117381 17.3707949,1.202297 7.186496,0.709679 0.286276,1.225773 0.325785,-1.113523 16.996822,1.912617 0.382133,1.132766 0.285628,-1.130842 19.888144,2.381823 0.895033,0.686523 0.640947,2.156937 -0.04767,13.69537 c 1.310302,11.427521 2.365858,15.702858 2.365858,15.702858 l 1.03707,10.53356 1.002613,1.51057 2.556729,10.05558 -21.422117,2.68951 -1.357673,2.48773 -1.228396,3.56571 0.241586,2.55225 23.989079,1.26092 c 0.316795,2.70452 1.114144,4.22562 1.114144,4.22562 0.942653,-1.95158 0.996137,-4.16289 0.996137,-4.16289 l 24.097051,-1.29569 0.24132,-2.45077 -1.22898,-3.69515 -1.35747,-2.48798 -21.419012,-2.68926 c 1.06656,-4.22428 2.451545,-10.04006 2.451545,-10.04006 l 1.093742,-1.53314 1.002225,-10.52651 c 1.462837,-7.701803 2.329457,-16.103111 2.329457,-16.103111 l -9.07e-4,-13.574203 0.491462,-1.975798 1.061812,-0.588384 19.863726,-2.335191 0.42812,1.084018 0.19106,-1.130842 17.07319,-1.895042 0.29729,1.329685 0.50519,-1.370095 6.86926,-0.845982 17.4947,-1.157076 1.78061,-0.07049 1.95211,0.84925 2.31275,1.508195 0.0622,-3.988413 -3.55321,-7.245596 c -1.07923,-2.001236 -2.53814,-2.388237 -3.39632,-2.404273 L 88.637338,48.278242 87.270403,47.974075 87.102135,35.021354 86.104185,25.485467 C 84.887101,17.89936 83.381326,12.926262 83.381326,12.926262 c 1.322167,1.131917 2.194868,2.546029 2.194868,2.546029 -0.521042,-1.81351 0.700145,-2.509275 0.700145,-2.509275 C 84.648758,9.8327359 82.743358,10.169395 82.743358,10.169395 82.167513,7.2904327 81.510364,7.0605751 81.510364,7.0605751 l 0.261081,-0.6011484 c 0.626602,0.3506788 0.950892,-0.1146079 0.950892,-0.1146079 1.536857,-0.7331791 8.077226,0.047065 8.077226,0.047065 L 90.701113,5.4367945 C 90.096427,4.2298792 86.063942,4.1452816 86.063942,4.1452816 83.526015,4.1170826 81.115277,5.1529616 81.115277,5.1529616 80.367705,0.32429073 79.106876,0.66218011 79.106876,0.66218011 Z",
		scale:   0.2,
	},
	"RV9": {
		name:    "VAN'S RV-9",
		svgPath: "m 187.29883,2.5117188 c -3.97579,0 -7.71875,12.9699732 -7.71875,17.2773432 h -7.5918 c -1.67444,0 -1.55433,-0.05396 -2.15234,3.294922 l -4.78321,32.650391 -0.83789,17.822266 H 21.050781 c -4.493133,0 -12.1451985,2.83044 -13.0859372,11.234375 L 2.1953125,131.32812 H 165.94727 l 13.79687,85.23047 h -54.79102 c -2.71635,0 -4.90234,2.71613 -4.90234,4.90235 v 32.0664 h 57.93555 l 5.43164,-12.67773 1.55273,4.14062 1.16407,15.7793 -1.80086,0.33536 1.15437,1.86386 0.25977,4.6582 1.54449,2.27967 1.54535,-2.27967 0.25781,-4.6582 1.23822,-1.86386 -1.8847,-0.33536 1.16406,-15.7793 1.55274,-4.14062 5.43164,12.67773 h 57.9375 v -32.0664 c 0,-2.18622 -2.18599,-4.90235 -4.90235,-4.90235 h -54.79297 l 13.79883,-85.23047 h 163.75 L 366.61914,84.791016 C 365.6784,76.387081 358.02829,73.556641 353.53516,73.556641 H 210.36914 L 209.5332,55.734375 204.74805,23.083984 c -0.59802,-3.348873 -0.4779,-3.294922 -2.15235,-3.294922 h -7.58984 c 0,-4.307372 -3.73124,-17.2773432 -7.70703,-17.2773432 z",
		scale:   0.07,
	},
	"SF34": {
		name:    "SAAB-FAIRCHILD SF-340",
		svgPath: "M 94.720832,2.06375 C 84.084583,6.4558334 84.666666,31.485416 84.666666,31.485416 V 68.58 L 69.267917,70.379167 V 66.83375 64.187917 l 0.79375,-2.06375 -1.058334,-1.005417 c -0.370417,-7.037917 -0.846666,-7.9375 -0.846666,-7.9375 1.957916,0.211666 3.280833,0.158754 3.280833,0.158754 7.9375,-0.158755 7.884584,-0.846671 7.884584,-0.846671 C 77.893333,51.59375 71.4375,51.964167 71.4375,51.964167 L 67.839167,52.07 C 66.622083,48.365834 65.7225,48.365833 65.7225,48.365833 c -1.27,0.3175 -2.275417,3.757084 -2.275417,3.757084 L 60.0075,52.07 c -7.9375,0.05292 -7.884583,0.635 -7.884583,0.635 0.105833,0.846667 7.831666,0.687917 7.831666,0.687917 1.27,0.211666 3.227917,-0.05292 3.227917,-0.05292 -1.5875,5.60917 -1.27,13.440836 -1.27,13.440836 V 71.27875 L 4.7095833,78.052083 C 2.2225,78.84583 2.2754167,82.126667 2.2754167,82.126667 L 1.905,82.920417 2.2754167,83.449583 v 3.968751 l 60.8541663,3.598333 c 0.05292,2.38125 2.38125,2.328333 2.38125,2.328333 2.116667,0.211666 2.38125,-2.116667 2.38125,-2.116667 l 16.615834,1.005417 v 40.11083 c 0,8.6785 1.481666,13.07042 1.481666,13.07042 l -26.564166,4.445 c -2.328333,0.58209 -2.434167,2.8575 -2.434167,2.8575 l 0.05292,5.18583 c -4e-6,0.635 1.05833,0.68792 1.05833,0.68792 l 30.215417,4.1275 2.487083,2.16958 c 1.5875,8.36083 3.386667,8.36084 3.386667,8.36084 l 0.211666,1.16416 0.211667,-1.11125 c 2.010833,-0.15875 3.4925,-8.57249 3.4925,-8.57249 l 2.43416,-2.27542 30.63876,-4.445 c 0.89958,-0.21168 0.84666,-0.635 0.84666,-0.635 v -4.18042 c -0.10583,-3.28084 -2.27541,-3.65125 -2.27541,-3.65125 l -26.77584,-4.33917 c 1.64042,-8.14917 1.48167,-12.91166 1.48167,-12.91166 V 91.969167 l 16.29833,-0.9525 c 0.26457,2.2225 2.43417,2.116666 2.43417,2.116666 2.38125,0.05292 2.43417,-2.487083 2.43417,-2.487083 l 61.11875,-3.704167 v -3.65125 l 0.68791,-0.687916 -0.635,-0.740834 C 186.47833,77.522916 184.52042,77.47 184.52042,77.47 l -57.41459,-6.561667 v -4.28625 -2.963333 c 0.42334,-0.47625 0.42334,-1.164167 0.42334,-1.164167 0.0529,-1.375833 -0.68792,-1.74625 -0.68792,-1.74625 0,-4.339166 -1.00542,-7.725833 -1.00542,-7.725833 1.32292,0.370417 3.12209,0.05292 3.12209,0.05292 8.20208,0.634997 7.99041,-0.740837 7.99041,-0.740837 -0.0529,-0.793749 -8.04333,-0.423333 -8.04333,-0.423333 -1.32292,-0.05292 -3.38667,0.05292 -3.38667,0.05292 -1.27,-3.810003 -2.27541,-3.810003 -2.27541,-3.810003 -0.95251,0.05292 -2.16959,3.862916 -2.16959,3.862916 -11.64167,-0.423333 -11.27125,0.582084 -11.27125,0.582084 0,0.740833 7.72584,0.740833 7.72584,0.740833 1.69333,0.105833 3.175,-0.15875 3.175,-0.15875 C 119.22125,58.89625 119.38,66.675 119.38,66.675 v 3.439583 L 104.61625,68.42125 V 31.538333 c 0,-26.0905456 -9.895418,-29.474583 -9.895418,-29.474583 z",
		scale:   0.17,
	},
	"SW3": {
		name:    "SWEARINGEN Merlin 3",
		svgPath: "m 318.56445,3.2109375 c -2.24808,0 -10.18066,7.6147715 -18.80468,33.7988285 -16.04996,48.730484 -18.10352,79.129474 -18.10352,148.906254 v 22.40234 l -51.1543,3.12109 c 0.69324,-7.62565 1.21485,-16.29198 1.21485,-27.73047 0,-29.63607 -4.50737,-54.94034 -7.62696,-63.08593 5.89255,2.07972 31.19709,1.9082 41.59571,1.9082 5.54593,0 5.71924,-7.45312 0,-7.45313 l -43.1543,1.2129 c -1.7331,-5.02601 -4.54459,-12.49024 -6.86523,-12.49024 -2.14733,0 -4.74737,7.29092 -6.48047,12.49024 l -44.88672,-1.38672 c -5.71925,0 -5.89359,7.10547 -0.34766,7.10547 6.93242,0 38.12842,0.86723 43.84766,-1.73243 -2.59966,9.18545 -8.83789,33.79609 -8.83789,63.60547 0,9.70538 -0.51993,21.8358 0,29.63477 L 9.7050781,225.30273 c -3.1195865,0.34662 -7.1063863,3.81309 -6.7597656,7.2793 L 5.71875,271.75 274.00391,315.94531 c 2.07972,7.97228 7.79882,36.91376 7.79882,46.61914 0,39.86139 8.35764,83.93335 12.13086,90.46875 l -71.57617,62.39258 c -5.71924,5.02602 -9.25,17.02449 -9.25,27.57617 v 15.92383 l 95.54883,-27.72851 v 3.8125 c 0,1.21965 1.81884,2.33398 5.5957,2.33398 0.34337,21.03027 3.09124,47.64062 4.29297,47.64063 1.20174,-10e-6 3.94746,-26.61036 4.29102,-47.64063 3.77684,0 5.63476,-1.11433 5.63476,-2.33398 v -3.8125 l 95.54883,27.72851 v -15.92383 c 0,-10.55168 -3.53074,-22.55015 -9.25,-27.57617 L 343.19336,453.0332 c 3.77322,-6.5354 12.13086,-50.60736 12.13086,-90.46875 0,-9.70537 5.7191,-38.64686 7.79883,-46.61914 L 631.4082,271.75 l 2.77344,-39.16797 c 0.34662,-3.46621 -3.64018,-6.93268 -6.75976,-7.2793 L 438.16797,213.51758 c 0.51995,-7.79897 0,-19.92939 0,-29.63477 0,-29.80938 -6.24021,-54.42002 -8.83985,-63.60547 5.71922,2.59966 36.91525,1.73243 43.84766,1.73243 5.54597,0 5.37355,-7.10547 -0.3457,-7.10547 l -44.88867,1.38672 c -1.73311,-5.19932 -4.33313,-12.49024 -6.48047,-12.49024 -2.32063,0 -5.13018,7.46423 -6.86328,12.49024 l -43.1543,-1.2129 c -5.71922,1e-5 -5.54593,7.45313 0,7.45313 10.39861,0 35.7012,0.17152 41.59375,-1.9082 -3.11958,8.14559 -7.625,33.44986 -7.625,63.08593 0,11.43849 0.51969,20.10482 1.21289,27.73047 l -51.1543,-3.12109 v -22.40234 c 0,-69.77678 -2.05161,-100.175771 -18.10156,-148.906254 C 328.74513,10.825707 320.81255,3.2109375 318.56445,3.2109375 Z",
		scale:   0.05,
	},
}

// for lookups from IATA to ICAO if needed, likely faster than iterating over Aircraft to find? maybe?
// map key-type is IATA
var IATAtoICAO = map[string]string{
	"320": "A320",
	"388": "A388",
	"100": "F100",
	"E90": "E190",
	"PL2": "PC12",
	"SF3": "SF34",
}

// Colour object
type RGBA struct {
	r, g, b, a float64
}

// Marker object
type Marker struct {
	Img     *ebiten.Image
	CentreX float64
	CentreY float64
	icao    string
	poly    [][][]float64
}

func (m *Marker) PointInsideMarker(x, y float64) bool {

	var pip bool

	points := make([][]float64, 0, 9)

	// make a 3x3 around the point to make this function a bit less "pixel perfect"
	points = append(points, []float64{x - 1, y - 1})
	points = append(points, []float64{x, y - 1})
	points = append(points, []float64{x + 1, y - 1})

	points = append(points, []float64{x - 1, y})
	points = append(points, []float64{x, y})
	points = append(points, []float64{x + 1, y})

	points = append(points, []float64{x - 1, y + 1})
	points = append(points, []float64{x, y + 1})
	points = append(points, []float64{x + 1, y + 1})

	for _, p := range points {
		pip = piper.Pip(p, m.poly) || pip
	}
	return pip
}

func GetAircraft(icao string, aircraftMarkers *map[string]Marker) (aircraftMarker Marker) { //*Marker {

	// determine image
	if _, ok := (*aircraftMarkers)[icao]; ok {
		// use marker that matches aircraft type if found
		aircraftMarker = (*aircraftMarkers)[icao]

	} else {

		switch icao {

		// close matches
		case "A359":
			aircraftMarker = (*aircraftMarkers)["A320"]

		case "A35K":
			aircraftMarker = (*aircraftMarkers)["A320"]

		case "B772":
			aircraftMarker = (*aircraftMarkers)["B77L"]

		case "B773":
			aircraftMarker = (*aircraftMarkers)["B77W"]

		case "BE55":
			aircraftMarker = (*aircraftMarkers)["SW3"]

		case "C208":
			aircraftMarker = (*aircraftMarkers)["RV9"]

		case "C210":
			aircraftMarker = (*aircraftMarkers)["RV9"]

		case "COL4":
			aircraftMarker = (*aircraftMarkers)["RV9"]

		case "DH8C":
			aircraftMarker = (*aircraftMarkers)["DH8D"]

		// catch-all
		default:
			aircraftMarker = (*aircraftMarkers)["B77L"]

			exists := false
			for _, item := range missingMarkers {
				if item == icao {
					exists = true
				}
			}
			if !exists {
				missingMarkers = append(missingMarkers, icao)
				log.Printf("Missing marker for aircraft type: %s", icao)
			}

		}
	}

	return aircraftMarker

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

	var wgInner sync.WaitGroup
	c := make(chan Marker, len(Aircraft))

	// Pre-render aircraft concurrently
	for k, v := range Aircraft {

		wgInner.Add(1)

		go func(k string, v marker) {
			defer wgInner.Done()
			log.Printf("Pre-rendering marker: %s (%s)", k, v.name)

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

			img, poly, err := imgFromSVG(r)
			if err != nil {
				log.Fatal(err)
			}
			c <- Marker{
				Img:     img,
				CentreX: float64(img.Bounds().Dx()) / 2,
				CentreY: float64(img.Bounds().Dy()) / 2,
				icao:    k,
				poly:    poly,
			}
		}(k, v)
	}

	wgInner.Wait()
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
