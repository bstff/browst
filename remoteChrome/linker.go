package remoteChrome

import (
	"browst/launcher"
	"context"
	"encoding/json"
	"fmt"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/input"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"

	"io/ioutil"
	"time"
)

var (
	viewWidth  = 720
	viewHeight = 640

	scrollX = 0
	scrollY = 50
)

type Linker struct {
	ctxt       context.Context
	cancelFunc context.CancelFunc
	devtools   *devtool.DevTools
	conn       *rpcc.Conn
}

func NewCDP(port int) *Linker {
	ctxt, cancel := context.WithCancel(context.Background())
	launchChrome(port, ctxt)

	time.Sleep(1 * time.Second)

	devt_url := fmt.Sprintf("http://localhost:%d", port)
	devt := devtool.New(devt_url)

	pageTarget, err := devt.Get(ctxt, devtool.Page)
	if err != nil {
		return nil
	}

	conn, err := rpcc.DialContext(ctxt, pageTarget.WebSocketDebuggerURL)

	if err != nil {
		cancel()
		return nil
	}

	return &Linker{
		ctxt,
		cancel,
		devt,
		conn,
	}
}

func NewCDPWithViewSize(port, w, h int) *Linker {
	ctxt, cancel := context.WithCancel(context.Background())
	viewWidth, viewHeight = w, h
	launchChrome(port, ctxt)

	time.Sleep(1 * time.Second)

	devt_url := fmt.Sprintf("http://localhost:%d", port)
	devt := devtool.New(devt_url)

	pageTarget, err := devt.Get(ctxt, devtool.Page)
	if err != nil {
		return nil
	}

	conn, err := rpcc.DialContext(ctxt, pageTarget.WebSocketDebuggerURL)

	if err != nil {
		cancel()
		return nil
	}

	return &Linker{
		ctxt,
		cancel,
		devt,
		conn,
	}
}

func (l *Linker) Close() {

	l.conn.Close()
	l.cancelFunc()
}

func launchChrome(port int, ctxt context.Context) {
	launcher.Run(ctxt,
		launcher.ExecPath("/usr/local/bin/chrome"),
		launcher.Flag("headless", true),
		launcher.Flag("no-first-run", true),
		launcher.Flag("no-default-browser-check", true),
		launcher.Flag("disable-gpu", true),
		// launcher.Flag("no-sandbox", true),
		launcher.Flag("hide-scrollbars", true),
		launcher.Flag("remote-debugging-port", port),
		launcher.Flag("window-size", fmt.Sprintf("%d,%d", viewWidth, viewHeight)),
	)

}

func (l *Linker) Client() (*cdp.Client, error) {
	conn := l.conn
	ctxt := l.ctxt

	c := cdp.NewClient(conn)
	err := c.Page.Enable(ctxt)
	if err != nil {
		return nil, err
	}

	err = c.DOM.Enable(ctxt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (l *Linker) Navigate(c *cdp.Client, url string) error {
	ctxt := l.ctxt

	// nav, err := c.Page.Navigate(ctxt, page.NewNavigateArgs(url))
	// fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

	_, err := c.Page.Navigate(ctxt, page.NewNavigateArgs(url))
	return err
}

func (l *Linker) Screenshot2Data(c *cdp.Client) ([]byte, error) {
	ctxt := l.ctxt

	screenshotArgs := page.NewCaptureScreenshotArgs().SetFormat("png")
	screenshot, err := c.Page.CaptureScreenshot(ctxt, screenshotArgs)

	return screenshot.Data, err
}

func (l *Linker) Screenshot2File(c *cdp.Client, path string) error {
	buf, err := l.Screenshot2Data(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf, 0644)
	return err
}

func (l *Linker) ListAllTargets(c *cdp.Client) ([]target.Info, error) {
	ctxt := l.ctxt

	reply, err := c.Target.GetTargets(ctxt)

	return reply.TargetInfos, err
}

func (l *Linker) ListPageTargets(c *cdp.Client) ([]target.Info, error) {
	return l.ListTargetsWithType(c, "page")
}

func (l *Linker) ListTargetsWithType(c *cdp.Client, typ string) ([]target.Info, error) {
	targets, err := l.ListAllTargets(c)
	if err != nil {
		return nil, err
	}

	var ret []target.Info
	for _, t := range targets {
		if t.Type == typ {
			ret = append(ret, t)
		}
	}

	return ret, nil
}

func (l *Linker) DVTListAllTargets() ([]*devtool.Target, error) {
	dvt := l.devtools
	ctxt := l.ctxt

	return dvt.List(ctxt)
}

func (l *Linker) DVTListTargetsWithType(typ string) ([]*devtool.Target, error) {
	targets, err := l.DVTListAllTargets()
	if err != nil {
		return nil, err
	}

	var ret []*devtool.Target
	for _, t := range targets {
		if string(t.Type) == typ {
			ret = append(ret, t)
		}
	}

	return ret, nil
}

func (l *Linker) DVTListPageTargets() ([]*devtool.Target, error) {
	return l.DVTListTargetsWithType("page")
}

func (l *Linker) ActivateTarget(c *cdp.Client, targetID target.ID) error {
	ctxt := l.ctxt

	return c.Target.ActivateTarget(ctxt, target.NewActivateTargetArgs(targetID))
}

func (l *Linker) DVTActivateTarget(target *devtool.Target) error {
	ctxt := l.ctxt
	dvt := l.devtools

	return dvt.Activate(ctxt, target)

}

func (l *Linker) DVTCloseTarget(target *devtool.Target) error {
	ctxt := l.ctxt
	dvt := l.devtools

	return dvt.Close(ctxt, target)

}

func (l *Linker) MouseClickXY(c *cdp.Client, x, y int) error {
	ctxt := l.ctxt

	pressLeft :=
		input.NewDispatchMouseEventArgs("mousePressed", float64(x), float64(y)).
			SetButton("left").SetClickCount(1)

	err := c.Input.DispatchMouseEvent(ctxt, pressLeft)
	if err != nil {
		return err
	}

	return l.releaseMouse(c, "left", x, y)
}

func (l *Linker) releaseMouse(c *cdp.Client, btn string, x, y int) error {
	ctxt := l.ctxt

	release :=
		input.NewDispatchMouseEventArgs("mouseReleased", float64(x), float64(y)).
			SetButton(btn)

	return c.Input.DispatchMouseEvent(ctxt, release)
}

func (l *Linker) MouseWheel(c *cdp.Client, delta int) error {
	ctxt := l.ctxt

	scrollJS := `(function(x, y) {
		window.scrollTo(x, y);
		return [window.scrollX, window.scrollY];
	})(%d, %d)`

	expression := fmt.Sprintf(scrollJS, scrollX, scrollY+delta)
	evalArgs :=
		runtime.NewEvaluateArgs(expression).SetAwaitPromise(true).SetReturnByValue(true)
	_, err := c.Runtime.Evaluate(ctxt, evalArgs)
	if err != nil {
		return err
	}
	scrollY += delta
	return l.releaseMouse(c, "middle", 0, delta)
}

func (l *Linker) NodeForLocation(c *cdp.Client,
	x, y int) (dom.NodeID, dom.BackendNodeID, error) {
	ctxt := l.ctxt

	reply, err := c.DOM.GetNodeForLocation(ctxt, dom.NewGetNodeForLocationArgs(x, y))

	if err != nil {
		return -1, -1, err
	}
	return *reply.NodeID, reply.BackendNodeID, err
}

func (l *Linker) DescribeNode(c *cdp.Client,
	nodeID dom.NodeID,
	backendNodeID dom.BackendNodeID) (dom.Node, error) {

	ctxt := l.ctxt
	desc := dom.NewDescribeNodeArgs()
	if nodeID > 0 {
		desc.SetNodeID(nodeID)
	}
	if backendNodeID > 0 {
		desc.SetBackendNodeID(backendNodeID)
	}

	reply, err := c.DOM.DescribeNode(ctxt, desc)

	return reply.Node, err
}

func (l *Linker) SelectNodes(c *cdp.Client, sel string) ([]dom.NodeID, error) {
	ctxt := l.ctxt

	doc, err := c.DOM.GetDocument(ctxt, nil)
	if err != nil {
		return nil, err
	}
	reply, err := c.DOM.QuerySelectorAll(ctxt,
		dom.NewQuerySelectorAllArgs(doc.Root.NodeID, sel))

	return reply.NodeIDs, err
}

func (l *Linker) NodeAttributes(c *cdp.Client, id int) ([]string, error) {
	ctxt := l.ctxt

	reply, err := c.DOM.GetAttributes(ctxt,
		dom.NewGetAttributesArgs(dom.NodeID(id)))

	return reply.Attributes, err
}

func (l *Linker) NodeBoxModel(c *cdp.Client,
	id dom.NodeID,
	backendNodeID dom.BackendNodeID) (dom.BoxModel, error) {

	ctxt := l.ctxt
	boxmodel := dom.NewGetBoxModelArgs()
	if id > 0 {
		boxmodel.SetNodeID(id)
	}
	if backendNodeID > 0 {
		boxmodel.SetBackendNodeID(backendNodeID)
	}
	reply, err := c.DOM.GetBoxModel(ctxt, boxmodel)

	return reply.Model, err

	// rect := domRect{
	// 	int(boxmodel.Model.Margin[0]),
	// 	int(boxmodel.Model.Margin[1]),
	// 	int(boxmodel.Model.Margin[4]),
	// 	int(boxmodel.Model.Margin[5]),
	// }
}

func (l *Linker) ResolveNode(c *cdp.Client,
	id dom.NodeID,
	backendNodeID dom.BackendNodeID) (runtime.RemoteObject, error) {

	ctxt := l.ctxt
	resolve := dom.NewResolveNodeArgs()
	if id > 0 {
		resolve.SetNodeID(id)
	}
	if backendNodeID > 0 {
		resolve.SetBackendNodeID(backendNodeID)
	}
	reply, err := c.DOM.ResolveNode(ctxt, resolve)
	return reply.Object, err
}

func (l *Linker) SetAttributeValue(c *cdp.Client, id int, name, value string) error {

	ctxt := l.ctxt

	return c.DOM.SetAttributeValue(ctxt,
		dom.NewSetAttributeValueArgs(dom.NodeID(id), name, value))
}

func (l *Linker) Location(c *cdp.Client) (string, error) {
	ctxt := l.ctxt

	expression := `document.location.toString()`
	evalArgs :=
		runtime.NewEvaluateArgs(expression).SetAwaitPromise(true).SetReturnByValue(true)
	reply, err := c.Runtime.Evaluate(ctxt, evalArgs)
	if err != nil {
		return "", err
	}

	var url string
	if err = json.Unmarshal(reply.Result.Value, &url); err != nil {
		fmt.Println(err)
		return "", err
	}
	return url, nil
}

func (l *Linker) EventFrameNavigated(c *cdp.Client) (*page.Frame, error) {
	ctxt := l.ctxt

	frameNavigated, err := c.Page.FrameNavigated(ctxt)
	if err != nil {
		return nil, err
	}
	defer frameNavigated.Close()

	var frame *page.Frame = nil
loop:
	for {
		select {
		case <-frameNavigated.Ready():
			reply, _ := frameNavigated.Recv()
			if err == nil {
				frame = &reply.Frame
			}
			break loop
		default:
			time.After(6 * time.Second)
		}
	}
	return frame, nil
}
