# Blog CRUD API with Go-Fiber

	This project is a simple Blog CRUD API built using the Go programming language and the Fiber web framework. It provides endpoints to create, read, update, and delete blog posts, and includes Swagger documentation for easy API exploration.

	## Features

	- RESTful API for managing blog posts
	- PostgreSQL database integration
	- Swagger documentation for API endpoints
	- Fiber web framework for fast and lightweight HTTP handling

	## Prerequisites

	- Go 1.18 or later
	- PostgreSQL database
	- `go mod` for dependency management

	## Installation

	1. Clone the repository:
		```bash
		git clone https://github.com/your-username/blog-crud-api.git
		cd blog-crud-api
		```

	2. Install dependencies:
		```bash
		go mod tidy
		```

	3. Set up the PostgreSQL database:
		- Create a database named `blog`.
		- Update the `dbURL` in the `connectDB` function in `main.go` with your PostgreSQL credentials.

	4. Run the application:
		```bash
		go run main.go
		```

	5. Access the API at `http://localhost:3000`.

	## API Endpoints

	### Blog Post Endpoints

	| Method | Endpoint                  | Description                |
	|--------|---------------------------|----------------------------|
	| POST   | `/api/blog-post`          | Create a new blog post     |
	| GET    | `/api/blog-post`          | Retrieve all blog posts    |
	| GET    | `/api/blog-post/{id}`     | Retrieve a blog post by ID |
	| PATCH  | `/api/blog-post/{id}`     | Update a blog post by ID   |
	| DELETE | `/api/blog-post/{id}`     | Delete a blog post by ID   |

	### Swagger Documentation

	Access the Swagger documentation at `http://localhost:3000/swagger/index.html`
 To access APIs which are running on AWS using Swagger: http://65.0.45.10:3000/swagger/index.html#/blog/get_api_blog_post__id_

Database Schema: blog ; Table: blog_posts
Column	    Type	      Constraints	      Description
id	        SERIAL	    PRIMARY KEY	      Auto-incremented blog post ID
title	      TEXT	      NOT NULL	        Title of the blog post
description	TEXT      	NOT NULL	        Short description or summary
body	      TEXT	      NOT NULL	        Main content of the blog post
created_at	TIMESTAMPTZ	DEFAULT NOW()	    Timestamp when created
updated_at	TIMESTAMPTZ	DEFAULT NOW()	    Timestamp when last updated
