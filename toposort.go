package grammar

import "fmt"

// ##################################################
// simple toposort algorithm
// not optimized for memory consumption and CPU time
// ##################################################

// DAG, nodes have links to other nodes

type (
	nodes map[ruleName]links
	links map[ruleName]struct{}
)

// toposort returns all dependent rules in topological sort order.
func (g *Grammar) toposort() ([]ruleName, error) {
	// fill dag datastruct, nodes with links aka rules with subrules
	dag := make(nodes)

	// for all rules do ...
	for node, rule := range g.rules {
		dag[node] = make(links)
		for _, link := range rule.subrules {
			dag[node][link] = struct{}{}
		}
	}

	var result []ruleName

	// do til break condition
	for len(dag) != 0 {
		nextNodes := nodesWithoutLinks(dag)

		// cyclic dependency!
		if len(nextNodes) == 0 {
			// collect remaining rules for error reporting
			var remaining []ruleName
			for ruleName := range dag {
				remaining = append(remaining, ruleName)
			}

			// unsuccessful return with cyclic dependency error
			return nil, fmt.Errorf("grammar %q, unresolved rules: %v, %w", g.name, remaining, errCyclic)
		}

		// handle the next nodes in topo sort order
		for _, node := range nextNodes {
			// push terminal node to result
			result = append(result, node)

			// delete terminal node from topo
			delete(dag, node)

			// delete terminal nodes from links
			for _, links := range dag {
				delete(links, node)
			}
		}
	}

	return result, nil
}

// nodesWithoutLinks returns terminal nodes.
func nodesWithoutLinks(topo nodes) []ruleName {
	var result []ruleName

	for node, links := range topo {
		if len(links) == 0 {
			result = append(result, node)
		}
	}

	return result
}
