package graph

import (
	"context"
	"fmt"
)

const (
	EndNode = "__END__"
)

// NodeHandler defines an ADK v2 node function
type NodeHandler func(ctx context.Context, state *AgentState) (*AgentState, error)

// ConditionRouter dynamically decides the next node
type ConditionRouter func(ctx context.Context, state *AgentState) (string, error)

// StateGraph coordinates the ADK 2.0 graph topology
type StateGraph struct {
	entryPoint       string
	nodes            map[string]NodeHandler
	staticEdges      map[string]string
	conditionalEdges map[string]ConditionRouter
	interrupts       map[string]bool
}

func NewStateGraph() *StateGraph {
	return &StateGraph{
		nodes:            make(map[string]NodeHandler),
		staticEdges:      make(map[string]string),
		conditionalEdges: make(map[string]ConditionRouter),
		interrupts:       make(map[string]bool),
	}
}

func (g *StateGraph) AddNode(name string, handler NodeHandler) {
	g.nodes[name] = handler
}

func (g *StateGraph) SetEntryPoint(name string) {
	g.entryPoint = name
}

func (g *StateGraph) AddEdge(from, to string) {
	g.staticEdges[from] = to
}

func (g *StateGraph) AddConditionalEdge(from string, router ConditionRouter) {
	g.conditionalEdges[from] = router
}

// AddInterrupt registers a Human-in-the-Loop checkpoint before a designated node
func (g *StateGraph) AddInterrupt(nodeName string) {
	g.interrupts[nodeName] = true
}

// Execute runs the graph from entry point to completion, respecting interrupts and conditional edges
func (g *StateGraph) Execute(ctx context.Context, state *AgentState) (*AgentState, error) {
	if g.entryPoint == "" {
		return nil, fmt.Errorf("state graph entry point not defined")
	}

	currentNode := g.entryPoint

	for {
		if currentNode == EndNode {
			state.ExecutionEnded = true
			return state, nil
		}

		// Check for Human-in-the-loop interrupt gate
		if g.interrupts[currentNode] && !state.HumanApproved {
			state.RequiresHuman = true
			return state, nil // Halts execution to wait for user approval
		}

		handler, exists := g.nodes[currentNode]
		if !exists {
			return nil, fmt.Errorf("node '%s' not found in graph", currentNode)
		}

		var err error
		state, err = handler(ctx, state)
		if err != nil {
			return nil, fmt.Errorf("error executing node '%s': %w", currentNode, err)
		}

		// Determine next node
		if router, hasCond := g.conditionalEdges[currentNode]; hasCond {
			nextNode, condErr := router(ctx, state)
			if condErr != nil {
				return nil, fmt.Errorf("conditional routing error from '%s': %w", currentNode, condErr)
			}
			currentNode = nextNode
		} else if nextNode, hasStatic := g.staticEdges[currentNode]; hasStatic {
			currentNode = nextNode
		} else {
			// No outbound edge implies termination
			currentNode = EndNode
		}
	}
}
