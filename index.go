package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	prisma "github.com/maneaionutdev/prisma-go/generated/prisma-client"
)

type todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}

var client = prisma.New(nil)
var ctx = context.TODO()

func addTodo(c *gin.Context) {
	var newTodo todo

	if err := c.BindJSON(&newTodo); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "datele nu sunt introduse corect"})
		return
	}

	todo, err := client.CreateTodo(prisma.TodoCreateInput{
		Title:     newTodo.Title,
		Body:      newTodo.Body,
		Completed: &newTodo.Completed,
	}).Exec(ctx)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "ceva nu a mers bine, incerca iar"})
		return
	}

	c.IndentedJSON(http.StatusOK, todo)
}

func getTodos(c *gin.Context) {

	todos, err := client.Todoes(nil).Exec(ctx)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "ceva nu a mers bine, incerca iar"})

	}

	c.IndentedJSON(http.StatusOK, todos)
}

func editTodo(c *gin.Context) {
	var todo todo
	id := c.Param("id")

	if err := c.BindJSON(&todo); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "datele nu sunt introduse corect"})
		return
	}

	fmt.Println(todo)

	updatedTodo, err := client.UpdateTodo(prisma.TodoUpdateParams{
		Where: prisma.TodoWhereUniqueInput{
			ID: &id,
		},
		Data: prisma.TodoUpdateInput{
			Title: &todo.Title,
			Body:  &todo.Body,
		},
	}).Exec(ctx)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Todo updated", "todo": updatedTodo})
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	_, err := client.DeleteTodo(prisma.TodoWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Todo deleted", "id": id})
}

func setTodoDone(c *gin.Context) {
	id := c.Param("id")

	todo, err := client.Todo(prisma.TodoWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})

	}

	isCompleted := !todo.Completed

	updatedTodo, err := client.UpdateTodo(prisma.TodoUpdateParams{
		Where: prisma.TodoWhereUniqueInput{
			ID: &id,
		},
		Data: prisma.TodoUpdateInput{
			Completed: &isCompleted,
		},
	}).Exec(ctx)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "something went wrong", "error": err})

	}

	c.IndentedJSON(http.StatusOK, updatedTodo)
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/todos", getTodos)
	router.PATCH("/editTodo/:id", editTodo)
	router.PATCH("/setDone/:id", setTodoDone)
	router.POST("/todo", addTodo)
	router.DELETE("/todos/:id", deleteTodo)

	router.Run("localhost:5000")

}
