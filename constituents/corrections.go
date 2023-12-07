package constituents

import (
	"math"

	astro "github.com/ryan-lang/tides/astronomy"
)

// Constants
const DEG_TO_RAD = math.Pi / 180
const RAD_TO_DEG = 180 / math.Pi

func fUnity(a *astro.Astro) float64 {
	return 1
}

// Schureman equations 73, 65
func fMm(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	mean := (2/3.0 - math.Pow(math.Sin(omega), 2)) * (1 - (3/2.0)*math.Pow(math.Sin(i), 2))
	return (2/3.0 - math.Pow(math.Sin(I), 2)) / mean
}

// Schureman equations 74, 66
func fMf(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	mean := math.Pow(math.Sin(omega), 2) * math.Pow(math.Cos(0.5*i), 4)
	return math.Pow(math.Sin(I), 2) / mean
}

// Schureman equations 75, 67
func fO1(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	mean := math.Sin(omega) * math.Pow(math.Cos(0.5*omega), 2) * math.Pow(math.Cos(0.5*i), 4)
	return math.Sin(I) * math.Pow(math.Cos(0.5*I), 2) / mean
}

// Schureman equations 76, 68
func fJ1(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	mean := math.Sin(2*omega) * (1 - (3/2.0)*math.Pow(math.Sin(i), 2))
	return math.Sin(2*I) / mean
}

// Schureman equations 77, 69
func fOO1(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	mean := math.Sin(omega) * math.Pow(math.Sin(0.5*omega), 2) * math.Pow(math.Cos(0.5*i), 4)
	return math.Sin(I) * math.Pow(math.Sin(0.5*I), 2) / mean
}

// Schureman equations 78, 70
func fM2(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	mean := math.Pow(math.Cos(0.5*omega), 4) * math.Pow(math.Cos(0.5*i), 4)
	return math.Pow(math.Cos(0.5*I), 4) / mean

}

// Schureman equations 227, 226, 68
// Should probably eventually include the derivations of the magic numbers (0.5023 etc).
func fK1(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	nu := DEG_TO_RAD * a.SolarAnomaly()
	sin2IcosnuMean := math.Sin(2*omega) * (1 - (3/2.0)*math.Pow(math.Sin(i), 2))
	mean := 0.5023*sin2IcosnuMean + 0.1681
	return math.Pow(0.2523*math.Pow(math.Sin(2*I), 2)+0.1689*math.Sin(2*I)*math.Cos(nu)+0.0283, 0.5) / mean
}

// Schureman equations 215, 213, 204
// It can be (and has been) confirmed that the exponent for R_a reads 1/2 via Schureman Table 7
func fL2(a *astro.Astro) float64 {
	P := DEG_TO_RAD * a.P()
	I := DEG_TO_RAD * a.InclinationAngle()
	rAInv := math.Pow(1-12*math.Pow(math.Tan(0.5*I), 2)*math.Cos(2*P)+36*math.Pow(math.Tan(0.5*I), 4), 0.5)
	return fM2(a) * rAInv
}

// Schureman equations 235, 234, 71
// Again, magic numbers
func fK2(a *astro.Astro) float64 {
	omega, _ := pairAsRad(a.TerrestrialObliquity())
	i, _ := pairAsRad(a.LunarInclination())
	I := DEG_TO_RAD * a.InclinationAngle()
	nu := DEG_TO_RAD * a.SolarAnomaly()
	sinsqIcos2nuMean := math.Pow(math.Sin(omega), 2) * (1 - (3/2.0)*math.Pow(math.Sin(i), 2))
	mean := 0.5023*sinsqIcos2nuMean + 0.0365
	return math.Pow(0.2523*math.Pow(math.Sin(I), 4)+0.0367*math.Pow(math.Sin(I), 2)*math.Cos(2*nu)+0.0013, 0.5) / mean
}

// Schureman equations 206, 207, 195
func fM1(a *astro.Astro) float64 {
	P := DEG_TO_RAD * a.P()
	I := DEG_TO_RAD * a.InclinationAngle()
	qAInv := math.Pow(
		0.25+
			1.5*math.Cos(I)*math.Cos(2*P)*math.Pow(math.Cos(0.5*I), -0.5)+
			2.25*math.Pow(math.Cos(I), 2)*math.Pow(math.Cos(0.5*I), -4),
		0.5,
	)
	return fO1(a) * qAInv
}

// See e.g. Schureman equation 149
func fModd(a *astro.Astro, n float64) float64 {
	return math.Pow(fM2(a), n/2.0)
}

// Node factors u, see Table 2 of Schureman.
func uZero(a *astro.Astro) float64 {
	return 0
}

func uMf(a *astro.Astro) float64 {
	return -2.0 * a.LunarElongation()
}

func uO1(a *astro.Astro) float64 {
	return 2.0*a.LunarElongation() - a.SolarAnomaly()
}

func uJ1(a *astro.Astro) float64 {
	return -a.SolarAnomaly()
}

func uOO1(a *astro.Astro) float64 {
	return -2.0*a.LunarElongation() - a.SolarAnomaly()
}

func uM2(a *astro.Astro) float64 {
	return 2.0*a.LunarElongation() - 2.0*a.SolarAnomaly()
}

func uK1(a *astro.Astro) float64 {
	return -a.LunarPerigeeAnomaly()
}

// Schureman 214
func uL2(a *astro.Astro) float64 {
	I := DEG_TO_RAD * a.InclinationAngle()
	P := DEG_TO_RAD * a.P()
	R := RAD_TO_DEG * math.Atan(math.Sin(2*P)/((1/6.0)*math.Pow(math.Tan(0.5*I), -2)-math.Cos(2*P)))
	return 2.0*a.LunarElongation() - 2.0*a.SolarAnomaly() - R
}

func uK2(a *astro.Astro) float64 {
	return -2.0 * a.SolarPerigeeAnomaly()
}

// Schureman 202
func uM1(a *astro.Astro) float64 {
	I := DEG_TO_RAD * a.InclinationAngle()
	P := DEG_TO_RAD * a.P()
	Q := RAD_TO_DEG * math.Atan(((5*math.Cos(I)-1)/(7*math.Cos(I)+1))*math.Tan(P))
	return a.LunarElongation() - a.SolarAnomaly() + Q
}

func uModd(a *astro.Astro, n float64) float64 {
	return (n / 2.0) * uM2(a)
}

func pairAsRad(a, b float64) (float64, float64) {
	return DEG_TO_RAD * a, b
}
