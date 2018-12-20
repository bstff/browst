package remoteChrome

import (
	"browst/common"

	"fmt"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"

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

	err := l.MouseClickXY(cur, x, y)
	if err != nil {
		return false
	}

	url, err := c.maybeTargetURL()
	if err != nil {
		return false
	}
	if len(url) > 0 {
		c.Navigate(url)
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

func (c *Chrome) isInputNode(left, top, right, bottom int) dom.NodeID {
	l := c.linker

	x := (left + right) / 2
	y := (top + bottom) / 2

	NodeID, BackendNodeID, err := l.NodeForLocation(cur, x, y)
	if err != nil {
		return -1
	}

	node, err := l.DescribeNode(cur, NodeID, BackendNodeID)
	if err != nil {
		return -1
	}
	if node.NodeName != "INPUT" {
		return -1
	}

	inputType := common.Attribute(node.Attributes, "type")
	if inputType == "" || inputType == "text" {
		return NodeID
	}

	return -1
}

func (c *Chrome) maybeInput(left, top, right, bottom int) bool {
	if -1 == c.isInputNode(left, top, right, bottom) {
		return false
	}

	ev := common.Event{
		ID: common.WaitInput,
		Payload: common.Region{
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		},
	}
	wait(ev)

	// l.SetAttributeValue(cur, NodeID, "value", "any")

	return true
}

func (c *Chrome) ABSInput(value string, left, top, right, bottom int) bool {
	l := c.linker

	NodeID := c.isInputNode(left, top, right, bottom)
	if int(NodeID) == -1 {
		fmt.Println("absinput no node")
		return false
	}

	if l.SetAttributeValue(cur, NodeID, "value", value) != nil {
		return false
	}
	return true
}
