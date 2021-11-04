package main

// import "fmt"
import "net/http"
import "github.com/gin-gonic/gin"

type httpResponse struct {
    Body        string  `json:"body"`
}

type battlesnakeDetails struct {
    ApiVersion  string  `json:"apiversion"`
    Author      string  `json:"author"`
    Color       string  `json:"color"`
    Head        string  `json:"head"`
    Tail        string  `json:"tail"`
    Version     string  `json:"version"`
}

// TODO: declare settings object
type settings struct {
    FoodSpawnChance int32   `json:"foodSpawnChance"`
}

type ruleset struct {
    Name        string      `json:"name"`
    Version     string      `json:"version"`
    Settings    settings    `json:"settings"`
}

type game struct {
    Id        string    `json:"game"`
    Ruleset   ruleset   `json:"ruleset"`
    Timeout   int32     `json:"timeout"`
    Source    string    `json:"source"`
}

type pos struct {
    X   int32
    Y   int32
}

// TODO: declare battlesnake object
type snake struct {
}

type board struct {
    Heigth      int32   `json:"height"`
    Width       int32   `json:"width"`
    Food        []pos   `json:"food"`
    Hazards     []pos   `json:"hazards"`
    Snakes      []snake `json:"snakes"`
}

type move struct {
    Game        game    `json:"game"`
    Turn        string  `json:"turn"`
    Board       board   `json:"board"`
    You         snake   `json:"you"`
}

type nextMove struct {
    Move        string  `json:"move"`
    Shout       string  `json:"shout"`
}

func moveSnake(c *gin.Context) {
    var game move

    if err := c.BindJSON(&game); err != nil {
        return
    }

    c.IndentedJSON(http.StatusOK, nextMove{
        Move: "up",
        Shout: "Moving up!",
    })
}

func startGame(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, httpResponse{
        Body: "This is a test response",
    })
}

func endGame(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, httpResponse{
        Body: "This is a test response",
    })
}

func getSnake(c *gin.Context) {
    apiVersion := "1"
    c.IndentedJSON(http.StatusOK, battlesnakeDetails{
        ApiVersion: apiVersion,
        Author:     "salasberryfin",
        Color:      "#888888",
        Head:       "default",
        Tail:       "default",
        Version:    "0.0.1-beta",
    })
}

func main() {
    router := gin.Default()

    router.GET("/", getSnake)
    router.POST("/move", moveSnake)
    router.POST("/end", endGame)

    router.Run("localhost:8080")
}
