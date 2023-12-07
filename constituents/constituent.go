package constituents

import (
	"fmt"
	"math"

	astro "github.com/ryan-lang/tides/astronomy"
)

type ConstituentName string

var (
	// Long Term
	CONSTITUENT_Z0  Constituent = Constituent{"Z0", []float64{0, 0, 0, 0, 0, 0, 0}, uZero, fUnity}
	CONSTITUENT_SA  Constituent = Constituent{"SA", []float64{0, 0, 1, 0, 0, 0, 0}, uZero, fUnity}
	CONSTITUENT_SSA Constituent = Constituent{"SSA", []float64{0, 0, 2, 0, 0, 0, 0}, uZero, fUnity}
	CONSTITUENT_MM  Constituent = Constituent{"MM", []float64{0, 1, 0, -1, 0, 0, 0}, uZero, fMm}
	CONSTITUENT_MF  Constituent = Constituent{"MF", []float64{0, 2, 0, 0, 0, 0, 0}, uMf, fMf}
	// dinurals
	CONSTITUENT_Q1  Constituent = Constituent{"Q1", []float64{1, -2, 0, 1, 0, 0, 1}, uO1, fO1}
	CONSTITUENT_O1  Constituent = Constituent{"O1", []float64{1, -1, 0, 0, 0, 0, 1}, uO1, fO1}
	CONSTITUENT_K1  Constituent = Constituent{"K1", []float64{1, 1, 0, 0, 0, 0, -1}, uK1, fK1}
	CONSTITUENT_J1  Constituent = Constituent{"J1", []float64{1, 2, 0, -1, 0, 0, -1}, uJ1, fJ1}
	CONSTITUENT_M1  Constituent = Constituent{"M1", []float64{1, 0, 0, 0, 0, 0, 1}, uM1, fM1}
	CONSTITUENT_P1  Constituent = Constituent{"P1", []float64{1, 1, -2, 0, 0, 0, 1}, uZero, fUnity}
	CONSTITUENT_S1  Constituent = Constituent{"S1", []float64{1, 1, -1, 0, 0, 0, 0}, uZero, fUnity}
	CONSTITUENT_OO1 Constituent = Constituent{"OO1", []float64{1, 3, 0, 0, 0, 0, -1}, uOO1, fOO1}
	// Semi diurnals
	CONSTITUENT_2N2  Constituent = Constituent{"2N2", []float64{2, -2, 0, 2, 0, 0, 0}, uM2, fM2}
	CONSTITUENT_N2   Constituent = Constituent{"N2", []float64{2, -1, 0, 1, 0, 0, 0}, uM2, fM2}
	CONSTITUENT_NU2  Constituent = Constituent{"NU2", []float64{2, -1, 2, -1, 0, 0, 0}, uM2, fM2}
	CONSTITUENT_M2   Constituent = Constituent{"M2", []float64{2, 0, 0, 0, 0, 0, 0}, uM2, fM2}
	CONSTITUENT_LAM2 Constituent = Constituent{"LAM2", []float64{2, 1, -2, 1, 0, 0, 2}, uM2, fM2}
	CONSTITUENT_L2   Constituent = Constituent{"L2", []float64{2, 1, 0, -1, 0, 0, 2}, uL2, fL2}
	CONSTITUENT_T2   Constituent = Constituent{"T2", []float64{2, 2, -3, 0, 0, 1, 0}, uZero, fUnity}
	CONSTITUENT_S2   Constituent = Constituent{"S2", []float64{2, 2, -2, 0, 0, 0, 0}, uZero, fUnity}
	CONSTITUENT_R2   Constituent = Constituent{"R2", []float64{2, 2, -1, 0, 0, -1, 2}, uZero, fUnity}
	CONSTITUENT_K2   Constituent = Constituent{"K2", []float64{2, 2, 0, 0, 0, 0, 0}, uK2, fK2}
	// Third diurnal
	CONSTITUENT_M3 Constituent = Constituent{"M3", []float64{3, 0, 0, 0, 0, 0, 0}, func(a *astro.Astro) float64 { return uModd(a, 3) }, func(a *astro.Astro) float64 { return fModd(a, 3) }}

	// COMPOUND ===
	CONSTITUENT_MSF CompoundConstituent = NewCompoundConstituent("MSF", []CompoundContituentMember{{CONSTITUENT_S2, 1}, {CONSTITUENT_M2, -1}})
	// Diurnal
	CONSTITUENT_2Q1 CompoundConstituent = NewCompoundConstituent("2Q1", []CompoundContituentMember{{CONSTITUENT_N2, 1}, {CONSTITUENT_J1, -1}})
	CONSTITUENT_RHO CompoundConstituent = NewCompoundConstituent("RHO", []CompoundContituentMember{{CONSTITUENT_NU2, 1}, {CONSTITUENT_K1, -1}})
	// Semi-Diurnal
	CONSTITUENT_MU2  CompoundConstituent = NewCompoundConstituent("MU2", []CompoundContituentMember{{CONSTITUENT_M2, 2}, {CONSTITUENT_S2, -1}})
	CONSTITUENT_2SM2 CompoundConstituent = NewCompoundConstituent("2SM2", []CompoundContituentMember{{CONSTITUENT_S2, 2}, {CONSTITUENT_M2, -1}})
	// Third-Diurnal
	CONSTITUENT_2MK3 CompoundConstituent = NewCompoundConstituent("2MK3", []CompoundContituentMember{{CONSTITUENT_M2, 1}, {CONSTITUENT_O1, 1}})
	CONSTITUENT_MK3  CompoundConstituent = NewCompoundConstituent("MK3", []CompoundContituentMember{{CONSTITUENT_M2, 1}, {CONSTITUENT_K1, 1}})
	// Quarter-Diurnal
	CONSTITUENT_MN4 CompoundConstituent = NewCompoundConstituent("MN4", []CompoundContituentMember{{CONSTITUENT_M2, 1}, {CONSTITUENT_N2, 1}})
	CONSTITUENT_M4  CompoundConstituent = NewCompoundConstituent("M4", []CompoundContituentMember{{CONSTITUENT_M2, 2}})
	CONSTITUENT_MS4 CompoundConstituent = NewCompoundConstituent("MS4", []CompoundContituentMember{{CONSTITUENT_M2, 1}, {CONSTITUENT_S2, 1}})
	CONSTITUENT_S4  CompoundConstituent = NewCompoundConstituent("S4", []CompoundContituentMember{{CONSTITUENT_S2, 2}})
	// Sixth-Diurnal
	CONSTITUENT_M6 CompoundConstituent = NewCompoundConstituent("M6", []CompoundContituentMember{{CONSTITUENT_M2, 3}})
	CONSTITUENT_S6 CompoundConstituent = NewCompoundConstituent("S6", []CompoundContituentMember{{CONSTITUENT_S2, 3}})
	// Eighth-Diurnals
	CONSTITUENT_M8 CompoundConstituent = NewCompoundConstituent("M8", []CompoundContituentMember{{CONSTITUENT_M2, 4}})
)

type (
	Constituent struct {
		Name           string
		Coefficients   []float64
		NodeFactorFunc func(a *astro.Astro) float64
		FormFactorFunc func(a *astro.Astro) float64
	}
	CompoundConstituent struct {
		Name    string
		Members []CompoundContituentMember
	}
	CompoundContituentMember struct {
		Constituent Constituent
		Factor      float64
	}
)

func NewCompoundConstituent(name string, members []CompoundContituentMember) CompoundConstituent {
	// // members have their coefficients multiplied by their factor
	// newMembers := make([]CompoundContituentMember, len(members))
	// copy(newMembers, members)

	// for i, member := range newMembers {
	// 	newCoefficients := make([]float64, len(member.Constituent.Coefficients))
	// 	for j, coefficient := range member.Constituent.Coefficients {
	// 		newCoefficients[j] = coefficient * float64(member.Factor)
	// 	}
	// 	newMembers[i].Constituent.Coefficients = newCoefficients
	// }

	return CompoundConstituent{
		Name:    name,
		Members: members,
	}
}

func (c *Constituent) GetName() string {
	return c.Name
}

func (c *Constituent) Speed(a *astro.Astro) float64 {
	_, astroSpeeds := DoodsonNumbers(a)
	return dotArray(c.Coefficients, astroSpeeds)
}

func (c *Constituent) Value(a *astro.Astro) float64 {
	astroValues, _ := DoodsonNumbers(a)
	return dotArray(c.Coefficients, astroValues)
}

// u
func (c *Constituent) NodeFactor(a *astro.Astro) float64 {
	return c.NodeFactorFunc(a)
}

// f
func (c *Constituent) FormFactor(a *astro.Astro) float64 {
	return c.FormFactorFunc(a)
}

func (c *CompoundConstituent) GetName() string {
	return c.Name
}

func (c *CompoundConstituent) Speed(a *astro.Astro) float64 {
	speed := 0.0
	for _, member := range c.Members {
		speed += member.Constituent.Speed(a) * member.Factor
	}
	return speed
}

func (c *CompoundConstituent) Value(a *astro.Astro) float64 {
	value := 0.0
	for _, member := range c.Members {
		value += member.Constituent.Value(a) * member.Factor
	}
	return value
}

// u
func (c *CompoundConstituent) NodeFactor(a *astro.Astro) float64 {
	nodeFactor := 0.0
	for _, member := range c.Members {
		nodeFactor += member.Constituent.NodeFactor(a) * member.Factor
	}
	return nodeFactor
}

// f
func (c *CompoundConstituent) FormFactor(a *astro.Astro) float64 {
	var f []float64
	for _, member := range c.Members {
		f = append(f, math.Pow(member.Constituent.FormFactor(a), math.Abs(member.Factor)))
	}
	product := 1.0
	for _, value := range f {
		product *= value
	}
	return product
}

func DoodsonNumbers(a *astro.Astro) ([]float64, []float64) {
	thsA, thsS := a.EquilibriumArgument()
	sA, sS := a.LunarLongitude()
	hA, hS := a.SolarLongitude()
	pA, pS := a.LunarPerigee()
	nA, nS := a.LunarNode()
	ppA, ppS := a.SolarPerigee()
	angle90A, angle90S := a.FixedAngle(90)
	return []float64{thsA, sA, hA, pA, nA, ppA, angle90A}, []float64{thsS, sS, hS, pS, nS, ppS, angle90S}
}

func dotArray(a, b []float64) float64 {
	if len(a) != len(b) {
		// Handle the error as appropriate
		fmt.Println("Error: Arrays must be of the same length")
		return 0
	}

	result := 0.0
	for i := range a {
		result += a[i] * b[i]
	}
	return result
}
