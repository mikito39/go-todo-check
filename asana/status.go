package asana

import "log"

// TaskStatus returns false if task is already closed
func TaskStatus(taskID string, taskListID string) bool {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.ListAllTasks(taskListID)
	if err != nil {
		log.Fatal(err)
	}
	for page := range res {
		for _, task := range page.Data {
			if task.ID == taskID && !task.Completed {
				return true
			}
		}
	}
	return false
}
