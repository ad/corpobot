package db

import (
	// s "database/sql"
	"fmt"
	// "reflect"
	"time"
	"strings"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	_ "github.com/mattn/go-sqlite3"
)

// User ...
type User struct {
	ID             int64     `sql:"id"`
	FirstName      string    `sql:"first_name"`
	LastName       string    `sql:"last_name"`
	UserName       string    `sql:"user_name"`
	TelegramID 	   int       `sql:"telegram_id"`
	IsBot		   bool		 `sql:"is_bot"`
	Role 	   	   string    `sql:"role"`
	CreatedAt      time.Time `sql:"created_at"`
}

func (u *User) String() string {
    return fmt.Sprintf("@%s %s %s [id %d] (%s)", u.UserName, u.FirstName, u.LastName, u.TelegramID, u.Role)
}

// AddUserIfNotExist ...
func AddUserIfNotExist(db *sql.DB, user *User) (*User, error) {
	var returnModel User

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM users WHERE telegram_id = ?;`, user.TelegramID)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*User); ok && returnModel.Role == "deleted" {
		return returnModel, fmt.Errorf(UserDeleted)
	}

	if returnModel, ok := result.Interface().(*User); ok && returnModel.Role == "blocked" {
		return returnModel, fmt.Errorf(UserBlocked)
	}

	if returnModel, ok := result.Interface().(*User); ok && returnModel.UserName != "" {
		return returnModel, fmt.Errorf(UserAlreadyExists)
	}

	res, err := db.Exec(
		"INSERT INTO users (first_name, last_name, user_name, telegram_id, is_bot) VALUES (?, ?, ?, ?, ?);",
		user.FirstName,
		user.LastName,
		user.UserName,
		user.TelegramID,
		user.IsBot,
	)

	if err != nil {
		return nil, err
	}

	user.ID, _ = res.LastInsertId()
	user.CreatedAt = time.Now()

	dlog.Debugf("%s (%d) added at %s\n", user.UserName, user.ID, user.CreatedAt)

	return user, nil
}


// GetUsers ...
func GetUsers(db *sql.DB, roles []string) (users []*User, err error) {
	if len(roles) == 0 {
		roles = []string{"owner", "admin", "member", "new"}
	}

	args := make([]interface{}, len(roles))
	for i, role := range roles {
	    args[i] = role
	}

	var returnModel User
	sql := `SELECT
	*
FROM
	users
WHERE
	role IN (?` + strings.Repeat(",?", len(args)-1) + `)
ORDER BY
	role, id;`

	result, err := QuerySQLList(db, returnModel, sql, args...)
	if err != nil {
		return users, err
	}

	for _, item := range result {
		if returnModel, ok := item.Interface().(*User); ok {
			users = append(users, returnModel)
		}
	}

	return users, err
}

// UpdateUserRole ...
func UpdateUserRole(db *sql.DB, user *User) (int64, error) {
	result, err := db.Exec(
		"UPDATE users SET role = ? WHERE user_name = ? AND role != ?;",
		user.Role,
		user.UserName,
		user.Role)

	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return -1, err
	}

	return rows, nil
}