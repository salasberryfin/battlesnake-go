package main

import (
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

type MoveDetails struct {
	IsDead      bool
	IsHealthy   bool
	AvoidWall   bool
	AvoidSnakes bool
	MoveDepth   int
}

func distanceTo(head Coordinates, target Coordinates) int32 {
	/*
		Calculate distance between my head and object in board
		Distance between two points in 2D space: sqrt((olTail.X-newHead.X)^2 + (oldTail.Y-newHead.Y)^2)
	*/

	return int32(math.Sqrt((math.Pow(2, float64(target.X-head.X)) + (math.Pow(2, float64(target.Y-head.Y))))))
}

func pathTo(head Coordinates, target Coordinates) int32 {
	/*
		Give extra score if snake moves closer to target.
	*/
	distance := distanceTo(head, target)
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

	return me.Health > distanceTo(me.Head, closest)+5
}

func isDead(me BattleSnake) bool {
	/*
		Check if BattleSnake's Health < 0 after given move
	*/

	return me.Health < 0
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
	newSnakes := board.Snakes
	for index, player := range board.Snakes {
		if player.Id == snake.Id {
			newSnakes[index] = nextBattleSnake
			break
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
		IsHealthy:   false,
		IsDead:      true,
		AvoidWall:   false,
		AvoidSnakes: false,
	}
	if avoidWall(me.Head, Coordinates{board.Width, board.Height}) {
		moveDetails.AvoidWall = true
	}
	if avoidBattleSnakes(me.Head, board.Snakes) {
		moveDetails.AvoidSnakes = true
	}
	if isHealthy(me, board.Food) {
		moveDetails.IsHealthy = true
	}
	if isDead(me) {
		moveDetails.IsDead = true
	}

	return moveDetails
}

func whatNext(me BattleSnake, board Board, searchDepth int) NextMove {
	/*
		Main decision making function.
	*/

	moves := []string{
		"up",
		"down",
		"right",
		"left",
	}

	for _, move := range moves {
		// Go over each move
		// Check avoid walls||self||snakes
		newHeadPosition := nextCoords(me.Head, move)
		ateFood := eatFood(newHeadPosition, board.Food)
		newMe, newBoard := nextTurn(me, newHeadPosition, ateFood, board)
		moveSituation := getMoveDetails(newMe, board)
		moveSituation.MoveDepth = searchDepth

		if searchDepth == 0 {
			return selectNextMove()
		}
		searchDepth -= 1
		return whatNext(newMe, newBoard, searchDepth)
	}

	return NextMove{}
}
