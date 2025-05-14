package encryption

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidatePassword(t *testing.T) {
	plain := "password"
	hashed, err := HashPassword(plain)
	log.Println(hashed)
	require.Nil(t, err)
	require.NotEmpty(t, hashed)

	err = ValidatePassword(hashed, plain)
	require.Nil(t, err)
}
