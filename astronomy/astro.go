package astronomy

import (
	"math"
	"time"

	"github.com/soniakeys/meeus/v3/julian"
)

const (
	JULIAN_CENTURIES_TO_DEG_PER_HOUR = 1 / (24 * 365.25 * 100)
)

type (
	Astro struct {
		Time time.Time
	}
)

// s
func (a *Astro) LunarLongitude() (float64, float64) {
	return calcValAndSpeed(LUNAR_LONGITUDE, a.Time)
}

// h
func (a *Astro) SolarLongitude() (float64, float64) {
	return calcValAndSpeed(SOLAR_LONGITUDE, a.Time)
}

// p
func (a *Astro) LunarPerigee() (float64, float64) {
	return calcValAndSpeed(LUNAR_PERIGEE, a.Time)
}

// N
func (a *Astro) LunarNode() (float64, float64) {
	return calcValAndSpeed(LUNAR_NODE, a.Time)
}

// pp
func (a *Astro) SolarPerigee() (float64, float64) {
	return calcValAndSpeed(SOLAR_PERIGEE, a.Time)
}

// omega
func (a *Astro) TerrestrialObliquity() (float64, float64) {
	return calcValAndSpeed(TERRESTRIAL_OBLIQUITY, a.Time)
}

// i
func (a *Astro) LunarInclination() (float64, float64) {
	return calcValAndSpeed(LUNAR_INCLINATION, a.Time)
}

// T + h - s
func (a *Astro) EquilibriumArgument() (float64, float64) {
	hourAngle, hourSpeed := a.hourAngle()
	sAngle, sSpeed := a.LunarLongitude()
	hAngle, hSpeed := a.SolarLongitude()
	v := hourAngle + hAngle - sAngle
	s := hourSpeed + hSpeed - sSpeed
	return v, s
}

func (a *Astro) FixedAngle(angle float64) (float64, float64) {
	return calcValAndSpeed([]float64{angle}, a.Time)
}

// I
func (a *Astro) InclinationAngle() float64 {
	N, _ := a.LunarNode()
	i, _ := a.LunarInclination()
	omega, _ := a.TerrestrialObliquity()
	return modulus(inclinationAngle(N, i, omega), 360)
}

// xi
func (a *Astro) LunarElongation() float64 {
	N, _ := a.LunarNode()
	i, _ := a.LunarInclination()
	omega, _ := a.TerrestrialObliquity()
	return modulus(lunarElongation(N, i, omega), 360)
}

// nu
func (a *Astro) SolarAnomaly() float64 {
	N, _ := a.LunarNode()
	i, _ := a.LunarInclination()
	omega, _ := a.TerrestrialObliquity()
	return modulus(solarAnomaly(N, i, omega), 360)
}

// nup
func (a *Astro) LunarPerigeeAnomaly() float64 {
	N, _ := a.LunarNode()
	i, _ := a.LunarInclination()
	omega, _ := a.TerrestrialObliquity()
	return modulus(lunarPerigeeAnomaly(N, i, omega), 360)
}

// nupp
func (a *Astro) SolarPerigeeAnomaly() float64 {
	N, _ := a.LunarNode()
	i, _ := a.LunarInclination()
	omega, _ := a.TerrestrialObliquity()
	return modulus(solarPerigeeAnomaly(N, i, omega), 360)
}

func (a *Astro) P() float64 { // TODO rename
	p, _ := a.LunarPerigee()
	xi := a.LunarElongation()
	return p - (modulus(xi, 360))
}

// hour
func (a *Astro) hourAngle() (float64, float64) {
	v := (JulianDate(a.Time) - math.Floor(JulianDate(a.Time))) * 360.0
	return v, 15.0
}

func calcValAndSpeed(coeffs []float64, t time.Time) (float64, float64) {
	v := modulus(polynomial(coeffs, julianCenturies(t)), 360)
	s := derivativePolynomial(coeffs, julianCenturies(t)) * JULIAN_CENTURIES_TO_DEG_PER_HOUR
	return v, s
}

func polynomial(coeffs []float64, x float64) float64 {
	result := 0.0
	for i := len(coeffs) - 1; i >= 0; i-- {
		result = result*x + coeffs[i]
	}
	return result
}

func derivativePolynomial(coeffs []float64, x float64) float64 {
	if len(coeffs) == 1 {
		return 0 // The derivative of a constant is zero
	}

	result := 0.0
	for i := len(coeffs) - 1; i > 0; i-- {
		result = result*x + float64(i)*coeffs[i]
	}
	return result
}

// i
func inclinationAngle(N, i, omega float64) float64 {
	N = DEG_TO_RAD * N
	i = DEG_TO_RAD * i
	omega = DEG_TO_RAD * omega
	cosI := math.Cos(i)*math.Cos(omega) - math.Sin(i)*math.Sin(omega)*math.Cos(N)
	return RAD_TO_DEG * math.Acos(cosI)
}

// xi
func lunarElongation(N, i, omega float64) float64 {
	N = DEG_TO_RAD * N
	i = DEG_TO_RAD * i
	omega = DEG_TO_RAD * omega
	e1 := (math.Cos(0.5*(omega-i)) / math.Cos(0.5*(omega+i))) * math.Tan(0.5*N)
	e2 := (math.Sin(0.5*(omega-i)) / math.Sin(0.5*(omega+i))) * math.Tan(0.5*N)
	e1 = math.Atan(e1)
	e2 = math.Atan(e2)
	e1 = e1 - 0.5*N
	e2 = e2 - 0.5*N
	return -(e1 + e2) * RAD_TO_DEG
}

// nu
func solarAnomaly(N, i, omega float64) float64 {
	N = DEG_TO_RAD * N
	i = DEG_TO_RAD * i
	omega = DEG_TO_RAD * omega
	e1 := (math.Cos(0.5*(omega-i)) / math.Cos(0.5*(omega+i))) * math.Tan(0.5*N)
	e2 := (math.Sin(0.5*(omega-i)) / math.Sin(0.5*(omega+i))) * math.Tan(0.5*N)
	e1 = math.Atan(e1)
	e2 = math.Atan(e2)
	e1 = e1 - 0.5*N
	e2 = e2 - 0.5*N
	return (e1 - e2) * RAD_TO_DEG
}

// Schureman equation 224
// nup
func lunarPerigeeAnomaly(N, i, omega float64) float64 {
	I := DEG_TO_RAD * inclinationAngle(N, i, omega)
	nu := DEG_TO_RAD * solarAnomaly(N, i, omega)
	return RAD_TO_DEG * math.Atan(
		(math.Sin(2*I)*math.Sin(nu))/
			(math.Sin(2*I)*math.Cos(nu)+0.3347),
	)
}

// Schureman equation 232
func solarPerigeeAnomaly(N, i, omega float64) float64 {
	I := DEG_TO_RAD * inclinationAngle(N, i, omega)
	nu := DEG_TO_RAD * solarAnomaly(N, i, omega)
	tan2nupp := (math.Pow(math.Sin(I), 2) * math.Sin(2*nu)) /
		(math.Pow(math.Sin(I), 2)*math.Cos(2*nu) + 0.0727)
	return RAD_TO_DEG * 0.5 * math.Atan(tan2nupp)
}

func JulianDate(t time.Time) float64 {
	return julian.TimeToJD(t)
}

// T - Time in Julian centuries from J2000.0 (Meeus formula 11.1)
func julianCenturies(t time.Time) float64 {
	return (JulianDate(t) - 2451545.0) / 36525
}

func modulus(a, b float64) float64 {
	return math.Mod(math.Mod(a, b)+b, b)
}
