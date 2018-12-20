package raspi

import (
	"fmt"
	"os"
	"syscall"
)

const (
	I2C_SLAVE = 0x0703
)

type I2C struct {
	addr uint8
	bus  int
	rc   *os.File
}

func NewI2C(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	v := &I2C{rc: f, bus: bus, addr: addr}
	return v, nil
}

func (v *I2C) GetBus() int {
	return v.bus
}

func (v *I2C) GetAddr() uint8 {
	return v.addr
}

func (v *I2C) write(buf []byte) (int, error) {
	return v.rc.Write(buf)
}

func (v *I2C) WriteBytes(buf []byte) (int, error) {
	return v.write(buf)
}

func (v *I2C) Write(buf []byte) (int, error) {
	return v.write(buf)
}

func (v *I2C) Read(buf []byte) (int, error) {
	return v.rc.Read(buf)
}

func (v *I2C) read(buf []byte) (int, error) {
	return v.rc.Read(buf)
}

func (v *I2C) ReadBytes(buf []byte) (int, error) {
	n, err := v.read(buf)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (v *I2C) Close() error {
	return v.rc.Close()
}

func (v *I2C) ReadRegBytes(reg byte, n int) ([]byte, int, error) {
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return nil, 0, err
	}
	buf := make([]byte, n)
	c, err := v.ReadBytes(buf)
	if err != nil {
		return nil, 0, err
	}
	return buf, c, nil

}

func (v *I2C) ReadRegU8(reg byte) (byte, error) {
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = v.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (v *I2C) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	_, err := v.WriteBytes(buf)
	if err != nil {
		return err
	}
	return nil
}

func (v *I2C) ReadRegU16BE(reg byte) (uint16, error) {
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = v.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	return w, nil
}

func (v *I2C) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := v.ReadRegU16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

func (v *I2C) ReadRegS16BE(reg byte) (int16, error) {
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = v.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	return w, nil
}

func (v *I2C) ReadRegS16LE(reg byte) (int16, error) {
	w, err := v.ReadRegS16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil

}

func (v *I2C) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := v.WriteBytes(buf)
	if err != nil {
		return err
	}
	return nil
}

func (v *I2C) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return v.WriteRegU16BE(reg, w)
}

func (v *I2C) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := v.WriteBytes(buf)
	if err != nil {
		return err
	}
	return nil
}

func (v *I2C) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return v.WriteRegS16BE(reg, w)
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
