package generator

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type testcase struct {
		desc    string
		charSet string
		expErr  bool
	}

	testcases := []testcase{
		{
			desc:    "good characterset",
			charSet: AlphaNumericChars,
			expErr:  false,
		},
		{
			desc:    "empty characterset",
			charSet: "",
			expErr:  true,
		},
		{
			desc:    "too long a characterset",
			charSet: strings.Repeat("a", 257),
			expErr:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := New(tc.charSet)
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, g)
		})
	}
}

func TestRandomBytes(t *testing.T) {
	g, err := New(AlphaNumericChars)
	require.NoError(t, err)
	require.NotNil(t, g)
	size := 36
	p, err := g.RandomBytes(size)
	require.NoError(t, err)
	require.Equal(t, len(p), size)
}
