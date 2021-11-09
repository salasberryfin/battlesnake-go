package main

import (
	"fmt"
	"math/rand"
	"time"
)

type MoveMatrix struct {
	MoveName  string
	HitWalls  bool
	HitBody   bool
	HitSnakes bool
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

	return me.Health > 10
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
		Check if BattleSnake avoids own body
	*/

	nextBody := myBody[1:] // Do not test against own head
	for _, square := range nextBody {
		fmt.Printf("Evaluating against my body item (%v, %v)\n", square.X, square.Y)
		if (newHeadPos.X == square.X) && (newHeadPos.Y == square.Y) {
			fmt.Println("Detected collision")
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
	fmt.Println("Avoid walls")

	return true
}

func nextBattleSnake(current BattleSnake, newHead Coordinates, ateFood bool) BattleSnake {
	/*
		Generate properties of my BattleSnake for given move
	*/

	// If move means eating food: increase BattleSnake length and health
	newLength := current.Length
	newHealth := current.Health
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
	for mvt, coords := range moves {
		fmt.Printf("Testing move: %v\n", mvt)
		fmt.Printf("Current latency: %v\n", me.Latency)
		ateFood := eatFood(coords, board.Food)
		afterMoveBattleSnake := nextBattleSnake(me, coords, ateFood)
		// If BattleSnake avoids walls and own body: consider the move safe
		if avoidWall(afterMoveBattleSnake.Head, Coordinates{X: board.Width, Y: board.Width}) && avoidSnake(afterMoveBattleSnake.Head, afterMoveBattleSnake.Body) && isHealthy((afterMoveBattleSnake)) {
			decision[mvt] = MoveMatrix{
				MoveName:  mvt,
				HitWalls:  false,
				HitBody:   false,
				HitSnakes: false,
			}
			if ateFood {
				safeMoves = append(safeMoves, decision[mvt])
			} else {
				safeMovesNoFood = append(safeMovesNoFood, decision[mvt])
			}
		}
	}

	var potentialMoves []MoveMatrix
	if len(safeMovesNoFood) > 0 {
		potentialMoves = safeMovesNoFood
	} else {
		potentialMoves = safeMoves
	}
	fmt.Printf("The following moves are considered safe (%v): %v\n", len(potentialMoves), potentialMoves)

	// Select a random move from the set of "valid" moves
	rand.Seed(time.Now().UnixNano())
	if len(potentialMoves) > 0 {
		randMove := potentialMoves[rand.Intn(len(potentialMoves))]
		fmt.Printf("MOVE: %v\n", randMove.MoveName)
		return NextMove{
			Move:  randMove.MoveName,
			Shout: "Let's go",
		}
	} else {
		return NextMove{ // If no safe moves were detected: something failed -> move up
			Move:  "up",
			Shout: "Failed to find a safe move: going up!",
		}
	}
}
