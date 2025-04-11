package main

import (
	"database/sql"
	"log"
	// "os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// BlogPost represents the data model (Single Responsibility: represents blog post structure only)
type BlogPost struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BlogRepository handles DB interactions (Single Responsibility)
type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

// CreateBlogPost inserts a new blog post
func (r *BlogRepository) CreateBlogPost(post *BlogPost) error {
	return r.db.QueryRow(
		`INSERT INTO blog_posts (title, description, body) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		post.Title, post.Description, post.Body,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

// GetAllBlogPosts retrieves all blog posts
func (r *BlogRepository) GetAllBlogPosts() ([]BlogPost, error) {
	rows, err := r.db.Query("SELECT id, title, description, body, created_at, updated_at FROM blog_posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.Body, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetBlogPost retrieves a single post by ID
func (r *BlogRepository) GetBlogPost(id string) (*BlogPost, error) {
	var post BlogPost
	err := r.db.QueryRow("SELECT id, title, description, body, created_at, updated_at FROM blog_posts WHERE id = $1", id).
		Scan(&post.ID, &post.Title, &post.Description, &post.Body, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// DeleteBlogPost deletes a blog post by ID
func (r *BlogRepository) DeleteBlogPost(id string) error {
	_, err := r.db.Exec("DELETE FROM blog_posts WHERE id = $1", id)
	return err
}

// UpdateBlogPost updates a post by ID
func (r *BlogRepository) UpdateBlogPost(id string, post *BlogPost) error {
	_, err := r.db.Exec(
		`UPDATE blog_posts SET title=$1, description=$2, body=$3, updated_at=NOW() WHERE id=$4`,
		post.Title, post.Description, post.Body, id,
	)
	return err
}

var repo *BlogRepository

// Connects to DB (Single Responsibility)
func connectDB() *sql.DB {
	// dbURL := os.Getenv("DATABASE_URL")
	dbURL := "postgres://postgres:krishna@localhost:5432/blog?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
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
	return db
}

// createPost handles HTTP POST request (Interface Segregation: interacts via fiber)
func createPost(c *fiber.Ctx) error {
	post := new(BlogPost)
	if err := c.BodyParser(post); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := repo.CreateBlogPost(post); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(post)
}

func getAllPosts(c *fiber.Ctx) error {
	posts, err := repo.GetAllBlogPosts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(posts)
}

func getPost(c *fiber.Ctx) error {
	id := c.Params("id")
	post, err := repo.GetBlogPost(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
	}
	return c.JSON(post)
}

func deletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repo.DeleteBlogPost(id); err != nil {
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
	if err := repo.UpdateBlogPost(id, post); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Post updated"})
}

func main() {
	db := connectDB()
	repo = NewBlogRepository(db)

	app := fiber.New()

	// RESTful routes
	app.Post("/api/blog-post", createPost)
	app.Get("/api/blog-post", getAllPosts)
	app.Get("/api/blog-post/:id", getPost)
	app.Delete("/api/blog-post/:id", deletePost)
	app.Patch("/api/blog-post/:id", updatePost)

	log.Fatal(app.Listen(":3000"))
}
