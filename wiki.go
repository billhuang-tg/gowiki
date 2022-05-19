package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

)

type Page struct {
	Title string `json:"title"`
	Body string `json:"body"`
}

var pages = []Page{}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, []byte(p.Body), 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title
	body, err := os.ReadFile(filename) // return []byte, error
	if err != nil {
		return nil, err
	}
	return &Page{title, string(body)}, nil
}

func createPage(c *gin.Context) {
	var newPage Page

	// Call BindJSON to bind the received JSON to newPage
	if err := c.BindJSON(&newPage); err != nil {
		return
	}

	// No error, add the new page to pages
	pages = append(pages, newPage)
	newPage.save()

	c.IndentedJSON(http.StatusCreated, newPage)
}

func getPages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, pages)
}

func collectLocalPages() {
	root := "."
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if filepath.Ext(path) == ".txt" {
			page, err := loadPage(path)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			pages = append(pages, *page)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("loading %d pages from local storages\n", len(pages))
	fmt.Println(pages)
}

func main() {

	// collect pages from local storages
	collectLocalPages()
	
	router := gin.Default()
	router.POST("/createPage", createPage)
	router.GET("/getPages", getPages)
	router.Run("localhost:8080")

}