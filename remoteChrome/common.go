package remoteChrome

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func deletePNG(pathname string) {

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

func findAttribute(array []string, attr string) string {
	res := ""
	for k, v := range array {
		if v == attr {
			res = array[k+1]
			break
		}
	}

	return res
}

func mendURL(url string) string {
	if len(url) > 0 {
		if strings.Index(url, "https://") == 0 {
			return url
		}
		if strings.Index(url, "//") == 0 {
			url = "https:" + url
		} else if strings.Index(url, "http://") != 0 {
			url = "http://" + url
		}
	}
	return url
}

func crossRect(l1, t1, r1, b1, l2, t2, r2, b2 int) bool {
	centerX := (l1 + r1) / 2
	centerY := (t1 + b1) / 2

	if centerX > l2 && centerX < r2 && centerY > t2 && centerY < b2 {
		return true
	}

	minx := max(l1, l2)
	miny := max(t1, t2)
	maxx := min(r1, r2)
	maxy := min(b1, b2)

	if minx < maxx && miny < maxy {
		return true
	}
	return false
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
