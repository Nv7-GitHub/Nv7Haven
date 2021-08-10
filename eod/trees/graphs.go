package trees

import (
	"bytes"
	"fmt"
	"image"
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type Graph struct {
	viz   *graphviz.Graphviz
	graph *cgraph.Graph
	nodes map[string]*cgraph.Node

	lock      *sync.RWMutex
	elemCache map[string]types.Element
}

func NewGraph(dat types.ServerData) (*Graph, error) {
	viz := graphviz.New()
	graph, err := viz.Graph()
	if err != nil {
		return nil, err
	}
	graph.SetSplines("ortho")
	return &Graph{
		viz:   viz,
		graph: graph,
		nodes: make(map[string]*cgraph.Node),

		elemCache: dat.ElemCache,
		lock:      dat.Lock,
	}, nil
}

func (g *Graph) AddElem(elem string) (string, bool) {
	elem = strings.ToLower(elem)

	// Already exists
	_, exists := g.nodes[elem]
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

	// Create Node
	node, err := g.graph.CreateNode(elem)
	if err != nil {
		return err.Error(), false
	}
	node.SetLabel(el.Name)
	node.SetShape("box")
	node.SetStyle("rounded")

	// Add parents and connections to parents
	for _, par := range el.Parents {
		lower := strings.ToLower(par)
		resp, suc := g.AddElem(par)
		if !suc {
			return resp, false
		}

		_, err := g.graph.CreateEdge(elem+"+"+lower, node, g.nodes[lower])
		if err != nil {
			return err.Error(), false
		}
	}

	// Finish
	g.nodes[elem] = node
	return "", true
}

func (g *Graph) RenderPNG() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	err := g.viz.Render(g.graph, graphviz.PNG, buf)
	return buf, err
}

func (g *Graph) Render() (image.Image, error) {
	return g.viz.RenderImage(g.graph)
}

func (g *Graph) Close() {
	g.graph.Close()
	g.viz.Close()
}
