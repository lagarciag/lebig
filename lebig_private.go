package lebig

import (
	"math/bits"
)

//-------------------------------------------
func sizeInWordsFromBits(sizeInBits uint) uint {
	lenInWords := sizeInBits / uint(64)
	if sizeInBits%64 != 0 {
		lenInWords++
	}
	return lenInWords
}

func sizeInWordsFromBytes(sizeInBytes uint) uint {
	lenInWords := sizeInBytes / uint(8)
	if sizeInBytes%8 != 0 {
		lenInWords++
	}
	return lenInWords
}

func sizeInBytes(sizeInBits uint) uint {
	lenInBytes := sizeInBits / 8
	if sizeInBits%8 != 0 {
		lenInBytes++
	}
	return lenInBytes
}

func recalcSizeInBytes(inWords []uint64) (outSizeInBytes uint) {
	lastLenBits := bits.Len64(inWords[len(inWords)-1])
	lastlenBytes := sizeInBytes(uint(lastLenBits))

	if len(inWords) > 0 {
		if lastLenBits > 0 {
			outSizeInBytes = uint((len(inWords) * 8) - 8 + int(lastlenBytes))
		} else {
			outSizeInBytes = uint(len(inWords) * 8)
		}
	} else {
		outSizeInBytes = 0
	}
	return outSizeInBytes
}

func RemoveMostSignificantZeroesFromBytes(in []byte) (out []byte) {

	removeMostSignificantZeros := true
	removeCounter := 0
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] == 0 && removeMostSignificantZeros {
			removeCounter++
		} else {
			removeMostSignificantZeros = false
		}
	}
	out = in[0 : len(in)-removeCounter]

	if len(out) == 0 {
		return in[0:1]
	}

	return out
}

func RemoveMostSignificantZeroesFromWords(in []uint64) (out []uint64) {
	removeMostSignificantZeros := true
	removeCounter := 0
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] == 0 && removeMostSignificantZeros {
			removeCounter++
		} else {
			removeMostSignificantZeros = false
		}
	}

	out = in[0 : len(in)-removeCounter]
	return out
}
