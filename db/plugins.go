package db

import (
	// s "database/sql"
	"fmt"
	// "reflect"
	"time"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	_ "github.com/mattn/go-sqlite3"
)

// User ...
type Plugin struct {
	ID        int64     `sql:"id"`
	Name      string    `sql:"name"`
	State     string    `sql:"state"`
	CreatedAt time.Time `sql:"created_at"`
}

func (p *Plugin) String() string {
	return fmt.Sprintf("%s — %s", p.Name, p.State)
}

// AddPluginIfNotExist ...
func AddPluginIfNotExist(db *sql.DB, plugin *Plugin) (*Plugin, error) {
	var returnModel Plugin

	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM plugins WHERE name = ?;`, plugin.Name)
	if err != nil {
		return nil, err
	}

	if returnModel, ok := result.Interface().(*Plugin); ok && returnModel.Name != "" {
		return returnModel, nil
	}

	res, err := db.Exec(
		"INSERT INTO plugins (name, state) VALUES (?, ?);",
		plugin.Name,
		plugin.State,
	)
	if err != nil {
		return nil, err
	}

	plugin.ID, _ = res.LastInsertId()
	plugin.CreatedAt = time.Now()

	dlog.Debugf("%s (%d) added at %s\n", plugin.Name, plugin.State, plugin.CreatedAt)

	return plugin, nil
}

// UpdatePluginState ...
func UpdatePluginState(db *sql.DB, plugin *Plugin) (int64, error) {
	result, err := db.Exec(
		"UPDATE plugins SET state = ? WHERE name = ? AND state != ?;",
		plugin.State,
		plugin.Name,
		plugin.State)
	if err != nil {
		return -1, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rows, nil
}
