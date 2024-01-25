// Package bitreader takes byte array and provides abilitiy to read bits and various bit encoding
package bitreader

import (
	"errors"
)

type BitReader struct {
	data      []byte
	byteIndex uint64
	bitIndex  uint8
}

func New(data []byte) *BitReader {
	return &BitReader{
		data: data,
	}
}

func (b *BitReader) Append(data []byte) {
	b.data = append(b.data, data...)
}

// Trims the underlying slice by removing the read bytes
func (b *BitReader) Trim() {
	b.data = b.data[b.byteIndex:]
	b.byteIndex = 0
}

func (b *BitReader) ReadBit() (uint8, error) {

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

func (b *BitReader) ReadFlag() (bool, error) {
	v, err := b.ReadBit()
	return v == 1, err
}

func (b *BitReader) ReadBits(count int) (uint64, error) {
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

func (b *BitReader) ReadUint8() (uint8, error) {
	v, err := b.ReadBits(8)
	return uint8(v), err
}

func (b *BitReader) ReadUint16() (uint16, error) {
	v, err := b.ReadBits(16)
	return uint16(v), err
}

func (b *BitReader) ReadUint32() (uint32, error) {
	v, err := b.ReadBits(32)
	return uint32(v), err
}

func (b *BitReader) ReadUint64() (uint64, error) {
	v, err := b.ReadBits(64)
	return uint64(v), err
}

func (b *BitReader) Skip(bits int) error {
	_, err := b.ReadBits(bits)
	return err
}

func (b *BitReader) ReadInt8() (int8, error) {
	v, err := b.ReadUint8()
	return int8(v), err
}

func (b *BitReader) ReadInt16() (int16, error) {
	v, err := b.ReadUint16()
	return int16(v), err
}

func (b *BitReader) ReadInt32() (int32, error) {
	v, err := b.ReadUint32()
	return int32(v), err
}

func (b *BitReader) ReadInt64() (int64, error) {
	v, err := b.ReadUint64()
	return int64(v), err
}

func (b *BitReader) CurrentIndex() int {
	return int(b.byteIndex)*8 + int(b.bitIndex)
}

func (b *BitReader) Reverse(count int) error {
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
func (reader *BitReader) ReadUev() (uint64, error) {
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
func (reader *BitReader) ReadEv() (int64, error) {
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
