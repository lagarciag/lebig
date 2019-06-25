package lebig_test

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/lagarciag/lebig"

)

const globalRepeat = 1000

func TestMain(t *testing.M) {
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)
	//seed = int64(1541178078900106424)
	fmt.Println(" >>>>>>>>>>>>>> SEED: ", seed)
	v := t.Run()
	os.Exit(v)

}

func TestSetBytes(t *testing.T) {
	t.Parallel()
	t.Log(t.Name())
	t.Log("repeats: ", globalRepeat)
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 1
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = byte(rand.Intn(253)) + 1
		}
		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		outBytes := anInt.Bytes()

		randBytes = lebig.RemoveMostSignificantZeroesFromBytes(randBytes)

		checkSlices(t, randBytes, outBytes, x)
	}
}

func BenchmarkSetBytes(b *testing.B) {

	sizeInBytes := 512
	theSlice := make([]byte, sizeInBytes)

	for i := range theSlice {
		theSlice[i] = byte(i)
	}
	anInt := lebig.Int{}
	for n := 0; n < b.N; n++ {
		anInt.SetBytes(theSlice)
	}
}

func BenchmarkSetBytesBigInt(b *testing.B) {

	sizeInBytes := 512
	theSlice := make([]byte, sizeInBytes)

	for i := range theSlice {
		theSlice[i] = byte(i)
	}
	anInt := big.Int{}
	for n := 0; n < b.N; n++ {
		anInt.SetBytes(theSlice)
	}
}

func TestSetBytesSpec1(t *testing.T) {
	t.Parallel()
	t.Log(t.Name())
	inputBytes := []byte{0, 186, 66, 17, 232, 6, 170, 143, 86, 147}
	anInt := lebig.Int{}
	anInt.SetBytes(inputBytes)
	outBytes := anInt.Bytes()

	if !reflect.DeepEqual(inputBytes, outBytes) {
		t.Error("not Equal")
		t.Error(inputBytes)
		t.Error(outBytes)
		t.FailNow()
	}
}

func TestSetBytesZeroLow(t *testing.T) {
	t.Parallel()
	t.Log(t.Name())
	t.Log("repeats: ", globalRepeat)
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 2
		zeroLowBytes := rand.Intn(sizeInBytes - 1)
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = byte(rand.Intn(254))
		}

		for i := range randBytes {
			if i <= zeroLowBytes {
				randBytes[i] = 0
			}
		}

		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		outBytes := anInt.Bytes()

		randBytes = lebig.RemoveMostSignificantZeroesFromBytes(randBytes)

		checkSlices(t, randBytes, outBytes, x)

	}
}

func TestSetBytesZeroHigh(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	t.Log("repeats: ", globalRepeat)
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 2
		zeroLowHigh := rand.Intn(sizeInBytes - 1)
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = byte(rand.Intn(254))
		}

		for i := range randBytes {
			if i >= zeroLowHigh {
				randBytes[i] = 0
			}
		}

		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		outBytes := anInt.Bytes()

		randBytes = lebig.RemoveMostSignificantZeroesFromBytes(randBytes)

		checkSlices(t, randBytes, outBytes, x)

	}
}

func TestSetUint64CheckUint(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 1000
	t.Log(t.Name())
	t.Log("repeats: ", globalRepeat)
	for x := 0; x < globalRepeat; x++ {
		buf := make([]byte, 8)
		rand.Read(buf) // Always succeeds, no need to check error
		randUint := binary.LittleEndian.Uint64(buf)

		anInt := lebig.Int{}
		anInt.SetUint64(randUint)
		outInt := anInt.Uint64()

		if randUint != outInt {
			t.Error("not Equal")
			t.Error(randUint)
			t.Error(outInt)
			t.FailNow()
		}
	}
}

func TestShiftLeftSmall(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 1
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = byte(rand.Intn(254))
		}
		//toShift := uint(rand.Intn((sizeInBytes * 8) + 1))
		toShift := uint(rand.Intn(64))

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		anInt.SmallShiftLeft(toShift)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		lebig.ReverseSliceOfBytes(randBytes)
		aBigInt.SetBytes(randBytes)
		aBigInt.Lsh(&aBigInt, toShift)
		newBytes := aBigInt.Bytes()

		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)

	}
}

func TestShiftLeftBig(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 1
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = byte(rand.Intn(254))
		}
		toShift := uint(rand.Intn((sizeInBytes * 8) + 64))

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		anInt.ShiftLeft(toShift)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		lebig.ReverseSliceOfBytes(randBytes)
		aBigInt.SetBytes(randBytes)
		aBigInt.Lsh(&aBigInt, toShift)
		newBytes := aBigInt.Bytes()

		lebig.ReverseSliceOfBytes(newBytes)

		checkSlices(t, newBytes, outBytes, x)
	}
}

func BenchmarkShiftLeft(b *testing.B) {

	sizeInBytes := 512
	sl := uint(88)
	theSlice := make([]byte, sizeInBytes)

	for i := range theSlice {
		theSlice[i] = byte(i)
	}
	anInt := lebig.Int{}
	for n := 0; n < b.N; n++ {
		anInt.SetBytes(theSlice)
		anInt.ShiftLeft(sl)
	}
}

func TestShiftRightSmall(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 1
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = byte(rand.Intn(254))
		}
		//toShift := uint(rand.Intn((sizeInBytes * 8) + 1))
		toShift := uint(rand.Intn(64))

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		anInt.ShiftRight(toShift)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		lebig.ReverseSliceOfBytes(randBytes)
		aBigInt.SetBytes(randBytes)
		aBigInt.Rsh(&aBigInt, toShift)
		newBytes := aBigInt.Bytes()

		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)
	}
}

func TestShiftRightBig(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 1
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = byte(rand.Intn(0xFF))
		}
		toShift := uint(rand.Intn((sizeInBytes * 8)) + 1)

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		anInt.ShiftRight(toShift)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		lebig.ReverseSliceOfBytes(randBytes)
		aBigInt.SetBytes(randBytes)
		aBigInt.Rsh(&aBigInt, toShift)
		newBytes := aBigInt.Bytes()
		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)
	}
}

func TestShiftRightBigWhenZero(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytes := rand.Intn(1000) + 1
		randBytes := make([]byte, sizeInBytes)
		for i := range randBytes {
			randBytes[i] = 0
		}
		toShift := uint(rand.Intn((sizeInBytes * 8)) + 1)

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytes)
		anInt.ShiftRight(toShift)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		lebig.ReverseSliceOfBytes(randBytes)
		aBigInt.SetBytes(randBytes)
		aBigInt.Rsh(&aBigInt, toShift)
		newBytes := aBigInt.Bytes()
		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)
	}
}

func TestAndBytes(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytesOperan := rand.Intn(1000) + 1
		randBytesInital := make([]byte, sizeInBytesOperan)
		for i := range randBytesInital {
			randBytesInital[i] = byte(rand.Intn(254))
		}

		sizeInBytesOperand := rand.Intn(1000) + 1
		randBytesOperand := make([]byte, sizeInBytesOperand)
		for i := range randBytesOperand {
			randBytesOperand[i] = byte(rand.Intn(254))
		}

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytesInital)
		anInt.AndBytes(randBytesOperand)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		aBigIntOperand := big.Int{}
		lebig.ReverseSliceOfBytes(randBytesInital)
		lebig.ReverseSliceOfBytes(randBytesOperand)
		aBigInt.SetBytes(randBytesInital)
		aBigIntOperand.SetBytes(randBytesOperand)
		_ = aBigInt.And(&aBigInt, &aBigIntOperand)
		newBytes := aBigInt.Bytes()

		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)
	}
}

func TestAndBytes1(t *testing.T) {
	t.Parallel()
	t.Log(t.Name())
	bytesInital := []byte{178, 177, 61, 248, 118, 5, 165, 90, 54}

	bytesOperand := []byte{72, 221, 190, 70, 169, 101, 67, 39, 132}

	// set int bytes
	anInt := lebig.Int{}
	anInt.SetBytes(bytesInital)
	anInt.AndBytes(bytesOperand)
	outBytes := anInt.Bytes()

	// set big int bytes
	aBigInt := big.Int{}
	aBigIntOperand := big.Int{}
	lebig.ReverseSliceOfBytes(bytesInital)
	lebig.ReverseSliceOfBytes(bytesOperand)
	aBigInt.SetBytes(bytesInital)
	aBigIntOperand.SetBytes(bytesOperand)
	_ = aBigInt.And(&aBigInt, &aBigIntOperand)
	newBytes := aBigInt.Bytes()

	lebig.ReverseSliceOfBytes(newBytes)
	checkSlices(t, newBytes, outBytes, 0)

}

func TestOrBytes(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 100000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytesOperan := rand.Intn(1000) + 1
		randBytesInital := make([]byte, sizeInBytesOperan)
		for i := range randBytesInital {
			randBytesInital[i] = byte(rand.Intn(254))
		}

		sizeInBytesOperand := rand.Intn(10000) + 1
		randBytesOperand := make([]byte, sizeInBytesOperand)
		for i := range randBytesOperand {
			randBytesOperand[i] = byte(rand.Intn(254))
		}

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytesInital)
		anInt.OrBytes(randBytesOperand)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		aBigIntOperand := big.Int{}
		lebig.ReverseSliceOfBytes(randBytesInital)
		lebig.ReverseSliceOfBytes(randBytesOperand)
		aBigInt.SetBytes(randBytesInital)
		aBigIntOperand.SetBytes(randBytesOperand)
		_ = aBigInt.Or(&aBigInt, &aBigIntOperand)
		newBytes := aBigInt.Bytes()

		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)
	}
}

func TestAndUint64(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 10000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {

		buf := make([]byte, 8)
		rand.Read(buf) // Always succeeds, no need to check error
		randUintInitial := binary.LittleEndian.Uint64(buf)

		buf2 := make([]byte, 8)
		rand.Read(buf2) // Always succeeds, no need to check error
		randUintOperand := binary.LittleEndian.Uint64(buf)

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetUint64(randUintInitial)
		anInt.AndUint64(randUintOperand)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		aBigIntOperand := big.Int{}
		aBigInt.SetUint64(randUintInitial)
		aBigIntOperand.SetUint64(randUintOperand)
		_ = aBigInt.And(&aBigInt, &aBigIntOperand)
		newBytes := aBigInt.Bytes()

		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)
	}
}

func TestIntOrUint64(t *testing.T) {
	t.Parallel()
	//const globalRepeat = 1000
	t.Log(t.Name())
	for x := 0; x < globalRepeat; x++ {
		sizeInBytesOperan := rand.Intn(1000) + 1
		randBytesInital := make([]byte, sizeInBytesOperan)
		for i := range randBytesInital {
			randBytesInital[i] = byte(rand.Intn(254))
		}

		sizeInBytesOperand := rand.Intn(10000) + 1
		randBytesOperand := make([]byte, sizeInBytesOperand)
		for i := range randBytesOperand {
			randBytesOperand[i] = byte(rand.Intn(254))
		}

		// set int bytes
		anInt := lebig.Int{}
		anInt.SetBytes(randBytesInital)
		anInt.OrBytes(randBytesOperand)
		outBytes := anInt.Bytes()

		// set big int bytes
		aBigInt := big.Int{}
		aBigIntOperand := big.Int{}
		lebig.ReverseSliceOfBytes(randBytesInital)
		lebig.ReverseSliceOfBytes(randBytesOperand)
		aBigInt.SetBytes(randBytesInital)
		aBigIntOperand.SetBytes(randBytesOperand)
		_ = aBigInt.Or(&aBigInt, &aBigIntOperand)
		newBytes := aBigInt.Bytes()

		lebig.ReverseSliceOfBytes(newBytes)
		checkSlices(t, newBytes, outBytes, x)
	}
}

func TestSimpleSet2(t *testing.T) {
	t.Parallel()

	theBytes := []byte{182, 157, 18, 73, 149, 160, 239, 154, 183, 63, 80, 239, 0}

	anInt := lebig.Int{}
	anInt.SetBytes(theBytes)

	newBytes := anInt.Bytes()

	theBytes = lebig.RemoveMostSignificantZeroesFromBytes(theBytes)

	if !reflect.DeepEqual(newBytes, theBytes) {
		t.Error(theBytes, len(theBytes))
		t.Error(newBytes, len(newBytes))
	}

}

func TestReverseSliceOfBytes(t *testing.T) {
	in := []byte{1, 2, 3, 4, 5}
	out := []byte{5, 4, 3, 2, 1}
	lebig.ReverseSliceOfBytes(in)
	if !reflect.DeepEqual(out, in) {
		t.Error("not equal", in, out)
	}

}

func checkSlices(t *testing.T, sIn, sOut []byte, count int) {
	if !reflect.DeepEqual(sIn, sOut) {
		diffOk := false
		if len(sOut) == 1 && sOut[0] == 0 {
			if len(sIn) == 0 {
				diffOk = true
			}
		} else if len(sIn) == 1 && sIn[0] == 0 {
			if len(sOut) == 0 {
				diffOk = true
			}
		}

		if !diffOk {
			t.Error("not Equal on repetetion: ", count)
			t.Error("sIn bytes       :", sIn, len(sIn))
			t.Error("sOut bytes      :", sOut, len(sOut))
			t.FailNow()
		}

	}

}
