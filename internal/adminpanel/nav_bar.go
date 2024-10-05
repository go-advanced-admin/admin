package adminpanel

import "fmt"

// NavBarItem represents an item in the navigation bar.
type NavBarItem struct {
	Name              string
	Link              string
	Bold              bool
	NavBarAppendSlash bool
}

// HTML returns the HTML representation of the navigation bar item.
func (i *NavBarItem) HTML() string {
	finalHTML := i.Name
	if i.Bold {
		finalHTML = `<h2 class="text-2x1 lg:block hidden font-semibold m1-2 pr-2">` + finalHTML + "</h2>"
	}
	if i.Link != "" {
		finalHTML = fmt.Sprintf(`<a class="link" href="%s">%s</a>`, i.Link, finalHTML)
	}
	return finalHTML
}

// NavBarGenerator defines a function type for generating navigation bar items.
type NavBarGenerator = func(ctx interface{}) NavBarItem
