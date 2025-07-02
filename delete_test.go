package cqb

import (
	"reflect"
	"testing"
)

func TestDeleteBuilder(t *testing.T) {
	t.Run("basic delete", func(t *testing.T) {
		q := Delete("users").Where("id = ?", 1)
		sql, args, err := q.Build()
		wantSQL := "DELETE FROM users WHERE id = ?"
		wantArgs := []interface{}{1}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if !reflect.DeepEqual(args, wantArgs) {
			t.Errorf("got args %v, want %v", args, wantArgs)
		}
	})

	t.Run("delete with raw where", func(t *testing.T) {
		q := Delete("users").Where(Raw("id = 1"))
		sql, args, err := q.Build()
		wantSQL := "DELETE FROM users WHERE id = 1"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != wantSQL {
			t.Errorf("got SQL %q, want %q", sql, wantSQL)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("error on missing table", func(t *testing.T) {
		q := Delete("").Where("id = ?", 1)
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on invalid where type", func(t *testing.T) {
		q := Delete("users").Where(123)
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})
}

func TestPostgresDeleteBuilder_Returning(t *testing.T) {
	pq := NewPostgresDelete("users")
	pq.DeleteBuilder = pq.DeleteBuilder.Where("id = ?", 1)
	pq = pq.Returning("id")
	sql, args, err := pq.Build()
	wantSQL := "DELETE FROM users WHERE id = ? RETURNING id"
	wantArgs := []interface{}{1}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sql != wantSQL {
		t.Errorf("got SQL %q, want %q", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("got args %v, want %v", args, wantArgs)
	}
}
