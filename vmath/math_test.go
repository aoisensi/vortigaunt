package vmath

import (
	"reflect"
	"testing"
)

func TestVecToGL(t *testing.T) {
	v := [3]float32{1, 2, 3}
	expected := [3]float32{2, 3, 1}
	result := VecToGL(v)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("VecToGL(%v) = %v; want %v", v, result, expected)
	}
}

func TestQuatToGL(t *testing.T) {
	q := [4]float32{1, 2, 3, 4}
	expected := [4]float32{2, 3, 1, 4}
	result := QuatToGL(q)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("QuatToGL(%v) = %v; want %v", q, result, expected)
	}
}

func TestVecMulScalar(t *testing.T) {
	v := [3]float32{1, 2, 3}
	s := float32(2)
	expected := [3]float32{2, 4, 6}
	result := VecMulScalar(v, s)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("VecMulScalar(%v, %v) = %v; want %v", v, s, result, expected)
	}
}

func TestFtoD3(t *testing.T) {
	f := [3]float32{1, 2, 3}
	expected := [3]float64{1, 2, 3}
	result := FtoD3(f)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FtoD3(%v) = %v; want %v", f, result, expected)
	}
}

func TestFtoD4(t *testing.T) {
	f := [4]float32{1, 2, 3, 4}
	expected := [4]float64{1, 2, 3, 4}
	result := FtoD4(f)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("FtoD4(%v) = %v; want %v", f, result, expected)
	}
}

func TestMakeTranslateMat(t *testing.T) {
	tVec := [3]float32{1, 2, 3}
	expected := [4][4]float32{
		{1, 0, 0, 1},
		{0, 1, 0, 2},
		{0, 0, 1, 3},
		{0, 0, 0, 1},
	}
	result := MakeTranslateMat(tVec)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MakeTranslateMat(%v) = %v; want %v", tVec, result, expected)
	}
}

func TestMakeRotateMat(t *testing.T) {
	q := [4]float32{0, 0, 0, 1}
	expected := [4][4]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	result := MakeRotateMat(q)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MakeRotateMat(%v) = %v; want %v", q, result, expected)
	}
}

func TestMultiplyMat(t *testing.T) {
	a := [4][4]float32{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	b := [4][4]float32{
		{17, 18, 19, 20},
		{21, 22, 23, 24},
		{25, 26, 27, 28},
		{29, 30, 31, 32},
	}
	expected := [4][4]float32{
		{250, 260, 270, 280},
		{618, 644, 670, 696},
		{986, 1028, 1070, 1112},
		{1354, 1412, 1470, 1528},
	}
	result := MultiplyMat(a, b)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MultiplyMat(%v, %v) = %v; want %v", a, b, result, expected)
	}
}

func TestInverseMat(t *testing.T) {
	m := [4][4]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	expected := [4][4]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	result := InverseMat(m)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InverseMat(%v) = %v; want %v", m, result, expected)
	}
}
