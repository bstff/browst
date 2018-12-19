// +build darwin

package launcher

const (
	// DefaultChromePath is the default path to use for Chrome if the
	// executable is not in $PATH.
	DefaultChromePath = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`
)

// DefaultChromeNames are the default Chrome executable names to look for in
// $PATH.
var DefaultChromeNames []string
