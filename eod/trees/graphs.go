package trees

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/goccy/go-graphviz"
)

type Graph struct {
	added    map[string]types.Empty
	body     *strings.Builder
	dot      *strings.Builder
	finished bool
	special  *strings.Builder

	lock      *sync.RWMutex
	elemCache map[string]types.Element
}

func NewGraph(dat types.ServerData) (*Graph, error) {
	dot := &strings.Builder{}
	dot.WriteString("digraph tree {\n")
	dot.WriteString("\tnode [ fontname=\"Arial\", shape=\"box\", style=\"rounded\" ];\n")
	dot.WriteString("\tedge [ dir=\"none\" ];\n")

	return &Graph{
		added:    make(map[string]types.Empty),
		dot:      dot,
		body:     &strings.Builder{},
		special:  &strings.Builder{},
		finished: false,

		elemCache: dat.ElemCache,
		lock:      dat.Lock,
	}, nil
}

func escapeGraphNode(txt string) string {
	return strings.ReplaceAll(txt, "\"", "\\\"")
}

func (g *Graph) AddElem(elem string, special bool) (string, bool) {
	elem = strings.ToLower(elem)

	// Already exists
	_, exists := g.added[elem]
	if exists {
		return "", true
	}

	// Get Element
	g.lock.RLock()
	el, exists := g.elemCache[elem]
	g.lock.RUnlock()
	if !exists {
		return fmt.Sprintf("Element **%s** doesn't exist!", elem), false
	}

	// Create Node Style because top level
	if special {
		fmt.Fprintf(g.special, "\t\"%s\" [ style=\"rounded,filled\", fontcolor=\"#ffffff\", fillcolor=\"#000000\" ];\n", escapeGraphNode(el.Name))
	}

	// Add parents and connections to parents
	for _, par := range el.Parents {
		g.AddElem(par, false)

		g.lock.RLock()
		parEl, exists := g.elemCache[strings.ToLower(par)]
		g.lock.RUnlock()
		if !exists {
			return fmt.Sprintf("Element **%s** doesn't exist!", elem), false
		}

		fmt.Fprintf(g.body, "\t\"%s\" -> \"%s\"\n", escapeGraphNode(el.Name), escapeGraphNode(parEl.Name))
	}

	// Finish
	g.added[elem] = types.Empty{}
	return "", true
}

func (g *Graph) Close(special bool, splines string) {
	if !g.finished {
		fmt.Fprintf(g.dot, "\tgraph [ splines=%s ];\n", splines)
		g.dot.WriteString(g.body.String())
		if special {
			g.dot.WriteString(g.special.String())
		}
		g.dot.WriteString("}")
		g.finished = true
	}
}

func (g *Graph) String(special bool, splines string) string {
	g.Close(special, splines)
	return g.dot.String()
}

func (g *Graph) Render(special bool, layout graphviz.Layout, format graphviz.Format) (*bytes.Buffer, error) {
	splines := "ortho"
	if layout == graphviz.TWOPI {
		splines = "false"
	}
	g.Close(special, splines)
	buf := bytes.NewBuffer(nil)

	graph, err := graphviz.ParseBytes([]byte(g.dot.String()))
	if err != nil {
		return nil, err
	}

	viz := graphviz.New()
	viz.SetLayout(layout)
	err = viz.Render(graph, format, buf)

	graph.Close()
	viz.Close()

	return buf, err
}

func (g *Graph) NodeCount() int {
	return len(g.added)
}
