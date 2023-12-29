package endpoints

import "log"

type TaskEventModel struct {
	EventType string
	Timestamp string
}

type TaskModel struct {
	Id           int              `json:"id"`
	ParentId     int              `json:"parent_id"`
	ChildrenIDs  []int            `json:"children_ids"`
	Text         string           `json:"text"`
	CreatedAt    string           `json:"created_at"`
	TaskStatus   string           `json:"task_status"`
	EventHistory []TaskEventModel `json:"event_history"`
}

type LinearizedTaskReference struct {
	TaskId                int        `json:"task_id"`
	Depth                 int        `json:"depth"`
	NumberOfChildren      int        `json:"number_of_children"`
	NumberOfAllNodesBelow int        `json:"number_of_all_nodes_below"`
	Task                  *TaskModel `json:"task"`
}

var visited_nodes, in_progress []bool
var linearized_tasks []LinearizedTaskReference
var current_depth int

func find_index(tasks *[]TaskModel, index int) int {
	for i, task := range *tasks {
		if task.Id == index {
			return i
		}
	}
	return -1
}

// Returns the number of nodes below including itself for given node
func dfs_helper(i int, tasks *[]TaskModel) int {
	if visited_nodes[i] {
		log.Println("A node visited multiple times, probably referenced by multiple parents. Check 'number of nodes below'.")
		return 1
	}

	visited_nodes[i] = true
	number_of_nodes_below := 0

	current_depth += 1
	for _, child_id := range (*tasks)[i].ChildrenIDs {
		slice_index := find_index(tasks, child_id)
		number_of_nodes_below += dfs_helper(slice_index, tasks)
	}
	current_depth -= 1

	number_of_nodes_below += 1
	linearized_tasks = append(linearized_tasks, LinearizedTaskReference{
		TaskId:                (*tasks)[i].Id,
		Task:                  &(*tasks)[i],
		Depth:                 current_depth,
		NumberOfChildren:      len((*tasks)[i].ChildrenIDs),
		NumberOfAllNodesBelow: number_of_nodes_below,
	})
	return number_of_nodes_below
}

func DFS(tasks []TaskModel) []LinearizedTaskReference {

	current_depth = 0

	// Empty initialize visited_nodes
	visited_nodes = make([]bool, len(tasks))
	for i := range visited_nodes {
		visited_nodes[i] = false
	}

	// Empty initialize visited_nodes
	in_progress = make([]bool, len(tasks))
	for i := range in_progress {
		in_progress[i] = false
	}

	for i := range tasks {
		dfs_helper(i, &tasks)
	}

	return linearized_tasks
}
