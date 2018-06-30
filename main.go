package main

import (
	cr "helloWorld/customrenderer"
	"helloWorld/util"
	"image/color"

	"github.com/golang/geo/r3"
)

func customPrint(a string) {
	println(a)
}

func main() {
	width, height := 1600, 1600
	// width, height := 800, 500
	// white := color.RGBA{255, 255, 255, 255}
	// red := color.NRGBA{255, 0, 0, 255}
	// green := color.NRGBA{0, 255, 0, 255}
	// blue := color.NRGBA{0, 0, 255, 255}
	clear := color.NRGBA{0, 0, 0, 0}
	finalRes := cr.NewRenderer(width, height)

	// draw african head
	lightDir := r3.Vector{X: 0, Y: 0, Z: -1.0}
	obj, err := util.NewModel("african_head.obj")
	// texture, err := util.LoadTexture("african_head_diffuse.tga")
	if err != nil {
		panic(err)
	}
	for _, face := range obj.Faces {
		v0, v1, v2 := obj.Verts[face[0]], obj.Verts[face[1]], obj.Verts[face[2]]
		// vt0, vt1, vt2 := obj.VTs[face[0]], obj.Verts[face[1]], obj.Verts[face[2]]
		d1, d2 := v2.Sub(v0), v1.Sub(v0)
		faceNorm := d1.Cross(d2).Normalize()
		intensity := faceNorm.Dot(lightDir) * 255
		uint8tensity := uint8(intensity)
		if intensity > 0 {
			// multiply z values by height to scale from float to int better
			p0 := r3.Vector{
				X: (v0.X + 1.0) * float64(width) / 2,
				Y: (v0.Y + 1.0) * float64(height) / 2,
				Z: v0.Z * float64(width),
			}
			p1 := r3.Vector{
				X: (v1.X + 1.0) * float64(width) / 2,
				Y: (v1.Y + 1.0) * float64(height) / 2,
				Z: v1.Z * float64(width),
			}
			p2 := r3.Vector{
				X: (v2.X + 1.0) * float64(width) / 2,
				Y: (v2.Y + 1.0) * float64(height) / 2,
				Z: v2.Z * float64(width),
			}
			finalRes.VecTriangle(p0, p1, p2, color.NRGBA{uint8tensity, uint8tensity, uint8tensity, 255})
		}
	}
	util.DrawFile(finalRes.I, "output.png")

	zIndexMap := util.InitImage(width, height, clear)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if finalRes.ZBuf[i][j] > cr.MinInt {
				zIndexMap.Set(i, height-j, color.NRGBA{uint8(finalRes.ZBuf[i][j]), 0, 0, 255})
			}
		}
	}
	util.DrawFile(zIndexMap, "zmap.png")

}
