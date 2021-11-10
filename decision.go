package main

import (
	"fmt"
	"math"
)

type MoveMatrix struct {
	MoveName  string
	HitWalls  bool
	HitBody   bool
	HitSnakes bool
	MoveScore int32
}

/*
TODO:
- Calculate distance to closest food to avoid health issues
- Check multiple turns before deciding which one to choose
*/

func isHealthy(me BattleSnake) bool {
	/*
		Check if BattleSnake's Health > `int` after given move
	*/
	fmt.Printf("Next health will be: %v\n", me.Health)

	return me.Health > 5
}

func isAlive(me BattleSnake) bool {
	/*
		Check if BattleSnake's Health > 0 after given move
	*/
	fmt.Printf("Next health will be: %v\n", me.Health)

	return me.Health > 0
}

func eatFood(newHeadPos Coordinates, boardFood []Coordinates) bool {
	/*
		Check if BattleSnake eats food for given move
	*/

	for _, food := range boardFood {
		if (newHeadPos.X == food.X) && (newHeadPos.Y == food.Y) {
			fmt.Printf("Eat food in (%v, %v)\n", newHeadPos.X, newHeadPos.Y)
			return true
		}
	}

	return false
}

func avoidSnake(newHeadPos Coordinates, myBody []Coordinates) bool {
	/*
		Check if BattleSnake avoids own body and other BattleSnakes
	*/

	nextBody := myBody[1:] // Do not test against own head
	for _, square := range nextBody {
		//fmt.Printf("Evaluating against my body item (%v, %v)\n", square.X, square.Y)
		if (newHeadPos.X == square.X) && (newHeadPos.Y == square.Y) {
			//fmt.Println("Detected collision")
			return false
		}
	}

	return true
}

func avoidWall(headPos Coordinates, boardSize Coordinates) bool {
	/*
		Check if BattleSnake avoids walls
	*/

	if (headPos.X > boardSize.X-1) || (headPos.Y > boardSize.Y-1) {
		return false
	}
	if (headPos.X < 0) || (headPos.Y < 0) {
		return false
	}
	//fmt.Println("Avoid walls")

	return true
}

func nextBattleSnake(current BattleSnake, newHead Coordinates, ateFood bool) BattleSnake {
	/*
		Generate properties of my BattleSnake for given move
	*/

	// If move means eating food: increase BattleSnake length and health
	// Else health decreases
	newLength := current.Length
	newHealth := current.Health - 1
	if ateFood {
		newLength = current.Length + 1
		newHealth = 100
	}

	// BattleSnake body coordinates after the move
	newBody := []Coordinates{
		newHead,
	}
	movedBody := current.Body[:len(current.Body)-1]
	newBody = append(newBody, movedBody...)

	// Update Body, head and length accordingly
	nextBattleSnake := BattleSnake{
		Id:      current.Id,
		Name:    current.Name,
		Health:  newHealth,
		Body:    newBody,
		Latency: current.Latency,
		Head:    newHead,
		Length:  newLength,
		Shout:   current.Shout,
		Squad:   current.Squad,
	}
	//fmt.Printf("Current BattleSnake: %v\n", current)
	//fmt.Printf("Next BattleSnake: %v\n", nextBattleSnake)

	return nextBattleSnake
}

func pathToTail(newHead Coordinates, oldTail Coordinates) int32 {
	/*
		Give extra score if snake moves closer to tail.
		Distance between two points in 2D space: sqrt((olTail.X-newHead.X)^2 + (oldTail.Y-newHead.Y)^2)
	*/
	distance := math.Sqrt((math.Pow(2, float64(oldTail.X-newHead.X)) + (math.Pow(2, float64(oldTail.Y-newHead.Y)))))
	fmt.Println("Distance between head and tail: ", distance)
	if distance < 2 {
		return 20
	} else if distance < 4 {
		return 10
	} else if distance < 6 {
		return 6
	} else {
		return 0
	}
}

func checkFuture(me BattleSnake, board Board, moves map[string]Coordinates, searchDepth int32) int32 {
	/*
		Check future possible moves
	*/

	var afterMoveBattleSnake BattleSnake
	var nextMoveScore int32
	for _, coords := range moves {
		ateFood := eatFood(coords, board.Food)
		afterMoveBattleSnake = nextBattleSnake(me, coords, ateFood)
		// If BattleSnake avoids walls and own body: add to move score
		if avoidWall(afterMoveBattleSnake.Head, Coordinates{X: board.Width, Y: board.Width}) && avoidSnake(afterMoveBattleSnake.Head, afterMoveBattleSnake.Body) && isHealthy((afterMoveBattleSnake)) {
			nextMoveScore += 1 + pathToTail(afterMoveBattleSnake.Head, me.Body[len(me.Body)-1])
		}
	}
	if searchDepth == 0 {
		return nextMoveScore
	} else {
		searchDepth -= 1
		return nextMoveScore + checkFuture(afterMoveBattleSnake, board, moves, searchDepth)
	}
}

func avoidObstacles(me BattleSnake, board Board) NextMove {
	/*
		Main decision making function.
	*/

	// Basic decision matrix
	decision := make(map[string]MoveMatrix)

	// Next Moves and Coordinates
	moves := make(map[string]Coordinates)
	moves["up"] = Coordinates{
		X: me.Head.X,
		Y: me.Head.Y + 1,
	}
	moves["down"] = Coordinates{
		X: me.Head.X,
		Y: me.Head.Y - 1,
	}
	moves["left"] = Coordinates{
		X: me.Head.X - 1,
		Y: me.Head.Y,
	}
	moves["right"] = Coordinates{
		X: me.Head.X + 1,
		Y: me.Head.Y,
	}

	var safeMoves []MoveMatrix
	var safeMovesNoFood []MoveMatrix
	var lastChance []MoveMatrix
	for mvt, coords := range moves {
		fmt.Printf("Testing move: %v\n", mvt)
		//fmt.Printf("Current latency: %v\n", me.Latency)
		ateFood := eatFood(coords, board.Food)
		afterMoveBattleSnake := nextBattleSnake(me, coords, ateFood)
		// If BattleSnake avoids walls and own body: consider the move safe
		if avoidWall(afterMoveBattleSnake.Head, Coordinates{X: board.Width, Y: board.Width}) && avoidSnake(afterMoveBattleSnake.Head, afterMoveBattleSnake.Body) {
			decision[mvt] = MoveMatrix{
				MoveName:  mvt,
				HitWalls:  false,
				HitBody:   false,
				HitSnakes: false,
				MoveScore: checkFuture(afterMoveBattleSnake, board, moves, 10),
			}
			if isHealthy((afterMoveBattleSnake)) {
				if ateFood {
					safeMoves = append(safeMoves, decision[mvt])
				} else {
					safeMovesNoFood = append(safeMovesNoFood, decision[mvt])
				}
			} else {
				// BattleSnake is about to die, but we'll keep moving!
				lastChance = append(lastChance, decision[mvt])
			}
		}
	}

	var potentialMoves []MoveMatrix
	if len(safeMovesNoFood) > 0 {
		potentialMoves = safeMovesNoFood
	} else if len(safeMoves) > 0 {
		potentialMoves = safeMoves
	} else {
		potentialMoves = lastChance
	}
	fmt.Printf("The following moves are considered safe (%v): %v\n", len(potentialMoves), potentialMoves)

	bestScore := int32(-1)
	var selectedMove MoveMatrix
	for _, move := range potentialMoves {
		fmt.Println(move.MoveName, " has a score of ", move.MoveScore)
		if move.MoveScore > bestScore {
			bestScore = int32(move.MoveScore)
			selectedMove = move
		}
	}
	fmt.Println("BattleSnake is moving ", selectedMove.MoveName, " with a score of: ", selectedMove.MoveScore)

	// Select a random move from the set of "valid" moves
	//rand.Seed(time.Now().UnixNano())
	if len(potentialMoves) > 0 {
		//randMove := potentialMoves[rand.Intn(len(potentialMoves))]
		//fmt.Printf("MOVE: %v\n", randMove.MoveName)
		fmt.Printf("MOVE: %v\n", selectedMove.MoveName)
		return NextMove{
			Move:  selectedMove.MoveName,
			Shout: "Let's go",
		}
	} else {
		fmt.Println("[CRITICAL] No moves are considered safe.")
		return NextMove{ // If no safe moves were detected: something failed -> move up
			Move:  "up",
			Shout: "Failed to find a safe move: going up!",
		}
	}
}
