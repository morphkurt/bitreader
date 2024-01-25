[![Test Coverage](https://raw.githubusercontent.com/wiki/morphkurt/bitreader/coverage.svg)](https://raw.githack.com/wiki/morphkurt/bitreader/coverage.html)
![CI](https://github.com/morphkurt/bitreader/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/morphkurt/bitreader)](https://goreportcard.com/report/github.com/morphkurt/bitreader)

# Bitreader

Go module to read bits

# Sample usage

```
data := []byte{0x1a, 0xfc, 0xda}

br := New(data)

br.Skip(8) // allow skipping of bits
r, err := br.ReadUint16()
```

# Supported Types

The following types are allowed

`uint8`, `uint16`, `uint32`, `uint64`, `int8`, `int16`, `int32`, `int64`

# Supported Functions

1. Allows to read signed and unsigned Exp-Golomb codes through `ReadUev()` and `ReadEv()` functions
2. Allows to Skip Bits via `Skip(bits int)`
