package path

import (
    "fmt"
    "testing"
)

// TestAStar validate functionality of best path
// search algorithm A*
func TestAStar(t *testing.T) {
    testStart := Coordinates{
        X: 0,
        Y: 3,
    }
    testGoal := Coordinates{
        X: 5,
        Y: 2,
    }

    result := astar(testStart, testGoal)
    fmt.Println(result)
}
