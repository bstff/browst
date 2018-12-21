package remoteChrome

import (
	"browst/common"

	"fmt"

	"github.com/mafredri/cdp"

	// "io/ioutil"
	"time"
)

type WaitInput func(ev common.Event)

var (
	cur  *cdp.Client = nil
	wait WaitInput
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

func (c *Chrome) Close() {
	l := c.linker
	l.Close()
}

func SetWaitInputFunc(f WaitInput) {
	wait = f
}

func (c *Chrome) Start(url string) {

	err := c.Navigate(url)
	if err != nil {
		return
	}

	err = c.keepTarget()
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Chrome) Navigate(url string) error {
	l := c.linker
	if cur == nil {
		c, err := l.Client()
		if err != nil {
			return err
		}
		cur = c
	}

	return l.Navigate(cur, url)
}

func (c *Chrome) Wheel(delta int) error {
	l := c.linker

	return l.MouseWheel(cur, delta)
}

func (c *Chrome) RunScreenshot2Data(quit chan struct{}, ch chan []byte, delay int) {
	l := c.linker

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

func (c *Chrome) RunScreenshot2File(
	quit chan struct{},
	path string,
	ch chan string,
	delay int) {

	l := c.linker

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

func (c *Chrome) maybeNavigate(x, y int) bool {
	l := c.linker
	urlBefore, err := l.Location(cur)
	if err != nil {
		fmt.Println(urlBefore)
		return false
	}

	err = l.MouseClickXY(cur, x, y)
	if err != nil {
		return false
	}

	url, err := c.maybeNewTarget()
	if err != nil {
		return false
	}
	if len(url) > 0 {
		c.Navigate(url)
		return true
	}
	urlAfter, err := l.Location(cur)
	if err != nil {
		fmt.Println(urlAfter)
		return false
	}
	if urlAfter != urlBefore {
		c.Navigate(urlAfter)
		return true
	}
	return false
}

func (c *Chrome) Clicked(col, row, left, top, right, bottom int) int {

	x := (left + right) / 2
	y := (top + bottom) / 2

	if c.maybeNavigate(x, y) {
		return 1
	}
	if c.maybeInput(left, top, right, bottom) {
		return 2
	}

	return 0
}

func (c *Chrome) maybeInput(left, top, right, bottom int) bool {
	l := c.linker

	x := (left + right) / 2
	y := (top + bottom) / 2

	NodeID, BackendNodeID, err := l.NodeForLocation(cur, x, y)
	if err != nil {
		fmt.Println("can't click input")
		fmt.Println(BackendNodeID)
		return false
	}

	node, err := l.DescribeNode(cur, NodeID, BackendNodeID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if node.NodeName != "INPUT" {
		return false
	}

	// nodeAttributes, err := l.NodeAttributes(cur, int(NodeID))
	// if err != nil {
	// 	return false
	// }
	inputType := common.Attribute(node.Attributes, "type")
	if inputType != "" && inputType != "text" {
		return false
	}

	ev := common.Event{
		ID: common.WaitInput,
		Payload: common.BuffWaited{
			ID: int(NodeID),
		},
	}
	wait(ev)

	return true
}

func (c *Chrome) ABSInput(value string, id int) bool {
	l := c.linker

	err := l.SetAttributeValue(cur, id, "value", value)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// nodeAttributes, err := l.NodeAttributes(cur, id)
	// if err != nil {
	// 	return false
	// }
	// inputType := common.Attribute(nodeAttributes, "value")
	// path := "log.txt"
	// ioutil.WriteFile(path, []byte(inputType), 0644)

	return true
}
