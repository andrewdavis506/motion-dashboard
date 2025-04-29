package models

import (
	"time"
)

// Task represents a task from the Motion API
type Task struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	ScheduledStart time.Time `json:"scheduledStart"`
	ScheduledEnd   time.Time `json:"scheduledEnd"`
	Chunks         []Chunk   `json:"chunks"`
}

// Chunk represents a time chunk for a task
type Chunk struct {
	ID             string    `json:"id"`
	ScheduledStart time.Time `json:"scheduledStart"`
	ScheduledEnd   time.Time `json:"scheduledEnd"`
}

// TaskWithTiming represents a task with additional timing information
type TaskWithTiming struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// TasksResponse represents the API response from Motion
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
	Meta  Meta   `json:"meta"`
}

// Meta contains pagination information
type Meta struct {
	NextCursor string `json:"nextCursor"`
}

// DashboardData represents the data needed for the dashboard UI
type DashboardData struct {
	CurrentTask *TaskWithTiming `json:"currentTask"`
	NextTask    *TaskWithTiming `json:"nextTask"`
}
