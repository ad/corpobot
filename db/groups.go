package db

import (
	// s "database/sql"
	"fmt"
	// "reflect"
	"strings"
	"time"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"

	_ "github.com/mattn/go-sqlite3" // Register some sql
)

// Group ...
type Group struct {
	ID        int64     `sql:"id"`
	Name      string    `sql:"name"`
	State     string    `sql:"state"`
	CreatedAt time.Time `sql:"created_at"`
}

func (g *Group) String() string {
	return g.Name
}

// AddGroupIfNotExist ...
func AddGroupIfNotExist(db *sql.DB, group *Group) (*Group, error) {
	var returnModel Group

	if group.State == "" {
		group.State = "active"
	}

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM groups WHERE name = ?;`, group.Name)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*Group); ok && returnModel.State == "deleted" {
		return returnModel, fmt.Errorf(GroupDeleted)
	}

	if returnModel, ok := result.Interface().(*Group); ok && returnModel.Name != "" {
		return returnModel, fmt.Errorf(GroupAlreadyExists)
	}

	res, err := db.Exec(
		"INSERT INTO groups (name, state) VALUES (?, ?);",
		group.Name,
		group.State,
	)
	if err != nil {
		return nil, err
	}

	group.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	group.CreatedAt = time.Now()

	dlog.Debugf("%s (%d) added at %s\n", group.Name, group.ID, group.CreatedAt)

	return group, nil
}

// GetGroups ...
func GetGroups(db *sql.DB, states []string) (groups []*Group, err error) {
	if len(states) == 0 {
		states = []string{"active"}
	}

	args := make([]interface{}, len(states))
	for i, state := range states {
		args[i] = state
	}

	var returnModel Group
	sql := `SELECT
	*
FROM
	groups
WHERE
	state IN (?` + strings.Repeat(",?", len(args)-1) + `)
ORDER BY
	state, name;`

	result, err := QuerySQLList(db, returnModel, sql, args...)
	if err != nil {
		return groups, err
	}

	for _, item := range result {
		if returnModel, ok := item.Interface().(*Group); ok {
			groups = append(groups, returnModel)
		}
	}

	return groups, err
}

// UpdateGroupState ...
func UpdateGroupState(db *sql.DB, group *Group) (int64, error) {
	result, err := db.Exec(
		"UPDATE groups SET state = ? WHERE name = ? AND state != ?;",
		group.State,
		group.Name,
		group.State)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

// UpdateGroupName ...
func UpdateGroupName(db *sql.DB, oldName, newName string) (int64, error) {
	result, err := db.Exec(
		"UPDATE groups SET name = ? WHERE name = ? AND name != ?;",
		newName,
		oldName,
		newName)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

// GetGroupByName ...
func GetGroupByName(db *sql.DB, group *Group) (*Group, error) {
	var returnModel Group

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM groups WHERE name = ?;`, group.Name)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*Group); ok && returnModel.Name != "" {
		return returnModel, nil
	}

	return nil, fmt.Errorf(GroupNotFound)
}

// AddGroupGroupChatIfNotExist ...
func AddGroupGroupChatIfNotExist(db *sql.DB, group *Group, groupchat *Groupchat) (bool, error) {
	res, err := db.Exec(
		"INSERT INTO groups_groupchats (group_id, groupchat_id) VALUES (?, ?);",
		group.ID,
		groupchat.ID,
	)
	if err != nil {
		return false, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteGroupGroupChat ...
func DeleteGroupGroupChat(db *sql.DB, group *Group, groupchat *Groupchat) (bool, error) {
	res, err := db.Exec(
		"DELETE FROM groups_groupchats WHERE group_id = ? AND groupchat_id = ?;",
		group.ID,
		groupchat.ID,
	)
	if err != nil {
		return false, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return false, err
	}

	return true, nil
}

// AddGroupUserIfNotExist ...
func AddGroupUserIfNotExist(db *sql.DB, group *Group, user *User) (bool, error) {
	res, err := db.Exec(
		"INSERT INTO groups_Users (group_id, user_id) VALUES (?, ?);",
		group.ID,
		user.ID,
	)
	if err != nil {
		return false, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteGroupUser ...
func DeleteGroupUser(db *sql.DB, group *Group, user *User) (bool, error) {
	res, err := db.Exec(
		"DELETE FROM groups_Users WHERE group_id = ? AND user_id = ?;",
		group.ID,
		user.ID,
	)
	if err != nil {
		return false, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return false, err
	}

	return true, nil
}
