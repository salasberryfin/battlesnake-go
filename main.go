package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	//"encoding/json"
)

/////////////////////
// BattleSnake API //
/////////////////////

// RequestPayload: /start, /move & /end
type RequestPayload struct {
	Game  Game        `json:"game"`
	Turn  int32       `json:"turn"`
	Board Board       `json:"board"`
	You   BattleSnake `json:"you"`
}

// Details response: /
type BattleSnakeDetails struct {
	ApiVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
	Version    string `json:"version"`
}

// Next move response: /move
type NextMove struct {
	Move  string `json:"move"`
	Shout string `json:"shout"`
}

/////////////////////////
// BattleSnake Objects //
/////////////////////////

// Game Oject
type Game struct {
	Id      string  `json:"game"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
	Source  string  `json:"source"`
}

// Ruleset Object
type Ruleset struct {
	Name     string          `json:"name"`
	Version  string          `json:"version"`
	Settings RulesetSettings `json:"settings"`
}

// RulesetSettings
type RulesetSettings struct {
	FoodSpawnChance     int32          `json:"foodSpawnChance"`
	MinimumFood         int32          `json:"minimumFood"`
	HazardDamagePerTurn int32          `json:"hazardDamagePerTurn"`
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"royale.shrinkEveryNTurns"`
}

// RoyaleSettings - this settings are specific to Royale games
type RoyaleSettings struct {
	ShrinkEveryNTurns int32 `json:"shrinkEveryNTurns"`
}

// SquadSettings - this settings are specific to Squad games
type SquadSettings struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

// BattleSnake Object
type BattleSnake struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Health  int32         `json:"health"`
	Body    []Coordinates `json:"body"`
	Latency string        `json:"latency"`
	Head    Coordinates   `json:"head"`
	Length  int32         `json:"length"`
	Shout   string        `json:"shout"`
	Squad   string        `json:"squad"`
}

// Board Object
type Board struct {
	Height  int32         `json:"height"`
	Width   int32         `json:"width"`
	Food    []Coordinates `json:"food"`
	Hazards []Coordinates `json:"hazards"`
	Snakes  []BattleSnake `json:"snakes"`
}

// (X, Y) position in the board
type Coordinates struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

//////////////
// GAME API //
//////////////

func parseJsonRequest(c *gin.Context) RequestPayload {
	var request RequestPayload

	if err := c.BindJSON(&request); err != nil {
		return RequestPayload{}
	}

	return request
}

// POST /move
func moveSnake(c *gin.Context) {
	request := parseJsonRequest(c)
	move := avoidObstacles(request.You, request.Board)

	//move := NextMove{
	//    Move: "down",
	//    Shout: "heeeey",
	//}

	c.IndentedJSON(http.StatusOK, move)
}

// POST /start
func startGame(c *gin.Context) {
	// request := parseJsonRequest(c)

	c.IndentedJSON(http.StatusOK,
		"This response will be ignored anyway!",
	)
}

// POST /end
func endGame(c *gin.Context) {
	// request := parseJsonRequest(c)

	c.IndentedJSON(http.StatusOK,
		"This response will be ignored anyway!",
	)
}

// GET /
func getSnake(c *gin.Context) {
	apiVersion := "1"
	fmt.Println("Received request to /")
	c.IndentedJSON(http.StatusOK, BattleSnakeDetails{
		ApiVersion: apiVersion,      // required
		Author:     "salasberryfin", // optional
		Color:      "#888888",       // optional
		Head:       "default",       // optional
		Tail:       "default",       // optional
		Version:    "0.0.1-beta",    // optional
	})
}

func main() {
	router := gin.Default()

	router.GET("/", getSnake)
	router.POST("/start", startGame)
	router.POST("/move", moveSnake)
	router.POST("/end", endGame)

	router.Run("localhost:80")
}
