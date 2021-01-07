package db

import (
	"strings"
	"time"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"

	_ "github.com/mattn/go-sqlite3" // Register some sql
)

// User ...
type Plugin struct {
	ID        int64     `sql:"id"`
	Name      string    `sql:"name"`
	State     string    `sql:"state"`
	CreatedAt time.Time `sql:"created_at"`
}

func (p *Plugin) String() string {
	var b strings.Builder
	b.WriteString(p.Name)
	b.WriteRune(' ')
	b.WriteRune('—')
	b.WriteRune(' ')
	b.WriteString(p.State)
	return b.String()
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

	plugin.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	plugin.CreatedAt = time.Now()

	dlog.Debugf("%s (%s) added at %s\n", plugin.Name, plugin.State, plugin.CreatedAt)

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

// GetPlugins ...
func GetPlugins(db *sql.DB) (plugins []*Plugin, err error) {
	var returnModel Plugin
	sql := `SELECT
	*
FROM
	plugins
ORDER BY
	name;`

	result, err := QuerySQLList(db, returnModel, sql, nil)
	if err != nil {
		return plugins, err
	}

	for _, item := range result {
		if returnModel, ok := item.Interface().(*Plugin); ok {
			plugins = append(plugins, returnModel)
		}
	}

	return plugins, err
}

func (p *Plugin) IsEnabled() bool {
	return p.State == "enabled"
}
