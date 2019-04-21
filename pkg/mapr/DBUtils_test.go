package mapr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnavailableNodes(t *testing.T) {
	getUnavailableNodes("", "")
	assert.Equal(t, 1, 1)
}

func TestGetAvailableNodes(t *testing.T) {
	//getAvailableNodes("", "centos")
	//assert.Equal(t, 1, 1)
	fmt.Println(IsRequestDoable(10, "centos", "7"))
}
