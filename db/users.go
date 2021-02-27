package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	_ "github.com/mattn/go-sqlite3" // Register some sql
)

// User ...
type User struct {
	ID         int64     `sql:"id"`
	FirstName  string    `sql:"first_name"`
	LastName   string    `sql:"last_name"`
	UserName   string    `sql:"user_name"`
	TelegramID int64     `sql:"telegram_id"`
	IsBot      bool      `sql:"is_bot"`
	Role       string    `sql:"role"`
	Birthday   time.Time `sql:"birthday"`
	CreatedAt  time.Time `sql:"created_at"`
}

func (u *User) String() string {
	var b strings.Builder
	if u.UserName != "" {
		b.WriteRune('@')
		b.WriteString(u.UserName)
		b.WriteRune(' ')
	}
	if u.FirstName != "" {
		b.WriteString(u.FirstName)
		b.WriteRune(' ')
	}
	if u.LastName != "" {
		b.WriteString(u.LastName)
		b.WriteRune(' ')
	}
	if u.TelegramID != 0 {
		b.WriteRune('[')
		b.WriteString(strconv.FormatInt(u.TelegramID, 10))
		b.WriteRune(']')
		b.WriteRune(' ')
	}
	if u.Role != "" {
		b.WriteRune('(')
		b.WriteString(u.Role)
		b.WriteRune(')')
	}

	if b.Len() == 0 {
		b.WriteString("ID")
		b.WriteString(strconv.FormatInt(u.ID, 10))
	}

	return b.String()
}

func (u *User) Short() string {
	var b strings.Builder
	if u.UserName != "" {
		b.WriteRune('@')
		b.WriteString(u.UserName)
	}
	if u.FirstName != "" {
		if b.Len() != 0 {
			b.WriteRune(' ')
		}
		b.WriteString(u.FirstName)
	}
	if u.LastName != "" {
		if b.Len() != 0 {
			b.WriteRune(' ')
		}
		b.WriteString(u.LastName)
	}
	if b.Len() == 0 {
		if u.TelegramID != 0 {
			b.WriteRune('[')
			b.WriteString(strconv.FormatInt(u.TelegramID, 10))
			b.WriteRune(']')
		}
	}

	if b.Len() == 0 {
		b.WriteString("ID")
		b.WriteString(strconv.FormatInt(u.ID, 10))
	}

	return b.String()
}

func (u *User) Paragraph() string {
	var b strings.Builder
	if u.FirstName != "" {
		b.WriteString(u.FirstName)
		b.WriteRune(' ')
	}
	if u.LastName != "" {
		b.WriteString(u.LastName)
		b.WriteRune(' ')
	}
	if u.FirstName != "" || u.LastName != "" {
		b.WriteString("\n")
	}
	if u.UserName != "" {
		b.WriteRune('@')
		b.WriteString(u.UserName)
		b.WriteString("\n")
	}
	if u.TelegramID != 0 {
		b.WriteString("id: ")
		b.WriteString(strconv.FormatInt(u.TelegramID, 10))
		b.WriteString("\n")
	}
	b.WriteString("Role: ")
	b.WriteString(u.Role)
	b.WriteString("\n")

	if !u.Birthday.IsZero() {
		b.WriteString("Birthday: ")
		b.WriteString(u.Birthday.Format("2006.01.02"))
	}
	return b.String()
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

	if returnModel, ok := result.Interface().(*User); ok && returnModel.TelegramID > 0 {
		return returnModel, fmt.Errorf(UserAlreadyExists)
	}

	res, err := db.Exec(
		"INSERT INTO users (first_name, last_name, user_name, telegram_id, is_bot, role) VALUES (?, ?, ?, ?, ?, ?);",
		user.FirstName,
		user.LastName,
		user.UserName,
		user.TelegramID,
		user.IsBot,
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	user.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

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
		AND
	is_bot = False
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
		"UPDATE users SET role = ? WHERE telegram_id = ? AND role != ? AND role != 'owner';",
		user.Role,
		user.TelegramID,
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

// UpdateUserBirthday ...
func UpdateUserBirthday(db *sql.DB, user *User) (int64, error) {
	result, err := db.Exec("UPDATE users SET birthday = ? WHERE telegram_id = ?;", user.Birthday, user.TelegramID)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

// GetUserByID ...
func GetUserByID(db *sql.DB, user *User) (*User, error) {
	var returnModel User

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM users WHERE id = ?;`, user.ID)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*User); ok && returnModel.Role != "" {
		return returnModel, nil
	}

	return nil, fmt.Errorf(UserNotFound)
}

// GetUserByTelegramID ...
func GetUserByTelegramID(db *sql.DB, user *User) (*User, error) {
	var returnModel User

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM users WHERE telegram_id = ?;`, user.TelegramID)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*User); ok && returnModel.Role != "" {
		return returnModel, nil
	}

	return nil, fmt.Errorf(UserNotFound)
}
