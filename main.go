package main

import "browst/ui"
import "browst/common"
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

func runCDP() {

	b := remoteChrome.NewWithViewSize(9223, width, height)

	remoteChrome.SetWaitInputFunc(func(ev common.Event) {
		handlerChromeEvent(ev)
	})
	// url := "http://192.168.0.156:8080"
	// url := "http://www.bilibili.com"
	url := "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=golang&rsv_pq=cc657c92000278a4&rsv_t=efebRYwNRIqUH76tr8pdIqkxrMqHTFSeU2jg9v7pL2Il33nuomqpiPSfO3k&rqlang=cn&rsv_enter=1&rsv_sug3=7&rsv_sug1=6&rsv_sug7=100&rsv_sug2=0&inputT=1832&rsv_sug4=2450"
	// url := `https://www.2345.com/?39291`
	b.Start(url)

	screenshot_quit = make(chan struct{})
	b.RunScreenshot2Data(screenshot_quit, screenshot_data, 38)

	ui.SetHandlerInputFunc(func(ev common.Event) {
		handlerUIEvent(b, ev)
	})

	ui.Run(screenshot_data, width, height)
	b.Close()
	close(screenshot_quit)
}

func main() {

	runCDP()
	fmt.Println("quit")
}

func handleClickEvent(b *remoteChrome.Chrome, ret int) {
	switch ret {
	case 1:
		// startScreenshot(b)
	default:
		break
	}
}

func handlerUIEvent(b *remoteChrome.Chrome, ev common.Event) {
	switch ev.ID {
	case common.InputURL:
		url := ev.Payload.(string)
		// fmt.Print(url)
		b.Navigate(url)

	case common.InputMouseClick:
		r := ev.Payload.(common.Region)
		ret := b.Clicked(r.X, r.Y, r.Left, r.Top, r.Right, r.Bottom)
		handleClickEvent(b, ret)

	case common.InputMouseWheel:
		r := ev.Payload.(common.Region)
		b.Wheel(r.Y)

	case common.PageReset:
		b.PageReset()

	case common.InputWaited:
		r := ev.Payload.(common.BuffWaited)
		content := string(r.Cont)
		b.ABSInput(content, r.ID)
	}
}

func handlerChromeEvent(ev common.Event) {
	switch ev.ID {
	case common.WaitInput:
		r := ev.Payload.(common.BuffWaited)
		ui.WaitInput(r.ID)
	}
}
