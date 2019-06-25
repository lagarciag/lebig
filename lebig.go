// +build internal

package lebig

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"math/bits"
)

type Int struct {
	sizeInBytes uint
	sizeInBits  uint
	sizeInWords uint
	abs         []uint64
}

func (this *Int) Words() []uint64 {
	return this.abs
}

func (this *Int) SetBytes(in []byte) {
	_ = in[0]

	in = RemoveMostSignificantZeroesFromBytes(in)
	this.sizeInBytes = uint(len(in))
	lenInWords := sizeInWordsFromBytes(this.sizeInBytes)
	this.sizeInWords = lenInWords
	this.abs = make([]uint64, lenInWords)

	// convert to slice of words
	tmpIn := make([]byte, lenInWords*8)
	copy(tmpIn, in)
	for i := range this.abs {
		this.abs[i] = binary.LittleEndian.Uint64(tmpIn[i*8 : i*8+8])
	}

	leadingZeroes := bits.LeadingZeros64(this.abs[len(this.abs)-1])
	this.sizeInBits = uint(64*len(this.abs)) - uint(leadingZeroes)

	if this.sizeInBits == 0 {
		this.sizeInBits = 1
	}
	this.sizeInWords = uint(len(this.abs))
}

func (this *Int) SetUint64(in uint64) {
	inBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(inBytes, in)
	this.SetBytes(inBytes)
}

func (this *Int) Uint64() uint64 {
	if len(this.abs) > 0 {
		return this.abs[0]
	} else {
		return 0
	}

}

func (this *Int) Bytes() []byte {
	if len(this.abs) == 0 {
		return []byte{}
	}
	this.sizeInBytes = recalcSizeInBytes(this.abs)
	tmpOut := make([]byte, len(this.abs)*8)
	for i, word := range this.abs {
		binary.LittleEndian.PutUint64(tmpOut[i*8:i*8+8], word)
	}

	return RemoveMostSignificantZeroesFromBytes(tmpOut[0:this.sizeInBytes])
}

func (this *Int) SmallShiftRight(sr uint) {
	if len(this.abs) > 0 {
		if this.abs[len(this.abs)-1] > 0 {
			if sr <= this.sizeInBits {
				this.sizeInBits = uint(this.sizeInBits - sr)
				this.sizeInBytes = sizeInBytes(this.sizeInBits)
				this.sizeInWords = sizeInWordsFromBits(this.sizeInBits)
				for i := range this.abs {
					if i == 0 {
						this.abs[i] = this.abs[i] >> sr
					} else {
						sl := 64 - sr
						srWord := this.abs[i] << sl
						this.abs[i-1] |= srWord
						this.abs[i] = this.abs[i] >> sr
					}
				}
				this.abs = this.abs[0:this.sizeInWords]
			} else {
				this.abs = []uint64{}
			}
			this.abs = RemoveMostSignificantZeroesFromWords(this.abs)
		}
	}
}

func (this *Int) SmallShiftLeft(sl uint) {
	if len(this.abs) > 0 {

		// if the most significant word is not > 0 it means that the slice is of size 1 & the value is 0.
		// this is true because all zeroes to the left of the most significant 1 were removed in the SetBytes function
		// if the value to shift is zero, the shifting is unnecessary.
		//if this.abs[len(this.abs)-1] > 0 {
		this.sizeInBits = uint(this.sizeInBits + sl)
		this.sizeInBytes = sizeInBytes(this.sizeInBits)
		this.sizeInWords = sizeInWordsFromBits(this.sizeInBits)

		if uint(len(this.abs)) < this.sizeInWords {
			diffSize := this.sizeInWords - uint(len(this.abs))
			diffSlice := make([]uint64, diffSize)
			this.abs = append(this.abs, diffSlice...)
		}

		this.sizeInWords = uint(len(this.abs))
		this.sizeInBytes = recalcSizeInBytes(this.abs)

		for i := len(this.abs) - 1; i >= 0; i-- {
			if i == len(this.abs)-1 {
				this.abs[i] = this.abs[i] << sl
			} else {
				sr := 64 - sl
				srWord := this.abs[i] >> sr
				this.abs[i+1] |= srWord
				this.abs[i] = this.abs[i] << sl
			}
		}
		//}
	}
}

func (this *Int) ShiftLeft(sl uint) {
	if sl <= 64 {
		this.SmallShiftLeft(sl)
	} else {
		count := sl / 64
		if sl%64 > 0 {
			count++
		}
		nsl := uint(64)
		for i := 0; i < int(count); i++ {
			if sl > 64 {
				sl -= 64
				nsl = 64
			} else {
				nsl = sl
			}
			this.SmallShiftLeft(nsl)
		}
	}
}

func (this *Int) ShiftRight(sl uint) {
	lastValue := 64
	if sl <= 64 {
		this.SmallShiftRight(sl)
	} else {
		count := sl / 64
		if sl%64 != 0 {
			tmpVal := 64 * count
			lastValue = int(sl - tmpVal)
			count++
		}

		nsl := uint(64)
		for i := 0; i < int(count); i++ {
			if i == int(count-1) {
				nsl = uint(lastValue)
			}
			this.SmallShiftRight(nsl)
		}
	}
}

func (this *Int) AndUint64(in uint64) {
	if len(this.abs) > 0 {
		this.abs[0] &= in
		this.abs = RemoveMostSignificantZeroesFromWords(this.abs)
	}
}

func (this *Int) AndBytes(in []byte) {
	if len(in) > 0 {
		intIn := Int{}
		intIn.SetBytes(in)
		intInWords := intIn.Words()
		if len(this.abs) < len(intInWords) {
			for i := range this.abs {
				this.abs[i] &= intInWords[i]
			}

		} else {
			for i := range intInWords {

				this.abs[i] &= intInWords[i]
			}
			this.abs = RemoveMostSignificantZeroesFromWords(this.abs[0:len(intInWords)])
		}
	} else {
		for i := range this.abs {
			this.abs[i] = 0
		}
	}
}

func (this *Int) OrUint64(in uint64) {
	this.abs[0] |= in
	this.abs = RemoveMostSignificantZeroesFromWords(this.abs)
}

func (this *Int) OrBytes(in []byte) {
	if len(in) != 0 {
		intIn := Int{}
		intIn.SetBytes(in)
		intInWords := intIn.Words()

		if len(this.abs) < len(intInWords) {
			for i := range this.abs {
				intInWords[i] |= this.abs[i]
			}
			this.abs = intInWords
		} else {
			for i := range intInWords {
				this.abs[i] |= intInWords[i]
			}
		}
		this.abs = RemoveMostSignificantZeroesFromWords(this.abs)
	}
}

func ReverseSliceOfBytes(in []byte) {
	for i := len(in)/2 - 1; i >= 0; i-- {
		opp := len(in) - 1 - i
		//in[i], in[opp] = bits.Reverse8(in[opp]), bits.Reverse8(in[i])
		in[i], in[opp] = in[opp], in[i]
	}
}

//---------
func tmp() {
	a := big.Int{}
	fmt.Println(a)
}
