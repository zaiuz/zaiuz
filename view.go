package zaiuz

// TODO: Support for text/template for some kind of views (i.e. json/xml)
import tmpl "html/template"

// TODO: interface View {}, then type JsonView and type XmlView : )
// type View interface{ Render(c *Context) error }

// The html/template template name to starts rendering from. Usually this is the topmost
// {{define "root"}} block in your template.
const RootTemplateName = "root"

// Represents a single html/template template. Encapsulate the template pathname from
// controller action code and allows further subviews based on this view.
type HtmlView struct {
	template  *tmpl.Template
	filenames []string
}

// Creates a new html view from the specified template name which should be a
// html/template-compatible html template file.
func NewHtmlView(filenames ...string) *HtmlView {
	if len(filenames) < 1 {
		panic("needs at least 1 filename.")
	}

	t := tmpl.New(RootTemplateName)
	t, e := t.ParseFiles(filenames...)
	if e != nil {
		panic(e) // better to failfast here since views are pre-loaded at startup.
	}

	return &HtmlView{t, filenames}
}

// Creates a subview from the receiving view. Subview templates contains all templates
// defined in the parent view.
func (view *HtmlView) Subview(filenames ...string) *HtmlView {
	if len(filenames) < 1 {
		panic("need at least 1 filename.")
	}

	return NewHtmlView(append(view.filenames, filenames...)...)
}

// Renders the view to the response in the supplied Context with the given view data
// context.
func (view *HtmlView) Render(c *Context, data interface{}) error {
	// TODO: Configurable/overridable content type support
	w := c.ResponseWriter
	w.Header()["Content-Type"] = []string{"text/html"}
	return view.template.Execute(w, data)
}

// Same as calling HtmlView.Render(*Context, interface{}) but will panic if there is an
// error.
func (view *HtmlView) MustRender(c *Context, data interface{}) {
	e := view.Render(c, data)
	if e != nil {
		panic(e)
	}
}
