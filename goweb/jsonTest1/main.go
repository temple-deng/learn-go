package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Post struct {
	Id  int `json:"id"`
	Content string `json:"content"`
	Author Author `json:"author"`
	Comments []Comment `json:"comments"`
}

type Author struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	Id int `json:"id"`
	Content string `json:"content"`
	Author string `json:"author"`
}

func main() {
	_, err := decode("post.json")
	if err != nil {
		fmt.Println("Error:", err)
	}
}


func decode(filename string) (post Post, err error){
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("There is an error on Opending File:", err)
		return
	}

	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON Data:", err)
		return
	}

	json.Unmarshal(jsonData, &post)
	return
}

