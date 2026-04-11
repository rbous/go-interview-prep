package polymorphism

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}

func TestCircleArea(t *testing.T) {
	c := Circle{Radius: 5}
	want := math.Pi * 25
	if !almostEqual(c.Area(), want) {
		t.Errorf("Circle.Area() = %f, want %f", c.Area(), want)
	}
}

func TestRectangleArea(t *testing.T) {
	r := Rectangle{Width: 4, Height: 6}
	want := 24.0
	if !almostEqual(r.Area(), want) {
		t.Errorf("Rectangle.Area() = %f, want %f", r.Area(), want)
	}
}

func TestTriangleArea(t *testing.T) {
	tr := Triangle{Base: 10, Height: 5}
	want := 25.0
	if !almostEqual(tr.Area(), want) {
		t.Errorf("Triangle.Area() = %f, want %f", tr.Area(), want)
	}
}

func TestTotalAreaPolymorphism(t *testing.T) {
	c := Circle{Radius: 1}       // pi
	r := Rectangle{Width: 2, Height: 3} // 6
	tr := Triangle{Base: 4, Height: 5}  // 10

	got := TotalArea(c, r, tr)
	want := math.Pi + 6 + 10
	if !almostEqual(got, want) {
		t.Errorf("TotalArea() = %f, want %f", got, want)
	}
}

func TestDescribe(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	got := Describe(r)
	want := "Rectangle: area=12.00"
	if got != want {
		t.Errorf("Describe() = %q, want %q", got, want)
	}
}

func TestShapeInterfaceSatisfaction(t *testing.T) {
	// This test verifies that all types satisfy Shape at compile time.
	var shapes []Shape
	shapes = append(shapes, Circle{Radius: 1})
	shapes = append(shapes, Rectangle{Width: 1, Height: 1})
	shapes = append(shapes, Triangle{Base: 1, Height: 1})

	if len(shapes) != 3 {
		t.Errorf("expected 3 shapes, got %d", len(shapes))
	}

	for _, s := range shapes {
		if s.Name() == "" {
			t.Errorf("shape %T has empty name", s)
		}
		if s.Area() <= 0 {
			t.Errorf("shape %T has non-positive area", s)
		}
	}
}
