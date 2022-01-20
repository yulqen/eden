// check https://gosamples.dev/sqlite-intro/

package repository

type Entry struct {
	ID      int64
	Time    string
	Content string
}

// sqlite date type equates to fmt.Println(time.Now().Format(time.RFC3339)) in Go (according to :https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang)

// 2.2. Date and Time Datatype

//    SQLite does not have a storage class set aside for storing dates and/or times. Instead, the built-in Date And Time Functions of SQLite are
//    capable of storing dates and times as TEXT, REAL, or INTEGER values:
//      * TEXT as ISO8601 strings ("YYYY-MM-DD HH:MM:SS.SSS").
//      * REAL as Julian day numbers, the number of days since noon in Greenwich on November 24, 4714 B.C. according to the proleptic Gregorian
//        calendar.
//      * INTEGER as Unix Time, the number of seconds since 1970-01-01 00:00:00 UTC.

// This is a good guide to dates and times in go: https://qvault.io/golang/golang-date-time/
