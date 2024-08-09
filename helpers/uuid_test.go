package helpers

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	regex := regexp.MustCompile(
		`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`,
	)

	assert.True(t, regex.MatchString(uuid))

	uuid2 := GenerateUUID()
	assert.NotEqual(t, uuid, uuid2)
}
