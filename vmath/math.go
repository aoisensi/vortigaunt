package vmath

func VecToGL(v [3]float32) [3]float32 {
	return [3]float32{v[1], v[2], v[0]}
}

func QuatToGL(q [4]float32) [4]float32 {
	return [4]float32{q[1], q[2], q[0], q[3]}
}

func VecMulScalar(v [3]float32, s float32) [3]float32 {
	return [3]float32{v[0] * s, v[1] * s, v[2] * s}
}

func FtoD3(f [3]float32) [3]float64 {
	return [3]float64{float64(f[0]), float64(f[1]), float64(f[2])}
}

func FtoD4(f [4]float32) [4]float64 {
	return [4]float64{float64(f[0]), float64(f[1]), float64(f[2]), float64(f[3])}
}
