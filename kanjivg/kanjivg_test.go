package kanjivg

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Shape struct {
	R, G, B   float64
	LineWidth float64
}

type Arc struct {
	Shape
	X, Y, R, Angle1, Angle2 float64
}

type Curve struct {
	Shape
	X, Y, X1, Y1, X2, Y2, X3, Y3 float64
}

type Drawing struct {
	Shape interface{}
	Fill  bool
}

type canvasState struct {
	r, g, b   float64
	x, y      float64
	lineWidth float64
}

type Canvas struct {
	canvasState
	Drawings  []Drawing
	undrawn   []interface{}
	oldStates []canvasState
}

func (c *Canvas) Arc(xc, yc, r, angle1, angle2 float64) {
	c.undrawn = append(c.undrawn, Arc{
		Shape:  c.currentShape(),
		X:      xc,
		Y:      yc,
		R:      r,
		Angle1: angle1,
		Angle2: angle2,
	})
	c.MoveTo(xc+r*math.Cos(angle2), yc+r*math.Sin(angle2))
}

func (c *Canvas) CurveTo(x1, y1, x2, y2, x3, y3 float64) {
	c.undrawn = append(c.undrawn, Curve{
		Shape: c.currentShape(),
		X:     c.x,
		Y:     c.y,
		X1:    x1,
		Y1:    y1,
		X2:    x2,
		Y2:    y2,
		X3:    x3,
		Y3:    y3,
	})
	c.MoveTo(x3, y3)
}

func (c *Canvas) MoveTo(x, y float64) {
	c.x = x
	c.y = y
}

func (c *Canvas) GetCurrentPoint() (x, y float64) {
	return c.x, c.y
}

func (c *Canvas) GetLineWidth() float64 {
	return c.lineWidth
}

func (c *Canvas) SetSourceRGB(r, g, b float64) {
	c.r = r
	c.g = g
	c.b = b
}

func (c *Canvas) currentShape() Shape {
	return Shape{R: c.r, G: c.g, B: c.b, LineWidth: c.lineWidth}
}

func (c *Canvas) Fill() {
	c.draw(true)
}

func (c *Canvas) Stroke() {
	c.draw(false)
}

func (c *Canvas) Save() {
	c.oldStates = append(c.oldStates, c.canvasState)
}

func (c *Canvas) Restore() {
	c.canvasState = c.oldStates[len(c.oldStates)-1]
	c.oldStates = c.oldStates[:len(c.oldStates)-1]
}

func (c *Canvas) draw(fill bool) {
	for _, shape := range c.undrawn {
		c.Drawings = append(c.Drawings, Drawing{shape, fill})
	}
	c.undrawn = nil
}

func TestDrawTo(t *testing.T) {
	tests := []struct {
		stroke    Stroke
		markStart bool
		drawings  []Drawing
	}{
		{"M10,10C15,15,20,15,25,10", false, []Drawing{
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				X:     10,
				Y:     10,
				X1:    15,
				Y1:    15,
				X2:    20,
				Y2:    15,
				X3:    25,
				Y3:    10,
			}, Fill: false},
		}},
		{"M10,10C15,15,20,15,25,10", true, []Drawing{
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				X:     10,
				Y:     10,
				X1:    15,
				Y1:    15,
				X2:    20,
				Y2:    15,
				X3:    25,
				Y3:    10,
			}, Fill: false},
			{Shape: Arc{
				Shape:  Shape{LineWidth: 1, R: 1, G: 0, B: 0},
				X:      10,
				Y:      10,
				R:      1.5,
				Angle1: 0,
				Angle2: 2 * math.Pi,
			}, Fill: true},
		}},
	}

	closeTo := func(x, y float64) bool {
		return math.Abs(x-y) < 0.01
	}

	for _, test := range tests {
		c := new(Canvas)
		c.lineWidth = 1
		test.stroke.DrawTo(c, test.markStart)

		if diff := cmp.Diff(test.drawings, c.Drawings, cmp.Comparer(closeTo)); diff != "" {
			t.Errorf("for stroke %q: %v", test.stroke, diff)
		}
	}
}
