package graphics

type Color struct {
	R byte
	G byte
	B byte
}

type Pixel struct {
	X     int
	Y     int
	Color Color
}
