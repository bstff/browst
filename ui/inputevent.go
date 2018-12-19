package ui

type InputEvent struct {
	ID      string
	Payload interface{}
}

type Mouse struct {
	X      int
	Y      int
	Left   int
	Top    int
	Right  int
	Bottom int
}
