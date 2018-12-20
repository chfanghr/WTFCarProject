package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

const epochStart = 122192928000000000

type epochFunc func() time.Time
type hwAddrFunc func() (net.HardwareAddr, error)

var (
	global = newRFC4122Generator()

	posixUID = uint32(os.Getuid())
	posixGID = uint32(os.Getgid())
)

func NewV1() (UUID, error) {
	return global.NewV1()
}

func NewV2(domain byte) (UUID, error) {
	return global.NewV2(domain)
}

func NewV3(ns UUID, name string) UUID {
	return global.NewV3(ns, name)
}

func NewV4() (UUID, error) {
	return global.NewV4()
}

func NewV5(ns UUID, name string) UUID {
	return global.NewV5(ns, name)
}

type Generator interface {
	NewV1() (UUID, error)
	NewV2(domain byte) (UUID, error)
	NewV3(ns UUID, name string) UUID
	NewV4() (UUID, error)
	NewV5(ns UUID, name string) UUID
}

type rfc4122Generator struct {
	clockSequenceOnce sync.Once
	hardwareAddrOnce  sync.Once
	storageMutex      sync.Mutex

	rand io.Reader

	epochFunc     epochFunc
	hwAddrFunc    hwAddrFunc
	lastTime      uint64
	clockSequence uint16
	hardwareAddr  [6]byte
}

func newRFC4122Generator() Generator {
	return &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
}

func (g *rfc4122Generator) NewV1() (UUID, error) {
	u := UUID{}

	timeNow, clockSeq, err := g.getClockSequence()
	if err != nil {
		return Nil, err
	}
	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	hardwareAddr, err := g.getHardwareAddr()
	if err != nil {
		return Nil, err
	}
	copy(u[10:], hardwareAddr)

	u.SetVersion(V1)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

func (g *rfc4122Generator) NewV2(domain byte) (UUID, error) {
	u, err := g.NewV1()
	if err != nil {
		return Nil, err
	}

	switch domain {
	case DomainPerson:
		binary.BigEndian.PutUint32(u[:], posixUID)
	case DomainGroup:
		binary.BigEndian.PutUint32(u[:], posixGID)
	}

	u[9] = domain

	u.SetVersion(V2)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

func (g *rfc4122Generator) NewV3(ns UUID, name string) UUID {
	u := newFromHash(md5.New(), ns, name)
	u.SetVersion(V3)
	u.SetVariant(VariantRFC4122)

	return u
}

func (g *rfc4122Generator) NewV4() (UUID, error) {
	u := UUID{}
	if _, err := io.ReadFull(g.rand, u[:]); err != nil {
		return Nil, err
	}
	u.SetVersion(V4)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

func (g *rfc4122Generator) NewV5(ns UUID, name string) UUID {
	u := newFromHash(sha1.New(), ns, name)
	u.SetVersion(V5)
	u.SetVariant(VariantRFC4122)

	return u
}

func (g *rfc4122Generator) getClockSequence() (uint64, uint16, error) {
	var err error
	g.clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		if _, err = io.ReadFull(g.rand, buf); err != nil {
			return
		}
		g.clockSequence = binary.BigEndian.Uint16(buf)
	})
	if err != nil {
		return 0, 0, err
	}

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	timeNow := g.getEpoch()

	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	return timeNow, g.clockSequence, nil
}

func (g *rfc4122Generator) getHardwareAddr() ([]byte, error) {
	var err error
	g.hardwareAddrOnce.Do(func() {
		if hwAddr, err := g.hwAddrFunc(); err == nil {
			copy(g.hardwareAddr[:], hwAddr)
			return
		}

		if _, err = io.ReadFull(g.rand, g.hardwareAddr[:]); err != nil {
			return
		}
		g.hardwareAddr[0] |= 0x01
	})
	if err != nil {
		return []byte{}, err
	}
	return g.hardwareAddr[:], nil
}

func (g *rfc4122Generator) getEpoch() uint64 {
	return epochStart + uint64(g.epochFunc().UnixNano()/100)
}

func newFromHash(h hash.Hash, ns UUID, name string) UUID {
	u := UUID{}
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))

	return u
}

func defaultHWAddrFunc() (net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []byte{}, err
	}
	for _, iface := range ifaces {
		if len(iface.HardwareAddr) >= 6 {
			return iface.HardwareAddr, nil
		}
	}
	return []byte{}, fmt.Errorf("uuid: no HW address found")
}
