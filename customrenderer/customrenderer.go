package customrenderer

import (
	"image"
	"image/color"
	"math"

	"github.com/6palace/CustomRenderer/util"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
)

const MinInt = -int(^uint(0)>>1) - 1

// CustomRenderer is a custom 3d render type
type CustomRenderer struct {
	I    *image.NRGBA
	ZBuf [][]int
	// drawChan chan CustomRendererPixel
}

// NewRenderer returns a new renderer
func NewRenderer(width, height int) *CustomRenderer {
	black := color.NRGBA{0, 0, 0, 255}
	res := CustomRenderer{I: util.InitImage(width, height, black)}
	// initialize z-buffer
	res.ZBuf = make([][]int, width)
	for i := range res.ZBuf {
		res.ZBuf[i] = make([]int, height)
		for j := range res.ZBuf[i] {
			res.ZBuf[i][j] = MinInt
		}
	}
	return &res
}

func barycentric(p0, p1, p2, test image.Point) r3.Vector {
	//find cross of (01x 02x t0x) X (01y 02y t0y)
	xPoint := r3.Vector{
		X: float64(p2.X - p0.X),
		Y: float64(p1.X - p0.X),
		Z: float64(p0.X - test.X),
	}
	yPoint := r3.Vector{
		X: float64(p2.Y - p0.Y),
		Y: float64(p1.Y - p0.Y),
		Z: float64(p0.Y - test.Y),
	}
	res := yPoint.Cross(xPoint)
	if math.Abs(res.Z) < 1 {
		return r3.Vector{X: -1, Y: 1, Z: 1}
	} else {
		return r3.Vector{
			X: 1 - (res.X+res.Y)/res.Z,
			Y: res.Y / res.Z,
			Z: res.X / res.Z,
		}
	}
}

func mapVts(vt []r2.Point, barycentricCoords r3.Vector, texture image.Image, intensity float64) color.Color {
	textBounds := texture.Bounds()
	mappedX := vt[0].X*barycentricCoords.X + vt[1].X*barycentricCoords.Y + vt[2].X*barycentricCoords.Z
	mappedX = mappedX * float64(textBounds.Dx())
	// flip coordinate system around vertically
	mappedY := vt[0].Y*barycentricCoords.X + vt[1].Y*barycentricCoords.Y + vt[2].Y*barycentricCoords.Z
	mappedY = (1 - mappedY) * float64(textBounds.Dx())
	origR, origG, origB, origA := texture.At(int(mappedX), int(mappedY)).RGBA()
	newR, newG, newB := uint8(uint32(float64(origR)*intensity)>>8), uint8(uint32(float64(origG)*intensity)>>8), uint8(uint32(float64(origB)*intensity)>>8)
	return color.NRGBA{
		R: newR,
		G: newG,
		B: newB,
		A: uint8(origA >> 8),
	}
}

// TexturedTriangle renders a textured triangle with zBuf
func (rend *CustomRenderer) TexturedTriangle(v []r3.Vector, vt []r2.Point, texture image.Image, intensity float64) {
	boundingBox := findBBox(v[0], v[1], v[2])

	// ops := boundingBox.Dx() * boundingBox.Dy()
	// sem := make(chan bool, ops)
	for i := boundingBox.Min.X; i < boundingBox.Max.X; i++ {
		for j := boundingBox.Min.Y; j < boundingBox.Max.Y; j++ {
			ip0 := image.Point{int(v[0].X), int(v[0].Y)}
			ip1 := image.Point{int(v[1].X), int(v[1].Y)}
			ip2 := image.Point{int(v[2].X), int(v[2].Y)}
			// currently barycentric does not take z axis into account
			barycentricCoords := barycentric(ip0, ip1, ip2, image.Point{i, j})
			if barycentricCoords.X < 0 || barycentricCoords.Y < 0 || barycentricCoords.Z < 0 {
			} else {
				zCoord := v[0].Z*barycentricCoords.X + v[1].Z*barycentricCoords.Y + v[2].Z*barycentricCoords.Z
				rawTextureColor := mapVts(vt, barycentricCoords, texture, intensity)
				rend.blSet(image.Point{i, j}, rawTextureColor, int(zCoord))
			}

		}
	}
}

// VecTriangle renders a filled-in triangle with zBuf
func (rend *CustomRenderer) VecTriangle(vertexes []r3.Vector, color color.Color) {
	v0, v1, v2 := vertexes[0], vertexes[1], vertexes[2]
	boundingBox := findBBox(v0, v1, v2)
	// ops := boundingBox.Dx() * boundingBox.Dy()
	// sem := make(chan bool, ops)
	for i := boundingBox.Min.X; i < boundingBox.Max.X; i++ {
		for j := boundingBox.Min.Y; j < boundingBox.Max.Y; j++ {
			// go func(i, j int) {
			ip0 := image.Point{int(v0.X), int(v0.Y)}
			ip1 := image.Point{int(v1.X), int(v1.Y)}
			ip2 := image.Point{int(v2.X), int(v2.Y)}
			// currently barycentric does not take z axis into account
			barycentricCoords := barycentric(ip0, ip1, ip2, image.Point{i, j})
			if barycentricCoords.X < 0 || barycentricCoords.Y < 0 || barycentricCoords.Z < 0 {
			} else {
				zCoord := v0.Z*barycentricCoords.X + v1.Z*barycentricCoords.Y + v2.Z*barycentricCoords.Z
				rend.blSet(image.Point{i, j}, color, int(zCoord))
			}
			// sem <- true
			// }(i, j)
		}
	}
	// for i := 0; i < ops; i++ {
	// 	<-sem
	// }
}

func findBBox(p0, p1, p2 r3.Vector) image.Rectangle {
	x0 := math.Min(float64(p0.X), float64(p1.X))
	x0 = math.Min(x0, float64(p2.X))
	x1 := math.Max(float64(p0.X), float64(p1.X))
	x1 = math.Max(x1, float64(p2.X))
	y0 := math.Min(float64(p0.Y), float64(p1.Y))
	y0 = math.Min(y0, float64(p2.Y))
	y1 := math.Max(float64(p0.Y), float64(p1.Y))
	y1 = math.Max(y1, float64(p2.Y))
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

// Line draws a single line
func (rend *CustomRenderer) Line(p0, p1 image.Point, color color.Color) {
	dx := math.Abs(float64(p1.X - p0.X))
	dy := math.Abs(float64(p1.Y - p0.Y))
	steep := false
	if dy > dx {
		// swap x and y of both points
		p0.X, p0.Y = p0.Y, p0.X
		p1.X, p1.Y = p1.Y, p1.X
		steep = true
	}
	if p0.X > p1.X {
		// swap starting point
		p0, p1 = p1, p0
	}
	for i := p0.X; i < p1.X; i++ {
		frac := float64(i-p0.X) / float64(p1.X-p0.X)
		j := p0.Y + int(frac*float64(p1.Y-p0.Y))
		if steep {
			rend.blSet(image.Point{j, i}, color, 0)
		} else {
			rend.blSet(image.Point{i, j}, color, 0)
		}
	}
}

func (rend *CustomRenderer) blSet(p image.Point, val color.Color, z int) {
	realY := rend.I.Bounds().Dy() - p.Y
	if z > rend.ZBuf[p.X][p.Y] {
		rend.ZBuf[p.X][p.Y] = z
		rend.I.Set(p.X, realY, val)
	}
}
