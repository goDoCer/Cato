package parser

import "time"

//An enum representing the types of tasks
const (
	Group = iota + 1
	Individual
	Unassessed
	UnassessedSub
)

//Module represents a module for example Reasoning about Programs
type Module struct {
	name  string
	tasks []Task
}

//Task represents a block in the cate timetable
type Task struct {
	name     string
	class    int
	info     []string //The links to notes for the task
	deadline time.Time
}
