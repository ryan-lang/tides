package constituents

import (
	"math"
	"testing"
	"time"

	astro "github.com/ryan-lang/tides/astronomy"
)

var TEST_DATE = time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC)

const VAL_TOLERANCE = 0.000001

func TestFUnity(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.0
	result := fUnity(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fUnity() failed, expected %f but got %f", expected, result)
	}
}

func TestFMm(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 0.8941124780426585
	result := fMm(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fMm() failed, expected %f but got %f", expected, result)
	}
}

func TestFMf(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.380589713035402
	result := fMf(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fMf() failed, expected %f but got %f", expected, result)
	}
}

func TestFO1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.1571433842381185
	result := fO1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fO1() failed, expected %f but got %f", expected, result)
	}
}

func TestFJ1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.1443477229145904
	result := fJ1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fJ1() failed, expected %f but got %f", expected, result)
	}
}

func TestFOO1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.647183902791902
	result := fOO1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fOO1() failed, expected %f but got %f", expected, result)
	}
}

func TestFM2(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 0.9698615012436436
	result := fM2(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fM2() failed, expected %f but got %f", expected, result)
	}
}

func TestFK1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.0973117472647704
	result := fK1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fK1() failed, expected %f but got %f", expected, result)
	}
}

func TestFL2(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.147067972436108
	result := fL2(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fL2() failed, expected %f but got %f", expected, result)
	}
}

func TestFK2(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.2608955649232407
	result := fK2(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fK2() failed, expected %f but got %f", expected, result)
	}
}

func TestFM1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 1.5305197393171937
	result := fM1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fM1() failed, expected %f but got %f", expected, result)
	}
}

func TestFModd(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	n := 3.0
	expected := 0.9551346058944368
	result := fModd(a, n)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("fModd() failed, expected %f but got %f", expected, result)
	}
}

func TestUZero(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 0.0
	result := uZero(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uZero() failed, expected %f but got %f", expected, result)
	}
}

func TestUMf(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := -11.423502593010198
	result := uMf(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uMf() failed, expected %f but got %f", expected, result)
	}
}

func TestUO1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 5.094839139355827
	result := uO1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uO1() failed, expected %f but got %f", expected, result)
	}
}

func TestUJ1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := -6.3286634536543716
	result := uJ1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uJ1() failed, expected %f but got %f", expected, result)
	}
}

func TestUOO1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := -17.75216604666457
	result := uOO1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uOO1() failed, expected %f but got %f", expected, result)
	}
}

func TestUM2(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := -1.2338243142985448
	result := uM2(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uM2() failed, expected %f but got %f", expected, result)
	}
}

func TestUK1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := -4.503444637802829
	result := uK1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uK1() failed, expected %f but got %f", expected, result)
	}
}

func TestUL2(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := 15.623387911836264
	result := uL2(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uL2() failed, expected %f but got %f", expected, result)
	}
}

func TestUK2(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := -9.496278929046412
	result := uK2(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uK2() failed, expected %f but got %f", expected, result)
	}
}

func TestUM1(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	expected := -35.33609082266406
	result := uM1(a)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uM1() failed, expected %f but got %f", expected, result)
	}
}

func TestUModd(t *testing.T) {
	a := &astro.Astro{Time: TEST_DATE}
	n := 3.0
	expected := -1.8507364714478172
	result := uModd(a, n)
	if math.Abs(result-expected) > VAL_TOLERANCE {
		t.Errorf("uModd() failed, expected %f but got %f", expected, result)
	}
}
