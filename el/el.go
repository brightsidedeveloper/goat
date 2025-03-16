package el

import (
	"syscall/js"

	"github.com/brightsidedeveloper/goat"
)

type NodeBuilder struct {
	node goat.GoatNode
}

func NewNode(tag string) *NodeBuilder {
	return &NodeBuilder{
		node: goat.GoatNode{
			Tag:      tag,
			Attrs:    make(map[string]string),
			Events:   make(map[string]js.Func),
			Children: []goat.GoatNode{},
		},
	}
}

func (b *NodeBuilder) Attr(key, value string) *NodeBuilder {
	b.node.Attrs[key] = value
	return b
}

func (b *NodeBuilder) Event(event string, handler js.Func) *NodeBuilder {
	b.node.Events[event] = handler
	return b
}

func (b *NodeBuilder) Child(child goat.GoatNode) *NodeBuilder {
	b.node.Children = append(b.node.Children, child)
	return b
}

func (b *NodeBuilder) Text(text string) *NodeBuilder {
	if len(b.node.Children) == 0 {
		b.node.Text = text
	} else {
		b.node.Children = append(b.node.Children, goat.GoatNode{Text: text})
	}
	return b
}

func (b *NodeBuilder) Build() goat.GoatNode {
	return b.node
}

// Pre made Node Builders

// Commons
func Div() *NodeBuilder {
	return NewNode("div")
}

func Span() *NodeBuilder {
	return NewNode("span")
}

func P() *NodeBuilder {
	return NewNode("p")
}

func Article() *NodeBuilder {
	return NewNode("article")
}

func Section() *NodeBuilder {
	return NewNode("section")
}

func Header() *NodeBuilder {
	return NewNode("header")
}

func Footer() *NodeBuilder {
	return NewNode("footer")
}

func Nav() *NodeBuilder {
	return NewNode("nav")
}

func Main() *NodeBuilder {
	return NewNode("main")
}

func Aside() *NodeBuilder {
	return NewNode("aside")
}

func A() *NodeBuilder {
	return NewNode("a")
}

func Ul() *NodeBuilder {
	return NewNode("ul")
}

func Ol() *NodeBuilder {
	return NewNode("ol")
}

func Li() *NodeBuilder {
	return NewNode("li")
}

func Dl() *NodeBuilder {
	return NewNode("dl")
}

func Dt() *NodeBuilder {
	return NewNode("dt")
}

func Dd() *NodeBuilder {
	return NewNode("dd")
}

func Table() *NodeBuilder {
	return NewNode("table")
}

func Thead() *NodeBuilder {
	return NewNode("thead")
}

func Tbody() *NodeBuilder {
	return NewNode("tbody")
}

func Tfoot() *NodeBuilder {
	return NewNode("tfoot")
}

func Tr() *NodeBuilder {
	return NewNode("tr")
}

func Th() *NodeBuilder {
	return NewNode("th")
}

func Td() *NodeBuilder {
	return NewNode("td")
}

func Form() *NodeBuilder {
	return NewNode("form")
}

func Button() *NodeBuilder {
	return NewNode("button")
}

func Input() *NodeBuilder {
	return NewNode("input")
}

func Label() *NodeBuilder {
	return NewNode("label")
}

func Fieldset() *NodeBuilder {
	return NewNode("fieldset")
}

func Legend() *NodeBuilder {
	return NewNode("legend")
}

func Figure() *NodeBuilder {
	return NewNode("figure")
}

func Figcaption() *NodeBuilder {
	return NewNode("figcaption")
}

func Br() *NodeBuilder {
	return NewNode("br")
}

func Hr() *NodeBuilder {
	return NewNode("hr")
}

func Img() *NodeBuilder {
	return NewNode("img")
}

// Text

func Text(text string) goat.GoatNode {
	return NewNode("").Text(text).Build()
}

func H1() *NodeBuilder {
	return NewNode("h1")
}

func H2() *NodeBuilder {
	return NewNode("h2")
}

func H3() *NodeBuilder {
	return NewNode("h3")
}

func H4() *NodeBuilder {
	return NewNode("h4")
}

func H5() *NodeBuilder {
	return NewNode("h5")
}

func H6() *NodeBuilder {
	return NewNode("h6")
}

func Blockquote() *NodeBuilder {
	return NewNode("blockquote")
}

func Pre() *NodeBuilder {
	return NewNode("pre")
}

func Code() *NodeBuilder {
	return NewNode("code")
}

func Em() *NodeBuilder {
	return NewNode("em")
}

func Strong() *NodeBuilder {
	return NewNode("strong")
}

func I() *NodeBuilder {
	return NewNode("i")
}

func B() *NodeBuilder {
	return NewNode("b")
}

func U() *NodeBuilder {
	return NewNode("u")
}

func S() *NodeBuilder {
	return NewNode("s")
}

func Small() *NodeBuilder {
	return NewNode("small")
}

func Sub() *NodeBuilder {
	return NewNode("sub")
}

func Sup() *NodeBuilder {
	return NewNode("sup")
}

func Q() *NodeBuilder {
	return NewNode("q")
}

func Cite() *NodeBuilder {
	return NewNode("cite")
}

func Abbr() *NodeBuilder {
	return NewNode("abbr")
}

func Time() *NodeBuilder {
	return NewNode("time")
}

func Mark() *NodeBuilder {
	return NewNode("mark")
}

func Del() *NodeBuilder {
	return NewNode("del")
}

func Ins() *NodeBuilder {
	return NewNode("ins")
}

// Special

func Style() *NodeBuilder {
	return NewNode("style")
}

func Script() *NodeBuilder {
	return NewNode("script")
}
