package main

import (
    "database/sql"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
    var err error
    connStr := "postgres://postgres:user@localhost/go_crud_db?sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Connected to PostgreSQL database!")
}

func main() {
    initDB()
    r := gin.Default()
// to avoid 404 
	r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to the CRUD API!",
        })
    })

    // CRUD routes
    r.GET("/items", getItems)
    r.POST("/items", createItem)
    r.GET("/items/:id", getItem)
    r.PUT("/items/:id", updateItem)
    r.DELETE("/items/:id", deleteItem)

    r.Run(":8080")
}

// CRUD functions will be added here
type Item struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// Get all items
func getItems(c *gin.Context) {
    rows, err := db.Query("SELECT id, name FROM items")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var items []Item
    for rows.Next() {
        var item Item
        if err := rows.Scan(&item.ID, &item.Name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        items = append(items, item)
    }
    c.JSON(http.StatusOK, items)
}

// Create an item
func createItem(c *gin.Context) {
    var item Item
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := db.QueryRow("INSERT INTO items(name) VALUES($1) RETURNING id", item.Name).Scan(&item.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, item)
}

// Get a single item
func getItem(c *gin.Context) {
    id := c.Param("id")
    var item Item
    err := db.QueryRow("SELECT id, name FROM items WHERE id = $1", id).Scan(&item.ID, &item.Name)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
        return
    }
    c.JSON(http.StatusOK, item)
}

// Update an item
func updateItem(c *gin.Context) {
    id := c.Param("id")
    var item Item
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := db.Exec("UPDATE items SET name=$1 WHERE id=$2", item.Name, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Item updated"})
}

// Delete an item
func deleteItem(c *gin.Context) {
    id := c.Param("id")
    _, err := db.Exec("DELETE FROM items WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
