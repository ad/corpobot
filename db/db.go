package db

import (
	s "database/sql"
	"fmt"
	"reflect"
	"time"

	// "strings"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"

	_ "github.com/mattn/go-sqlite3" // Register some sql
)

// TelegramMessage ...
type TelegramMessage struct {
	ID         int
	TelegramID int
	FirstName  string
	LastName   string
	UserName   string
	Message    string
	IsBot      bool
	Date       time.Time
}

// InitDB ...
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "db/corpobot.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS telegram_messages (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"telegram_id" INTEGER NOT NULL,
		"message" TEXT DEFAULT "",
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "telegram_messages_user_id" FOREIGN KEY ("telegram_id") REFERENCES "users" ("telegram_id")
	);`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "users" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"first_name" VARCHAR(32) NOT NULL,
		"last_name" VARCHAR(32) NOT NULL,
		"user_name" VARCHAR(32) NOT NULL,
		"telegram_id" INTEGER NOT NULL,
		"role" VARCHAR(32) NOT NULL DEFAULT "new",
		"is_bot" bool NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "users_telegram_id" UNIQUE ("telegram_id") ON CONFLICT IGNORE
	  );

		CREATE TRIGGER IF NOT EXISTS users_updated_at_Trigger
		AFTER UPDATE On users
		BEGIN
		   UPDATE users SET updated_at = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
		END;`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "groups" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" VARCHAR(32) NOT NULL,
		"state" VARCHAR(32) NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "groups_name" UNIQUE ("name") ON CONFLICT IGNORE
	  );

		CREATE TRIGGER IF NOT EXISTS groups_updated_at_Trigger
		AFTER UPDATE On groups
		BEGIN
		   UPDATE groups SET updated_at = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
		END;`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "groupchats" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"title" TEXT DEFAULT "",
		"telegram_id" INTEGER NOT NULL,
		"invite_link" TEXT DEFAULT "",
		"state" VARCHAR(32) NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "groupchats_telegram_id" UNIQUE ("telegram_id") ON CONFLICT IGNORE
	  );

		CREATE TRIGGER IF NOT EXISTS groupchats_updated_at_Trigger
		AFTER UPDATE On groupchats
		BEGIN
		   UPDATE groupchats SET updated_at = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
		END;`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "plugins" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" text NOT NULL,
		"state" VARCHAR(32) NOT NULL DEFAULT "enabled",
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "plugins_name" UNIQUE ("name") ON CONFLICT IGNORE
	  );

		CREATE TRIGGER IF NOT EXISTS plugins_updated_at_Trigger
		AFTER UPDATE On plugins
		BEGIN
		   UPDATE plugins SET updated_at = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
		END;`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "groups_groupchats" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"group_id" INTEGER NOT NULL,
		"groupchat_id" INTEGER NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "groups_groupchats_group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE,
		CONSTRAINT "groups_groupchats_groupchat_id" FOREIGN KEY ("groupchat_id") REFERENCES "groupchats" ("id") ON DELETE CASCADE,
		CONSTRAINT "groups_groupchats_pair" UNIQUE ("group_id" ASC, "groupchat_id" ASC) ON CONFLICT IGNORE
	  );`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	return db, nil
}

// ExecSQL ...
func ExecSQL(db *sql.DB, sql string) error {
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), sql)
	}

	return nil
}

// QuerySQLObject ...
func QuerySQLObject(db *sql.DB, returnModel interface{}, sql string, args ...interface{}) (reflect.Value, error) {
	t := reflect.TypeOf(returnModel)
	u := reflect.New(t)

	err := db.QueryRow(sql, args...).Scan(u.Interface())
	switch {
	case err == s.ErrNoRows:
		return u, nil
	case err != nil:
		return u, fmt.Errorf("%s: %s", err.Error(), sql)
	}

	return u, nil
}

// QuerySQLList ...
func QuerySQLList(db *sql.DB, returnModel interface{}, sql string, args ...interface{}) ([]reflect.Value, error) {
	var result []reflect.Value

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Error(), sql)
	}

	t := reflect.TypeOf(returnModel)

	for rows.Next() {
		u := reflect.New(t)
		if err = rows.Scan(u.Interface()); err != nil {
			return nil, fmt.Errorf("%s: %s", err.Error(), sql)
		}
		result = append(result, u)
	}

	return result, nil
}

// StoreTelegramMessage ...
func StoreTelegramMessage(db *sql.DB, message *TelegramMessage) error {
	user := &User{
		TelegramID: message.TelegramID,
		FirstName:  message.FirstName,
		LastName:   message.LastName,
		UserName:   message.UserName,
		IsBot:      message.IsBot,
	}

	_, err := AddUserIfNotExist(db, user)
	if err != nil && err.Error() != UserAlreadyExists {
		return err
	}

	_, err2 := db.Exec(
		"INSERT INTO telegram_messages (telegram_id, message, created_at) VALUES (?, ?, ?);",
		message.TelegramID,
		message.Message,
		message.Date)

	if err2 != nil {
		return err2
	}

	return nil
}
