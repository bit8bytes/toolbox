package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"errors"
	"fmt"
)

type UUID [16]byte

func New() (UUID, error) {
	var u UUID
	_, err := rand.Read(u[:])
	if err != nil {
		return u, errors.New("failed to create cryptographically secure UUID")
	}

	// Set version 4 (random) and variant bits
	u[6] = (u[6] & 0x0f) | 0x40 // version 4
	u[8] = (u[8] & 0x3f) | 0x80 // variant 10

	return u, nil
}

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func (u *UUID) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok || len(bytes) != 16 {
		return fmt.Errorf("UUID: cannot convert %T to UUID, expected 16 bytes", value)
	}
	copy(u[:], bytes)
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return u[:], nil
}
