package sqltk

import (
	"reflect"
	"testing"
)

func init() {
	SetDialect(NoQuoteIdent())
}

func TestConditionBuilder_Basic(t *testing.T) {
	t.Run("empty condition", func(t *testing.T) {
		cond := NewCond()
		sql, args, err := cond.Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if sql != "" {
			t.Errorf("got SQL %q, want empty string", sql)
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})

	t.Run("simple equality", func(t *testing.T) {
		cond := NewCond().Equal("id", 1)
		sql, args, err := cond.Build()
		wantSQL := "id = ?"
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

	t.Run("simple inequality", func(t *testing.T) {
		cond := NewCond().NotEqual("active", false)
		sql, args, err := cond.Build()
		wantSQL := "active != ?"
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

	t.Run("null equality", func(t *testing.T) {
		cond := NewCond().Equal("deleted_at", nil)
		sql, args, err := cond.Build()
		wantSQL := "deleted_at IS NULL"
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

	t.Run("null inequality", func(t *testing.T) {
		cond := NewCond().NotEqual("created_at", nil)
		sql, args, err := cond.Build()
		wantSQL := "created_at IS NOT NULL"
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
}

func TestConditionBuilder_Comparison(t *testing.T) {
	t.Run("greater than", func(t *testing.T) {
		cond := NewCond().GreaterThan("age", 18)
		sql, args, err := cond.Build()
		wantSQL := "age > ?"
		wantArgs := []interface{}{18}
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

	t.Run("greater than or equal", func(t *testing.T) {
		cond := NewCond().GreaterThanOrEqual("score", 100)
		sql, args, err := cond.Build()
		wantSQL := "score >= ?"
		wantArgs := []interface{}{100}
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

	t.Run("less than", func(t *testing.T) {
		cond := NewCond().LessThan("price", 50.0)
		sql, args, err := cond.Build()
		wantSQL := "price < ?"
		wantArgs := []interface{}{50.0}
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

	t.Run("less than or equal", func(t *testing.T) {
		cond := NewCond().LessThanOrEqual("quantity", 10)
		sql, args, err := cond.Build()
		wantSQL := "quantity <= ?"
		wantArgs := []interface{}{10}
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
}

func TestConditionBuilder_Like(t *testing.T) {
	t.Run("like pattern", func(t *testing.T) {
		cond := NewCond().Like("name", "%john%")
		sql, args, err := cond.Build()
		wantSQL := "name LIKE ?"
		wantArgs := []interface{}{"%john%"}
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

	t.Run("not like pattern", func(t *testing.T) {
		cond := NewCond().NotLike("email", "%@spam.com")
		sql, args, err := cond.Build()
		wantSQL := "email NOT LIKE ?"
		wantArgs := []interface{}{"%@spam.com"}
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
}

func TestConditionBuilder_In(t *testing.T) {
	t.Run("in values", func(t *testing.T) {
		cond := NewCond().In("status", "active", "pending", "approved")
		sql, args, err := cond.Build()
		wantSQL := "status IN (?, ?, ?)"
		wantArgs := []interface{}{"active", "pending", "approved"}
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

	t.Run("not in values", func(t *testing.T) {
		cond := NewCond().NotIn("category", "deleted", "archived")
		sql, args, err := cond.Build()
		wantSQL := "category NOT IN (?, ?)"
		wantArgs := []interface{}{"deleted", "archived"}
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

	t.Run("in empty list error", func(t *testing.T) {
		_, _, err := NewCond().In("status").Build()
		if err == nil {
			t.Errorf("expected error for empty IN list, got none")
		}
	})
}

func TestConditionBuilder_Between(t *testing.T) {
	t.Run("between values", func(t *testing.T) {
		cond := NewCond().Between("age", 18, 65)
		sql, args, err := cond.Build()
		wantSQL := "age BETWEEN ? AND ?"
		wantArgs := []interface{}{18, 65}
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

	t.Run("not between values", func(t *testing.T) {
		cond := NewCond().NotBetween("price", 10.0, 100.0)
		sql, args, err := cond.Build()
		wantSQL := "price NOT BETWEEN ? AND ?"
		wantArgs := []interface{}{10.0, 100.0}
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
}

func TestConditionBuilder_Null(t *testing.T) {
	t.Run("is null", func(t *testing.T) {
		cond := NewCond().IsNull("deleted_at")
		sql, args, err := cond.Build()
		wantSQL := "deleted_at IS NULL"
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

	t.Run("is not null", func(t *testing.T) {
		cond := NewCond().IsNotNull("created_at")
		sql, args, err := cond.Build()
		wantSQL := "created_at IS NOT NULL"
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
}

func TestConditionBuilder_Exists(t *testing.T) {
	t.Run("exists subquery", func(t *testing.T) {
		subq := Select("1").From("orders").Where(NewStringCondition("orders.user_id = users.id"))
		cond := NewCond().Exists(subq)
		sql, args, err := cond.Build()
		wantSQL := "EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
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

	t.Run("not exists subquery", func(t *testing.T) {
		subq := Select("1").From("deleted_users").
			Where(NewCond().Raw("deleted_users.id = users.id"))
		cond := NewCond().NotExists(subq)
		sql, args, err := cond.Build()
		wantSQL := "NOT EXISTS (SELECT 1 FROM deleted_users WHERE deleted_users.id = users.id)"
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

	t.Run("exists raw", func(t *testing.T) {
		cond := NewCond().Exists(Raw("SELECT 1 FROM orders WHERE user_id = 1"))
		sql, args, err := cond.Build()
		wantSQL := "EXISTS (SELECT 1 FROM orders WHERE user_id = 1)"
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
}

func TestConditionBuilder_Combination(t *testing.T) {
	t.Run("and combination", func(t *testing.T) {
		cond1 := NewCond().Equal("active", true)
		cond2 := NewCond().GreaterThan("age", 18)
		combined := cond1.And(cond2)

		sql, args, err := combined.Build()
		wantSQL := "active = ? AND age > ?"
		wantArgs := []interface{}{true, 18}
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

	t.Run("or combination", func(t *testing.T) {
		cond1 := NewCond().Equal("status", "active")
		cond2 := NewCond().Equal("status", "pending")
		combined := cond1.Or(cond2)

		sql, args, err := combined.Build()
		wantSQL := "(status = ?) OR (status = ?)"
		wantArgs := []interface{}{"active", "pending"}
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

	t.Run("complex combination", func(t *testing.T) {
		cond1 := NewCond().Equal("active", true).And(NewCond().GreaterThan("age", 18))
		cond2 := NewCond().Equal("vip", true).And(NewCond().GreaterThan("age", 16))
		combined := cond1.Or(cond2)

		sql, args, err := combined.Build()
		wantSQL := "(active = ? AND age > ?) OR (vip = ? AND age > ?)"
		wantArgs := []interface{}{true, 18, true, 16}
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
}

func TestConditionBuilder_Case(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		cond := NewCond().Case().
			When(NewCond().Equal("status", "active"), "Active User").
			When(NewCond().Equal("status", "pending"), "Pending User").
			Else("Unknown User").
			End()

		sql, args, err := cond.Build()
		wantSQL := "CASE WHEN status = ? THEN ? WHEN status = ? THEN ? ELSE ? END"
		wantArgs := []interface{}{"active", "Active User", "pending", "Pending User", "Unknown User"}
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

	t.Run("case without else", func(t *testing.T) {
		cond := NewCond().Case().
			When(NewCond().GreaterThan("score", 90), "A").
			When(NewCond().GreaterThan("score", 80), "B").
			End()

		sql, args, err := cond.Build()
		wantSQL := "CASE WHEN score > ? THEN ? WHEN score > ? THEN ? END"
		wantArgs := []interface{}{90, "A", 80, "B"}
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
}

func TestConditionBuilder_Integration(t *testing.T) {
	t.Run("with select builder", func(t *testing.T) {
		cond := NewCond().
			Equal("active", true).
			And(NewCond().GreaterThan("age", 18)).
			And(NewCond().In("status", "active", "pending"))

		q := Select("id", "name").From("users").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "SELECT id, name FROM users WHERE active = ? AND age > ? AND status IN (?, ?)"
		wantArgs := []interface{}{true, 18, "active", "pending"}
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

	t.Run("with having clause", func(t *testing.T) {
		cond := NewCond().GreaterThan("COUNT(*)", 1)

		q := Select("category").From("products").GroupBy("category").Having(cond)
		sql, args, err := q.Build()
		wantSQL := "SELECT category FROM products GROUP BY category HAVING COUNT(*) > ?"
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
}

func TestConditionBuilder_Dialect(t *testing.T) {
	t.Run("with mysql dialect", func(t *testing.T) {
		cond := NewCond().WithDialect(MySQL()).Equal("user.name", "john")
		sql, args, err := cond.Build()
		wantSQL := "`user`.`name` = ?"
		wantArgs := []interface{}{"john"}
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

	t.Run("with postgres dialect", func(t *testing.T) {
		cond := NewCond().WithDialect(Postgres()).Equal("user.name", "john")
		sql, args, err := cond.Build()
		wantSQL := "\"user\".\"name\" = ?"
		wantArgs := []interface{}{"john"}
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
}

func TestConditionBuilder_String(t *testing.T) {
	t.Run("debug string", func(t *testing.T) {
		cond := NewCond().Equal("name", "john").And(NewCond().GreaterThan("age", 18))
		debugStr := cond.GetUnsafeString()
		wantStr := "name = 'john' AND age > 18"
		if debugStr != wantStr {
			t.Errorf("got debug string %q, want %q", debugStr, wantStr)
		}
	})
}

func TestConditionInterface(t *testing.T) {
	t.Run("string condition", func(t *testing.T) {
		cond := NewStringCondition("active = ? AND age > ?", true, 18)
		sql, args, err := cond.BuildCondition()
		wantSQL := "active = ? AND age > ?"
		wantArgs := []interface{}{true, 18}
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

	t.Run("raw condition", func(t *testing.T) {
		cond := NewRawCondition(Raw("id = 1"))
		sql, args, err := cond.BuildCondition()
		wantSQL := "id = 1"
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

	t.Run("condition builder as condition", func(t *testing.T) {
		cond := NewCond().Equal("active", true).And(NewCond().GreaterThan("age", 18))
		sql, args, err := cond.BuildCondition()
		wantSQL := "active = ? AND age > ?"
		wantArgs := []interface{}{true, 18}
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
}

func TestTypeSafeWhere(t *testing.T) {
	t.Run("where with string condition", func(t *testing.T) {
		cond := NewStringCondition("active = ? AND age > ?", true, 18)
		q := Select("id").From("users").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "SELECT id FROM users WHERE active = ? AND age > ?"
		wantArgs := []interface{}{true, 18}
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

	t.Run("where with raw condition", func(t *testing.T) {
		cond := NewRawCondition(Raw("id = 1"))
		q := Select("id").From("users").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "SELECT id FROM users WHERE id = 1"
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

	t.Run("where with condition builder", func(t *testing.T) {
		cond := NewCond().Equal("active", true).And(NewCond().GreaterThan("age", 18))
		q := Select("id").From("users").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "SELECT id FROM users WHERE active = ? AND age > ?"
		wantArgs := []interface{}{true, 18}
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

	t.Run("raw where now requires AsCondition", func(t *testing.T) {
		cond := AsCondition(Raw("id = 1"))
		q := Select("id").From("users").Where(cond)
		sql, args, err := q.Build()
		wantSQL := "SELECT id FROM users WHERE id = 1"
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

	t.Run("invalid type now requires proper condition", func(t *testing.T) {
		// This test demonstrates that the compiler will catch invalid types
		// We can't test this at runtime since it's a compile-time error
		t.Skip("This is now a compile-time error, not a runtime error")
	})
}
