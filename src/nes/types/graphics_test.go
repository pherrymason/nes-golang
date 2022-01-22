package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var black = Color{0, 0, 0}
var white = Color{255, 255, 255}
var checkTile = Tile{
	Pixels: [8 * 8]Color{
		black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black,
		black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black,
		black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black,
		black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black,
	},
}

func TestFrame_PushTile(t *testing.T) {
	frame := Frame{}
	frame.PushTile(checkTile, 0, 0)

	assert.Equal(t, checkTile.Pixels[0:7], frame.Pixels[0:7])
	assert.Equal(t, checkTile.Pixels[8:16], frame.Pixels[256:256+8])
}
