package polymorphism

import (
	"math"
	"fmt"
)

// Shapes exercise — Polymorphism via interfaces.
//
// You have a Shape interface and several concrete types.
// There are bugs preventing the code from compiling and
// the tests from passing. Fix the concrete types so they
// properly satisfy the Shape interface.
//
// Rules:
// - Do NOT modify the Shape interface or the TotalArea function.
// - Do NOT modify the test file.

// Shape describes any 2-D shape that can report its area and name.
type Shape interface {
	Area() float64
	Name() string
}

// TotalArea returns the sum of areas from any number of shapes.
// This is polymorphism: it works on *any* Shape without knowing the concrete type.
func TotalArea(shapes ...Shape) float64 {
	var total float64
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// Describe returns "ShapeName: area=X.XX" for any shape.
func Describe(s Shape) string {
	return s.Name() + ": area=" + formatFloat(s.Area())
}

func formatFloat(f float64) string {
	// Round to 2 decimal places for consistent output.
	return fmt.Sprintf("%.2f", math.Floor(f*100) / 100)
}

// --- Concrete types below. Fix them. ---

// Circle should satisfy Shape.
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Name() string {
	return "Circle"
}

// Rectangle should satisfy Shape.
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Name() string {
	return "Rectangle"
}

// Triangle should satisfy Shape.
type Triangle struct {
	Base, Height float64
}

func (t Triangle) Area() float64 {
	return t.Base * t.Height / 2
}

func (t Triangle) Name() string {
	return "Triangle"
}
