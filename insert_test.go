package sqltk

import (
	"encoding/json"
	"reflect"
	"testing"
)

func init() {
	SetDialect(NoQuoteIdent())
}

func TestInsertBuilder(t *testing.T) {
	t.Run("single row insert", func(t *testing.T) {
		q := Insert("users").Columns("id", "name").Values(1, "Alice")
		sql, args, err := q.Build()
		wantSQL := "INSERT INTO users (id, name) VALUES (?, ?)"
		wantArgs := []interface{}{1, "Alice"}
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

	t.Run("multi row insert", func(t *testing.T) {
		q := Insert("users").Columns("id", "name").Values(1, "Alice").Values(2, "Bob")
		sql, args, err := q.Build()
		wantSQL := "INSERT INTO users (id, name) VALUES (?, ?), (?, ?)"
		wantArgs := []interface{}{1, "Alice", 2, "Bob"}
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

	t.Run("error on mismatched values", func(t *testing.T) {
		q := Insert("users").Columns("id", "name").Values(1)
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on missing table", func(t *testing.T) {
		q := Insert("").Columns("id").Values(1)
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on missing columns", func(t *testing.T) {
		q := Insert("users").Values(1)
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("error on missing values", func(t *testing.T) {
		q := Insert("users").Columns("id")
		_, _, err := q.Build()
		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("debugsql single row", func(t *testing.T) {
		q := Insert("users").Columns("id", "name").Values(1, "Alice")
		got := q.DebugSQL()
		want := "INSERT INTO users (id, name) VALUES (1, 'Alice')"
		if got != want {
			t.Errorf("DebugSQL got %q, want %q", got, want)
		}
	})

	t.Run("debugsql multi row", func(t *testing.T) {
		q := Insert("users").Columns("id", "name").Values(1, "Alice").Values(2, "Bob")
		got := q.DebugSQL()
		want := "INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob')"
		if got != want {
			t.Errorf("DebugSQL got %q, want %q", got, want)
		}
	})
}

func TestPostgresInsertBuilder_Returning(t *testing.T) {
	pq := NewPostgresInsert("users")
	pq.InsertBuilder = pq.InsertBuilder.Columns("name", "age").Values("Alice", 30)
	pq = pq.Returning("id", "age")
	sql, args, err := pq.Build()
	wantSQL := "INSERT INTO \"users\" (\"name\", \"age\") VALUES ($1, $2) RETURNING id, age"
	wantArgs := []interface{}{"Alice", 30}
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

func TestPostgresInsertBuilder_PGJSON(t *testing.T) {
	pq := NewPostgresInsert("users")
	jsonVal := map[string]interface{}{"foo": 1, "bar": []int{2, 3}}
	pq.InsertBuilder = pq.InsertBuilder.Columns("name", "data").Values("Alice", PGJSON{V: jsonVal})
	pq = pq.Returning("id")
	sql, args, err := pq.Build()
	wantSQL := "INSERT INTO \"users\" (\"name\", \"data\") VALUES ($1, $2) RETURNING id"
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sql != wantSQL {
		t.Errorf("got SQL %q, want %q", sql, wantSQL)
	}
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	jsonArg, ok := args[1].(PGJSON)
	if !ok {
		t.Fatalf("expected PGJSON for arg, got %T", args[1])
	}
	jsonBytes, err := jsonArg.Value()
	if err != nil {
		t.Fatalf("PGJSON.Value() error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(jsonBytes.([]byte), &got); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if got["foo"].(float64) != 1 {
		t.Errorf("expected foo=1, got %v", got["foo"])
	}
}

func TestPostgresInsertBuilder_PGArray(t *testing.T) {
	pq := NewPostgresInsert("users")
	arrVal := []string{"foo", "bar"}
	pq.InsertBuilder = pq.InsertBuilder.Columns("tags").Values(PGArray{V: arrVal})
	pq = pq.Returning("id")
	sql, args, err := pq.Build()
	wantSQL := "INSERT INTO \"users\" (\"tags\") VALUES ($1) RETURNING id"
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sql != wantSQL {
		t.Errorf("got SQL %q, want %q", sql, wantSQL)
	}
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(args))
	}
	pgArr, ok := args[0].(PGArray)
	if !ok {
		t.Fatalf("expected PGArray for arg, got %T", args[0])
	}
	if got, want := pgArr.V, arrVal; !reflect.DeepEqual(got, want) {
		t.Errorf("expected array %v, got %v", want, got)
	}
}
