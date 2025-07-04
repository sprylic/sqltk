package sqltk

import (
	"reflect"
	"testing"
)

func init() {
	SetDialect(Standard())
}

func TestUpdateBuilder(t *testing.T) {
	t.Run("basic update", func(t *testing.T) {
		q := Update("users").Set("name", "Alice").Where(NewStringCondition("id = ?", 1))
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET name = ? WHERE id = ?"
		wantArgs := []interface{}{"Alice", 1}
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

	t.Run("multiple sets", func(t *testing.T) {
		q := Update("users").Set("name", "Alice").Set("age", 30).Where(NewStringCondition("id = ?", 1))
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET name = ?, age = ? WHERE id = ?"
		wantArgs := []interface{}{"Alice", 30, 1}
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

	t.Run("raw set", func(t *testing.T) {
		q := Update("users").SetRaw("updated_at = NOW()")
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET updated_at = NOW()"
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

	t.Run("where raw", func(t *testing.T) {
		q := Update("users").Set("name", "Alice").Where(AsCondition(Raw("id = 1")))
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET name = ? WHERE id = 1"
		wantArgs := []interface{}{"Alice"}
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

	t.Run("where with condition builder", func(t *testing.T) {
		cond := NewCond().
			Equal("active", true).
			And(NewCond().GreaterThan("age", 18)).
			And(NewCond().In("status", "active", "pending"))

		q := Update("users").Set("name", "Alice").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET name = ? WHERE active = ? AND age > ? AND status IN (?, ?)"
		wantArgs := []interface{}{"Alice", true, 18, "active", "pending"}
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

	t.Run("where with complex condition builder", func(t *testing.T) {
		cond := NewCond().
			Equal("active", true).
			Or(NewCond().Equal("vip", true)).
			And(NewCond().GreaterThan("age", 16))

		q := Update("users").Set("name", "Alice").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET name = ? WHERE (active = ?) OR (vip = ?) AND age > ?"
		wantArgs := []interface{}{"Alice", true, true, 16}
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

	t.Run("where equal", func(t *testing.T) {
		q := Update("users").Set("active", true).WhereEqual("id", 1)
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET active = ? WHERE id = ?"
		wantArgs := []interface{}{true, 1}
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
		q := Update("users").Set("active", false).WhereNotEqual("id", 2)
		sql, args, err := q.Build()
		wantSQL := "UPDATE users SET active = ? WHERE id != ?"
		wantArgs := []interface{}{false, 2}
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
		q := Update("").Set("name", "Alice")
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on no sets", func(t *testing.T) {
		q := Update("users")
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on invalid where type", func(t *testing.T) {
		// This test demonstrates that the compiler will catch invalid types
		// We can't test this at runtime since it's a compile-time error
		t.Skip("This is now a compile-time error, not a runtime error")
	})
}

func TestPostgresUpdateBuilder_Returning(t *testing.T) {
	pq := NewPostgresUpdate("users")
	pq.UpdateBuilder = pq.UpdateBuilder.Set("name", "Alice").Where(NewStringCondition("id = ?", 1))
	pq = pq.Returning("id", "name")
	sql, args, err := pq.Build()
	wantSQL := "UPDATE users SET name = ? WHERE id = ? RETURNING id, name"
	wantArgs := []interface{}{"Alice", 1}
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
