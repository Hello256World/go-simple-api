package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgrald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func returnTheBook(c *gin.Context) {
	id := c.Query("id")

	myBook, err := giveMeBook(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	myBook.Quantity += 1
	c.IndentedJSON(http.StatusAccepted, myBook)
}

func checkOutTheBook(c *gin.Context) {
	id := c.Query("id")

	myBook, err := giveMeBook(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if myBook.Quantity <= 0 {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The book is out"})
		return
	}

	myBook.Quantity -= 1
	c.IndentedJSON(http.StatusAccepted, myBook)
}

func getBookById(c *gin.Context) {
	id := c.Param("id")

	myBook, err := giveMeBook(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusFound, myBook)
}

func giveMeBook(id string) (*book, error) {
	for index, val := range books {
		if val.ID == id {
			return &books[index], nil
		}
	}

	return nil, errors.New("there is nothing to match")
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book

	// dataByte, err := io.ReadAll(c.Request.Body)

	// if err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Somthing went Wrong in sec 1"})
	// 	return
	// }

	// if err := json.Unmarshal(dataByte, &newBook); err != nil {
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Somthing went Wrong in sec 2"})
	// 	return
	// }

	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Somthing went Wrong"})
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, books)

}

func main() {
	r := gin.Default()
	r.GET("/books", getBooks)
	r.GET("/books/:id", getBookById)
	r.POST("/books", createBook)
	r.PATCH("/books/check", checkOutTheBook)
	r.PATCH("/books/return", returnTheBook)
	r.Run()
}
