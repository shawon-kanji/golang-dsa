package main

import (
	"fmt"
	"math"
)

// Shape interface defines the contract for all geometric shapes
type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

// Circle represents a circle with a radius
type Circle struct {
	Radius float64
}

// Area calculates the area of the circle (πr²)
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference of the circle (2πr)
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Name returns a descriptive name for the circle
func (c Circle) Name() string {
	return fmt.Sprintf("Circle (r=%.1f)", c.Radius)
}

// Rectangle represents a rectangle with width and height
type Rectangle struct {
	Width  float64
	Height float64
}

// Area calculates the area of the rectangle (width × height)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter of the rectangle 2(width + height)
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Name returns a descriptive name for the rectangle
func (r Rectangle) Name() string {
	return fmt.Sprintf("Rectangle (%.0fx%.0f)", r.Width, r.Height)
}

// Triangle represents a triangle with three sides
type Triangle struct {
	A, B, C float64 // three sides
}

// Area calculates the area using Heron's formula
// s = (a + b + c) / 2
// Area = √(s(s-a)(s-b)(s-c))
func (t Triangle) Area() float64 {
	s := t.Perimeter() / 2 // semi-perimeter
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Perimeter calculates the perimeter of the triangle (a + b + c)
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// Name returns a descriptive name for the triangle
func (t Triangle) Name() string {
	return fmt.Sprintf("Triangle (%.0f,%.0f,%.0f)", t.A, t.B, t.C)
}

// IsValid checks if the triangle satisfies the triangle inequality
func (t Triangle) IsValid() bool {
	return t.A+t.B > t.C && t.B+t.C > t.A && t.A+t.C > t.B
}

// Square represents a square (embeds Rectangle)
type Square struct {
	Rectangle // embedding Rectangle to demonstrate composition
}

// NewSquare creates a new square with the given side length
func NewSquare(side float64) Square {
	return Square{
		Rectangle: Rectangle{
			Width:  side,
			Height: side,
		},
	}
}

// Name returns a descriptive name for the square
func (s Square) Name() string {
	return fmt.Sprintf("Square (side=%.0f)", s.Width)
}

// Side returns the side length of the square
func (s Square) Side() float64 {
	return s.Width
}

// PrintShapeInfo prints detailed information about a shape
func PrintShapeInfo(s Shape) {
	fmt.Printf("%s - Area: %.2f, Perimeter: %.2f\n", s.Name(), s.Area(), s.Perimeter())
}

// TotalArea calculates the total area of all shapes
func TotalArea(shapes []Shape) float64 {
	var total float64
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// TotalPerimeter calculates the total perimeter of all shapes
func TotalPerimeter(shapes []Shape) float64 {
	var total float64
	for _, s := range shapes {
		total += s.Perimeter()
	}
	return total
}

// LargestShape returns the shape with the largest area
func LargestShape(shapes []Shape) Shape {
	if len(shapes) == 0 {
		return nil
	}

	largest := shapes[0]
	for _, s := range shapes[1:] {
		if s.Area() > largest.Area() {
			largest = s
		}
	}
	return largest
}

// FilterByMinArea returns shapes with area >= minArea
func FilterByMinArea(shapes []Shape, minArea float64) []Shape {
	var result []Shape
	for _, s := range shapes {
		if s.Area() >= minArea {
			result = append(result, s)
		}
	}
	return result
}

// GetShapeType demonstrates type assertions
func GetShapeType(s Shape) string {
	switch v := s.(type) {
	case Circle:
		return fmt.Sprintf("Circle with radius %.2f", v.Radius)
	case Rectangle:
		return fmt.Sprintf("Rectangle with dimensions %.2f x %.2f", v.Width, v.Height)
	case Triangle:
		return fmt.Sprintf("Triangle with sides %.2f, %.2f, %.2f", v.A, v.B, v.C)
	case Square:
		return fmt.Sprintf("Square with side %.2f", v.Side())
	default:
		return "Unknown shape"
	}
}

func main() {
	fmt.Println("=== Shape Calculator Demo ===")
	fmt.Println()

	// Create various shapes
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{A: 3, B: 4, C: 5},
		NewSquare(4),
	}

	// Print info for each shape
	fmt.Println("Individual Shape Details:")
	fmt.Println("--------------------------")
	for _, shape := range shapes {
		PrintShapeInfo(shape)
	}

	fmt.Println()
	fmt.Printf("Total Area: %.2f\n", TotalArea(shapes))
	fmt.Printf("Total Perimeter: %.2f\n", TotalPerimeter(shapes))

	// Find largest shape
	fmt.Println()
	fmt.Println("Largest Shape:")
	fmt.Println("--------------")
	largest := LargestShape(shapes)
	if largest != nil {
		PrintShapeInfo(largest)
	}

	// Filter shapes by minimum area
	fmt.Println()
	fmt.Println("Shapes with area >= 20:")
	fmt.Println("------------------------")
	largeShapes := FilterByMinArea(shapes, 20)
	for _, shape := range largeShapes {
		PrintShapeInfo(shape)
	}

	// Demonstrate type assertions
	fmt.Println()
	fmt.Println("Type Assertions Demo:")
	fmt.Println("----------------------")
	for _, shape := range shapes {
		fmt.Println(GetShapeType(shape))
	}

	// Demonstrate triangle validity check
	fmt.Println()
	fmt.Println("Triangle Validity Check:")
	fmt.Println("------------------------")
	validTriangle := Triangle{A: 3, B: 4, C: 5}
	invalidTriangle := Triangle{A: 1, B: 2, C: 10}
	fmt.Printf("%s - Valid: %v\n", validTriangle.Name(), validTriangle.IsValid())
	fmt.Printf("%s - Valid: %v\n", invalidTriangle.Name(), invalidTriangle.IsValid())

	// Demonstrate polymorphism with a function
	fmt.Println()
	fmt.Println("Polymorphism Demo:")
	fmt.Println("------------------")
	processShape := func(s Shape) {
		fmt.Printf("Processing %s with area %.2f\n", s.Name(), s.Area())
	}

	for _, shape := range shapes {
		processShape(shape)
	}
}
