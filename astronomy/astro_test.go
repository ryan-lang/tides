package astronomy_test

import (
	"math"
	"testing"
	"time"

	"github.com/ryan-lang/tides/astronomy"
	"github.com/stretchr/testify/assert"
)

var TEST_DATE = time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)

const VAL_TOLERANCE = 0.000001
const SPEED_TOLERANCE = 0.000001
const JULIAN_TOLERANCE = 0.000001

func TestJulianDate(t *testing.T) {
	localDate := time.Date(2023, 4, 10, 7, 0, 0, 0, time.Local)
	utcDate := time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)

	localDateExpected := 2460045.083333
	utcDateExpected := 2460044.500000

	localJulian := astronomy.JulianDate(localDate)
	utcJulian := astronomy.JulianDate(utcDate)

	assert.LessOrEqual(t, math.Abs(localJulian-localDateExpected), JULIAN_TOLERANCE)
	assert.LessOrEqual(t, math.Abs(utcJulian-utcDateExpected), JULIAN_TOLERANCE)
}

func TestValAndSpeed(t *testing.T) {

	a := astronomy.Astro{Time: TEST_DATE}

	testSet := []struct {
		Name          string
		Func          func() (float64, float64)
		ExpectedValue float64
		ExpectedSpeed float64
	}{
		{"LunarLongitude", a.LunarLongitude, 251.09824817610206, 0.5490165191936777},
		{"SolarLongitude", a.SolarLongitude, 17.976204134796717, 0.04106864016781501},
		{"LunarPerigee", a.LunarPerigee, 310.2269353387635, 0.004641808013309299},
		{"LunarNode", a.LunarNode, 34.96416085537629, -0.0022064056791100965},
		{"SolarPerigee", a.SolarPerigee, 283.33739705676953, 0.000001961252269116578},
		{"TerrestrialObliquity", a.TerrestrialObliquity, 23.436265471296906, -1.4832892045076382e-8},
		{"LunarInclination", a.LunarInclination, 5.144999999999982, 0},
		{"EquilibriumArgument", a.EquilibriumArgument, -53.122044041305344, 14.492052120974137},
	}

	for _, test := range testSet {

		value, speed := test.Func()

		if math.Abs(value-test.ExpectedValue) > VAL_TOLERANCE {
			t.Errorf("%s value = %v, want %v", test.Name, value, test.ExpectedValue)
		}
		if math.Abs(speed-test.ExpectedSpeed) > SPEED_TOLERANCE {
			t.Errorf("%s speed = %v, want %v", test.Name, speed, test.ExpectedSpeed)
		}
	}
}

func TestVal(t *testing.T) {

	a := astronomy.Astro{Time: TEST_DATE}

	testSet := []struct {
		Name          string
		Func          func() float64
		ExpectedValue float64
	}{
		{"InclinationAngle", a.InclinationAngle, 27.800493973844425},
		{"LunarElongation", a.LunarElongation, 5.709397347280401},
		{"SolarAnomaly", a.SolarAnomaly, 6.326072701951034},
		{"LunarPerigeeAnomaly", a.LunarPerigeeAnomaly, 4.501621177918025},
		{"SolarPerigeeAnomaly", a.SolarPerigeeAnomaly, 4.746245039723647},
		{"P", a.P, 304.5175379914831},
	}

	for _, test := range testSet {

		value := test.Func()

		if math.Abs(value-test.ExpectedValue) > VAL_TOLERANCE {
			t.Errorf("%s value = %v, want %v", test.Name, value, test.ExpectedValue)
		}
	}
}

func TestFixedAngle(t *testing.T) {
	expectedValue := 90.0
	expectedSpeed := 0.0

	a := astronomy.Astro{Time: TEST_DATE}
	value, speed := a.FixedAngle(90)

	if math.Abs(value-expectedValue) > VAL_TOLERANCE {
		t.Errorf("FixedAngle() value = %v, want %v", value, expectedValue)
	}
	if math.Abs(speed-expectedSpeed) > SPEED_TOLERANCE {
		t.Errorf("FixedAngle() speed = %v, want %v", speed, expectedSpeed)
	}
}
