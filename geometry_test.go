package pixel_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/faiface/pixel"
	"github.com/stretchr/testify/assert"
)

func TestRect_Resize(t *testing.T) {
	type rectTestTransform struct {
		name string
		f    func(pixel.Rect) pixel.Rect
	}

	// rectangles
	squareAroundOrigin := pixel.R(-10, -10, 10, 10)
	squareAround2020 := pixel.R(10, 10, 30, 30)
	rectangleAroundOrigin := pixel.R(-20, -10, 20, 10)
	rectangleAround2020 := pixel.R(0, 10, 40, 30)

	// resize transformations
	resizeByHalfAroundCenter := rectTestTransform{"by half around center", func(rect pixel.Rect) pixel.Rect {
		return rect.Resized(rect.Center(), rect.Size().Scaled(0.5))
	}}
	resizeByHalfAroundMin := rectTestTransform{"by half around Min", func(rect pixel.Rect) pixel.Rect {
		return rect.Resized(rect.Min, rect.Size().Scaled(0.5))
	}}
	resizeByHalfAroundMax := rectTestTransform{"by half around Max", func(rect pixel.Rect) pixel.Rect {
		return rect.Resized(rect.Max, rect.Size().Scaled(0.5))
	}}
	resizeByHalfAroundMiddleOfLeftSide := rectTestTransform{"by half around middle of left side", func(rect pixel.Rect) pixel.Rect {
		return rect.Resized(pixel.V(rect.Min.X, rect.Center().Y), rect.Size().Scaled(0.5))
	}}
	resizeByHalfAroundOrigin := rectTestTransform{"by half around the origin", func(rect pixel.Rect) pixel.Rect {
		return rect.Resized(pixel.ZV, rect.Size().Scaled(0.5))
	}}

	testCases := []struct {
		input     pixel.Rect
		transform rectTestTransform
		answer    pixel.Rect
	}{
		{squareAroundOrigin, resizeByHalfAroundCenter, pixel.R(-5, -5, 5, 5)},
		{squareAround2020, resizeByHalfAroundCenter, pixel.R(15, 15, 25, 25)},
		{rectangleAroundOrigin, resizeByHalfAroundCenter, pixel.R(-10, -5, 10, 5)},
		{rectangleAround2020, resizeByHalfAroundCenter, pixel.R(10, 15, 30, 25)},

		{squareAroundOrigin, resizeByHalfAroundMin, pixel.R(-10, -10, 0, 0)},
		{squareAround2020, resizeByHalfAroundMin, pixel.R(10, 10, 20, 20)},
		{rectangleAroundOrigin, resizeByHalfAroundMin, pixel.R(-20, -10, 0, 0)},
		{rectangleAround2020, resizeByHalfAroundMin, pixel.R(0, 10, 20, 20)},

		{squareAroundOrigin, resizeByHalfAroundMax, pixel.R(0, 0, 10, 10)},
		{squareAround2020, resizeByHalfAroundMax, pixel.R(20, 20, 30, 30)},
		{rectangleAroundOrigin, resizeByHalfAroundMax, pixel.R(0, 0, 20, 10)},
		{rectangleAround2020, resizeByHalfAroundMax, pixel.R(20, 20, 40, 30)},

		{squareAroundOrigin, resizeByHalfAroundMiddleOfLeftSide, pixel.R(-10, -5, 0, 5)},
		{squareAround2020, resizeByHalfAroundMiddleOfLeftSide, pixel.R(10, 15, 20, 25)},
		{rectangleAroundOrigin, resizeByHalfAroundMiddleOfLeftSide, pixel.R(-20, -5, 0, 5)},
		{rectangleAround2020, resizeByHalfAroundMiddleOfLeftSide, pixel.R(0, 15, 20, 25)},

		{squareAroundOrigin, resizeByHalfAroundOrigin, pixel.R(-5, -5, 5, 5)},
		{squareAround2020, resizeByHalfAroundOrigin, pixel.R(5, 5, 15, 15)},
		{rectangleAroundOrigin, resizeByHalfAroundOrigin, pixel.R(-10, -5, 10, 5)},
		{rectangleAround2020, resizeByHalfAroundOrigin, pixel.R(0, 5, 20, 15)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Resize %v %s", testCase.input, testCase.transform.name), func(t *testing.T) {
			testResult := testCase.transform.f(testCase.input)
			if testResult != testCase.answer {
				t.Errorf("Got: %v, wanted: %v\n", testResult, testCase.answer)
			}
		})
	}
}

func TestMatrix_Unproject(t *testing.T) {
	const delta = 1e-15
	t.Run("for rotated matrix", func(t *testing.T) {
		matrix := pixel.IM.
			Rotated(pixel.ZV, math.Pi/2)
		unprojected := matrix.Unproject(pixel.V(0, 1))
		assert.InDelta(t, unprojected.X, 1, delta)
		assert.InDelta(t, unprojected.Y, 0, delta)
	})
	t.Run("for moved matrix", func(t *testing.T) {
		matrix := pixel.IM.
			Moved(pixel.V(1, 2))
		unprojected := matrix.Unproject(pixel.V(2, 5))
		assert.InDelta(t, unprojected.X, 1, delta)
		assert.InDelta(t, unprojected.Y, 3, delta)
	})
	t.Run("for scaled matrix", func(t *testing.T) {
		matrix := pixel.IM.
			Scaled(pixel.ZV, 2)
		unprojected := matrix.Unproject(pixel.V(2, 4))
		assert.InDelta(t, unprojected.X, 1, delta)
		assert.InDelta(t, unprojected.Y, 2, delta)
	})
	t.Run("for scaled, rotated and moved matrix", func(t *testing.T) {
		matrix := pixel.IM.
			Scaled(pixel.ZV, 2).
			Rotated(pixel.ZV, math.Pi/2).
			Moved(pixel.V(2, 2))
		unprojected := matrix.Unproject(pixel.V(-2, 6))
		assert.InDelta(t, unprojected.X, 2, delta)
		assert.InDelta(t, unprojected.Y, 2, delta)
	})
	t.Run("for rotated and moved matrix", func(t *testing.T) {
		matrix := pixel.IM.
			Rotated(pixel.ZV, math.Pi/2).
			Moved(pixel.V(1, 1))
		unprojected := matrix.Unproject(pixel.V(1, 2))
		assert.InDelta(t, unprojected.X, 1, delta)
		assert.InDelta(t, unprojected.Y, 0, delta)
	})
	t.Run("for projected vertices using all kinds of matrices", func(t *testing.T) {
		namedMatrices := map[string]pixel.Matrix{
			"IM":                        pixel.IM,
			"Scaled":                    pixel.IM.Scaled(pixel.ZV, 0.5),
			"Scaled x 2":                pixel.IM.Scaled(pixel.ZV, 2),
			"Rotated":                   pixel.IM.Rotated(pixel.ZV, math.Pi/4),
			"Moved":                     pixel.IM.Moved(pixel.V(0.5, 1)),
			"Moved 2":                   pixel.IM.Moved(pixel.V(-1, -0.5)),
			"Scaled and Rotated":        pixel.IM.Scaled(pixel.ZV, 0.5).Rotated(pixel.ZV, math.Pi/4),
			"Scaled, Rotated and Moved": pixel.IM.Scaled(pixel.ZV, 0.5).Rotated(pixel.ZV, math.Pi/4).Moved(pixel.V(1, 2)),
			"Rotated and Moved":         pixel.IM.Rotated(pixel.ZV, math.Pi/4).Moved(pixel.V(1, 2)),
		}
		vertices := [...]pixel.Vec{
			pixel.V(0, 0),
			pixel.V(5, 0),
			pixel.V(5, 10),
			pixel.V(0, 10),
			pixel.V(-5, 10),
			pixel.V(-5, 0),
			pixel.V(-5, -10),
			pixel.V(0, -10),
			pixel.V(5, -10),
		}
		for matrixName, matrix := range namedMatrices {
			for _, vertex := range vertices {
				testCase := fmt.Sprintf("for matrix %s and vertex %v", matrixName, vertex)
				t.Run(testCase, func(t *testing.T) {
					projected := matrix.Project(vertex)
					unprojected := matrix.Unproject(projected)
					assert.InDelta(t, vertex.X, unprojected.X, delta)
					assert.InDelta(t, vertex.Y, unprojected.Y, delta)
				})
			}
		}
	})
	t.Run("for singular matrix", func(t *testing.T) {
		matrix := pixel.Matrix{0, 0, 0, 0, 0, 0}
		unprojected := matrix.Unproject(pixel.ZV)
		assert.True(t, math.IsNaN(unprojected.X))
		assert.True(t, math.IsNaN(unprojected.Y))
	})
}

func TestC(t *testing.T) {
	type args struct {
		radius float64
		center pixel.Vec
	}
	tests := []struct {
		name string
		args args
		want pixel.Circle
	}{
		{
			name: "C(): positive radius",
			args: args{radius: 10, center: pixel.ZV},
			want: pixel.Circle{Radius: 10, Center: pixel.ZV},
		},
		{
			name: "C(): zero radius",
			args: args{radius: 0, center: pixel.ZV},
			want: pixel.Circle{Radius: 0, Center: pixel.ZV},
		},
		{
			name: "C(): negative radius",
			args: args{radius: -5, center: pixel.ZV},
			want: pixel.Circle{Radius: -5, Center: pixel.ZV},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pixel.C(tt.args.center, tt.args.radius); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_String(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Circle.String(): positive radius",
			fields: fields{radius: 10, center: pixel.ZV},
			want:   "Circle(Vec(0, 0), 10.00)",
		},
		{
			name:   "Circle.String(): zero radius",
			fields: fields{radius: 0, center: pixel.ZV},
			want:   "Circle(Vec(0, 0), 0.00)",
		},
		{
			name:   "Circle.String(): negative radius",
			fields: fields{radius: -5, center: pixel.ZV},
			want:   "Circle(Vec(0, 0), -5.00)",
		},
		{
			name:   "Circle.String(): irrational radius",
			fields: fields{radius: math.Pi, center: pixel.ZV},
			want:   "Circle(Vec(0, 0), 3.14)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(tt.fields.center, tt.fields.radius)
			if got := c.String(); got != tt.want {
				t.Errorf("Circle.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Norm(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   pixel.Circle
	}{
		{
			name:   "Circle.Norm(): positive radius",
			fields: fields{radius: 10, center: pixel.ZV},
			want:   pixel.C(pixel.ZV, 10),
		},
		{
			name:   "Circle.Norm(): zero radius",
			fields: fields{radius: 0, center: pixel.ZV},
			want:   pixel.C(pixel.ZV, 0),
		},
		{
			name:   "Circle.Norm(): negative radius",
			fields: fields{radius: -5, center: pixel.ZV},
			want:   pixel.C(pixel.ZV, 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(tt.fields.center, tt.fields.radius)
			if got := c.Norm(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Norm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Area(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "Circle.Area(): positive radius",
			fields: fields{radius: 10, center: pixel.ZV},
			want:   100 * math.Pi,
		},
		{
			name:   "Circle.Area(): zero radius",
			fields: fields{radius: 0, center: pixel.ZV},
			want:   0,
		},
		{
			name:   "Circle.Area(): negative radius",
			fields: fields{radius: -5, center: pixel.ZV},
			want:   25 * math.Pi,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(tt.fields.center, tt.fields.radius)
			if got := c.Area(); got != tt.want {
				t.Errorf("Circle.Area() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Moved(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	type args struct {
		delta pixel.Vec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pixel.Circle
	}{
		{
			name:   "Circle.Moved(): positive movement",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{delta: pixel.V(10, 20)},
			want:   pixel.C(pixel.V(10, 20), 10),
		},
		{
			name:   "Circle.Moved(): zero movement",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{delta: pixel.ZV},
			want:   pixel.C(pixel.V(0, 0), 10),
		},
		{
			name:   "Circle.Moved(): negative movement",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{delta: pixel.V(-5, -10)},
			want:   pixel.C(pixel.V(-5, -10), 10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(tt.fields.center, tt.fields.radius)
			if got := c.Moved(tt.args.delta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Moved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Resized(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	type args struct {
		radiusDelta float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pixel.Circle
	}{
		{
			name:   "Circle.Resized(): positive delta",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{radiusDelta: 5},
			want:   pixel.C(pixel.V(0, 0), 15),
		},
		{
			name:   "Circle.Resized(): zero delta",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{radiusDelta: 0},
			want:   pixel.C(pixel.V(0, 0), 10),
		},
		{
			name:   "Circle.Resized(): negative delta",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{radiusDelta: -5},
			want:   pixel.C(pixel.V(0, 0), 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(tt.fields.center, tt.fields.radius)
			if got := c.Resized(tt.args.radiusDelta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Resized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Contains(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	type args struct {
		u pixel.Vec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Circle.Contains(): point on cicles' center",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{u: pixel.ZV},
			want:   true,
		},
		{
			name:   "Circle.Contains(): point offcenter",
			fields: fields{radius: 10, center: pixel.V(5, 0)},
			args:   args{u: pixel.ZV},
			want:   true,
		},
		{
			name:   "Circle.Contains(): point on circumference",
			fields: fields{radius: 10, center: pixel.V(10, 0)},
			args:   args{u: pixel.ZV},
			want:   true,
		},
		{
			name:   "Circle.Contains(): point outside circle",
			fields: fields{radius: 10, center: pixel.V(15, 0)},
			args:   args{u: pixel.ZV},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(tt.fields.center, tt.fields.radius)
			if got := c.Contains(tt.args.u); got != tt.want {
				t.Errorf("Circle.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Union(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	type args struct {
		d pixel.Circle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pixel.Circle
	}{
		{
			name:   "Circle.Union(): overlapping circles",
			fields: fields{radius: 5, center: pixel.ZV},
			args:   args{d: pixel.C(pixel.ZV, 5)},
			want:   pixel.C(pixel.ZV, 5),
		},
		{
			name:   "Circle.Union(): separate circles",
			fields: fields{radius: 1, center: pixel.ZV},
			args:   args{d: pixel.C(pixel.V(0, 2), 1)},
			want:   pixel.C(pixel.V(0, 1), 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(tt.fields.center, tt.fields.radius)
			if got := c.Union(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircle_Intersect(t *testing.T) {
	type fields struct {
		radius float64
		center pixel.Vec
	}
	type args struct {
		d pixel.Circle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pixel.Circle
	}{
		{
			name:   "Circle.Intersect(): intersecting circles",
			fields: fields{radius: 1, center: pixel.ZV},
			args:   args{d: pixel.C(pixel.V(1, 0), 1)},
			want:   pixel.C(pixel.V(0.5, 0), 1),
		},
		{
			name:   "Circle.Intersect(): non-intersecting circles",
			fields: fields{radius: 1, center: pixel.ZV},
			args:   args{d: pixel.C(pixel.V(3, 3), 1)},
			want:   pixel.C(pixel.V(1.5, 1.5), 0),
		},
		{
			name:   "Circle.Intersect(): first circle encompassing second",
			fields: fields{radius: 10, center: pixel.ZV},
			args:   args{d: pixel.C(pixel.V(3, 3), 1)},
			want:   pixel.C(pixel.ZV, 10),
		},
		{
			name:   "Circle.Intersect(): second circle encompassing first",
			fields: fields{radius: 1, center: pixel.V(-1, -4)},
			args:   args{d: pixel.C(pixel.ZV, 10)},
			want:   pixel.C(pixel.ZV, 10),
		},
		{
			name:   "Circle.Intersect(): matching circles",
			fields: fields{radius: 1, center: pixel.ZV},
			args:   args{d: pixel.C(pixel.ZV, 1)},
			want:   pixel.C(pixel.ZV, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pixel.C(
				tt.fields.center,
				tt.fields.radius,
			)
			if got := c.Intersect(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Circle.Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_IntersectCircle(t *testing.T) {
	type fields struct {
		Min pixel.Vec
		Max pixel.Vec
	}
	type args struct {
		c pixel.Circle
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   pixel.Vec
	}{
		{
			name:   "Rect.IntersectCircle(): no overlap",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(50, 50), 1)},
			want:   pixel.ZV,
		},
		{
			name:   "Rect.IntersectCircle(): circle contains rect",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(5, 5), 10)},
			want:   pixel.V(-15, 0),
		},
		{
			name:   "Rect.IntersectCircle(): rect contains circle",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(5, 5), 1)},
			want:   pixel.V(-6, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps bottom-left corner",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(0, 0), 1)},
			want:   pixel.V(1, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps top-left corner",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(0, 10), 1)},
			want:   pixel.V(1, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps bottom-right corner",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(10, 0), 1)},
			want:   pixel.V(-1, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps top-right corner",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(10, 10), 1)},
			want:   pixel.V(-1, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps two corners",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(0, 5), 6)},
			want:   pixel.V(6, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps left edge",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(0, 5), 1)},
			want:   pixel.V(1, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps bottom edge",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(5, 0), 1)},
			want:   pixel.V(0, 1),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps right edge",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(10, 5), 1)},
			want:   pixel.V(-1, 0),
		},
		{
			name:   "Rect.IntersectCircle(): circle overlaps top edge",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(5, 10), 1)},
			want:   pixel.V(0, -1),
		},
		{
			name:   "Rect.IntersectCircle(): edge is tangent of left side",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(-1, 5), 1)},
			want:   pixel.ZV,
		},
		{
			name:   "Rect.IntersectCircle(): edge is tangent of top side",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(5, -1), 1)},
			want:   pixel.ZV,
		},
		{
			name:   "Rect.IntersectCircle(): circle above rectangle",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(5, 12), 1)},
			want:   pixel.ZV,
		},
		{
			name:   "Rect.IntersectCircle(): circle below rectangle",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(5, -2), 1)},
			want:   pixel.ZV,
		},
		{
			name:   "Rect.IntersectCircle(): circle left of rectangle",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(-1, 5), 1)},
			want:   pixel.ZV,
		},
		{
			name:   "Rect.IntersectCircle(): circle right of rectangle",
			fields: fields{Min: pixel.ZV, Max: pixel.V(10, 10)},
			args:   args{c: pixel.C(pixel.V(11, 5), 1)},
			want:   pixel.ZV,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := pixel.Rect{
				Min: tt.fields.Min,
				Max: tt.fields.Max,
			}
			if got := r.IntersectCircle(tt.args.c); got != tt.want {
				t.Errorf("Rect.IntersectCircle() = %v, want %v", got, tt.want)
			}
		})
	}
}
