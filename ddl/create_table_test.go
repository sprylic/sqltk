package ddl

import (
	"strings"
	"testing"

	"github.com/sprylic/sqltk/raw"
	"github.com/sprylic/sqltk/sqldialect"

	"github.com/sprylic/sqltk/mysqlfunc"
)

func TestCreateTableBuilder(t *testing.T) {
	t.Run("basic create table", func(t *testing.T) {
		q := CreateTable("users").WithDialect(sqldialect.NoQuoteIdent()).
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("name").Type("VARCHAR").Size(255).NotNull())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL, name VARCHAR(255) NOT NULL)"

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

	t.Run("create table with if not exists", func(t *testing.T) {
		q := CreateTable("users").
			IfNotExists().
			AddColumn(Column("id").Type("INT").NotNull())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE IF NOT EXISTS users (id INT NOT NULL)"

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

	t.Run("create temporary table", func(t *testing.T) {
		q := CreateTable("temp_users").
			Temporary().
			AddColumn(Column("id").Type("INT").NotNull())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TEMPORARY TABLE temp_users (id INT NOT NULL)"

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

	t.Run("create table with primary key", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("name").Type("VARCHAR").Size(255)).
			PrimaryKey("id")

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL, name VARCHAR(255), PRIMARY KEY (id))"

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

	t.Run("create table with column primary key", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull().PrimaryKey()).
			AddColumn(Column("name").Type("VARCHAR").Size(255))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL, name VARCHAR(255), PRIMARY KEY (id))"

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

	t.Run("create table with multiple column primary keys", func(t *testing.T) {
		q := CreateTable("user_roles").
			AddColumn(Column("user_id").Type("INT").NotNull().PrimaryKey()).
			AddColumn(Column("role_id").Type("INT").NotNull().PrimaryKey())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE user_roles (user_id INT NOT NULL, role_id INT NOT NULL, PRIMARY KEY (user_id, role_id))"

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

	t.Run("create table with column unique constraint", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").AutoIncrement().NotNull().PrimaryKey()).
			AddColumn(Column("email").Type("VARCHAR").Size(255).NotNull().Unique())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL AUTO_INCREMENT, email VARCHAR(255) NOT NULL, PRIMARY KEY (id), UNIQUE (email))"

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

	t.Run("create table with multiple unique columns", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").AutoIncrement().NotNull().PrimaryKey()).
			AddColumn(Column("email").Type("VARCHAR").Size(255).NotNull().Unique()).
			AddColumn(Column("username").Type("VARCHAR").Size(100).NotNull().Unique())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL AUTO_INCREMENT, email VARCHAR(255) NOT NULL, username VARCHAR(100) NOT NULL, PRIMARY KEY (id), UNIQUE (email), UNIQUE (username))"

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

	t.Run("create table with unique constraint", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("email").Type("VARCHAR").Size(255)).
			Unique("idx_email", "email")

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL, email VARCHAR(255), CONSTRAINT idx_email UNIQUE (email))"

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

	t.Run("create table with check constraint", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("age").Type("INT")).
			Check("chk_age", "age >= 0")

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL, age INT, CONSTRAINT chk_age CHECK (age >= 0))"

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

	t.Run("create table with index", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("name").Type("VARCHAR").Size(255)).
			Index("idx_name", "name")

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL, name VARCHAR(255), INDEX idx_name (name))"

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

	t.Run("create table with foreign key", func(t *testing.T) {
		q := CreateTable("orders").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("user_id").Type("INT")).
			AddForeignKey(
				ForeignKey("fk_orders_user", "user_id").
					References("users", "id").
					OnDelete("CASCADE"),
			)

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE orders (id INT NOT NULL, user_id INT, CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE)"

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

	t.Run("create table with table options", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			Charset("utf8mb4").
			Collation("utf8mb4_unicode_ci").
			Comment("User accounts table").
			Engine("InnoDB")

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'User accounts table' ENGINE InnoDB"

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

	t.Run("create table with table-level foreign key", func(t *testing.T) {
		q := CreateTable("orders").
			AddColumn(Column("id").Type("INT").NotNull().PrimaryKey()).
			AddColumn(Column("user_id").Type("INT")).
			AddForeignKey(ForeignKey("fk_orders_user", "user_id").References("users", "id").OnDelete("CASCADE").OnUpdate("RESTRICT"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE orders (id INT NOT NULL, user_id INT, CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE RESTRICT, PRIMARY KEY (id))"

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

func TestColumnBuilder(t *testing.T) {
	t.Run("column with size", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("name").Type("VARCHAR").Size(100).NotNull())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (name VARCHAR(100) NOT NULL)"

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

	t.Run("column with precision and scale", func(t *testing.T) {
		q := CreateTable("products").
			AddColumn(Column("price").Type("DECIMAL").Precision(10, 2).NotNull())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE products (price DECIMAL(10,2) NOT NULL)"

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

	t.Run("column with default value", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("status").Type("VARCHAR").Size(20).Default("active")).
			AddColumn(Column("created_at").Type("TIMESTAMP").Default(raw.Raw("CURRENT_TIMESTAMP")))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (status VARCHAR(20) DEFAULT 'active', created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"

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

	t.Run("column with auto increment", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull().AutoIncrement())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL AUTO_INCREMENT)"

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

	t.Run("column with big auto increment", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("BIGINT").NotNull().AutoIncrement())

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id BIGINT NOT NULL AUTO_INCREMENT)"

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

	t.Run("column with charset and collation", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("name").Type("VARCHAR").Size(255).Charset("utf8mb4").Collation("utf8mb4_unicode_ci"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci)"

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

	t.Run("column with comment", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull().Comment("Primary key"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT NOT NULL COMMENT 'Primary key')"

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

	t.Run("column with raw SQL default", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("status").Type("VARCHAR").Default("active")).                           // String literal
			AddColumn(Column("created_at").Type("TIMESTAMP").Default(raw.Raw("CURRENT_TIMESTAMP"))). // Raw SQL
			AddColumn(Column("count").Type("INT").Default(0))                                        // Number literal

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (status VARCHAR DEFAULT 'active', created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, count INT DEFAULT 0)"

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

func TestCreateTableBuilder_Errors(t *testing.T) {
	t.Run("empty table name", func(t *testing.T) {
		_, _, err := CreateTable("").Build()
		if err == nil {
			t.Errorf("expected error for empty table name, got none")
		}
	})

	t.Run("no columns", func(t *testing.T) {
		_, _, err := CreateTable("users").Build()
		if err == nil {
			t.Errorf("expected error for no columns, got none")
		}
	})

	t.Run("column without type", func(t *testing.T) {
		q := CreateTable("users").AddColumn(Column("id"))
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for column without type, got none")
		}
	})

	t.Run("empty column name", func(t *testing.T) {
		q := CreateTable("users").AddColumn(Column(""))
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for empty column name, got none")
		}
	})

	t.Run("invalid size", func(t *testing.T) {
		q := CreateTable("users").AddColumn(Column("name").Type("VARCHAR").Size(0))
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for invalid size, got none")
		}
	})

	t.Run("invalid precision", func(t *testing.T) {
		q := CreateTable("users").AddColumn(Column("price").Type("DECIMAL").Precision(0, 2))
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for invalid precision, got none")
		}
	})

	t.Run("invalid scale", func(t *testing.T) {
		q := CreateTable("users").AddColumn(Column("price").Type("DECIMAL").Precision(10, 11))
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for invalid scale, got none")
		}
	})

	t.Run("primary key without columns", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			PrimaryKey()
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for primary key without columns, got none")
		}
	})

	t.Run("unique constraint without columns", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			Unique("idx_test")
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for unique constraint without columns, got none")
		}
	})

	t.Run("check constraint without expression", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			Check("chk_test", "")
		_, _, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		if err == nil {
			t.Errorf("expected error for check constraint without expression, got none")
		}
	})
}

func TestCreateTableBuilder_Dialect(t *testing.T) {
	t.Run("MySQL dialect", func(t *testing.T) {
		sqldialect.SetDialect(sqldialect.MySQL())
		defer sqldialect.SetDialect(sqldialect.NoQuoteIdent())

		q := CreateTable("users").
			AddColumns(
				Column("id").Type("int unsigned").NotNull(),
				Column("name").Type("VARCHAR").Size(255),
			)

		sql, args, err := q.Build()
		wantSQL := "CREATE TABLE `users` (`id` INT UNSIGNED NOT NULL, `name` VARCHAR(255))"

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

	t.Run("Postgres dialect", func(t *testing.T) {
		sqldialect.SetDialect(sqldialect.Postgres())
		defer sqldialect.SetDialect(sqldialect.NoQuoteIdent())

		q := CreateTable("users").
			AddColumn(Column("id").Type("INTEGER").NotNull()).
			AddColumn(Column("name").Type("VARCHAR").Size(255))

		sql, args, err := q.Build()
		wantSQL := "CREATE TABLE \"users\" (\"id\" INTEGER NOT NULL, \"name\" VARCHAR(255))"

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

func TestCreateTableBuilder_Complex(t *testing.T) {
	t.Run("complex table with all features", func(t *testing.T) {
		q := CreateTable("users").
			IfNotExists().
			AddColumn(Column("id").Type("INT").NotNull().AutoIncrement().Comment("Primary key")).
			AddColumn(Column("username").Type("VARCHAR").Size(50).NotNull()).
			AddColumn(Column("email").Type("VARCHAR").Size(255).NotNull()).
			AddColumn(Column("password_hash").Type("VARCHAR").Size(255).NotNull()).
			AddColumn(Column("age").Type("INT").Default(18)).
			AddColumn(Column("created_at").Type("TIMESTAMP").Default(raw.Raw("CURRENT_TIMESTAMP"))).
			AddColumn(Column("updated_at").Type("TIMESTAMP").Default(raw.Raw("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"))).
			PrimaryKey("id").
			Unique("idx_username", "username").
			Unique("idx_email", "email").
			Check("chk_age", "age >= 0 AND age <= 150").
			Index("idx_created_at", "created_at").
			Charset("utf8mb4").
			Collation("utf8mb4_unicode_ci").
			Comment("User accounts table").
			Engine("InnoDB")

		sql, args, err := q.Build()
		wantSQL := "CREATE TABLE IF NOT EXISTS users (id INT NOT NULL AUTO_INCREMENT COMMENT 'Primary key', username VARCHAR(50) NOT NULL, email VARCHAR(255) NOT NULL, password_hash VARCHAR(255) NOT NULL, age INT DEFAULT 18, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id), CONSTRAINT idx_username UNIQUE (username), CONSTRAINT idx_email UNIQUE (email), CONSTRAINT chk_age CHECK (age >= 0 AND age <= 150), INDEX idx_created_at (created_at)) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'User accounts table' ENGINE InnoDB"

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

func TestCreateTableBuilder_Postgres(t *testing.T) {
	t.Run("basic create table (postgres)", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("name").Type("VARCHAR").Size(255).NotNull()).
			WithDialect(sqldialect.Postgres())

		sql, args, err := q.Build()
		wantSQL := "CREATE TABLE \"users\" (\"id\" INT NOT NULL, \"name\" VARCHAR(255) NOT NULL)"

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

	t.Run("create table with table options (postgres)", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			Charset("utf8mb4").
			Collation("utf8mb4_unicode_ci").
			Comment("User accounts table").
			Engine("InnoDB").
			WithDialect(sqldialect.Postgres())

		sql, args, err := q.Build()
		wantSQL := "CREATE TABLE \"users\" (\"id\" INT NOT NULL) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'User accounts table' ENGINE InnoDB"

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

	t.Run("create table with unique constraint (postgres)", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("email").Type("VARCHAR").Size(255)).
			Unique("idx_email", "email").
			WithDialect(sqldialect.Postgres())

		sql, args, err := q.Build()
		wantSQL := "CREATE TABLE \"users\" (\"id\" INT NOT NULL, \"email\" VARCHAR(255), CONSTRAINT \"idx_email\" UNIQUE (\"email\"))"

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

	t.Run("create table with foreign key (postgres)", func(t *testing.T) {
		q := CreateTable("orders").
			AddColumn(Column("id").Type("INT").NotNull()).
			AddColumn(Column("user_id").Type("INT")).
			AddForeignKey(
				ForeignKey("fk_orders_user", "user_id").
					References("users", "id").
					OnDelete("CASCADE"),
			)

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE TABLE \"orders\" (\"id\" INT NOT NULL, \"user_id\" INT, CONSTRAINT \"fk_orders_user\" FOREIGN KEY (\"user_id\") REFERENCES \"users\" (\"id\") ON DELETE CASCADE)"

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

	t.Run("create table with auto increment (postgres)", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").NotNull().AutoIncrement())

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE TABLE \"users\" (\"id\" SERIAL NOT NULL)"

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

	t.Run("create table with big auto increment (postgres)", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("BIGINT").NotNull().AutoIncrement())

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := "CREATE TABLE \"users\" (\"id\" BIGSERIAL NOT NULL)"

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

func TestCreateTable_OnUpdateOnDeleteWithSqlFunc(t *testing.T) {
	t.Run("on update and on delete with string actions", func(t *testing.T) {
		q := CreateTable("orders").
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("user_id").Type("INT")).
			AddForeignKey(
				ForeignKey("fk_orders_user", "user_id").
					References("users", "id").
					OnDelete("CASCADE").
					OnUpdate("RESTRICT"),
			)

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE orders (id INT, user_id INT, CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE RESTRICT, PRIMARY KEY (id))"
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

	t.Run("on update and on delete with other string actions", func(t *testing.T) {
		q := CreateTable("orders").
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("user_id").Type("INT")).
			AddForeignKey(
				ForeignKey("fk_orders_user", "user_id").
					References("users", "id").
					OnDelete("SET NULL").
					OnUpdate("NO ACTION"),
			)

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE orders (id INT, user_id INT, CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL ON UPDATE NO ACTION, PRIMARY KEY (id))"
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

func TestCreateTable_ColumnOnUpdate(t *testing.T) {
	t.Run("column on update with string", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP"))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT, updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id))"
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

	t.Run("column on update with sqlfunc", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("updated_at").Type("TIMESTAMP").OnUpdate(mysqlfunc.CurrentTimestamp()))

		sql, args, err := q.WithDialect(sqldialect.NoQuoteIdent()).Build()
		wantSQL := "CREATE TABLE users (id INT, updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id))"
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

func TestCreateTable_OnUpdate_DialectSpecific(t *testing.T) {
	t.Run("MySQL OnUpdate", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP"))

		sql, args, err := q.WithDialect(sqldialect.MySQL()).Build()
		wantSQL := "CREATE TABLE `users` (`id` INT, `updated_at` TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`))"
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

	t.Run("PostgreSQL OnUpdate", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP"))

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		wantSQL := `CREATE TABLE "users" ("id" INT, "updated_at" TIMESTAMP, PRIMARY KEY ("id"));

CREATE OR REPLACE FUNCTION "users_updated_at_update_trigger"()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER "tr_users_updated_at_update"
    BEFORE UPDATE ON "users"
    FOR EACH ROW
    EXECUTE FUNCTION "users_updated_at_update_trigger"();`
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

	t.Run("PostgreSQL OnUpdate with multiple columns", func(t *testing.T) {
		q := CreateTable("users").
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP")).
			AddColumn(Column("modified_at").Type("TIMESTAMP").OnUpdate("NOW()"))

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		// Should contain both trigger functions and triggers
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(sql, `"users_updated_at_update_trigger"`) {
			t.Error("missing trigger function for updated_at")
		}
		if !strings.Contains(sql, `"users_modified_at_update_trigger"`) {
			t.Error("missing trigger function for modified_at")
		}
		if !strings.Contains(sql, `"tr_users_updated_at_update"`) {
			t.Error("missing trigger for updated_at")
		}
		if !strings.Contains(sql, `"tr_users_modified_at_update"`) {
			t.Error("missing trigger for modified_at")
		}
		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})
}

func TestCreateTable_OnUpdate_IfNotExists(t *testing.T) {
	t.Run("PostgreSQL OnUpdate with IfNotExists", func(t *testing.T) {
		q := CreateTable("users").
			IfNotExists().
			AddColumn(Column("id").Type("INT").PrimaryKey()).
			AddColumn(Column("updated_at").Type("TIMESTAMP").OnUpdate("CURRENT_TIMESTAMP"))

		sql, args, err := q.WithDialect(sqldialect.Postgres()).Build()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should contain IF NOT EXISTS in table creation
		if !strings.Contains(sql, "IF NOT EXISTS") {
			t.Error("missing IF NOT EXISTS in table creation")
		}

		// Should contain DO block for trigger creation
		if !strings.Contains(sql, "DO $$") {
			t.Error("missing DO block for trigger creation")
		}

		// Should contain trigger existence check
		if !strings.Contains(sql, "pg_trigger WHERE tgname =") {
			t.Error("missing trigger existence check")
		}

		if len(args) != 0 {
			t.Errorf("got args %v, want none", args)
		}
	})
}
