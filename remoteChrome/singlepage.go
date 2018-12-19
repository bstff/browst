package remoteChrome

import (
	"fmt"
	"github.com/mafredri/cdp/devtool"
)

var (
	onlyTarget *devtool.Target = nil
)

func (b *Chrome) saveSingleTargetID() error {
	l := b.linker

	targets, err := l.DVTListPageTargets()
	if err != nil {
		return err
	}
	if len(targets) != 1 {
		return fmt.Errorf("targets not only one")
	}
	for _, v := range targets {
		onlyTarget = v
	}

	return nil
}

func (b *Chrome) maybeTargetURL() (string, error) {
	l := b.linker

	url := ""

	targets, err := l.DVTListPageTargets()
	if err != nil {
		return url, err
	}

	var targetClose *devtool.Target
	if len(targets) == 2 {

		for _, v := range targets {
			if v.ID != onlyTarget.ID {
				url = v.URL
				targetClose = v
				break
			}
		}
	}
	if len(url) > 0 {
		go func() {
			for {
				err = l.DVTCloseTarget(targetClose)
				if err == nil {
					break
				}
			}
		}()
	}

	return url, err
}
