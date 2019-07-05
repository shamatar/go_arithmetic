package representation

import (
	"fmt"
	"math/big"
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

func TestBN254BaseField(t *testing.T) {
	modulus := big.NewInt(0)
	modulus.SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)
	montBits := uint(256)
	montR := big.NewInt(1)
	montR = big.NewInt(0).Lsh(montR, montBits)
	montR = big.NewInt(0).Mod(montR, modulus)
	montR2 := big.NewInt(0).Mul(montR, montR)
	montR2 = big.NewInt(0).Mod(montR2, modulus)

	var u256Modulus U256
	var u256R U256
	var u256R2 U256
	for i := 0; i < M; i++ {
		u256Modulus[i] = modulus.Uint64()
		u256R[i] = montR.Uint64()
		u256R2[i] = montR2.Uint64()
		modulus = modulus.Rsh(modulus, 64)
		montR = montR.Rsh(montR, 64)
		montR2 = montR2.Rsh(montR2, 64)
	}
	inv := uint64(1)
	for i := 0; i < 63; i++ {
		inv = inv * inv
		inv = inv * u256Modulus[0]
	}
	inv = (^inv) + 1
	var two U256
	two[0] = 2

	var three U256
	three[0] = 3

	params := FieldParams{
		modulus: &u256Modulus,
		montR:   &u256R,
		montR2:  &u256R2,
		montInv: inv,
	}

	fe_two := two.IntoFp(&params)
	fe_three := three.IntoFp(&params)
	fe_two.MulAssign(&fe_three)
	result := fe_two.IntoRepr()
	res := result.(U256)
	if res[0] != 6 {
		t.Fail()
	}
}
