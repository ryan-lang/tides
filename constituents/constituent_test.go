package constituents_test

import (
	"math"
	"testing"
	"time"

	astro "github.com/ryan-lang/tides/astronomy"
	"github.com/ryan-lang/tides/constituents"
	"github.com/stretchr/testify/assert"
)

var TEST_DATE = time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)

const VAL_TOLERANCE = 0.000001

func TestDoodsonNumbers(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}

	valuesActual, speedsActual := constituents.DoodsonNumbers(a)
	valuesExpected := []float64{-154.5664088341564, 247.25513254381076, 17.688723653775014, 310.1944426826883, 34.97960569512168, 283.33738332800283, 90}
	speedsExpected := []float64{14.492052120974137, 0.5490165191936777, 0.04106864016781501, 0.004641808013309299, -0.0022064056791100965, 0.000001961252269116578, 0}

	for i := range valuesExpected {
		assert.LessOrEqual(t, math.Abs(valuesExpected[i]-valuesActual[i]), VAL_TOLERANCE)
		assert.LessOrEqual(t, math.Abs(speedsExpected[i]-speedsActual[i]), VAL_TOLERANCE)
	}
}

func TestConstituentValue(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	c := constituents.CONSTITUENT_M2

	expectedValue := -309.1328176683128
	actualValue := c.Value(a)

	assert.LessOrEqual(t, math.Abs(expectedValue-actualValue), 0.000001)
}

func TestConstituentSpeed(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	c := constituents.CONSTITUENT_M2

	expectedSpeed := 28.984104241948273
	actualSpeed := c.Speed(a)

	assert.LessOrEqual(t, math.Abs(expectedSpeed-actualSpeed), VAL_TOLERANCE)
}

func TestCompoundConstituentValue(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	cc := constituents.NewCompoundConstituent("M6", []constituents.CompoundContituentMember{{constituents.CONSTITUENT_M2, 3}})

	expectedValue := -927.3984530049383
	actualValue := cc.Value(a)

	assert.LessOrEqual(t, math.Abs(expectedValue-actualValue), VAL_TOLERANCE)
}

func TestCompoundConstituentSpeed(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	cc := constituents.NewCompoundConstituent("M6", []constituents.CompoundContituentMember{{constituents.CONSTITUENT_M2, 3}})

	expectedSpeed := 86.95231272583933
	actualSpeed := cc.Speed(a)

	assert.LessOrEqual(t, math.Abs(expectedSpeed-actualSpeed), VAL_TOLERANCE)
}
