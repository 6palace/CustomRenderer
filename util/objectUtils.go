package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang/geo/r3"
)

// Model collection of vectors
type Model struct {
	VTs   [][2]float64
	Verts []r3.Vector
	Faces [][3]int
}

func fail(msg, line string, linenumb int) error {
	return fmt.Errorf(msg+" at %d: %s", linenumb, line)
}

// NewModel creates a model from an obj file
func NewModel(filename string) (*Model, error) {
	res := Model{
		Verts: make([]r3.Vector, 0, 1000),
		Faces: make([][3]int, 0, 1000),
	}
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fi)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, " ") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "v":

			x, err := strconv.ParseFloat(fields[1], 64)
			y, err := strconv.ParseFloat(fields[2], 64)
			z, err := strconv.ParseFloat(fields[3], 64)
			if err != nil {
				panic(err)
			}
			res.Verts = append(res.Verts, r3.Vector{
				X: x,
				Y: y,
				Z: z,
			})
		case "f":
			var f [3]int
			for i := 0; i < 3; i++ {
				face := strings.Split(fields[i+1], "/")
				if len(face) != 3 {
					return nil, fail("unsupported face shape (not a triangle)", line, lineNum)
				}
				vi, err := strconv.Atoi(face[0])
				if err != nil {
					return nil, fail("unsupported face vertex index", line, lineNum)
				}
				f[i] = vi - 1
			}
			res.Faces = append(res.Faces, f)
		case "vt":
			x, err := strconv.ParseFloat(fields[1], 64)
			y, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				panic(err)
			}
			res.VTs = append(res.VTs, [2]float64{
				x,
				y,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &res, nil
}
