package vmath

func Vec3ToGL(v [3]float32) [3]float32 {
	return [3]float32{v[1], v[2], v[0]}
}
