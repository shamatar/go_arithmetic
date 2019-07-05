package representation

import (
	"fmt"
	"math/bits"
	"strings"
)

type Representation interface {
	num_libs() uint64
	sub_noborrow(Representation)
	add_nocarry(Representation)
	num_bits() uint32
	is_zero() bool
	is_odd() bool
	is_even() bool
	div2()
	// shr(uint32)
	mul2()
	// shl(uint32)
	mont_mul_assign(Representation, Representation, uint64)
	mont_square(Representation, uint64)
	into_normal_repr(Representation, uint64) Representation
}

// M ...
const M = 4

type U256 [M]uint64
type U256MulResult [2 * M]uint64

func (repr *U256) String() string {
	var sb strings.Builder
	sb.WriteString("0x")
	for i := M - 1; i >= 0; i-- {
		fmt.Fprintf(&sb, "%016x", repr[i])
	}
	return sb.String()
}

func (repr *U256MulResult) String() string {
	var sb strings.Builder
	sb.WriteString("0x")
	for i := 2*M - 1; i >= 0; i-- {
		fmt.Fprintf(&sb, "%016x", repr[i])
	}
	return sb.String()
}

func (repr *U256) asVec() []uint64 {
	return repr[:]
}

func (repr *U256) num_libs() uint64 {
	return M
}

func (repr *U256) sub_noborrow(oth Representation) {
	other, _ := oth.(*U256)
	borrow := uint64(0)
	for i := 0; i < M; i++ {
		repr[i], borrow = bits.Sub64(repr[i], other[i], borrow)
	}
}

func (repr *U256) add_nocarry(oth Representation) {
	other, _ := oth.(*U256)
	carry := uint64(0)
	for i := 0; i < M; i++ {
		repr[i], carry = bits.Add64(repr[i], other[i], carry)
	}
}

func (repr *U256) num_bits() uint32 {
	num := uint32(M * 64)
	for i := M - 1; i >= 0; i-- {
		leading := bits.LeadingZeros64(repr[i])
		num -= uint32(leading)
		if leading != 64 {
			return num
		}
	}
	return num
}

func (repr *U256) is_zero() bool {
	for i := 0; i < M; i++ {
		if repr[i] != 0 {
			return false
		}
	}
	return true
}

func (repr *U256) is_odd() bool {
	return repr[0]&1 == 1
}

func (repr *U256) is_even() bool {
	return repr[0]&1 == 0
}

func (repr *U256) div2() {
	t := uint64(0)
	t2 := uint64(0)
	for i := M - 1; i >= 0; i-- {
		t2 = repr[i] << 63
		repr[i] >>= 1
		repr[i] |= t
		t = t2
	}
}

func (repr *U256) mul2() {
	t := uint64(0)
	t2 := uint64(0)
	for i := 0; i < M; i++ {
		t2 = repr[i] >> 63
		repr[i] <<= 1
		repr[i] |= t
		t = t2
	}
}

func adc(a, b, carry uint64) (uint64, uint64) {
	res, res_carry := bits.Add64(a, b, 0)
	res, carry_carry := bits.Add64(res, carry, 0)
	return res, res_carry + carry_carry
}

func mac_with_carry(a, b, c, carry uint64) (uint64, uint64) {
	m_hi, m_lo := bits.Mul64(b, c)
	sum_ab, sum_ab_carry := bits.Add64(m_lo, a, 0)
	sum_ab_and_carry, m_hi_carry := bits.Add64(sum_ab, carry, 0)
	m_hi, _ = bits.Add64(m_hi, sum_ab_carry, m_hi_carry)

	return sum_ab_and_carry, m_hi
}

func (repr *U256) mont_mul_assign(oth Representation, modul Representation, mont_inv uint64) {
	other, _ := oth.(*U256)
	modulus, _ := modul.(*U256)
	interm := repr.MulImpl(other)
	repr.mont_reduce(interm, modulus, mont_inv)
}

func (repr *U256) mont_square(modul Representation, mont_inv uint64) {
	modulus, _ := modul.(*U256)
	interm := repr.MulImpl(repr)
	repr.mont_reduce(interm, modulus, mont_inv)
}

func (repr *U256) MulImpl(oth Representation) U256MulResult {
	other, _ := oth.(*U256)
	var result U256MulResult
	for k := 0; k < M; k++ {
		carry := uint64(0)
		thisLimb := repr[k]
		for i := 0; i < M; i++ {
			result[k+i], carry = mac_with_carry(result[k+i], thisLimb, other[i], carry)
		}
		result[M+k] = carry
	}
	return result
}

func (repr *U256) mont_reduce(interm U256MulResult, modulus *U256, mont_inv uint64) {
	carry2 := uint64(0)
	for j := 0; j < M; j++ {
		carry := uint64(0)
		k := interm[j] * mont_inv
		for i := 0; i < M; i++ {
			interm[i+j], carry = mac_with_carry(interm[i+j], k, modulus[i], carry)
		}
		interm[M+j], carry = adc(interm[M+j], carry2, carry)
		carry2 = carry
	}
	for i := 0; i < M; i++ {
		repr[i] = interm[M+i]
	}
	repr.reduce(modulus)
}

func (repr *U256) compare(other *U256) int {
	for i := M - 1; i >= 0; i-- {
		if repr[i] > other[i] {
			return 1
		} else if repr[i] < other[i] {
			return -1
		}
	}
	return 0
}

func (repr *U256) isValid(modulus *U256) bool {
	if repr.compare(modulus) >= -1 {
		return false
	}
	return true
}

func (repr *U256) reduce(modulus *U256) {
	if repr.isValid(modulus) {
		repr.sub_noborrow(modulus)
	}
}

func (repr *U256) SquareImpl() U256MulResult {
	var result U256MulResult
	for k := 0; k < M; k++ {
		carry := uint64(0)
		thisLimb := repr[k]
		for i := k + 1; i < M; i++ {
			result[k+i], carry = mac_with_carry(result[k+i], thisLimb, repr[i], carry)
		}
		result[M+k] = carry
	}
	result[2*M-1] = result[2*M-2] >> 63
	for k := 2*M - 2; k >= 2; k-- {
		result[k] = (result[k] << 1) | (result[k-1] >> 63)
	}
	result[1] = result[1] << 1
	carry := uint64(0)
	for k := 0; k < M; k++ {
		thisLimb := repr[k]
		idx := 2 * k
		result[idx], carry = mac_with_carry(result[idx], thisLimb, thisLimb, carry)
		result[idx+1], carry = bits.Add64(result[idx+1], carry, 0)
	}

	return result
}

func (repr *U256) into_normal_repr(modul Representation, mont_inv uint64) Representation {
	modulus, _ := modul.(*U256)
	var interm U256MulResult
	var result U256
	for i := 0; i < M; i++ {
		interm[i] = repr[i]
	}
	result.mont_reduce(interm, modulus, mont_inv)

	return &result
}

type Fp struct {
	repr  Representation
	field *FieldParams
}

type FieldParams struct {
	modulus Representation
	montR   Representation
	montR2  Representation
	montInv uint64
}

func (repr U256) IntoFp(field *FieldParams) Fp {
	modulus := field.modulus.(*U256)
	montR2 := field.montR2.(*U256)
	mont_inv := field.montInv
	repr.mont_mul_assign(montR2, modulus, mont_inv)
	fe := Fp{
		repr:  &repr,
		field: field,
	}

	return fe
}

func (fe *Fp) MulAssign(other *Fp) {
	fe.repr.mont_mul_assign(other.repr, fe.field.modulus, fe.field.montInv)
}

func (fe *Fp) IntoRepr() Representation {
	result := fe.repr.into_normal_repr(fe.field.modulus, fe.field.montInv)
	return result
}
