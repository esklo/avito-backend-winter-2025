package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgon2_Hash(t *testing.T) {
	t.Parallel()

	t.Run("successful hash", func(t *testing.T) {
		t.Parallel()
		h := NewArgon2()

		hash1, salt1, err := h.Hash("password")
		require.NoError(t, err)
		require.NotEmpty(t, hash1)
		require.NotEmpty(t, salt1)

		hash2, salt2, err := h.Hash("password")
		require.NoError(t, err)
		require.NotEmpty(t, hash2)
		require.NotEmpty(t, salt2)

		assert.NotEqual(t, hash1, hash2)
		assert.NotEqual(t, salt1, salt2)
	})

	t.Run("empty password", func(t *testing.T) {
		t.Parallel()
		h := NewArgon2()

		hash, salt, err := h.Hash("")
		require.NoError(t, err)
		require.NotEmpty(t, hash)
		require.NotEmpty(t, salt)
	})
}

func TestArgon2_Verify(t *testing.T) {
	t.Parallel()

	t.Run("successful verification", func(t *testing.T) {
		t.Parallel()
		h := NewArgon2()
		password := "test-password"

		hash, salt, err := h.Hash(password)
		require.NoError(t, err)

		assert.True(t, h.Verify(password, hash, salt))
	})

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
		h := NewArgon2()
		password := "test-password"

		hash, salt, err := h.Hash(password)
		require.NoError(t, err)

		assert.False(t, h.Verify("wrong-password", hash, salt))
	})

	t.Run("empty password", func(t *testing.T) {
		t.Parallel()
		h := NewArgon2()

		hash, salt, err := h.Hash("")
		require.NoError(t, err)

		assert.True(t, h.Verify("", hash, salt))
		assert.False(t, h.Verify("some-password", hash, salt))
	})
}
