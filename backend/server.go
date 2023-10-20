package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/tictactoefan/tictactoe/backend/game_manager"
	"github.com/tictactoefan/tictactoe/backend/tictactoe"
	"net/http"
	"strconv"
)

func main() {

	e := echo.New()

	e.GET("/games", GetGames)
	e.GET("/games/:id", GetGame)
	e.POST("/games", CreateGames)
	e.POST("/games/move", MoveGame)

	e.Logger.Fatal(e.Start(":8080"))
}

// TODO: Update model to use different data types to avoid so many type conversions.
type GameJSON struct {
	ID       string `json:"id"`
	Board    string `json:"board"`
	NextTurn string `json:"nextTurn"`
	GameOver string `json:"gameOver"`
	Winner   string `json:"winner"`
}

type IdJSON struct {
	ID string `json:"id"`
}

type ErrorJSON struct {
	Error string `json:"error"`
}

var gm = game_manager.NewGameManager()

func GetGame(c echo.Context) error {
	id := c.Param("id")

	gameId, convErr := strconv.ParseInt(id, 10, 32)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, ErrorJSON{"bad id provided"})
	}

	game, findGameErr := gm.GetGameById(int(gameId))
	if findGameErr != nil {
		return c.JSON(http.StatusBadRequest, ErrorJSON{"game does not exist"})
	}

	var result = GameJSON{
		ID:       id,
		Board:    game.GetBoard(),
		NextTurn: string(game.GetNextTurn()),
		GameOver: strconv.FormatBool(game.IsGameOver()),
		Winner:   string(game.GetWinner()),
	}

	return c.JSON(http.StatusOK, result)
}

func GetGames(c echo.Context) error {
	games := gm.GetAllGames()
	gamesJSON := make([]GameJSON, 0, len(games))
	for _, game := range games {
		gamesJSON = append(gamesJSON, GameJSON{ // TODO: refactor into own function.
			ID:       strconv.FormatInt(int64(game.GetId()), 10),
			Board:    game.GetBoard(),
			NextTurn: string(game.GetNextTurn()),
			GameOver: strconv.FormatBool(game.IsGameOver()),
			Winner:   string(game.GetWinner()),
		})
	}

	return c.JSON(http.StatusOK, gamesJSON)
}

func CreateGames(c echo.Context) error {
	gameId := gm.CreateNewGame()

	var result = IdJSON{ID: strconv.FormatInt(int64(gameId), 10)}
	return c.JSON(http.StatusCreated, result)
}

type MakeMoveJSON struct {
	GameID   string `json:"gameId"`
	Location string `json:"location"`
	Player   string `json:"player"`
}

func MoveGame(c echo.Context) error {
	var makeMoveJSON MakeMoveJSON
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&makeMoveJSON)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorJSON{err.Error()})
	} else if makeMoveJSON.GameID == "" || makeMoveJSON.Location == "" || makeMoveJSON.Player == "" { // TODO: refactor into own function to use on all API calls
		return c.JSON(http.StatusBadRequest, ErrorJSON{"'gameId', 'location' and 'player' fields required"})
	}

	gameId, convErr := strconv.ParseInt(makeMoveJSON.GameID, 10, 32) // TODO: refactor into own function.
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, ErrorJSON{"bad gameId provided"})
	}
	game, findGameErr := gm.GetGameById(int(gameId)) // TODO: refactor into own function.
	if findGameErr != nil {
		return c.JSON(http.StatusBadRequest, ErrorJSON{"game does not exist"})
	}
	location, convErr := strconv.ParseInt(makeMoveJSON.Location, 10, 32)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, ErrorJSON{"bad location provided"})
	}

	err = gm.MakeMove(game.GetId(), tictactoe.Cell(makeMoveJSON.Player), tictactoe.Location(location))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorJSON{err.Error()})
	}
	return c.NoContent(200)
}
