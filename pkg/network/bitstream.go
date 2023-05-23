package network

import (
	"fmt"
	"reflect"
)

const Version = 3
const StackAllocationSize = 256
const CompressedVecMagnitudeEpsilon = 0.00001

type BitStream struct {
	numberOfBitsUsed      int
	numberOfBitsAllocated int
	readOffset            int
	copyData              bool
	data                  *[]byte
	stackData             [StackAllocationSize]byte
}

type BitStreamOptions struct {
	InitialBytesToAllocate int32
}

func NewBitStream() *BitStream {
	data := new([]byte)
	stackData := [StackAllocationSize]byte{}

	*data = stackData[:]

	numberOfBitsUsed := 0
	numberOfBitsAllocated := StackAllocationSize * 8
	readOffset := 0
	copyData := true

	return &BitStream{
		numberOfBitsUsed:      numberOfBitsUsed,
		numberOfBitsAllocated: numberOfBitsAllocated,
		readOffset:            readOffset,
		copyData:              copyData,
		data:                  data,
		stackData:             stackData,
	}
}

func (bitStream *BitStream) AddBitsAndReallocate(numberOfBitsToWrite int) {
	if numberOfBitsToWrite <= 0 {
		return
	}

	newNumberOfBitsAllocated := numberOfBitsToWrite + bitStream.numberOfBitsUsed

	if numberOfBitsToWrite+bitStream.numberOfBitsUsed > 0 && ((bitStream.numberOfBitsAllocated-1)>>3) < ((newNumberOfBitsAllocated-1)>>3) {

		newNumberOfBitsAllocated = (numberOfBitsToWrite + bitStream.numberOfBitsUsed) * 2
		amountToAllocate := bitsToBytes(newNumberOfBitsAllocated)

		if reflect.DeepEqual(*bitStream.data, bitStream.stackData[:]) {
			if amountToAllocate > StackAllocationSize {
				*bitStream.data = make([]byte, amountToAllocate)
				copy(*bitStream.data, bitStream.stackData[:bitsToBytes(bitStream.numberOfBitsAllocated)])
			}
		} else {
			if bitStream.copyData {
				*bitStream.data = make([]byte, amountToAllocate)
			} else {
				var newData []byte

				bitStream.copyData = true

				if amountToAllocate < StackAllocationSize {
					newData = bitStream.stackData[:]
				}

				if amountToAllocate >= StackAllocationSize {
					newData = make([]byte, amountToAllocate)
				}

				copy(newData, (*bitStream.data)[:bitsToBytes(bitStream.numberOfBitsAllocated)])

				if reflect.DeepEqual(newData, bitStream.stackData[:]) {
					bitStream.numberOfBitsAllocated = bytesToBits(StackAllocationSize)
				} else {
					bitStream.numberOfBitsAllocated = amountToAllocate
				}

				bitStream.data = &newData
			}
		}
	}

	if newNumberOfBitsAllocated > bitStream.numberOfBitsAllocated {
		bitStream.numberOfBitsAllocated = newNumberOfBitsAllocated
	}
}

func (bitStream *BitStream) SetNumberOfBitsAllocated(lengthInBits uint) {
	bitStream.numberOfBitsAllocated = int(lengthInBits)
}

func (bitStream *BitStream) Write(input *[]byte, numberOfBytes int) {
	if numberOfBytes == 0 {
		return
	}

	if (bitStream.numberOfBitsUsed & 7) != 0 {
		bitStream.WriteBits(input, numberOfBytes*8, true)
		return
	}

	bitStream.AddBitsAndReallocate(bytesToBits(numberOfBytes))
	copy((*bitStream.data)[bitsToBytes(bitStream.numberOfBitsUsed):numberOfBytes], *input)
	bitStream.numberOfBitsUsed += bytesToBits(numberOfBytes)
}

func (bitStream *BitStream) WriteBits(input *[]byte, numberOfBitsToWrite int, rightAlignedBits bool) {
	if numberOfBitsToWrite <= 0 {
		return
	}

	bitStream.AddBitsAndReallocate(numberOfBitsToWrite)

	var dataByte byte
	var numberOfBitsUsedMod8 int

	offset := 0

	numberOfBitsUsedMod8 = bitStream.numberOfBitsUsed & 7

	for numberOfBitsToWrite > 0 {
		dataByte = (*input)[offset]

		if numberOfBitsToWrite < 8 && rightAlignedBits {
			dataByte <<= 8 - numberOfBitsToWrite
		}

		if numberOfBitsUsedMod8 == 0 {
			(*bitStream.data)[bitStream.numberOfBitsUsed>>3] = dataByte
		} else {
			(*bitStream.data)[bitStream.numberOfBitsUsed>>3] |= dataByte >> (numberOfBitsUsedMod8)

			if 8-(numberOfBitsUsedMod8) < 8 && 8-(numberOfBitsUsedMod8) < numberOfBitsToWrite {
				(*bitStream.data)[(bitStream.numberOfBitsUsed>>3)+1] = (dataByte << (8 - (numberOfBitsUsedMod8)))
			}
		}

		if numberOfBitsToWrite >= 8 {
			bitStream.numberOfBitsUsed += 8
		} else {
			bitStream.numberOfBitsUsed += numberOfBitsToWrite
		}

		numberOfBitsToWrite -= 8

		offset++
	}
}

func (bitStream *BitStream) Reset() {
	bitStream.numberOfBitsUsed = 0
	bitStream.readOffset = 0
}

func (bitStream *BitStream) PrintBits() {
	if bitStream.numberOfBitsUsed <= 0 {
		fmt.Println("No bits")

		return
	}

	for i := 0; i < bitsToBytes(bitStream.numberOfBitsUsed); i++ {
		var stop int

		if i == (bitStream.numberOfBitsUsed-1)>>3 {
			stop = 8 - (((bitStream.numberOfBitsUsed - 1) & 7) + 1)
		} else {
			stop = 0
		}

		for l := 7; l >= stop; l-- {
			if ((*bitStream.data)[i]>>l)&1 == 1 {
				fmt.Print('1')
				continue
			}

			fmt.Print('0')
		}

		fmt.Print(' ')
	}

	fmt.Print('\n')
}
