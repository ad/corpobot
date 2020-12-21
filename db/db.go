package db

import (
	s "database/sql"
	"fmt"
	"reflect"
	"time"

	dlog "github.com/amoghe/distillog"
	sql "github.com/lazada/sqle"
	_ "github.com/mattn/go-sqlite3" // ...
)

// AlreadyExists ...
const (
	AlreadyExists = "already exists"
	UserNotFound  = "user not found"
	RepoNotFound  = "repo not found"
)

// TelegramMessage ...
type TelegramMessage struct {
	ID       int
	UserID   int
	UserName string
	Message  string
	Date     time.Time
}

// // GithubUser ...
// type GithubUser struct {
// 	ID             int64     `sql:"id"`
// 	Name           string    `sql:"name"`
// 	UserName       string    `sql:"user_name"`
// 	TelegramUserID string    `sql:"telegram_user_id"`
// 	Token          string    `sql:"token"`
// 	CreatedAt      time.Time `sql:"created_at"`
// }

// // GithubRepo ...
// type GithubRepo struct {
// 	ID        int64     `sql:"id"`
// 	Name      string    `sql:"name"`
// 	RepoName  string    `sql:"repo_name"`
// 	CreatedAt time.Time `sql:"created_at"`
// }

// // UserRepo ...
// type UserRepo struct {
// 	ID        int64     `sql:"id"`
// 	UserID    int64     `sql:"user_id"`
// 	RepoID    int64     `sql:"repo_id"`
// 	CreatedAt time.Time `sql:"created_at"`
// 	UpdatedAt time.Time `sql:"updated_at"`
// }

// // UsersReposResult ...
// type UsersReposResult struct {
// 	UserID         int64
// 	RepoID         int64
// 	TelegramUserID string
// 	Token          string
// 	RepoName       string
// 	UpdatedAt      time.Time
// }

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
		"user_id" INTEGER NOT NULL,
		"user_name" VARCHAR(32) DEFAULT "",
		"message" TEXT DEFAULT "",
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "github_users" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" text NOT NULL,
		"user_name" text NOT NULL,
		"token" text NOT NULL,
		"telegram_user_id" INTEGER NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "github_users_user_name" UNIQUE ("user_name") ON CONFLICT IGNORE
	  );`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "github_repos" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" text NOT NULL,
		"repo_name" text NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT "github_repos_repo_name" UNIQUE ("repo_name") ON CONFLICT IGNORE
	  );`)
	if err != nil {
		dlog.Errorf("%s", err)
	}

	err = ExecSQL(db, `CREATE TABLE IF NOT EXISTS "users_repos" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"user_id" INTEGER NOT NULL,
		"repo_id" INTEGER NOT NULL,
		"created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
		"updated_at" timestamp,
		CONSTRAINT "repos_user_id" FOREIGN KEY ("user_id") REFERENCES "github_users" ("id"),
		CONSTRAINT "repos_repo_id" FOREIGN KEY ("repo_id") REFERENCES "github_repos" ("id"),
		CONSTRAINT "repos_repo_id_user_id" UNIQUE ("user_id", "repo_id") ON CONFLICT IGNORE
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

// // AddUserIfNotExist ...
// func AddUserIfNotExist(db *sql.DB, user *GithubUser) (*GithubUser, error) {
// 	var returnModel GithubUser

// 	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM github_users WHERE user_name = ?;`, user.UserName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if returnModel, ok := result.Interface().(*GithubUser); ok && returnModel.UserName != "" {
// 		// dlog.Debugf("already exists: %#v", returnModel)
// 		return returnModel, fmt.Errorf(AlreadyExists)
// 	}

// 	res, err := db.Exec(
// 		"INSERT INTO github_users (name, user_name, token, telegram_user_id) VALUES (?, ?, ?, ?);",
// 		user.Name,
// 		user.UserName,
// 		user.Token,
// 		user.TelegramUserID,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	user.ID, _ = res.LastInsertId()
// 	user.CreatedAt = time.Now()

// 	dlog.Debugf("%s (%d) added at %s\n", user.UserName, user.ID, user.CreatedAt)

// 	return user, nil
// }

// // AddRepoIfNotExist ...
// func AddRepoIfNotExist(db *sql.DB, repo *GithubRepo) (*GithubRepo, error) {
// 	var returnModel GithubRepo

// 	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM github_repos WHERE repo_name = ?;`, repo.RepoName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if returnModel, ok := result.Interface().(*GithubRepo); ok && returnModel.RepoName != "" {
// 		// dlog.Debugf("already exists: %#v", returnModel)
// 		return returnModel, fmt.Errorf(AlreadyExists)
// 	}
// 	res, err := db.Exec(
// 		"INSERT INTO github_repos (name, repo_name) VALUES (?, ?);",
// 		repo.Name,
// 		repo.RepoName,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	repo.ID, _ = res.LastInsertId()
// 	repo.CreatedAt = time.Now()

// 	dlog.Debugf("%s (%d) added at %s\n", repo.RepoName, repo.ID, repo.CreatedAt)

// 	return repo, nil
// }

// // AddRepoLinkIfNotExist ...
// func AddRepoLinkIfNotExist(db *sql.DB, user *GithubUser, repo *GithubRepo, updatedAt time.Time) error {
// 	var returnModel UserRepo

// 	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM users_repos WHERE user_id = ? AND repo_id = ?;`, user.ID, repo.ID)
// 	if err != nil {
// 		return err
// 	}
// 	if returnModel, ok := result.Interface().(*UserRepo); ok && returnModel.UserID > 0 && returnModel.RepoID > 0 {
// 		// dlog.Debugf("already exists: %#v", returnModel)
// 		return fmt.Errorf(AlreadyExists)
// 	}

// 	res, err := db.Exec(
// 		"INSERT INTO users_repos (user_id, repo_id, updated_at) VALUES (?, ?, ?);",
// 		user.ID,
// 		repo.ID,
// 		updatedAt,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	id, _ := res.LastInsertId()

// 	dlog.Debugf("link %s <-> %s (%d) added at %s\n", user.Name, repo.RepoName, id, time.Now())

// 	return nil
// }

// StoreTelegramMessage ...
func StoreTelegramMessage(db *sql.DB, message TelegramMessage) error {
	_, err := db.Exec(
		"INSERT INTO telegram_messages (user_id, user_name, message, created_at) VALUES (?, ?, ?, ?);",
		message.UserID,
		message.UserName,
		message.Message,
		message.Date)

	if err != nil {
		return err
	}

	return nil
}

// // UpdateUserRepoLink ...
// func UpdateUserRepoLink(db *sql.DB, userRepoResult *UsersReposResult) error {
// 	_, err := db.Exec(
// 		"UPDATE users_repos SET updated_at = ? WHERE user_id = ? AND repo_id = ?;",
// 		userRepoResult.UpdatedAt,
// 		userRepoResult.UserID,
// 		userRepoResult.RepoID)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // GetUserRepos ...
// func GetUserRepos(db *sql.DB) (usersRepos []*UsersReposResult, err error) {
// 	var returnModel UsersReposResult
// 	sql := `select
// 	users_repos.user_id as user_id,
// 	users_repos.repo_id as repo_id,
// 	github_users.telegram_user_id as telegram_user_id,
// 	github_users.token as token,
// 	github_repos.repo_name as repo_name,
// 	users_repos.updated_at as updated_at
// FROM
// 	github_repos
// 	INNER JOIN users_repos ON github_repos.id = users_repos.repo_id
// 	INNER JOIN github_users ON github_users.id = users_repos.user_id
// WHERE
// 	users_repos.updated_at < DATETIME('now',  '-15 minutes');`

// 	result, err := QuerySQLList(db, returnModel, sql)
// 	if err != nil {
// 		return usersRepos, err
// 	}
// 	for _, item := range result {
// 		if returnModel, ok := item.Interface().(*UsersReposResult); ok {
// 			usersRepos = append(usersRepos, returnModel)
// 		}
// 	}

// 	return usersRepos, err
// }

// // GetUsers ...
// func GetUsers(db *sql.DB) (users []*GithubUser, err error) {
// 	var returnModel GithubUser
// 	sql := `select
// 	*
// FROM
// 	github_users;`

// 	result, err := QuerySQLList(db, returnModel, sql)
// 	if err != nil {
// 		return users, err
// 	}
// 	for _, item := range result {
// 		if returnModel, ok := item.Interface().(*GithubUser); ok {
// 			users = append(users, returnModel)
// 		}
// 	}

// 	return users, err
// }

// // GetGithubUserFromDB ...
// func GetGithubUserFromDB(db *sql.DB, id string) (*GithubUser, error) {
// 	var returnModel GithubUser

// 	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM github_users WHERE telegram_user_id = ?;`, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if returnModel, ok := result.Interface().(*GithubUser); ok && returnModel.UserName != "" {
// 		return returnModel, nil
// 	}

// 	return nil, fmt.Errorf(UserNotFound)
// }

// // GetGithubRepoByNameFromDB ...
// func GetGithubRepoByNameFromDB(db *sql.DB, repo_name string) (*GithubRepo, error) {
// 	var returnModel GithubRepo

// 	result, err := QuerySQLObject(db, returnModel, `SELECT * FROM github_repos WHERE repo_name = ?;`, repo_name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if returnModel, ok := result.Interface().(*GithubRepo); ok && returnModel.RepoName != "" {
// 		return returnModel, nil
// 	}

// 	return nil, fmt.Errorf(RepoNotFound)
// }

// // DeleteRepoUserLinkDB ...
// func DeleteRepoUserLinkDB(db *sql.DB, user *GithubUser, repo *GithubRepo) error {
// 	_, err := db.Exec(
// 		"DELETE FROM users_repos WHERE user_id = ? AND repo_id = ?;",
// 		user.ID,
// 		repo.ID)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
