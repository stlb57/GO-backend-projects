package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type todo struct {
	Id        int    `json:"Id"`
	Title     string `json:"Title"`
	Completed bool   `json:"Completed"`
}

func main() {
	//add/done/delete/list
	data, err := os.ReadFile("./todo.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var tasks []todo
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch os.Args[1] {
	case "list":
		for _, task := range tasks {
			fmt.Println(task.Id, task.Title, task.Completed)
		}
	case "done":
		id_done, _ := strconv.Atoi(os.Args[2])
		for i, task := range tasks {
			if task.Id == id_done {
				tasks[i].Completed = true
			}
		}
	case "delete":
		id_done, _ := strconv.Atoi(os.Args[2])
		for i, task := range tasks {
			if task.Id == id_done {
				tasks = append(tasks[:i], tasks[i+1:]...)
				break
			}
		}
	case "add":
		newTask := todo{
			Id:        len(tasks) + 1,
			Title:     os.Args[2],
			Completed: false,
		}
		tasks = append(tasks, newTask)
	}

	data, err = json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile("todo.json", data, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}
