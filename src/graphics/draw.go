package graphics

import r "github.com/gen2brain/raylib-go/raylib"

var font *r.Font
var fontSize float32 = 16
var fontSpacing float32 = 0

func DrawText(text string, x int, y int, color r.Color, fontSize float32) {
	r.DrawTextEx(*font, text, r.Vector2{X: float32(x), Y: float32(y)}, fontSize, fontSpacing, color)
}

func DrawArrow(x int32, y int32, width int32) {
	r.DrawRectangle(x, y, width, width, r.White)
	r.DrawRectangle(x+width, y, width, width, r.White)
	r.DrawRectangle(x+width*2, y, width, width, r.White)

	r.DrawRectangle(x+width, y+width, width, width, r.White)
}
