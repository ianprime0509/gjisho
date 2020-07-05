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
		{"M5,10c5,5 5,0 10 -5C8 2 6.7 8.5 2 5 6 10 15 20 20 15", false, []Drawing{
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				X:     5,
				Y:     10,
				X1:    10,
				Y1:    15,
				X2:    10,
				Y2:    10,
				X3:    15,
				Y3:    5,
			}, Fill: false},
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				X:     15,
				Y:     5,
				X1:    8,
				Y1:    2,
				X2:    6.7,
				Y2:    8.5,
				X3:    2,
				Y3:    5,
			}, Fill: false},
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				X:     2,
				Y:     5,
				X1:    6,
				Y1:    10,
				X2:    15,
				Y2:    20,
				X3:    20,
				Y3:    15,
			}, Fill: false},
		}},
		// The first stroke of kvg:kanji_057f6
		{"m16.186819,29.104494c1.331878,0.571588,2.950798,0.368357,4.340084,0.254039,5.499736,-0.495376,15.592153,-2.095821,21.252632,-2.514985,1.722256,-0.12702,3.478956,-0.203231,5.178248,0.241337", true, []Drawing{
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				// c1.331878,0.571588,2.950798,0.368357,4.340084,0.254039
				X:  16.186819,
				Y:  29.104494,
				X1: 17.518697,
				Y1: 29.676082,
				X2: 19.137617,
				Y2: 29.472851,
				X3: 20.526903,
				Y3: 29.358533,
			}, Fill: false},
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				// c5.499736,-0.495376,15.592153,-2.095821,21.252632,-2.514985
				X:  20.526903,
				Y:  29.358533,
				X1: 26.026639,
				Y1: 28.863157,
				X2: 36.119056,
				Y2: 27.262712,
				X3: 41.779535,
				Y3: 26.843548,
			}, Fill: false},
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				// c1.722256,-0.12702,3.478956,-0.203231,5.178248,0.241337
				X:  41.779535,
				Y:  26.843548,
				X1: 43.501791,
				Y1: 26.716528,
				X2: 45.258491,
				Y2: 26.640317,
				X3: 46.957783,
				Y3: 27.084885,
			}, Fill: false},
			{Shape: Arc{
				Shape:  Shape{LineWidth: 1, R: 1, G: 0, B: 0},
				X:      16.186819,
				Y:      29.104494,
				R:      1.5,
				Angle1: 0,
				Angle2: 2 * math.Pi,
			}, Fill: true},
		}},
		// The sixth stroke of kvg:kanji_05b66
		{"M37.25,46.5c1,0.25,3.75,0.25,5.5-0.25s18.25-4,20-4s2.75,0.75,1,2.25S54.5,53.5,53,54.75", false, []Drawing{
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				// c1,0.25,3.75,0.25,5.5-0.25
				X:  37.25,
				Y:  46.5,
				X1: 38.25,
				Y1: 46.75,
				X2: 41.0,
				Y2: 46.75,
				X3: 42.75,
				Y3: 46.25,
			}, Fill: false},
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				// s18.25-4,20-4
				X:  42.75,
				Y:  46.25,
				X1: 44.5,
				Y1: 45.75,
				X2: 61.0,
				Y2: 42.25,
				X3: 62.75,
				Y3: 42.25,
			}, Fill: false},
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				// s2.75,0.75,1,2.25
				X:  62.75,
				Y:  42.25,
				X1: 64.5,
				Y1: 42.25,
				X2: 65.5,
				Y2: 43.0,
				X3: 63.75,
				Y3: 44.5,
			}, Fill: false},
			{Shape: Curve{
				Shape: Shape{LineWidth: 1},
				// S54.5,53.5,53,54.75
				X:  63.75,
				Y:  44.5,
				X1: 62.0,
				Y1: 46.0,
				X2: 54.5,
				Y2: 53.5,
				X3: 53.0,
				Y3: 54.75,
			}},
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
			t.Errorf("for stroke %q (-want, +got):\n%v", test.stroke, diff)
		}
	}
}
