package crawler

const (
	httpKey  = "http://"
	httpsKey = "https://"
	wwwKey   = "www."
)

// List of HTML tag attributes which have a URL value:
// Source: https://stackoverflow.com/questions/2725156/complete-list-of-html-tag-attributes-which-have-a-url-value
var htmlElements = map[string][]string{
	"a":          {"href"},
	"applet":     {"codebase"},
	"area":       {"href"},
	"base":       {"href"},
	"blockquote": {"cite"},
	"body":       {"background"},
	"del":        {"cite"},
	"form":       {"action"},
	"frame":      {"longdesc", "src"},
	"head":       {"profile"},
	"iframe":     {"longdesc", "src"},
	"img":        {"longdesc", "src", "usemap"},
	"input":      {"src", "usemap"},
	"ins":        {"cite"},
	"link":       {"href"},
	"object":     {"classid", "codebase", "data", "usemap"},
	"q":          {"cite"},
	"script":     {"src"},
}
