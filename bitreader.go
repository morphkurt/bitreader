// Package bitreader takes byte array and provides abilitiy to read bits and various bit encoding
package bitreader

import (
	"errors"
)

type Bitreader struct {
	data      []byte
	byteIndex uint64
	bitIndex  uint8
}

func New(data []byte) *Bitreader {
	return &Bitreader{
		data: data,
	}
}

func (b *Bitreader) Append(data []byte) {
	b.data = append(b.data, data...)
}

// Trims the underlying slice by removing the read bytes
func (b *Bitreader) Trim() {
	b.data = b.data[b.byteIndex:]
	b.byteIndex = 0
}

func (b *Bitreader) ReadBit() (uint8, error) {

	if int(b.byteIndex) >= len(b.data) {
		return 0, errors.New("read exceeds the given byte array")
	}
	result := b.data[b.byteIndex] >> (7 - b.bitIndex) & 1
	b.bitIndex++
	if b.bitIndex == 8 {
		b.byteIndex++
		b.bitIndex = 0
	}
	return result, nil
}

func (b *Bitreader) ReadFlag() (bool, error) {
	v, err := b.ReadBit()
	return v == 1, err
}

func (b *Bitreader) ReadBits(count int) (uint64, error) {
	result := uint64(0)
	for i := 0; i < count; i++ {
		v, err := b.ReadBit()
		if err != nil {
			return 0, err
		}
		result = uint64(v) | result<<1
	}
	return result, nil
}

func (b *Bitreader) ReadUint8() (uint8, error) {
	v, err := b.ReadBits(8)
	return uint8(v), err
}

func (b *Bitreader) ReadUint16() (uint16, error) {
	v, err := b.ReadBits(16)
	return uint16(v), err
}

func (b *Bitreader) ReadUint32() (uint32, error) {
	v, err := b.ReadBits(32)
	return uint32(v), err
}

func (b *Bitreader) ReadUint64() (uint64, error) {
	v, err := b.ReadBits(64)
	return uint64(v), err
}

func (b *Bitreader) Skip(bits int) error {
	_, err := b.ReadBits(bits)
	return err
}

func (b *Bitreader) ReadInt8() (int8, error) {
	v, err := b.ReadUint8()
	return int8(v), err
}

func (b *Bitreader) ReadInt16() (int16, error) {
	v, err := b.ReadUint16()
	return int16(v), err
}

func (b *Bitreader) ReadInt32() (int32, error) {
	v, err := b.ReadUint32()
	return int32(v), err
}

func (b *Bitreader) ReadInt64() (int64, error) {
	v, err := b.ReadUint64()
	return int64(v), err
}

func (b *Bitreader) CurrentIndex() int {
	return int(b.byteIndex)*8 + int(b.bitIndex)
}

func (b *Bitreader) Reverse(count int) error {
	if count > b.CurrentIndex() {
		return errors.New("index reached < 0")
	}
	for i := 0; i < count; i++ {
		if b.bitIndex == 0 {
			b.bitIndex = 7
			b.byteIndex = b.byteIndex - 1
		} else {
			b.bitIndex = b.bitIndex - 1
		}
	}
	return nil
}

// ReadUev reads unsigned exp-golomb codes
func (reader *Bitreader) ReadUev() (uint64, error) {
	countZero := 0
	for true {
		v, err := reader.ReadBit()
		if err != nil {
			return 0, err
		}
		if v == 1 {
			break
		}
		countZero++
	}
	v, err := reader.ReadBits(countZero)
	if err != nil {
		return 0, err
	}
	return (1 << countZero) - 1 + v, nil
}

// ReadEv reads signed exp-golomb codes
func (reader *Bitreader) ReadEv() (int64, error) {
	v, err := reader.ReadUev()
	if err != nil {
		return 0, err
	}
	if v&1 == 1 {
		return int64((v + 1) / 2), nil
	} else {
		return int64(-(v / 2)), nil
	}
}
