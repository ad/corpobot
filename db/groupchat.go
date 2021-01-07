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

// Groupchat ...
type Groupchat struct {
	ID         int64     `sql:"id"`
	Title      string    `sql:"title"`
	TelegramID int64     `sql:"telegram_id"`
	State      string    `sql:"state"`
	InviteLink string    `sql:"invite_link"`
	CreatedAt  time.Time `sql:"created_at"`
}

func (gc *Groupchat) String() string {
	var b strings.Builder
	b.WriteString(gc.Title)
	b.WriteRune(' ')
	b.WriteRune('[')
	b.WriteString(strconv.FormatInt(gc.TelegramID, 10))
	b.WriteRune(']')
	b.WriteRune(' ')
	b.WriteString(gc.InviteLink)
	return b.String()
}

// GetGroupchats ...
func GetGroupchats(db *sql.DB, states []string) (groupchats []*Groupchat, err error) {
	if len(states) == 0 {
		states = []string{Active}
	}

	args := make([]interface{}, len(states))
	for i, state := range states {
		args[i] = state
	}

	var returnModel Groupchat
	sql := `SELECT
	*
FROM
	groupchats
WHERE
	state IN (?` + strings.Repeat(",?", len(args)-1) + `)
ORDER BY
	state, title;`

	result, err := QuerySQLList(db, returnModel, sql, args...)
	if err != nil {
		return groupchats, err
	}

	for _, item := range result {
		if returnModel, ok := item.Interface().(*Groupchat); ok {
			groupchats = append(groupchats, returnModel)
		}
	}

	return groupchats, err
}

// GetGroupchatsByGroupID ...
func GetGroupchatsByGroupID(db *sql.DB, groupID int64) (groupchats []*Groupchat, err error) {
	var returnModel Groupchat
	sql := `SELECT
	*
FROM
	groupchats
WHERE
	id IN (SELECT groupchat_id FROM groups_groupchats WHERE group_id = ?)
ORDER BY
	state, title;`

	result, err := QuerySQLList(db, returnModel, sql, groupID)
	if err != nil {
		return groupchats, err
	}

	for _, item := range result {
		if returnModel, ok := item.Interface().(*Groupchat); ok {
			groupchats = append(groupchats, returnModel)
		}
	}

	return groupchats, err
}

// AddGroupChatIfNotExist ...
func AddGroupChatIfNotExist(db *sql.DB, groupchat *Groupchat) (*Groupchat, error) {
	var returnModel Groupchat

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM groupchats WHERE telegram_id = ?;`, groupchat.TelegramID)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*Groupchat); ok && returnModel.State != "" {
		return returnModel, fmt.Errorf(GroupChatAlreadyExists)
	}

	res, err := db.Exec(
		"INSERT INTO groupchats (title, telegram_id, invite_link, state) VALUES (?, ?, ?, ?);",
		groupchat.Title,
		groupchat.TelegramID,
		groupchat.InviteLink,
		groupchat.State,
	)
	if err != nil {
		return nil, err
	}

	groupchat.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	groupchat.CreatedAt = time.Now()

	dlog.Debugf("%s (%d) added at %s\n", groupchat.Title, groupchat.ID, groupchat.CreatedAt)

	return groupchat, nil
}

// UpdateGroupChatInviteLink ...
func UpdateGroupChatInviteLink(db *sql.DB, groupchat *Groupchat) (int64, error) {
	result, err := db.Exec(
		"UPDATE groupchats SET invite_link = ? WHERE telegram_id = ?;",
		groupchat.InviteLink,
		groupchat.TelegramID)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

// UpdateGroupChatTitle ...
func UpdateGroupChatTitle(db *sql.DB, groupchat *Groupchat) (int64, error) {
	result, err := db.Exec(
		"UPDATE groupchats SET title = ? WHERE telegram_id = ?;",
		groupchat.Title,
		groupchat.TelegramID)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

// GetGroupChatByTelegramID ...
func GetGroupChatByTelegramID(db *sql.DB, groupchat *Groupchat) (*Groupchat, error) {
	var returnModel Groupchat

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM groupchats WHERE telegram_id = ?;`, groupchat.TelegramID)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*Groupchat); ok && returnModel.State != "" {
		return returnModel, nil
	}

	return nil, fmt.Errorf(GroupChatNotFound)
}

// GroupChatDelete ...
func GroupChatDelete(db *sql.DB, groupchat *Groupchat) (bool, error) {
	_, err := db.Exec(
		"DELETE FROM groupchats WHERE telegram_id = ?;",
		groupchat.TelegramID,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
