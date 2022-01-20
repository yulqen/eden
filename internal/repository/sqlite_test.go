package repository

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testConfigPath = "/tmp/CONFIG/"

func testConfigDir() string {
	d, err := SetUp(mockConfigDir)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func mockConfigDir() (string, error) {
	return testConfigPath, nil
}

func TestGetTodayTime(t *testing.T) {
	// This is a good guide to dates and times in go: https://qvault.io/golang/golang-date-time/
	// This test mainly tests the Go standard library, which is not what I want, but this
	// is useful for learning purposes.
	refTime := time.Date(1990, time.April, 10, 15, 0, 0, 0, time.UTC)
	isoDate := refTime.Format(time.RFC3339)
	// the reference format for all go dates is  Mon Jan 2 15:04:05 -0700 MST 2006 - you just pull the elements from this into time.Format()
	fixedTime := time.Date(1990, time.April, 10, 15, 0, 0, 0, time.UTC).Format("15:04")
	if refTime.Year() != 1990 {
		t.Error("Could not get the correct year")
	}
	if isoDate != "1990-04-10T15:00:00Z" {
		t.Errorf("Expected 1990-04-10T15:00:00Z got %s", isoDate)
	}
	if fixedTime != "15:00" {
		t.Errorf("Expected 15:00 but got %s", fixedTime)
	}
}

func TestMockedConfig(t *testing.T) {
	d := testConfigDir()
	if d != "/tmp/CONFIG/eden" {
		t.Errorf("Expected /tmp/CONFIG/eden but got %s", d)
	}

	defer func() {
		os.RemoveAll(filepath.Join("/tmp", "CONFIG"))
	}()

	dbpc := NewDBPathChecker(mockConfigDir)
	h := dbpc.Check()
	if !h {
		t.Error("the db config directory should be found but isn't")
	}
}

func TestCanAddEntry(t *testing.T) {
	d := testConfigDir()
	defer func() {
		os.RemoveAll(filepath.Join("/tmp", "CONFIG"))
	}()

	db, err := sql.Open("sqlite3", filepath.Join(d, dbName))
	if err != nil {
		t.Fatal(err)
	}

	r := NewSQLiteRepository(db)

	err = r.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	refTime := time.Date(2021, time.April, 10, 15, 0, 0, 0, time.UTC).Format(time.RFC3339)
	ent := Entry{Content: "Smash!", Time: refTime}

	retE, err := r.Create(ent)
	if err != nil {
		t.Error(err)
	}
	if retE.Content != "Smash!" {
		t.Errorf("Expected to get Smash! but got %s", retE.Content)
	}
	if retE.Time != "2021-04-10T15:00:00Z" {
		t.Errorf("Expected time to be \"2021-04-10T15:00:00Z\" but got %s", retE.Time)
	}
}
