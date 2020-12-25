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

// Group ...
type Group struct {
	ID             int64     `sql:"id"`
	Name      	   string    `sql:"name"`
	State          string    `sql:"state"`
	CreatedAt      time.Time `sql:"created_at"`
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

	group.ID, _ = res.LastInsertId()
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
	sql := `select
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
