package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/elastic/crd-ref-docs/config"
	"github.com/elastic/crd-ref-docs/types"
	"github.com/russross/blackfriday/v2"
)

const htmlAnchorPrefix = ""

type HtmlRenderer struct {
	conf *config.Config
	gvd []types.GroupVersionDetails
	*Functions
}

func NewHtmlRenderer(conf *config.Config) (*HtmlRenderer, error) {
	baseFuncs, err := NewFunctions(conf)
	if err != nil {
		return nil, err
	}
	return &HtmlRenderer{conf: conf, Functions: baseFuncs}, nil
}
func (adr *HtmlRenderer) Render(gvd []types.GroupVersionDetails) error {
	funcMap := combinedFuncMap(funcMap{prefix: "html", funcs: adr.ToFuncMap()}, funcMap{funcs: sprig.TxtFuncMap()})
	tmpl, err := loadTemplate(adr.conf.TemplatesDir, funcMap)
	if err != nil {
		return err
	}

	adr.gvd = gvd

	outputFile := adr.conf.OutputPath
	finfo, err := os.Stat(outputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if finfo != nil && finfo.IsDir() {
		outputFile = filepath.Join(outputFile, "out.html")
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.ExecuteTemplate(f, mainTemplate, gvd)
}

func (adr *HtmlRenderer) ToFuncMap() template.FuncMap {
	return template.FuncMap{
		"GroupVersionID":     adr.GroupVersionID,
		"RenderAnchorID":     adr.RenderAnchorID,
		"RenderExternalLink": adr.RenderExternalLink,
		"RenderDoc":          adr.RenderDoc,
		"RenderGVLink":       adr.RenderGVLink,
		"RenderLocalLink":    adr.RenderLocalLink,
		"RenderType":         adr.RenderType,
		"RenderTypeLink":     adr.RenderTypeLink,
		"SafeID":             adr.SafeID,
		"ShouldRenderType":   adr.ShouldRenderType,
		"TypeID":             adr.TypeID,
	}
}

func (adr *HtmlRenderer) ShouldRenderType(t *types.Type) bool {
	return t != nil && (t.GVK != nil || len(t.References) > 0)
}

// ToDo add function to render markdown doc to html

func (adr *HtmlRenderer) RenderType(t *types.Type) string {
	var sb strings.Builder
	switch t.Kind {
	case types.MapKind:
		sb.WriteString("object (")
		sb.WriteString("keys:")
		sb.WriteString(adr.RenderTypeLink(t.KeyType))
		sb.WriteString(", values:")
		sb.WriteString(adr.RenderTypeLink(t.ValueType))
		sb.WriteString(")")
	case types.ArrayKind, types.SliceKind:
		sb.WriteString(adr.RenderTypeLink(t.UnderlyingType))
		sb.WriteString(" array")
	default:
		sb.WriteString(adr.RenderTypeLink(t))
	}

	return sb.String()
}

func (adr *HtmlRenderer) RenderTypeLink(t *types.Type) string {
	text := adr.SimplifiedTypeName(t)

	link, local := adr.LinkForType(t)
	if link == "" {
		return text
	}

	if local {
		return adr.RenderLocalLink(htmlAnchorPrefix, link, text)
	} else {
		return adr.RenderExternalLink(link, text)
	}
}

func (adr *HtmlRenderer) RenderLocalLink(prefix, link, text string) string {
	return fmt.Sprintf("<a href=\"#%s%s\">%s</a>", prefix, link, text)
}

func (adr *HtmlRenderer) RenderExternalLink(link, text string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", link, text)
}

func (adr *HtmlRenderer) RenderGVLink(gv types.GroupVersionDetails) string {
	return adr.RenderLocalLink(htmlAnchorPrefix, adr.GroupVersionID(gv), gv.GroupVersionString())
}

func (adr *HtmlRenderer) RenderAnchorID(id string) string {
	return fmt.Sprintf("%s%s", htmlAnchorPrefix, adr.SafeID(id))
}

func (adr *HtmlRenderer) RenderDoc(doc string) string {
	// Dirty little hack to render the bullet lists lost because of markers deleting carriage returns
    exp := regexp.MustCompile(`\s\*\s`)
    doc2 := exp.ReplaceAllString(doc, "\n * ")
	return string(blackfriday.Run([]byte(doc2)))
}