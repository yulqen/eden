package repository

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

const (
	edenConfigDirName = "eden"
	dbName            = "eden.db"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

// getUserConfigDir allows replaces os.UserConfigDir
// for testing purposes.
type getUserConfigDir func() (string, error)

// DBPathChecker contains the func used to create the user config dir.
type DBPathChecker struct {
	userConfig getUserConfigDir
}

type Options struct {
	// path to the database file
	DBPath string
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS entries(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time TEXT NOT NULL,
		content TEXT NOT NULL
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(entry Entry) (*Entry, error) {
	res, err := r.db.Exec("INSERT INTO entries(content, time) values(?,?)", entry.Content, entry.Time)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	entry.ID = id

	return &entry, nil
}

func (r *SQLiteRepository) All() ([]Entry, error) {
	rows, err := r.db.Query("SELECT * FROM entries")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.Content); err != nil {
			return nil, err
		}
		all = append(all, entry)
	}
	return all, nil
}

func (r *SQLiteRepository) GetByID(id int64) (*Entry, error) {
	row := r.db.QueryRow("SELECT * FROM entries WHERE id = ?", id)

	var entry Entry
	if err := row.Scan(&entry.ID, &entry.Time, &entry.Content); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &entry, nil
}

func (r *SQLiteRepository) Update(id int64, updated Entry) (*Entry, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}

	res, err := r.db.Exec("UPDATE entries SET content = ? WHERE id = ?", updated.Content, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	res, err := r.db.Exec("DELETE FROM entries WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}
	return err
}

// NewDBPathChecker creates a DBPathChecker using whatever
// func you want as the argument, as long as it matches the
// type os.UserConfigDir. This makes it convenient for testing
// and was done as an experiment here to practice mocking in Go.
func NewDBPathChecker(f getUserConfigDir) *DBPathChecker {
	return &DBPathChecker{userConfig: f}
}

// Check returns true if the necessary config files (including
// the database) are in place - false if not
func (db *DBPathChecker) Check() bool {
	userConfig, err := db.userConfig()
	if err != nil {
		log.Fatal(err)
	}
	dbPath := filepath.Join(userConfig, "eden", "eden.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		println("Seems that " + dbPath + " not there")
		return false
	}
	return true
}

type Repository interface {
	Migrate() error
	Create(entry Entry) (*Entry, error)
	All() ([]Entry, error)
	GetByID(id int64) (*Entry, error)
	Update(id int64, updated Entry) (*Entry, error)
	Delete(id int64) error
}

// SetUp creates the config directory and requisite files
func SetUp(d getUserConfigDir) (string, error) {
	sysConfigDir, err := d()
	if err != nil {
		return "", err
	}
	// check if config folder exists
	edenConfigPath := filepath.Join(sysConfigDir, edenConfigDirName)
	dbPath := filepath.Join(edenConfigPath, dbName)
	if _, err := os.Stat(edenConfigPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(edenConfigPath, edenConfigDirName), 0700); err != nil {
			return "", err
		}
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		_, err = os.Create(dbPath)
		if err != nil {
			return "", err
		}

		_, err := os.Create(dbPath)
		if err != nil {
			return "", err
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return "", err
		}

		r := NewSQLiteRepository(db)

		err = r.Migrate()
		if err != nil {
			return "", err
		}
		if err != nil {
			return "", err
		}
	}
	return edenConfigPath, nil
}
