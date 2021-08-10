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
	dot      *strings.Builder
	finished bool

	lock      *sync.RWMutex
	elemCache map[string]types.Element
}

func NewGraph(dat types.ServerData) (*Graph, error) {
	dot := &strings.Builder{}
	dot.WriteString("digraph tree {\n")
	dot.WriteString("\tnode [ fontname=\"Arial\", shape=\"box\", style=\"rounded\" ];\n")
	dot.WriteString("\tedge [ dir=\"none\" ];\n")
	dot.WriteString("\tgraph [ splines=ortho ];\n")

	return &Graph{
		added:    make(map[string]types.Empty),
		dot:      dot,
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
		fmt.Fprintf(g.dot, "\t\"%s\" [ style=\"rounded,filled\", fontcolor=\"#ffffff\", fillcolor=\"#000000\" ];\n", escapeGraphNode(el.Name))
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

		fmt.Fprintf(g.dot, "\t\"%s\" -> \"%s\"\n", escapeGraphNode(el.Name), escapeGraphNode(parEl.Name))
	}

	// Finish
	g.added[elem] = types.Empty{}
	return "", true
}

func (g *Graph) Close() {
	if !g.finished {
		g.dot.WriteString("}")
		g.finished = true
	}
}

func (g *Graph) String() string {
	g.Close()
	return g.dot.String()
}

func (g *Graph) RenderPNG() (*bytes.Buffer, error) {
	g.Close()
	buf := bytes.NewBuffer(nil)

	graph, err := graphviz.ParseBytes([]byte(g.String()))
	if err != nil {
		return nil, err
	}

	viz := graphviz.New()
	err = viz.Render(graph, graphviz.PNG, buf)
	if err != nil {
		return nil, err
	}

	graph.Close()
	viz.Close()

	return buf, nil
}

func (g *Graph) NodeCount() int {
	return len(g.added)
}
