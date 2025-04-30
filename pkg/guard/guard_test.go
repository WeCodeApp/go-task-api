package guard

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_guard "task/pkg/guard/mock"
	"testing"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_guard.NewMockKeyStore(ctrl)
	keygen := mock_guard.NewMockKeyGenerator(ctrl)

	tt := []struct {
		name   string
		db     KeyStore
		opts   []Option
		expErr bool
	}{
		{
			name:   "valid no opts",
			db:     db,
			expErr: false,
		},
		{
			name:   "valid w/ valid key gen",
			db:     db,
			opts:   []Option{WithKeyGenerator(keygen)},
			expErr: false,
		},
		{
			name:   "valid w/ valid API Key len",
			db:     db,
			opts:   []Option{WithKeyLen(4)},
			expErr: false,
		},
		{
			name:   "nil db",
			db:     nil,
			expErr: true,
		},
		{
			name:   "nil key gen",
			db:     db,
			opts:   []Option{WithKeyGenerator(nil)},
			expErr: true,
		},
		{
			name:   "bad API Key len",
			db:     db,
			opts:   []Option{WithKeyLen(3)},
			expErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			g, err := New(tc.db, tc.opts...)
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, g)
		})
	}
}

func TestCreateKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_guard.NewMockKeyStore(ctrl)
	keygen := mock_guard.NewMockKeyGenerator(ctrl)
	validKey := "some-api-key"
	userID := "123456"
	userIDB64 := base64.StdEncoding.EncodeToString([]byte(userID))
	expKey := []byte(userIDB64 + "." + validKey)
	ctx := context.TODO()

	db.EXPECT().CreateClient(gomock.Any(), userID, expKey).Return(nil)
	keygen.EXPECT().RandomBytes(DefaultKeyLength).Return([]byte(validKey), nil)
	g, err := New(db, WithKeyGenerator(keygen))
	require.NoError(t, err)
	ak, err := g.CreateKey(ctx, userID)
	require.NotNil(t, ak)
	require.EqualValues(t, ak, expKey)

	_, err = g.CreateKey(ctx, "")
	require.Error(t, err)
}

func TestIsKeyValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	db := mock_guard.NewMockKeyStore(ctrl)
	validUsrID := "12.34"
	db.EXPECT().CreateClient(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	g, err := New(db)
	require.NoError(t, err)
	validKey, err := g.CreateKey(ctx, validUsrID)
	require.NoError(t, err)
	require.NotNil(t, validKey)

	db.EXPECT().GetKeyById(gomock.Any(), validUsrID).Return(validKey, nil)

	_, v := g.IsKeyValid(ctx, validKey)
	require.True(t, v)
	_, v = g.IsKeyValid(ctx, make([]byte, 0))
	require.NotEqual(t, v, true)
	_, v = g.IsKeyValid(ctx, []byte("."))
	require.NotEqual(t, v, true)
	require.NotEqual(t, bytes.SplitN(validKey, []byte("."), 2)[1], true)
	require.NotEqual(t, bytes.SplitN(validKey, []byte("."), 2)[0], true)
	require.NotEqual(t, []byte("aW52YWxpZFVzZXJJRA==.anapikey"), true)
}
