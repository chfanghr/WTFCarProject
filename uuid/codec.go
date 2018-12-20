package uuid

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func FromBytes(input []byte) (u UUID, err error) {
	err = u.UnmarshalBinary(input)
	return
}

func FromBytesOrNil(input []byte) UUID {
	uuid, err := FromBytes(input)
	if err != nil {
		return Nil
	}
	return uuid
}

func FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

func FromStringOrNil(input string) UUID {
	uuid, err := FromString(input)
	if err != nil {
		return Nil
	}
	return uuid
}

func (u UUID) MarshalText() (text []byte, err error) {
	text = []byte(u.String())
	return
}

func (u *UUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 32:
		return u.decodeHashLike(text)
	case 36:
		return u.decodeCanonical(text)
	case 38:
		return u.decodeBraced(text)
	case 41:
		fallthrough
	case 45:
		return u.decodeURN(text)
	default:
		return fmt.Errorf("uuid: incorrect UUID length: %s", text)
	}
}

func (u *UUID) decodeCanonical(t []byte) (err error) {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	src := t[:]
	dst := u[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	return
}

func (u *UUID) decodeHashLike(t []byte) (err error) {
	src := t[:]
	dst := u[:]

	if _, err = hex.Decode(dst, src); err != nil {
		return err
	}
	return
}

func (u *UUID) decodeBraced(t []byte) (err error) {
	l := len(t)

	if t[0] != '{' || t[l-1] != '}' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	return u.decodePlain(t[1 : l-1])
}

func (u *UUID) decodeURN(t []byte) (err error) {
	total := len(t)

	urn_uuid_prefix := t[:9]

	if !bytes.Equal(urn_uuid_prefix, urnPrefix) {
		return fmt.Errorf("uuid: incorrect UUID format: %s", t)
	}

	return u.decodePlain(t[9:total])
}

func (u *UUID) decodePlain(t []byte) (err error) {
	switch len(t) {
	case 32:
		return u.decodeHashLike(t)
	case 36:
		return u.decodeCanonical(t)
	default:
		return fmt.Errorf("uuid: incorrrect UUID length: %s", t)
	}
}

func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != Size {
		err = fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)

	return
}
