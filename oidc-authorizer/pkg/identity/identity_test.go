package identity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecode(t *testing.T) {
	h := Header{
		User: "user",
	}

	b64, err := h.Base64()
	require.NoError(t, err)
	require.Equal(t, "eyJ1c2VyIjoidXNlciJ9", b64)

	h2, err := FromBase64(b64)
	require.NoError(t, err)
	require.NotNil(t, h2)
	require.Equal(t, h, *h2)
}
