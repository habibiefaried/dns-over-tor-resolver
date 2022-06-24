package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	c, err := ReadConfig("../")
	if err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, c.Tor.Port, "53", "The two words should be the same.")
		assert.Equal(t, c.Manual["pure.privatedomain."], "127.1.5.6", "The two words should be the same.")
		assert.Equal(t, c.CacheTTL, 3600)
	}
}
