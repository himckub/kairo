package core

import (
	"reflect"
	"testing"
)

func TestTaskPatch_ApplyTo(t *testing.T) {
	task := Task{
		ID:          "1",
		Title:       "Original",
		Description: "Desc",
		Project:     "Proj",
		Tags:        []string{"a"},
		Priority:    P1,
		Status:      StatusTodo,
	}

	newTitle := "Updated"
	newTags := []string{"b", "c"}
	newStatus := StatusDone
	patch := TaskPatch{
		Title:  &newTitle,
		Tags:   &newTags,
		Status: &newStatus,
	}

	updated := patch.ApplyTo(task)

	if updated.Title != "Updated" {
		t.Errorf("expected title Updated, got %s", updated.Title)
	}
	if !reflect.DeepEqual(updated.Tags, []string{"b", "c"}) {
		t.Errorf("expected tags [b c], got %v", updated.Tags)
	}
	if updated.Status != StatusDone {
		t.Errorf("expected status Done, got %s", updated.Status)
	}
	// Verify immutability of original
	if task.Title != "Original" {
		t.Error("original task title was mutated")
	}
}
