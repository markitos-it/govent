package domain_test

import (
	"regexp"
	"testing"

	"go-vents/internal/domain/shared"

	"github.com/stretchr/testify/require"
)

func TestCanCreateValidUUIDv4(t *testing.T) {
	uuidRegex := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	uuid := shared.UUIDv4()

	matched, err := regexp.MatchString(uuidRegex, uuid)
	require.NoError(t, err)
	require.True(t, matched)
}

func TestCanValidateUUIDv4(t *testing.T) {
	var validUUIDs = []string{
		"1220c35c-e8ce-4bdc-a214-c1c1e0d52615",
		"88c42097-a504-462b-b6bd-ba656ced7821",
		"07afbd87-24d1-4d0c-bee7-241b65fc1d7a",
		"48df8508-8666-437c-91f6-6364386a4878",
		"cfa087db-2d34-4997-8504-b16d72ab6727",
		"5ce376e3-76e2-472b-a226-666d7d0bfd67",
		"d6696cbe-b0eb-4012-8723-c5d4ec07873f",
		"0e6c5dbc-71c3-4424-a5b0-31f28bcb2ff5",
		"5862e376-3675-43c4-b75d-76631627dbbe",
		"5829300a-4674-4792-b1e1-85be372e63f4",
	}
	for _, validUUID := range validUUIDs {
		require.True(t, shared.IsUUIDv4(validUUID))
	}

	badUUIDs := []string{
		"123e4567e89b12d3a456426655440000x",
		"123e4567-e89b-12d3-a456-42665544",
		"123e4567-e89b-12d3-a456-4266",
		"123e4567-e89b-12d3-a456426655440000",
		"123e4567e89b12d3a45642665544000z",
		"g23e4567-e89b-12d3-a456-426655440000",
		"123e4567-e89b-12d3-a456-426655440000 ",
		" 123e4567-e89b-12d3-a456-426655440000",
		"123e4567-e89b-12d3-a4564266554400",
		"123e4567e89b12d3a4564266554",
	}
	for _, validUUID := range badUUIDs {
		require.False(t, shared.IsUUIDv4(validUUID))
	}
}
