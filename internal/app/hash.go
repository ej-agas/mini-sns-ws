package app

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2id")
)

type Hasher interface {
	Hash(text string) (string, error)
	Check(plain, hashed string) (bool, error)
}

type Argon2IDHash struct {
	Params  Argon2IDParams
	Version int
	Salt    []byte
	Value   []byte
}

func (hash Argon2IDHash) String() string {
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		hash.Version,
		hash.Params.Memory,
		hash.Params.Iterations,
		hash.Params.Parallelism,
		base64.RawStdEncoding.EncodeToString(hash.Salt),
		base64.RawStdEncoding.EncodeToString(hash.Value),
	)
}

func (hash Argon2IDHash) Create(params Argon2IDParams, version int, salt []byte, value string) *Argon2IDHash {
	instance := &Argon2IDHash{}
	instance.Params = params
	instance.Version = version
	instance.Salt = salt
	instance.Value = argon2.IDKey([]byte(value), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	return instance
}

type Argon2IDParams struct {
	Memory      uint32
	Iterations  uint32
	SaltLength  uint32
	KeyLength   uint32
	Parallelism uint8
}

type Argon2IDHasher struct {
	Params Argon2IDParams
}

func (hasher *Argon2IDHasher) Hash(text string) (string, error) {
	salt, err := hasher.makeSalt()
	if err != nil {
		return "", err
	}

	hash := Argon2IDHash{}.Create(hasher.Params, argon2.Version, salt, text)
	return hash.String(), nil
}

func (a *Argon2IDHasher) makeSalt() ([]byte, error) {
	bytes := make([]byte, a.Params.SaltLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (a *Argon2IDHasher) Check(plain, encodedHash string) (bool, error) {
	hash, err := a.FromString(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(plain), hash.Salt, hash.Params.Iterations, hash.Params.Memory, hash.Params.Parallelism, hash.Params.KeyLength)

	if subtle.ConstantTimeCompare(hash.Value, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func (a *Argon2IDHasher) FromString(encodedHash string) (Argon2IDHash, error) {
	values := strings.Split(encodedHash, "$")

	if len(values) != 6 {
		return Argon2IDHash{}, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return Argon2IDHash{}, err
	}

	if version != argon2.Version {
		return Argon2IDHash{}, ErrIncompatibleVersion
	}

	params := Argon2IDParams{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)

	if err != nil {
		return Argon2IDHash{}, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return Argon2IDHash{}, err
	}
	params.SaltLength = uint32(len(salt))

	hashValue, err := base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return Argon2IDHash{}, err
	}
	params.KeyLength = uint32(len(hashValue))

	hash := Argon2IDHash{}
	hash.Params = params
	hash.Version = version
	hash.Salt = salt
	hash.Value = hashValue

	return hash, nil
}

func NewDefaultArgon2IDHasher() *Argon2IDHasher {
	return &Argon2IDHasher{Params: Argon2IDParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}}
}
