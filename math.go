package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func src2fbxXYZ(v mgl32.Vec3) mgl64.Vec3 {
	return mgl64.Vec3{float64(v[0]), float64(v[2]), float64(v[1])}
}

func src2fbxUV(v mgl32.Vec2) mgl64.Vec2 {
	return mgl64.Vec2{float64(v[0]), -float64(v[1])}
}
