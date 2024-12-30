// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package math provides template functions for mathematical operations.
package math

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sync/atomic"

	"git.sr.ht/~mekyt/latex2mathml"
	_math "github.com/gohugoio/hugo/common/math"
	"github.com/spf13/cast"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2target"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

var (
	errMustTwoNumbersError = errors.New("must provide at least two numbers")
	errMustOneNumberError  = errors.New("must provide at least one number")
)

// New returns a new instance of the math-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "math" namespace.
type Namespace struct{}

// Abs returns the absolute value of n.
func (ns *Namespace) Abs(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("the math.Abs function requires a numeric argument")
	}

	return math.Abs(af), nil
}

// Acos returns the arccosine, in radians, of n.
func (ns *Namespace) Acos(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}
	return math.Acos(af), nil
}

// Add adds the multivalued addends n1 and n2 or more values.
func (ns *Namespace) Add(inputs ...any) (any, error) {
	return ns.doArithmetic(inputs, '+')
}

// Asin returns the arcsine, in radians, of n.
func (ns *Namespace) Asin(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}
	return math.Asin(af), nil
}

// Atan returns the arctangent, in radians, of n.
func (ns *Namespace) Atan(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}
	return math.Atan(af), nil
}

// Atan2 returns the arc tangent of n/m, using the signs of the two to determine the quadrant of the return value.
func (ns *Namespace) Atan2(n, m any) (float64, error) {
	afx, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires numeric arguments")
	}
	afy, err := cast.ToFloat64E(m)
	if err != nil {
		return 0, errors.New("requires numeric arguments")
	}
	return math.Atan2(afx, afy), nil
}

// Ceil returns the least integer value greater than or equal to n.
func (ns *Namespace) Ceil(n any) (float64, error) {
	xf, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Ceil operator can't be used with non-float value")
	}

	return math.Ceil(xf), nil
}

// Cos returns the cosine of the radian argument n.
func (ns *Namespace) Cos(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}
	return math.Cos(af), nil
}

// Div divides n1 by n2.
func (ns *Namespace) Div(inputs ...any) (any, error) {
	return ns.doArithmetic(inputs, '/')
}

// Floor returns the greatest integer value less than or equal to n.
func (ns *Namespace) Floor(n any) (float64, error) {
	xf, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Floor operator can't be used with non-float value")
	}

	return math.Floor(xf), nil
}

// Log returns the natural logarithm of the number n.
func (ns *Namespace) Log(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Log operator can't be used with non integer or float value")
	}

	return math.Log(af), nil
}

// Max returns the greater of all numbers in inputs. Any slices in inputs are flattened.
func (ns *Namespace) Max(inputs ...any) (maximum float64, err error) {
	return ns.applyOpToScalarsOrSlices("Max", math.Max, inputs...)
}

// Min returns the smaller of all numbers in inputs. Any slices in inputs are flattened.
func (ns *Namespace) Min(inputs ...any) (minimum float64, err error) {
	return ns.applyOpToScalarsOrSlices("Min", math.Min, inputs...)
}

// Mod returns n1 % n2.
func (ns *Namespace) Mod(n1, n2 any) (int64, error) {
	ai, erra := cast.ToInt64E(n1)
	bi, errb := cast.ToInt64E(n2)

	if erra != nil || errb != nil {
		return 0, errors.New("modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("the number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

// ModBool returns the boolean of n1 % n2.  If n1 % n2 == 0, return true.
func (ns *Namespace) ModBool(n1, n2 any) (bool, error) {
	res, err := ns.Mod(n1, n2)
	if err != nil {
		return false, err
	}

	return res == int64(0), nil
}

// Mul multiplies the multivalued numbers n1 and n2 or more values.
func (ns *Namespace) Mul(inputs ...any) (any, error) {
	return ns.doArithmetic(inputs, '*')
}

// Pi returns the mathematical constant pi.
func (ns *Namespace) Pi() float64 {
	return math.Pi
}

// Pow returns n1 raised to the power of n2.
func (ns *Namespace) Pow(n1, n2 any) (float64, error) {
	af, erra := cast.ToFloat64E(n1)
	bf, errb := cast.ToFloat64E(n2)

	if erra != nil || errb != nil {
		return 0, errors.New("Pow operator can't be used with non-float value")
	}

	return math.Pow(af, bf), nil
}

// Product returns the product of all numbers in inputs. Any slices in inputs are flattened.
func (ns *Namespace) Product(inputs ...any) (product float64, err error) {
	fn := func(x, y float64) float64 {
		return x * y
	}
	return ns.applyOpToScalarsOrSlices("Product", fn, inputs...)
}

// Rand returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0).
func (ns *Namespace) Rand() float64 {
	return rand.Float64()
}

// Round returns the integer nearest to n, rounding half away from zero.
func (ns *Namespace) Round(n any) (float64, error) {
	xf, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Round operator can't be used with non-float value")
	}

	return _round(xf), nil
}

// Sin returns the sine of the radian argument n.
func (ns *Namespace) Sin(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}
	return math.Sin(af), nil
}

// Sqrt returns the square root of the number n.
func (ns *Namespace) Sqrt(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("Sqrt operator can't be used with non integer or float value")
	}

	return math.Sqrt(af), nil
}

// Sub subtracts multivalued.
func (ns *Namespace) Sub(inputs ...any) (any, error) {
	return ns.doArithmetic(inputs, '-')
}

// Sum returns the sum of all numbers in inputs. Any slices in inputs are flattened.
func (ns *Namespace) Sum(inputs ...any) (sum float64, err error) {
	fn := func(x, y float64) float64 {
		return x + y
	}
	return ns.applyOpToScalarsOrSlices("Sum", fn, inputs...)
}

// Tan returns the tangent of the radian argument n.
func (ns *Namespace) Tan(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}
	return math.Tan(af), nil
}

// ToDegrees converts radians into degrees.
func (ns *Namespace) ToDegrees(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}

	return af * 180 / math.Pi, nil
}

// ToRadians converts degrees into radians.
func (ns *Namespace) ToRadians(n any) (float64, error) {
	af, err := cast.ToFloat64E(n)
	if err != nil {
		return 0, errors.New("requires a numeric argument")
	}

	return af * math.Pi / 180, nil
}

// Convert Latex to MathMl.
func (ns *Namespace) Latex2MathMl(s string) (string, error) {
	return latex2mathml.Convert(
		s,
		"http://www.w3.org/1998/Math/MathML",
		"inline",
		0,
	), nil
}

func (ns *Namespace) D2(s string) (string, error) {
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2dagrelayout.DefaultLayout, nil
	}

	N1 := "#eff0eb"
	N2 := "#f3f99d"
	N3 := "#676d91"
	N4 := "#9aedfe"
	N5 := "#282a36"
	N6 := "#676d91"
	N7 := "#282a36"
	B1 := "#f1f1f0"
	B2 := "#9aedfe"
	B3 := "#676d91"
	B4 := "#282a36"
	B5 := "#676d91"
	B6 := "#282a36"
	AA2 := "#57c7ff"
	AA4 := "#676d91"
	AA5 := "#282a36"
	AB4 := "#676d91"
	AB5 := "#282a36"

	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.DarkMauve.ID,
		ThemeOverrides: &d2target.ThemeOverrides{
			N1:  &N1,
			N2:  &N2,
			N3:  &N3,
			N4:  &N4,
			N5:  &N5,
			N6:  &N6,
			N7:  &N7,
			B1:  &B1,
			B2:  &B2,
			B3:  &B3,
			B4:  &B4,
			B5:  &B5,
			B6:  &B6,
			AA2: &AA2,
			AA4: &AA4,
			AA5: &AA5,
			AB4: &AB4,
			AB5: &AB5,
		},
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}
	diagram, _, _ := d2lib.Compile(context.Background(), s, compileOpts, renderOpts)
	out, _ := d2svg.Render(diagram, renderOpts)

	return string(out), nil
}

func (ns *Namespace) applyOpToScalarsOrSlices(opName string, op func(x, y float64) float64, inputs ...any) (result float64, err error) {
	var i int
	var hasValue bool
	for _, input := range inputs {
		var values []float64
		var isSlice bool
		values, isSlice, err = ns.toFloatsE(input)
		if err != nil {
			err = fmt.Errorf("%s operator can't be used with non-float values", opName)
			return
		}
		hasValue = hasValue || len(values) > 0 || isSlice
		for _, value := range values {
			i++
			if i == 1 {
				result = value
				continue
			}
			result = op(result, value)
		}
	}

	if !hasValue {
		err = errMustOneNumberError
		return
	}
	return
}

func (ns *Namespace) toFloatsE(v any) ([]float64, bool, error) {
	vv := reflect.ValueOf(v)
	switch vv.Kind() {
	case reflect.Slice, reflect.Array:
		var floats []float64
		for i := 0; i < vv.Len(); i++ {
			f, err := cast.ToFloat64E(vv.Index(i).Interface())
			if err != nil {
				return nil, true, err
			}
			floats = append(floats, f)
		}
		return floats, true, nil
	default:
		f, err := cast.ToFloat64E(v)
		if err != nil {
			return nil, false, err
		}
		return []float64{f}, false, nil
	}
}

func (ns *Namespace) doArithmetic(inputs []any, operation rune) (value any, err error) {
	if len(inputs) < 2 {
		return nil, errMustTwoNumbersError
	}
	value = inputs[0]
	for i := 1; i < len(inputs); i++ {
		value, err = _math.DoArithmetic(value, inputs[i], operation)
		if err != nil {
			return
		}
	}
	return
}

var counter uint64

// Counter increments and returns a global counter.
// This was originally added to be used in tests where now.UnixNano did not
// have the needed precision (especially on Windows).
// Note that given the parallel nature of Hugo, you cannot use this to get sequences of numbers,
// and the counter will reset on new builds.
// <docsmeta>{"identifiers": ["now.UnixNano"] }</docsmeta>
func (ns *Namespace) Counter() uint64 {
	return atomic.AddUint64(&counter, uint64(1))
}
