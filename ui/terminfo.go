package ui

import (
	"syscall"
	"unsafe"
)

var (
	xCellPixel = 0
	yCellPixel = 0
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWinSize() *winsize {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return ws
}

func getCellPixel() (x, y int) {
	if xCellPixel == 0 {
		ws := getWinSize()

		xCellPixel = int(ws.Xpixel / ws.Col)
		yCellPixel = int(ws.Ypixel / ws.Row)
	}

	return xCellPixel, yCellPixel
}
