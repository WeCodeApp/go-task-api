package guard

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"task/pkg/generator"
)

const (
	DefaultKeyLength = 56
)

type KeyStore interface {
	CreateClient(ctx context.Context, userID string, key []byte) error
	GetKeyById(ctx context.Context, userID string) ([]byte, error)
}

type KeyGenerator interface {
	RandomBytes(length int) ([]byte, error)
}

type Guard struct {
	db  KeyStore
	gen KeyGenerator
	len int
}

var badKeyWUsrErrf = "invalid API key (%s) for %s"
var badKeyErrf = "invalid API key (%s)"

type Option func(*Guard) error

func WithKeyGenerator(kg KeyGenerator) Option {
	return func(g *Guard) error {
		if kg == nil {
			return errors.New("KeyGenerator was nil")
		}
		g.gen = kg
		return nil
	}
}

func WithKeyLen(l int) Option {
	return func(g *Guard) error {
		if l < 4 {
			return errors.New("Key length was too small")
		}
		g.len = l
		return nil
	}
}

func New(db KeyStore, opts ...Option) (*Guard, error) {
	if db == nil {
		return nil, errors.New("db was nil")
	}
	g := &Guard{db: db, len: DefaultKeyLength}
	var err error
	g.gen, err = generator.New(generator.AlphaNumericChars)
	if err != nil {
		return nil, errors.New("error creating charset")
	}
	for _, f := range opts {
		if err := f(g); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (s *Guard) IsKeyValid(ctx context.Context, key []byte) (string, bool) {
	pair := bytes.SplitN(key, []byte("."), 2)
	if len(pair) < 2 || len(pair[0]) == 0 {
		return "", false
	}

	userIDB := make([]byte, len(pair[0]))
	n, err := base64.StdEncoding.Decode(userIDB, pair[0])
	if err != nil {
		return "", false
	}
	userID := string(userIDB[:n])
	k, err := s.db.GetKeyById(ctx, userID)
	if err != nil {
		return "", false
	}

	return userID, bytes.Equal(k, key)
}

func (s *Guard) CreateKey(ctx context.Context, userID string) ([]byte, error) {
	if userID == "" {
		return nil, errors.New("userID was empty")
	}
	b64UsrID := base64.StdEncoding.EncodeToString([]byte(userID))
	key, err := s.gen.RandomBytes(s.len)
	if err != nil {
		return nil, err
	}
	key = bytes.Join([][]byte{[]byte(b64UsrID), key}, []byte("."))

	if err := s.db.CreateClient(ctx, userID, key); err != nil {
		return nil, err
	}

	return key, nil
}
