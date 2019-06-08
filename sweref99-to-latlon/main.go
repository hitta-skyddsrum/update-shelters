package coordConv

import (
	"github.com/im7mortal/UTM"
)

func Sweref99ToLatLon(coordinates [2]float64) [2]float64 {
	lat, lon, err := UTM.ToLatLon(coordinates[0], coordinates[1], 33, "", true)

	if err != nil {
		panic(err)
	}

	return [2]float64{lat, lon}
}
