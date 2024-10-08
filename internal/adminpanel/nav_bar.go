package adminpanel

import "fmt"

type NavBarItem struct {
	Name              string
	Link              string
	Bold              bool
	NavBarAppendSlash bool
}

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

type NavBarGenerator = func(ctx interface{}) NavBarItem
