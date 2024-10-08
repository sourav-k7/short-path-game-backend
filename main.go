package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

type Coordinates struct {
	Start [2]int `json:"start"`
	End   [2]int `json:"end"`
}

func findPathDFS(start, end [2]int, grid [][]bool) [][]int {
	path := [][]int{}
	visited := make([][]bool, 20)
	for i := range visited {
		visited[i] = make([]bool, 20)
	}
	var dfs func(x, y int) bool
	dfs = func(x, y int) bool {
		if x < 0 || x >= 20 || y < 0 || y >= 20 || visited[x][y] {
			return false
		}
		visited[x][y] = true
		path = append(path, []int{x, y})
		if x == end[0] && y == end[1] {
			return true
		}
		directions := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
		for _, dir := range directions {
			if dfs(x+dir[0], y+dir[1]) {
				return true
			}
		}
		path = path[:len(path)-1]
		return false
	}
	dfs(start[0], start[1])
	return path
}

func findPathHandler(c *fiber.Ctx) error {
	var coords Coordinates
	if err := c.BodyParser(&coords); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	grid := make([][]bool, 20)
	for i := range grid {
		grid[i] = make([]bool, 20)
	}
	path := findPathDFS(coords.Start, coords.End, grid)
	return c.JSON(fiber.Map{"path": path})
}

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from backend")
	})
	app.Post("/find-path", findPathHandler)
	log.Fatal(app.Listen(":5000"))
}
