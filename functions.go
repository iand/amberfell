/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

import (
	"fmt"
	"github.com/kierdavis/go/amberfell/mm3dmodel"
	"math"
	"os"
)

func IntPosition(pos Vectorf) Vectori {
	return Vectori{uint16(Round(pos[XAXIS], 0)),
		uint16(Round(pos[YAXIS], 0)),
		uint16(Round(pos[ZAXIS], 0))}
}

// Round a float to given precision
func Round(val float64, prec int) float64 {

	var rounder float64
	intermed := val * math.Pow(10, float64(prec))

	if intermed > 0 {
		rounder = math.Floor(intermed + 0.5)
	} else {
		rounder = math.Ceil(intermed - 0.5)
	}

	return rounder / math.Pow(10, float64(prec))
}

// Function LoadModel loads and returns an MM3D model.
func LoadModel(filename string) (model *mm3dmodel.Model) {
	fmt.Printf("Loading model: %s\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	model, err = mm3dmodel.Read(f)
	if err != nil {
		panic(err)
	}

	if model.NDirtySegments() > 0 {
		fmt.Fprintf(os.Stderr, "***** MM3D Warning: found %d dirty segments in %s! Tell Kier to add more functionality!\n", model.NDirtySegments(), filename)
	}

	return model
}

func HeadingToCompass(heading float64) string {
	if heading > 337.5 || heading <= 22.5 {
		return "E"
	} else if heading > 22.5 && heading <= 67.5 {
		return "NE"
	} else if heading > 67.5 && heading <= 112.5 {
		return "N"
	} else if heading > 112.5 && heading <= 157.5 {
		return "NW"
	} else if heading > 157.5 && heading <= 202.5 {
		return "W"
	} else if heading > 202.5 && heading <= 247.5 {
		return "SW"
	} else if heading > 247.5 && heading <= 292.5 {
		return "S"
	}

	return "NW"

}
