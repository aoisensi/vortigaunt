package vmath

func Vec3ToGL(v [3]float32) [3]float32 {
	return [3]float32{v[1], v[2], v[0]}
}

func VecMulScalar(v [3]float32, s float32) [3]float32 {
	return [3]float32{v[0] * s, v[1] * s, v[2] * s}
}
