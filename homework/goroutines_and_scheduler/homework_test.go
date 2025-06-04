package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	queue []*Task
	tasks map[int]*Task
}

func NewScheduler() Scheduler {
	return Scheduler{
		queue: make([]*Task, 0),
		tasks: make(map[int]*Task, 0),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks[task.Identifier] = &task
	s.queue = append(s.queue, &task)

	if len(s.queue) == 1 {
		return
	}

	i := len(s.queue) - 1
	for i > 0 {
		parent := (i - 1) / 2

		if s.queue[parent].Priority >= s.queue[i].Priority {
			break
		}

		s.queue[i], s.queue[parent] = s.queue[parent], s.queue[i]

		i = parent
	}
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.tasks[taskID].Priority = newPriority

	s.heapify(0)
}

func (s *Scheduler) GetTask() Task {
	task := s.queue[0]
	s.queue[0] = s.queue[len(s.queue)-1]
	s.queue = s.queue[:len(s.queue)-1]

	s.heapify(0)

	return *task
}

func (s *Scheduler) heapify(parent int) {
	for {
		size := len(s.queue)

		left := 2*parent + 1
		right := 2*parent + 2
		largest := parent

		if left < size && s.queue[left].Priority > s.queue[largest].Priority {
			largest = left
		}

		if right < size && s.queue[right].Priority > s.queue[largest].Priority {
			largest = right
		}

		if largest == parent {
			break
		}

		s.queue[parent], s.queue[largest] = s.queue[largest], s.queue[parent]

		parent = largest
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	// Так как мы inplace поменяли Priority, объекты различаются
	assert.Equal(t, task1.Identifier, task.Identifier)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
