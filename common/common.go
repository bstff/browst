package common

import (
	"io/ioutil"
	"os"
	"path"
)

const (
	InputURL        = "navigate"
	InputMouseClick = "click"
	InputMouseWheel = "wheel"

	WaitInput   = "wait"
	InputWaited = "waited"

	PageReset = "reset"
)

type Event struct {
	ID      string
	Payload interface{}
}

type Region struct {
	X      int
	Y      int
	Left   int
	Top    int
	Right  int
	Bottom int
}

type BuffWaited struct {
	ID   int
	Cont []byte
}

func DeletePNG(pathname string) {

	rd, _ := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			// fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			// GetAllFile(pathname + fi.Name() + "\\")
		} else {
			ext := path.Ext(fi.Name())
			if ".png" == ext {
				os.Remove(fi.Name())
			}
		}
	}
}

func Attribute(attributes []string, key string) string {
	res := ""
	for k, v := range attributes {
		if v == key {
			res = attributes[k+1]
			break
		}
	}

	return res
}

func InRect(x, y, l, t, r, b int) bool {
	if x > l && x < r && y > t && y < b {
		return true
	}
	return false
}

func LooseCrossRect(l1, t1, r1, b1, l2, t2, r2, b2 int) bool {
	centerX := (l1 + r1) / 2
	centerY := (t1 + b1) / 2

	if (centerX > l2 && centerX < r2) || (centerY > t2 && centerY < b2) {
		return true
	}
	return false
}

func CrossRect(l1, t1, r1, b1, l2, t2, r2, b2 int) bool {
	centerX := (l1 + r1) / 2
	centerY := (t1 + b1) / 2

	if centerX > l2 && centerX < r2 && centerY > t2 && centerY < b2 {
		return true
	}

	minx := Max(l1, l2)
	miny := Max(t1, t2)
	maxx := Min(r1, r2)
	maxy := Min(b1, b2)

	if minx < maxx && miny < maxy {
		return true
	}
	return false
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
