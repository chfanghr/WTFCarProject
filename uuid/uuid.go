package uuid

import (
	"bytes"
	"encoding/hex"
)

const Size = 16

type UUID [Size]byte

const (
	_ byte = iota
	V1
	V2
	V3
	V4
	V5
)

const (
	VariantNCS byte = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

const (
	DomainPerson = iota
	DomainGroup
	DomainOrg
)

var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

var Nil = UUID{}

var (
	NamespaceDNS  = Must(FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceURL  = Must(FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceOID  = Must(FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8"))
	NamespaceX500 = Must(FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8"))
)

func Equal(u1 UUID, u2 UUID) bool {
	return bytes.Equal(u1[:], u2[:])
}

func (u UUID) Version() byte {
	return u[6] >> 4
}

func (u UUID) Variant() byte {
	switch {
	case (u[8] >> 7) == 0x00:
		return VariantNCS
	case (u[8] >> 6) == 0x02:
		return VariantRFC4122
	case (u[8] >> 5) == 0x06:
		return VariantMicrosoft
	case (u[8] >> 5) == 0x07:
		fallthrough
	default:
		return VariantFuture
	}
}

func (u UUID) Bytes() []byte {
	return u[:]
}

func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

func (u *UUID) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

func (u *UUID) SetVariant(v byte) {
	switch v {
	case VariantNCS:
		u[8] = (u[8]&(0xff>>1) | (0x00 << 7))
	case VariantRFC4122:
		u[8] = (u[8]&(0xff>>2) | (0x02 << 6))
	case VariantMicrosoft:
		u[8] = (u[8]&(0xff>>3) | (0x06 << 5))
	case VariantFuture:
		fallthrough
	default:
		u[8] = (u[8]&(0xff>>3) | (0x07 << 5))
	}
}

func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}
