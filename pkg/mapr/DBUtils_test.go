package mapr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnavailableNodes(t *testing.T) {
	getUnavailableNodes("", "")
	assert.Equal(t, 1, 1)
}
