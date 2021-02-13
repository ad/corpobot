package db

import (
	"fmt"
	"strings"
	"time"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	_ "github.com/mattn/go-sqlite3" // Register some sql
)

// Meetingroom ...
type Meetingroom struct {
	ID         int64     `sql:"id"`
	Name       string    `sql:"name"`
	State      string    `sql:"state"`
	CreatedAt  time.Time `sql:"created_at"`
	UpdateddAt time.Time `sql:"updated_at"`
}

func (m *Meetingroom) String() string {
	return m.Name
}

// AddMeetingroomIfNotExist ...
func AddMeetingroomIfNotExist(db *sql.DB, m *Meetingroom) (*Meetingroom, error) {
	var returnModel Meetingroom

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM meetingrooms WHERE name = ?;`, m.Name)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*Meetingroom); ok && returnModel.State == "deleted" {
		return returnModel, fmt.Errorf(MeetingroomDeleted)
	}

	if returnModel, ok := result.Interface().(*Meetingroom); ok && returnModel.State == "blocked" {
		return returnModel, fmt.Errorf(MeetingroomBlocked)
	}

	if returnModel, ok := result.Interface().(*Meetingroom); ok && returnModel.Name != "" {
		return returnModel, fmt.Errorf(MeetingroomAlreadyExists)
	}

	res, err := db.Exec(
		"INSERT INTO meetingrooms (name, state) VALUES (?, ?);",
		m.Name,
		m.State,
	)
	if err != nil {
		return nil, err
	}

	m.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	m.CreatedAt = time.Now()

	dlog.Debugf("%s (%d) added at %s\n", m.Name, m.ID, m.CreatedAt)

	return m, nil
}

// GetMeetingrooms ...
func GetMeetingrooms(db *sql.DB, states []string) (ms []*Meetingroom, err error) {
	if len(states) == 0 {
		states = []string{"active"}
	}

	args := make([]interface{}, len(states))
	for i, state := range states {
		args[i] = state
	}

	var returnModel Meetingroom
	sql := `SELECT
	*
FROM
	meetingrooms
WHERE
	state IN (?` + strings.Repeat(",?", len(args)-1) + `)
ORDER BY
	state, id;`

	result, err := QuerySQLList(db, returnModel, sql, args...)
	if err != nil {
		return ms, err
	}

	for _, item := range result {
		if returnModel, ok := item.Interface().(*Meetingroom); ok {
			ms = append(ms, returnModel)
		}
	}

	return ms, err
}

// UpdateMeetingroomState ...
func UpdateMeetingroomState(db *sql.DB, m *Meetingroom) (int64, error) {
	result, err := db.Exec(
		"UPDATE meetingrooms SET state = ? WHERE name = ? AND state != ?;",
		m.State,
		m.Name,
		m.State)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}

// UpdateMeetingroomName ...
func UpdateMeetingroomName(db *sql.DB, oldName, newName string) (int64, error) {
	result, err := db.Exec("UPDATE meetingrooms SET name = ? WHERE name = ? AND name != ?;", newName, oldName, newName)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}
