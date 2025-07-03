package stk

import (
	"reflect"
	"testing"
)

func init() {
	SetDialect(Standard())
}

func TestDeleteBuilder(t *testing.T) {
	t.Run("basic delete", func(t *testing.T) {
		q := Delete("users").Where(NewStringCondition("id = ?", 1))
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

	t.Run("delete with condition builder", func(t *testing.T) {
		cond := NewCond().
			Equal("active", false).
			Or(NewCond().IsNull("deleted_at")).
			And(NewCond().LessThan("created_at", "2023-01-01"))

		q := Delete("users").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "DELETE FROM users WHERE (active = ?) OR (deleted_at IS NULL) AND created_at < ?"
		wantArgs := []interface{}{false, "2023-01-01"}
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

	t.Run("delete with exists condition", func(t *testing.T) {
		subq := Select("1").From("orders").Where(NewStringCondition("orders.user_id = users.id"))
		cond := NewCond().Exists(subq)

		q := Delete("users").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "DELETE FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
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

	t.Run("where equal", func(t *testing.T) {
		q := Delete("users").WhereEqual("active", true)
		sql, args, err := q.Build()
		wantSQL := "DELETE FROM users WHERE active = ?"
		wantArgs := []interface{}{true}
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

	t.Run("where not equal", func(t *testing.T) {
		q := Delete("users").WhereNotEqual("active", false)
		sql, args, err := q.Build()
		wantSQL := "DELETE FROM users WHERE active != ?"
		wantArgs := []interface{}{false}
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

	t.Run("error on missing table", func(t *testing.T) {
		q := Delete("").Where(NewStringCondition("id = ?", 1))
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
	pq.DeleteBuilder = pq.DeleteBuilder.Where(NewStringCondition("id = ?", 1))
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
