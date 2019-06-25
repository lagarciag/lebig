// +build !internal

package lebig

import (
	"math/big"
)

type Int struct {
	anInt big.Int
}

func (this *Int) SetBytes(in []byte) {
	_ = in[0]
	newIn := make([]byte, len(in))
	copy(newIn, in)
	ReverseSliceOfBytes(newIn)
	this.anInt.SetBytes(newIn)
}

func (this *Int) SetUint64(in uint64) {
	this.anInt.SetUint64(in)
}

func (this *Int) Uint64() uint64 {
	return this.anInt.Uint64()
}

func (this *Int) Bytes() []byte {
	out := this.anInt.Bytes()
	ReverseSliceOfBytes(out)
	return out
}

func (this *Int) SmallShiftRight(sr uint) {
	this.anInt.Rsh(&this.anInt, sr)
}

func (this *Int) SmallShiftLeft(sl uint) {
	this.anInt.Lsh(&this.anInt, sl)
}

func (this *Int) ShiftLeft(sl uint) {
	this.anInt.Lsh(&this.anInt, sl)
}

func (this *Int) ShiftRight(sl uint) {
	this.anInt.Rsh(&this.anInt, sl)
}

func (this *Int) AndUint64(in uint64) {
	op := big.Int{}
	op.SetUint64(in)
	this.anInt.And(&this.anInt, &op)
}

func (this *Int) newBigIntFromBytes(in []byte) *big.Int {
	newBytes := make([]byte, len(in))
	copy(newBytes, in)
	ReverseSliceOfBytes(newBytes)
	newBigInt := big.Int{}
	newBigInt.SetBytes(newBytes)
	return &newBigInt
}

func (this *Int) AndBytes(in []byte) {
	op := this.newBigIntFromBytes(in)
	this.anInt.And(&this.anInt, op)
}

func (this *Int) OrUint64(in uint64) {
	op := big.Int{}
	op.SetUint64(in)
	this.anInt.Or(&this.anInt, &op)
}

func (this *Int) OrBytes(in []byte) {
	op := this.newBigIntFromBytes(in)
	this.anInt.Or(&this.anInt, op)
}

func ReverseSliceOfBytes(in []byte) {
	for i := len(in)/2 - 1; i >= 0; i-- {
		opp := len(in) - 1 - i
		//in[i], in[opp] = bits.Reverse8(in[opp]), bits.Reverse8(in[i])
		in[i], in[opp] = in[opp], in[i]
	}
}
