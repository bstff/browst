package remoteChrome

import (
	"fmt"

	"github.com/mafredri/cdp"

	"time"
)

var (
	cur *cdp.Client = nil
)

type Chrome struct {
	linker *Linker
}

func New(port int) *Chrome {
	return &Chrome{
		NewCDP(port),
	}
}

func NewWithViewSize(port, w, h int) *Chrome {
	return &Chrome{
		NewCDPWithViewSize(port, w, h),
	}
}

func (b *Chrome) Close() {
	l := b.linker
	l.Close()
}

func (b *Chrome) Start(url string) {

	err := b.Navigate(url)
	if err != nil {
		return
	}

	err = b.saveSingleTargetID()
	if err != nil {
		fmt.Println(err)
	}
}

func (b *Chrome) Navigate(url string) error {
	l := b.linker
	if cur == nil {
		c, err := l.Client()
		if err != nil {
			return err
		}
		cur = c
	}

	return l.Navigate(cur, url)
}

func (b *Chrome) Clicked(left, top, right, bottom int) bool {
	l := b.linker

	x := (left + right) / 2
	y := (top + bottom) / 2
	err := l.MouseClickXY(cur, x, y)
	if err != nil {
		return false
	}

	url, err := b.maybeTargetURL()
	if err != nil {
		return false
	}
	b.Navigate(url)
	return true
}

func (b *Chrome) Wheel(delta int) error {
	l := b.linker

	return l.MouseWheel(cur, delta)
}

func (b *Chrome) RunScreenshot2Data(quit chan struct{}, ch chan []byte, delay int) {
	l := b.linker

	go func() {
		for {
			select {
			case <-quit:
				break
			default:
				buf, err := l.Screenshot2Data(cur)
				if err == nil {
					ch <- buf
				}
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}()
}

func (b *Chrome) RunScreenshot2File(
	quit chan struct{},
	path string,
	ch chan string,
	delay int) {

	l := b.linker

	go func() {
		for {
			select {
			case <-quit:
				break
			default:
				err := l.Screenshot2File(cur, path)
				if err == nil {
					ch <- path
				}
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}()
}
