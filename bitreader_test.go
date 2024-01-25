package bitreader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadExeedGivenByteArray(t *testing.T) {
	data := []byte{0xff, 0xff}
	expectedError := "read exceeds the given byte array"
	br := New(data)
	_, err := br.ReadInt64()
	assert.ErrorContains(t, err, expectedError)
}

func TestIncorrectUevRead(t *testing.T) {
	data := []byte{0x00, 0x0f}
	expectedError := "read exceeds the given byte array"
	br := New(data)
	_, err := br.ReadUev()
	assert.ErrorContains(t, err, expectedError)
}

func TestReadUint8(t *testing.T) {
	data := []byte{0b10100111}
	expectedData := uint8(0b10100111)
	br := New(data)
	r, _ := br.ReadUint8()
	assert.Equal(t, expectedData, r)
}

func TestReadUint16(t *testing.T) {
	data := []byte{0x1a, 0xfc, 0xda}
	expectedData := uint16(0xfcda)
	br := New(data)
	br.Skip(8)
	r, _ := br.ReadUint16()
	assert.Equal(t, expectedData, r)
}

func TestReadUint32(t *testing.T) {
	data := []byte{0xfc, 0xda, 0x5f, 0xaf}
	expectedData := uint32(0xfcda5faf)
	br := New(data)
	r, _ := br.ReadUint32()
	assert.Equal(t, expectedData, r)
}

func TestReadUint64(t *testing.T) {
	data := []byte{0xfc, 0xda, 0x5f, 0xaf, 0xec, 0xda, 0x5f, 0xaf}
	expectedData := uint64(0xfcda5fafecda5faf)
	br := New(data)
	r, _ := br.ReadUint64()
	assert.Equal(t, expectedData, r)
}

func TestReadInt64(t *testing.T) {
	data := []byte{0xfc, 0xda, 0x5f, 0xaf, 0xec, 0xda, 0x5f, 0xaf}
	expectedData := int64(-226763622031138897)
	br := New(data)
	r, _ := br.ReadInt64()
	assert.Equal(t, expectedData, r)
}

func TestReadInt32(t *testing.T) {
	data := []byte{0xfc, 0xda, 0x5f, 0xaf}
	expectedData := int32(-52797521)
	br := New(data)
	r, _ := br.ReadInt32()
	assert.Equal(t, expectedData, r)
}

func TestReadInt16(t *testing.T) {
	data := []byte{0xfc, 0xda}
	expectedData := int16(-806)
	br := New(data)
	r, _ := br.ReadInt16()
	assert.Equal(t, expectedData, r)
}

func TestReadInt8(t *testing.T) {
	data := []byte{0xfc}
	expectedData := int8(-4)
	br := New(data)
	r, _ := br.ReadInt8()
	assert.Equal(t, expectedData, r)
}

func TestReadUev(t *testing.T) {
	testData := [][]byte{
		{0b00001000},
		{0b00000100},
		{0b00000110},
		{0b00000010, 0b00000000},
		{0b00000010, 0b10000000},
		{0b00000011, 0b00000000},
		{0b00000011, 0b10000000},
	}
	expectedData := []uint64{0, 1, 2, 3, 4, 5, 6}

	for i, e := range expectedData {
		reader := New(testData[i])
		reader.Skip(4)
		b, err := reader.ReadUev()
		assert.Nil(t, err)
		assert.Equal(t, e, b)
	}
}

func TestReadEv(t *testing.T) {
	testData := [][]byte{
		{0b00001000},
		{0b00000100},
		{0b00000110},
		{0b00000010, 0b00000000},
		{0b00000010, 0b10000000},
		{0b00000011, 0b00000000},
		{0b00000011, 0b10000000},
	}
	expectedData := []int64{0, 1, -1, 2, -2, 3, -3}

	for i, e := range expectedData {
		reader := New(testData[i])
		reader.Skip(4)
		b, err := reader.ReadEv()
		assert.Nil(t, err)
		assert.Equal(t, e, b)
	}
}
