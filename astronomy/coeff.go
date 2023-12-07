package astronomy

import "math"

var (
	TERRESTRIAL_OBLIQUITY []float64
	SOLAR_PERIGEE         = []float64{280.46645 - 357.5291, 36000.76932 - 35999.0503, 0.0003032 + 0.0001559, 0.00000048}
	SOLAR_LONGITUDE       = []float64{280.46645, 36000.76983, 0.0003032}
	LUNAR_INCLINATION     = []float64{5.145}
	LUNAR_LONGITUDE       = []float64{218.3164591, 481267.88134236, -0.0013268, 1/538841.0 - 1/65194000.0}
	LUNAR_NODE            = []float64{125.044555, -1934.1361849, 0.0020762, 1 / 467410.0, -1 / 60616000.0}
	LUNAR_PERIGEE         = []float64{83.353243, 4069.0137111, -0.0103238, -1 / 80053.0, 1 / 18999000.0}
)

func init() {
	rawValues := []struct {
		degrees float64
		arcmins float64
		arcsecs float64
	}{
		{23, 26, 21.448},
		{0, 0, -4680.93},
		{0, 0, -1.55},
		{0, 0, 1999.25},
		{0, 0, -51.38},
		{0, 0, -249.67},
		{0, 0, -39.05},
		{0, 0, 7.12},
		{0, 0, 27.87},
		{0, 0, 5.79},
		{0, 0, 2.45},
	}

	// adjust the coefficients for T rather than U
	TERRESTRIAL_OBLIQUITY = make([]float64, len(rawValues))
	for i, v := range rawValues {
		TERRESTRIAL_OBLIQUITY[i] = sexagesimalToDecimal(v.degrees, v.arcmins, v.arcsecs) * math.Pow(1e-2, float64(i))
	}
}

func sexagesimalToDecimal(degrees float64, arcmins float64, arcsecs float64) float64 {
	return degrees + arcmins/60.0 + arcsecs/(60.0*60.0)
}
