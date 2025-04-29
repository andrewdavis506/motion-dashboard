package service

import (
	"context"
	"log/slog"
	"sync"
	"task-dashboard/internal/api"
	"task-dashboard/internal/models"
	"time"
)

const defaultCacheExpiry = 1 * time.Minute

// TaskService handles business logic for tasks.
type TaskService struct {
	apiClient   *api.MotionClient
	cacheMutex  sync.RWMutex
	cachedTasks []models.Task
	lastUpdate  time.Time
	cacheExpiry time.Duration
}

// NewTaskService creates a new TaskService.
func NewTaskService(apiClient *api.MotionClient) *TaskService {
	return &TaskService{
		apiClient:   apiClient,
		cacheExpiry: defaultCacheExpiry,
	}
}

// GetTasks returns tasks from cache if valid, or refreshes from API.
func (s *TaskService) GetTasks() ([]models.Task, error) {
	s.cacheMutex.RLock()
	tasksValid := time.Since(s.lastUpdate) < s.cacheExpiry && len(s.cachedTasks) > 0
	cached := s.cachedTasks
	s.cacheMutex.RUnlock()

	if tasksValid {
		return cached, nil
	}
	return s.RefreshTasks()
}

// RefreshTasks refreshes the task cache.
func (s *TaskService) RefreshTasks() ([]models.Task, error) {
	tasks, err := s.apiClient.FetchTasks()
	if err != nil {
		return nil, err
	}

	s.cacheMutex.Lock()
	s.cachedTasks = tasks
	s.lastUpdate = time.Now()
	s.cacheMutex.Unlock()

	slog.Info("Tasks refreshed", "count", len(tasks))
	return tasks, nil
}

// GetDashboardData returns current and next task data.
func (s *TaskService) GetDashboardData() (*models.DashboardData, error) {
	tasks, err := s.GetTasks()
	if err != nil {
		return nil, err
	}

	currentTask, nextTask := s.findCurrentAndNextTask(tasks)

	return &models.DashboardData{
		CurrentTask: currentTask,
		NextTask:    nextTask,
	}, nil
}

// findCurrentAndNextTask identifies the current and next tasks.
func (s *TaskService) findCurrentAndNextTask(tasks []models.Task) (*models.TaskWithTiming, *models.TaskWithTiming) {
	now := time.Now()
	var currentTask *models.TaskWithTiming
	var nextTask *models.TaskWithTiming

	for _, task := range tasks {
		if task.ScheduledStart.IsZero() || task.ScheduledEnd.IsZero() {
			continue
		}

		if now.After(task.ScheduledStart) && now.Before(task.ScheduledEnd) {
			currentTask = &models.TaskWithTiming{
				ID:        task.ID,
				Name:      task.Name,
				StartDate: task.ScheduledStart,
				EndDate:   task.ScheduledEnd,
			}
		} else if now.Before(task.ScheduledStart) {
			if nextTask == nil || task.ScheduledStart.Before(nextTask.StartDate) {
				nextTask = &models.TaskWithTiming{
					ID:        task.ID,
					Name:      task.Name,
					StartDate: task.ScheduledStart,
					EndDate:   task.ScheduledEnd,
				}
			}
		}
	}

	for _, task := range tasks {
		for _, chunk := range task.Chunks {
			if chunk.ScheduledStart.IsZero() || chunk.ScheduledEnd.IsZero() {
				continue
			}

			if now.After(chunk.ScheduledStart) && now.Before(chunk.ScheduledEnd) {
				currentTask = &models.TaskWithTiming{
					ID:        chunk.ID,
					Name:      task.Name,
					StartDate: chunk.ScheduledStart,
					EndDate:   chunk.ScheduledEnd,
				}
			} else if now.Before(chunk.ScheduledStart) {
				if nextTask == nil || chunk.ScheduledStart.Before(nextTask.StartDate) {
					nextTask = &models.TaskWithTiming{
						ID:        chunk.ID,
						Name:      task.Name,
						StartDate: chunk.ScheduledStart,
						EndDate:   chunk.ScheduledEnd,
					}
				}
			}
		}
	}

	return currentTask, nextTask
}

// StartPeriodicRefresh refreshes tasks periodically.
func (s *TaskService) StartPeriodicRefresh(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if _, err := s.RefreshTasks(); err != nil {
				slog.Warn("Periodic task refresh failed", "error", err)
			}
		case <-ctx.Done():
			slog.Info("Stopping periodic refresh")
			return
		}
	}
}

// GetLastUpdateTime returns the last update timestamp.
func (s *TaskService) GetLastUpdateTime() time.Time {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	return s.lastUpdate
}
