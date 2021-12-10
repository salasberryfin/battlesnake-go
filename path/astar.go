package path

import (
    "fmt"
    "math"
)


// Coordinates (X, Y) position in the board
type Coordinates struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

// Node in the tree of A* possibilities
type Node struct {
    Position    Coordinates
    Parent      Coordinates
    Fscore      int32
    Gscore      int32
    Hscore      int32
}

func manhattanDistance(origin, target Coordinates) int32 {
	/*
		Calculate shortest distance between two points in 2D space
        Euclidean: sqrt((target.X-origin.X)^2 + (target.Y-origin.Y)^2)
        Manhattan: abs (origin.x – target.x) + abs (origin.y – target.y)
	*/
    // Euclidean
	// distance := int32(math.Sqrt((math.Pow(float64(target.X-origin.X), 2) + (math.Pow(float64(target.Y-origin.Y), 2)))))
    // Manhattan
    distance := int32(math.Abs(float64(origin.X - target.X)) + math.Abs(float64(origin.Y - target.Y)) )

	return distance
}

func isDestination(pos, target Coordinates) bool {
    /*
        Check if current position is the end goal
    */

    return pos == target
}

// getLowestF returns the lower cost entry from the list of nodes
func getLowestF(nodes map[Coordinates]Node) Node {
    var result Node
    lowest := int32(math.MaxInt32)
    for _, node := range nodes {
        if node.Fscore < lowest {
            lowest = node.Fscore
            result = node
        }
    }

    return result
}

func findNeighbors(position Coordinates) []Coordinates {
    // Generate all four potential moves
    incX := Coordinates{
        X: position.X + 1,
        Y: position.Y,
    }
    incY := Coordinates{
        X: position.X,
        Y: position.Y + 1,
    }
    decX := Coordinates{
        X: position.X - 1,
        Y: position.Y,
    }
    decY := Coordinates{
        X: position.X,
        Y: position.Y - 1,
    }
    
    neighbors :=  []Coordinates {
        incX,
        incY,
        decX,
        decY,
    }

    // TODO: check if generated positions are impossible
    // for example: out of the board or negative

    return neighbors
}

// reconstructPath returns chosen path once A* is done
func reconstructPath(pos, start Coordinates, parent map[Coordinates]Node) []Coordinates {
    var result []Coordinates
    result = append(result, pos)

    var cost int
    for {
        previousPos := parent[pos].Position
        pos = previousPos
        result = append(result, pos)
        cost = cost + 1
        if pos == start {
            break
        }
    }

    // Reverse order
    for i, j := 0, len(result) - 1; i < j; i, j = i + 1, j - 1 {
        result[i], result[j] = result[j], result[i]
    }

    fmt.Println(result)
    fmt.Println("Cost for the given path is: ", cost)

    return result
}

func astar(start, target Coordinates) bool {
    /*
        Get best path to target coordinates
    */
    if isDestination(start, target) {
        fmt.Println("Destination reached.")
    }

    var closedMap = make(map[Coordinates]Node)
    var openMap = make(map[Coordinates]Node)
    var fScore, gScore, hScore int32
    var parent = make(map[Coordinates]Node)

    // initialize f(n)=g(n)+h(n)
    // gScore is the cost from source to n
    // hScore is the estimated cost from n to target
    // fScore is the estimated cost from start to target
    gScore = 0
    hScore = manhattanDistance(start, target)
    fScore = gScore + hScore

    openMap[start] = Node{
        Position: start,
        Parent: Coordinates{},
        Fscore: fScore,
        Gscore: gScore,
        Hscore: hScore,
    }

    var current Node
    for len(openMap) > 0 {
        current = getLowestF(openMap)
        if isDestination(current.Position, target) {
            reconstructPath(current.Position, start, parent)
            break
        }
        // remove current from openMap
        delete(openMap, current.Position)
        // append to closedMap
        closedMap[current.Position] = current
        for _, neighbor := range findNeighbors(current.Position) {
            newGScore := current.Gscore + manhattanDistance(current.Position, neighbor)
            // check if distance to neighbor was calculated before
            if openEntry, ok := openMap[neighbor]; ok {
                // entry exists in openMap
                if newGScore < openEntry.Gscore {
                    delete(openMap, neighbor)
                }
            } else if closedEntry, ok := closedMap[neighbor]; ok {
                if newGScore < closedEntry.Gscore {
                    delete(closedMap, neighbor)
                }
            } else {
                neighborGScore := newGScore
                neighborHScore := manhattanDistance(neighbor, target)
                neighborFScore := neighborGScore + neighborHScore
                openMap[neighbor] = Node{
                    Position: neighbor,
                    Parent: current.Position,
                    Fscore: neighborFScore,
                    Gscore: neighborGScore,
                    Hscore: neighborGScore,
                }
                parent[neighbor] = current 
            }
        }
    }

    return true
}
