package zaiuz

import "io/ioutil"
import "net/http/httptest"
import "testing"
import "./testutil"
import "regexp"
import "code.google.com/p/go-uuid/uuid"
import a "github.com/stretchr/testify/assert"

const (
	singleFile         = "./testviews/single.html"
	singleFileOutput   = "./testviews/single.output.html"
	parentFile         = "./testviews/parent.html"
	childFile          = "./testviews/child.html"
	combinedFileOutput = "./testviews/combined.output.html"
)

type TestData string

func NewTestData() TestData {
	return TestData(uuid.New())
}

func TestNewHtmlView(t *testing.T) {
	test := func() { NewHtmlView() }
	a.Panics(t, test, "does not throw error even when no filename given.")

	result := NewHtmlView(singleFile)
	a.NotNil(t, result, "constructor returns nil errorneously.")
}

func TestSubview(t *testing.T) {
	parent := NewHtmlView(parentFile)
	test := func() { parent.Subview() }
	a.Panics(t, test, "does not throw error even when no filename given.")

	result := parent.Subview(childFile)
	a.NotNil(t, result, "subview is nil errorneously.")
}

func TestRenderSingle(t *testing.T) {
	singleView := NewHtmlView(singleFile)
	output, e := ioutil.ReadFile(singleFileOutput)
	a.NoError(t, e)

	renderEqual(t, singleView, output)
}

func TestRenderParent(t *testing.T) {
	renderFail(t, NewHtmlView(parentFile))
}

func TestRenderChild(t *testing.T) {
	renderFail(t, NewHtmlView(childFile))
}

func TestRenderCombined(t *testing.T) {
	output, e := ioutil.ReadFile(combinedFileOutput)
	a.NoError(t, e)

	parent := NewHtmlView(parentFile)
	child := parent.Subview(childFile)

	result, e := renderToString(child, nil)
	a.NoError(t, e)
	a.Equal(t, string(result), string(output), "combined result wrong.")
}

func renderEqual(t *testing.T, view *HtmlView, expected []byte) {
	result, e := renderToString(view, nil)
	a.NoError(t, e)
	a.Equal(t, string(result), string(expected), "render result mismatch.")
}

func renderMatch(t *testing.T, view *HtmlView, data interface{}, pattern string) {
	re := regexp.MustCompile(pattern)
	result, e := renderToString(view, data)

	a.NoError(t, e)
	a.NotNil(t, re.FindString(result), "render output does not match pattern.")
}

func renderFail(t *testing.T, view *HtmlView) {
	_, e := renderToString(view, nil)
	a.Error(t, e, "expected rendering to fail.")
}

func renderToString(view *HtmlView, data interface{}) (string, error) {
	response, request := testutil.NewTestRequestPair()
	context := NewContext(response, request)

	e := view.Render(context, data)
	if e != nil {
		return "", e
	}

	resp := context.ResponseWriter.(*httptest.ResponseRecorder)
	return string(resp.Body.Bytes()), nil
}
