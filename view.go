package zaiuz

// TODO: Support for text/template for some kind of views (i.e. json/xml)
import tmpl "html/template"

// TODO: interface View {}, then type JsonView and type XmlView : )
// type View interface{ Render(c *Context) error }

const RootTemplateName = "root"

type HtmlView struct{
	template *tmpl.Template
	filenames []string
}

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

func (view *HtmlView) Subview(filenames ...string) *HtmlView {
	if len(filenames) < 1 {
		panic("need at least 1 filename.")
	}

	return NewHtmlView(append(view.filenames, filenames...)...)
}

func (view *HtmlView) Render(c *Context) error {
	// TODO: Configurable/overridable content type support
	w := c.ResponseWriter
	w.Header()["Content-Type"] = []string{"text/html"}
	return view.template.Execute(w, view)
}

