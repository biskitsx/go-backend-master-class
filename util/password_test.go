package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)
	hash, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	match := CheckPassword(password, hash)
	require.NoError(t, match)

	wrongpassword := RandomString(11)
	notmatch := CheckPassword(wrongpassword, hash)
	require.Error(t, notmatch, bcrypt.ErrMismatchedHashAndPassword.Error())
}
