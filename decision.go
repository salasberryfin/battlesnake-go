package main

import (
	"fmt"
	"math"
)

/*
TODO:
	1. MoveScore struct
	2. Check Walls/Own/Snakes
	3. Multiply by layers to go full depth: closer turns have higher value
	4. Find path to follow tail
	5. If healthy && follow tail avoids obstacles -> follow tail
	6. If healthy && follow tail !safe -> move with highest score
	7. If !healthy -> find path to closest food
*/

type Decision struct {
	Move        NextMove
	ToFood      int32
	ToTail      int32
	FutureScore int32
}

type MoveDetails struct {
	Alive   bool
	Healthy bool
}

func distanceTo(head Coordinates, target Coordinates) int32 {
	/*
		Calculate distance between my head and object in board
		Distance between two points in 2D space: sqrt((target.X-newHead.X)^2 + (target.Y-newHead.Y)^2)
	*/

	return int32(math.Sqrt((math.Pow(2, float64(target.X-head.X)) + (math.Pow(2, float64(target.Y-head.Y))))))
}

func pathTo(head Coordinates, target Coordinates) int32 {
	/*
		Give extra score if snake moves closer to target.
	*/
	distance := distanceTo(head, target)

	return distance
}

func closestItem(head Coordinates, target []Coordinates) Coordinates {
	closest := int32(999)
	var coords Coordinates
	for _, fuel := range target {
		dist := distanceTo(head, fuel)
		if dist < closest {
			closest, coords = dist, fuel
		}
	}

	return coords
}

func isHealthy(me BattleSnake, food []Coordinates) bool {
	/*
		Check if BattleSnake's Health > `int` after given move
	*/
	closest := closestItem(me.Head, food)

	return me.Health > distanceTo(me.Head, closest)+1
}

func isDead(me BattleSnake, battleSnakes []BattleSnake, boardSize Coordinates) bool {
	/*
		Check if BattleSnake's Health < 0 after given move
	*/
	//fmt.Println("Did I starve?: ", me.Health < 0)
	//fmt.Println("I'm colliding with snakes: ", !avoidBattleSnakes(me.Head, battleSnakes))
	//fmt.Println("I'm colliding with walls: ", !avoidBoardLimits(me.Head, boardSize))
	//fmt.Println("Dead evaluation: ", (me.Health < 0) || (!avoidBattleSnakes(me.Head, battleSnakes)) || (!avoidBoardLimits(me.Head, boardSize)))

	return (me.Health < 1) || (!avoidBattleSnakes(me.Head, battleSnakes)) || (!avoidBoardLimits(me.Head, boardSize))
}

func eatFood(newHeadPos Coordinates, boardFood []Coordinates) bool {
	/*
		Check if BattleSnake eats food for given move
	*/

	for _, food := range boardFood {
		if (newHeadPos.X == food.X) && (newHeadPos.Y == food.Y) {
			return true
		}
	}

	return false
}

func avoidBattleSnakes(newHeadPos Coordinates, battleSnakes []BattleSnake) bool {
	for _, snake := range battleSnakes {
		for _, part := range snake.Body {
			if (newHeadPos.X == part.X) && (newHeadPos.Y == part.Y) {
				return false
			}
		}
	}
	return true
}

func avoidBoardLimits(headPos Coordinates, boardSize Coordinates) bool {
	/*
		Check if BattleSnake avoids walls
	*/

	if (headPos.X > boardSize.X-1) || (headPos.Y > boardSize.Y-1) {
		return false
	}
	if (headPos.X < 0) || (headPos.Y < 0) {
		return false
	}

	return true
}

func nextTurn(snake BattleSnake, newHead Coordinates, ateFood bool, board Board) (BattleSnake, Board) {
	/*
		Generate properties of my BattleSnake and resulting board for given move
	*/

	// If move means eating food: increase BattleSnake length and health
	// Else health decreases
	newLength := snake.Length
	newHealth := snake.Health - 1
	newFood := board.Food
	if ateFood {
		newLength = snake.Length + 1
		newHealth = 100
		for index, item := range board.Food {
			// If ateFood, remove the food item from the board
			if (item.X == newHead.X) && (item.Y == newHead.Y) {
				newFood = append(board.Food[:index], board.Food[index+1:]...)
				break
			}
		}
	}

	// BattleSnake body coordinates after the move
	newBody := []Coordinates{
		newHead,
	}
	movedBody := snake.Body[:len(snake.Body)-1]
	newBody = append(newBody, movedBody...)

	// Update Body, head and length accordingly
	nextBattleSnake := BattleSnake{
		Id:      snake.Id,
		Name:    snake.Name,
		Health:  newHealth,
		Body:    newBody,
		Latency: snake.Latency,
		Head:    newHead,
		Length:  newLength,
		Shout:   snake.Shout,
		Squad:   snake.Squad,
	}

	// Update Snakes in board
	var newSnakes []BattleSnake
	for index, player := range board.Snakes {
		if player.Id == snake.Id {
			newSnakes = append(newSnakes, nextBattleSnake)
		} else {
			newSnakes = append(newSnakes, board.Snakes[index])
		}
	}

	// Update game board
	nextBoard := Board{
		Height:  board.Height,
		Width:   board.Width,
		Food:    newFood,
		Hazards: board.Hazards,
		Snakes:  newSnakes,
	}

	return nextBattleSnake, nextBoard
}

func nextCoords(current Coordinates, move string) Coordinates {
	/*
		Return new position for a given board move:
			up, down, left, right
	*/

	var nextPosition Coordinates
	switch move {
	case "up":
		nextPosition = Coordinates{
			X: current.X,
			Y: current.Y + 1,
		}
	case "down":
		nextPosition = Coordinates{
			X: current.X,
			Y: current.Y - 1,
		}
	case "right":
		nextPosition = Coordinates{
			X: current.X + 1,
			Y: current.Y,
		}
	case "left":
		nextPosition = Coordinates{
			X: current.X - 1,
			Y: current.Y,
		}
	}

	return nextPosition
}

func selectNextMove() NextMove {
	/*
		Return selected next move
	*/

	return NextMove{Move: "up", Shout: "something random"}
}

func getMoveDetails(me BattleSnake, board Board) MoveDetails {
	/*
		Return move details
	*/

	moveDetails := MoveDetails{
		Healthy: false,
		Alive:   true,
	}
	//if avoidBoardLimits(me.Head, Coordinates{board.Width, board.Height}) {
	//	moveDetails.AvoidWall = true
	//}
	//if avoidBattleSnakes(me.Head, board.Snakes) {
	//	moveDetails.AvoidSnakes = true
	//}
	if isHealthy(me, board.Food) {
		moveDetails.Healthy = true
	}
	if isDead(me, board.Snakes, Coordinates{X: board.Width, Y: board.Height}) {
		moveDetails.Alive = false
	}

	return moveDetails
}

func whatNext(me BattleSnake, board Board, searchDepth int32) int32 {
	/*
		Main decision making function.
	*/

	moves := []string{
		"up",
		"down",
		"right",
		"left",
	}

	var score int32
	for _, move := range moves {
		// Give higher score to moves in which I'm alive for longer
		newHeadPosition := nextCoords(me.Head, move)
		ateFood := eatFood(newHeadPosition, board.Food)
		newMe, newBoard := nextTurn(me, newHeadPosition, ateFood, board)
		moveSituation := getMoveDetails(newMe, board)
		// If I'm alive: += 1*searchDepth
		if moveSituation.Alive {
			score += 1 * searchDepth
		}

		if searchDepth > 0 {
			searchDepth -= 1
			score = score + whatNext(newMe, newBoard, searchDepth)
		}
	}

	return score
}

func checkMoves(me BattleSnake, board Board) NextMove {
	/*
		If avoid walls && snakes:
			check future moves and give score
			populate moveDetails
			if healthy:
				tail chasing
			if !healthy:
				go for nearest food
	*/

	moves := []string{
		"up",
		"down",
		"right",
		"left",
	}

	//var safeMoves map[string][]NextMove
	safeMoves := make(map[string][]Decision)

	for _, move := range moves {
		fmt.Println("Checking moving ", move)
		firstHeadPosition := nextCoords(me.Head, move)
		newMe, _ := nextTurn(me, firstHeadPosition, eatFood(firstHeadPosition, board.Food), board)
		moveSituation := getMoveDetails(newMe, board)
		alive_score := whatNext(newMe, board, 5)
		if moveSituation.Alive {
			if moveSituation.Healthy {
				safeMoves["healthy"] = append(safeMoves["healthy"],
					Decision{
						Move:        NextMove{Move: move, Shout: "yuhu"},
						ToTail:      distanceTo(newMe.Head, newMe.Body[len(newMe.Body)-1]),
						FutureScore: alive_score,
					})
			} else {
				fmt.Println("First checking closest food item.")
				safeMoves["unhealthy"] = append(safeMoves["unhealthy"],
					Decision{
						Move:        NextMove{Move: move, Shout: "yuhu"},
						ToFood:      distanceTo(newMe.Head, closestItem(newMe.Head, board.Food)),
						FutureScore: alive_score,
					})
			}
		} else {
			fmt.Println("Move is not safe: ", move)
		}
	}

	var selectedMove NextMove
	if len(safeMoves["healthy"]) > 0 {
		fmt.Println("Should go tail chasing")
		//closest_to_tail := int32(999)
		//for _, next := range safeMoves["healthy"] {
		//	if next.ToTail < closest_to_tail {
		//		closest_to_tail = next.ToTail
		//		selectedMove = next.Move
		//	}
		//}
		highest_score := int32(0)
		for _, next := range safeMoves["healthy"] {
			if next.FutureScore > highest_score {
				highest_score = next.FutureScore
				selectedMove = next.Move
			}
		}
	} else if len(safeMoves["unhealthy"]) > 0 {
		fmt.Println("Should go for food")
		//closest_to_food := int32(999)
		//for _, next := range safeMoves["unhealthy"] {
		//	if next.ToFood < closest_to_food {
		//		closest_to_food = next.ToFood
		//		selectedMove = next.Move
		//	}
		//}
		highest_score := int32(0)
		for _, next := range safeMoves["unhealthy"] {
			if next.FutureScore < highest_score {
				highest_score = next.FutureScore
				selectedMove = next.Move
			}
		}
	} else {

		fmt.Println("I'm dead")
		selectedMove = NextMove{Move: "up", Shout: "failed"}
	}

	return selectedMove
}
