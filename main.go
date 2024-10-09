package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Coordinates struct {
	Start     [2]int   `json:"start"`
	End       [2]int   `json:"end"`
	Obstacles [][2]int `json:"obstacles"`
}

func findPath(start, end [2]int, grid [][]bool) [][]int {
	visited := make([][]bool, 20)
	for i := range visited {
		visited[i] = make([]bool, 20)
	}
	directions := [][2]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	type Node struct {
		Position [2]int
		Path     [][]int
	}
	queue := []Node{{start, [][]int{start[:]}}}
	visited[start[0]][start[1]] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.Position == end {
			return current.Path
		}
		for _, dir := range directions {
			next := [2]int{current.Position[0] + dir[0], current.Position[1] + dir[1]}
			if next[0] >= 0 && next[0] < 20 && next[1] >= 0 && next[1] < 20 &&
				!visited[next[0]][next[1]] && !grid[next[0]][next[1]] {
				visited[next[0]][next[1]] = true
				newPath := append([][]int{}, current.Path...)
				newPath = append(newPath, next[:])
				queue = append(queue, Node{next, newPath})
			}
		}
	}
	return [][]int{}
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
	for _, obstacle := range coords.Obstacles {
		if obstacle[0] >= 0 && obstacle[0] < 20 && obstacle[1] >= 0 && obstacle[1] < 20 {
			grid[obstacle[0]][obstacle[1]] = true
		}
	}
	path := findPath(coords.Start, coords.End, grid)
	return c.JSON(fiber.Map{"path": path})
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from backend")
	})

	app.Post("/find-path", findPathHandler)

	log.Fatal(app.Listen(":5000"))
}
