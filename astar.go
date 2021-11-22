package main

import (
    "fmt"
    "math"
    //"container/heap"
)

type Node struct {
    Position    Coordinates
    Preceding   Coordinates
    Fscore      int32
    Gscore      int32
    Hscore      int32
}

func heuristicFunction(origin Coordinates, target Coordinates) int32 {
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

func isDestination(pos Coordinates, target Coordinates) bool {
    /*
        Check if current position is the end goal
    */

    return pos == target
}

func restorePath() {
    /*
        Return chosen path once A* is done
    */
}

func astar(start Coordinates, target Coordinates) bool {
    /*
        Get best path to target coordinates
    */
    if isDestination(start, target) {
        fmt.Println("Destination reached.")
    }

    // TODO
    // openList: priority queue

    var closedList []Node
    var openList []Node
    var fScore, gScore, hScore int32

    closedList = []Node{}
    fmt.Println(closedList)

    // initialize f(n)=g(n)+h(n)
    // all set to 0 for origin node
    gScore = 0
    hScore = 0
    fScore = gScore + hScore

    openList = append(openList, Node{
        Position: start,
        Preceding: Coordinates{},
        Fscore: fScore,
        Gscore: gScore,
        Hscore: hScore,
    })

    for len(openList) > 0 {
        fmt.Println("Do stuff...")
        // TODO
    }

    return true
}
