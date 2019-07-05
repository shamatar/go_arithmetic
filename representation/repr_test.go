package representation

import (
	"fmt"
	"testing"
)

func TestMulImpl(t *testing.T) {
	var a U256
	var b U256
	uint64max := ^uint64(0)
	a[0] = uint64max
	a[1] = uint64max
	a[2] = uint64max
	a[3] = uint64max

	b[0] = uint64max
	b[1] = uint64max
	b[2] = uint64max
	b[3] = uint64max

	fmt.Println(a.String())

	mulRes := a.MulImpl(&b)
	if mulRes.String() != "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000000000000000001" {
		t.Fail()
	}
}

func TestSquareImpl(t *testing.T) {
	var a U256
	var b U256
	uint64max := ^uint64(0)
	a[0] = uint64max
	a[1] = uint64max
	a[2] = uint64max
	a[3] = uint64max

	b[0] = uint64max
	b[1] = uint64max
	b[2] = uint64max
	b[3] = uint64max

	fmt.Println(a.String())

	mulRes := a.MulImpl(&b)
	squareRes := a.SquareImpl()
	if mulRes != squareRes {
		t.Fail()
	}
	if mulRes.String() != "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000000000000000000000000000000000000000000000000001" {
		t.Fail()
	}
}

func TestMulImplByLimb(t *testing.T) {
	var a U256
	var b U256
	uint64max := ^uint64(0)
	a[0] = uint64max
	a[1] = uint64max
	a[2] = uint64max
	a[3] = uint64max

	b[0] = uint64max

	fmt.Println(a.String())
	fmt.Println(b.String())

	mulRes := a.MulImpl(&b)
	fmt.Println(mulRes.String())
	t.Log(mulRes)
	if mulRes.String() != "0x000000000000000000000000000000000000000000000000fffffffffffffffeffffffffffffffffffffffffffffffffffffffffffffffff0000000000000001" {
		t.Fail()
	}

}

func TestMulImplByTwoLimbs(t *testing.T) {
	var a U256
	var b U256
	uint64max := ^uint64(0)
	a[0] = uint64max
	a[1] = uint64max
	a[2] = uint64max
	a[3] = uint64max

	b[0] = uint64max
	b[1] = uint64max

	fmt.Println(a.String())
	fmt.Println(b.String())

	mulRes := a.MulImpl(&b)
	fmt.Println(mulRes.String())
	t.Log(mulRes)
	if mulRes.String() != "0x00000000000000000000000000000000fffffffffffffffffffffffffffffffeffffffffffffffffffffffffffffffff00000000000000000000000000000001" {
		t.Fail()
	}

}

func TestMAC(t *testing.T) {
	uint64max := ^uint64(0)
	a := uint64max
	b := uint64max
	c := uint64(0)
	carry := uint64(0)
	res, carry := mac_with_carry(a, b, c, carry)
	if res != uint64max || carry != 0 {
		t.Fail()
	}
}

func TestMAC2(t *testing.T) {
	uint64max := ^uint64(0)
	a := uint64(0)
	b := uint64max
	c := uint64max
	carry := uint64(0)
	res, carry := mac_with_carry(a, b, c, carry)
	if res != 1 || carry != uint64max-1 {
		t.Fail()
	}
}

func TestMAC3(t *testing.T) {
	uint64max := ^uint64(0)
	a := uint64max
	b := uint64max
	c := uint64max
	carry := uint64(0)
	res, carry := mac_with_carry(a, b, c, carry)
	if res != 0 || carry != uint64max {
		t.Fail()
	}
}
