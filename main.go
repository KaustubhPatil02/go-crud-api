package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"context"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *pgx.Conn

func initDB() {
	var err error
	dbURL := "postgres://username:password@localhost:5432/mydb"
	db, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	fmt.Println("Database connected!")
}

func getItems(c *gin.Context) {
	rows, _ := db.Query(context.Background(), "SELECT id, name FROM items")
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		rows.Scan(&item.ID, &item.Name)
		items = append(items, item)
	}
	c.JSON(http.StatusOK, items)
}

func createItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(context.Background(), "INSERT INTO items (name) VALUES ($1)", item.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create item"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func updateItem(c *gin.Context) {
	id := c.Param("id")
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(context.Background(), "UPDATE items SET name=$1 WHERE id=$2", item.Name, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update item"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func deleteItem(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec(context.Background(), "DELETE FROM items WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Item deleted"})
}

func main() {
	initDB()
	defer db.Close(context.Background())

	r := gin.Default()

	r.GET("/items", getItems)
	r.POST("/items", createItem)
	r.PUT("/items/:id", updateItem)
	r.DELETE("/items/:id", deleteItem)

	r.Run(":8080")
}
