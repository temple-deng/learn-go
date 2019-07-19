package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Post struct {
	Id		int
	Content  string
	Author	 string
}

func main() {
	csvFile, err := os.Create("posts.csv")
	if err != nil {
		panic(err)
	}

	defer csvFile.Close()

	allPosts := []Post{
		Post{1, "Hello World!", "Sau Sheong"},
		Post{2, "Bonjour Monde!", "Pierre"},
		Post{3, "Hola Mundo!", "Pedro"},
		Post{4, "Greetings Earthlings", "Sau Sheong"},	// 为什么这里也需要逗号
	}

	writer := csv.NewWriter(csvFile)

	for _, post := range allPosts {
		line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
		err := writer.Write(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		panic(nil)
	}

	var posts []Post
	for _, item := range record {
		id, _ := strconv.ParseInt(item[0], 0, 0)
		post := Post{Id: int(id), Content: item[1], Author: item[2],}
		posts = append(posts, post)
	}

	fmt.Println(posts[0].Id)
	fmt.Println(posts[0].Content)
	fmt.Println(posts[0].Author)
}