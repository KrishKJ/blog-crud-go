package main

import (
	"database/sql"
	"log"
	// "os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type BlogPost struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var db *sql.DB

func connectDB() {
	var err error
	// dbURL := os.Getenv("DATABASE_URL") // example: postgres://user:password@localhost:5432/blog
	dbURL := "postgres://postgres:krishna@localhost:5432/blog?sslmode=disable"

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS blog_posts (
			id SERIAL PRIMARY KEY,
			title TEXT,
			description TEXT,
			body TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		log.Fatal(err)
	}
}

func createPost(c *fiber.Ctx) error {
	post := new(BlogPost)
	if err := c.BodyParser(post); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	err := db.QueryRow(
		`INSERT INTO blog_posts (title, description, body) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		post.Title, post.Description, post.Body,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(post)
}

func getAllPosts(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, title, description, body, created_at, updated_at FROM blog_posts")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Body, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		posts = append(posts, post)
	}
	return c.JSON(posts)
}

func getPost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post BlogPost
	err := db.QueryRow("SELECT id, title, description, body, created_at, updated_at FROM blog_posts WHERE id = $1", id).
		Scan(&post.ID, &post.Title, &post.Description, &post.Body, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}
	return c.JSON(post)
}

func deletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Exec("DELETE FROM blog_posts WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func updatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	post := new(BlogPost)
	if err := c.BodyParser(post); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	_, err := db.Exec(
		`UPDATE blog_posts SET title=$1, description=$2, body=$3, updated_at=NOW() WHERE id=$4`,
		post.Title, post.Description, post.Body, id,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Post updated"})
}

func main() {
	connectDB()
	app := fiber.New()

	app.Post("/api/blog-post", createPost)
	app.Get("/api/blog-post", getAllPosts)
	app.Get("/api/blog-post/:id", getPost)
	app.Delete("/api/blog-post/:id", deletePost)
	app.Patch("/api/blog-post/:id", updatePost)

	log.Fatal(app.Listen(":3000"))
}
