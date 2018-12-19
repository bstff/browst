package main

import "browst/ui"
import "browst/remoteChrome"
import "fmt"

const (
	width  = 1280
	height = 720
)

var (
	screenshot_quit = make(chan struct{})
	screenshot_data = make(chan []byte)
)

func startScreenshot(b *remoteChrome.Chrome) {
	close(screenshot_quit)
	screenshot_quit = make(chan struct{})
	b.RunScreenshot2Data(screenshot_quit, screenshot_data, 38)
}

func runCDP() {

	b := remoteChrome.NewWithViewSize(9223, width, height)

	// curURL := "http://192.168.0.156:8080"
	// curURL := "http://www.bilibili.com"
	curURL := "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=golang&rsv_pq=cc657c92000278a4&rsv_t=efebRYwNRIqUH76tr8pdIqkxrMqHTFSeU2jg9v7pL2Il33nuomqpiPSfO3k&rqlang=cn&rsv_enter=1&rsv_sug3=7&rsv_sug1=6&rsv_sug7=100&rsv_sug2=0&inputT=1832&rsv_sug4=2450"
	b.Start(curURL)

	startScreenshot(b)

	ui.SetHandlerInput(func(ev ui.InputEvent) {
		handlerInput(b, ev)
	})

	ui.Run(screenshot_data, width, height)
	b.Close()
	close(screenshot_quit)
}

func main() {

	runCDP()
	fmt.Println("quit")
}

func handlerInput(b *remoteChrome.Chrome, ev ui.InputEvent) {
	switch ev.ID {
	case "navigate":
		url := ev.Payload.(string)
		// fmt.Print(url)
		b.Navigate(url)
	case "click":
		r := ev.Payload.(ui.Mouse)
		ret := b.Clicked(r.Left, r.Top, r.Right, r.Bottom)
		if ret {
			startScreenshot(b)
			// fmt.Println("restart screenshot")
		}
	case "wheel":
		r := ev.Payload.(ui.Mouse)
		b.Wheel(r.Y)
		// fmt.Println(r.Y)
	}
}
