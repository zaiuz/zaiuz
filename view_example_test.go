package zaiuz_test

import "io"
import "os"
import "net/http/httptest"
import . "github.com/zaiuz/zaiuz"
import "github.com/zaiuz/testutil"

var ParentView = NewHtmlView("./testviews/example-parent.html")

type ParentViewData struct {
	Title string
}

var ChildView = ParentView.Subview("./testviews/example-child.html")

type ChildViewData struct {
	*ParentViewData
	Content string
}

func ExampleHtmlView() {
	parent := &ParentViewData{Title: "ExampleRender Test Title"}
	child := &ChildViewData{
		ParentViewData: parent,
		Content:        "The quick brown fox jumps over the lazy dog.",
	}

	context := testutil.NewTestContext()

	ChildView.Render(context, child)

	recorder := context.ResponseWriter.(*httptest.ResponseRecorder)
	io.Copy(os.Stdout, recorder.Body)

	// Output:
	// <html>
	// <head>
	// <title>ExampleRender Test Title</title>
	// </head>
	// <body>
	// <h1>ExampleRender Test Title</h1>
	//
	// <p>The quick brown fox jumps over the lazy dog.</p>
	//
	// </body>
	// </html>
}
