package graphics

import r "github.com/lachee/raylib-goplus/raylib"

var font *r.Font
var fontSize float32 = 16
var fontSpacing float32 = 0

func InitDrawer() {
	font = r.LoadFont("./assets/Pixel_NES.otf")
	r.SetTextureFilter(font.Texture, r.FilterPoint)
}

func DrawText(text string, x int, y int, color r.Color) {
	r.DrawTextEx(*font, text, r.Vector2{X: float32(x), Y: float32(y)}, fontSize, fontSpacing, color)
}

func DrawArrow(x int, y int, width int) {
	r.DrawRectangle(x, y, width, width, r.White)
	r.DrawRectangle(x+width, y, width, width, r.White)
	r.DrawRectangle(x+width*2, y, width, width, r.White)

	r.DrawRectangle(x+width, y+width, width, width, r.White)
}
