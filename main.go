package main

import (
	"browst/common"
	"browst/remoteChrome"
	"browst/ui"
	"fmt"
	"os"
	"strings"
)

const (
	width  = 1280
	height = 720
)

var (
	screenshot_quit = make(chan struct{})
	screenshot_data = make(chan []byte)
)

func runCDP(url string) {
	if url == "" {
		url = `https://www.baidu.com`
	}

	b := remoteChrome.NewWithViewSize(9223, width, height)

	remoteChrome.SetWaitInputFunc(func(ev common.Event) {
		handlerChromeEvent(ev)
	})

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
	url := ""
	if len(os.Args) > 1 {
		url = os.Args[1]
		if strings.Index(url, "http://") != 0 {
			url = "http://" + url
		}
	}
	runCDP(url)
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

	case common.Page2Top:
		b.Page2Top()

	case common.InputWaited:
		r := ev.Payload.(common.BuffWaited)
		content := string(r.Cont)
		b.ABSInput(content, r.ID)

	case common.NaviHist:
		r := ev.Payload.(common.Region)
		b.NavigateHistory(r.Y)
	}
}

func handlerChromeEvent(ev common.Event) {
	switch ev.ID {
	case common.WaitInput:
		r := ev.Payload.(common.BuffWaited)
		ui.WaitInput(r.ID)
	}
}
