package post05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// Connection details
var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

// MSDSCourse is for holding full course data
// MSDSCourse table + cname
type MSDSCourse struct {
	ID      int
	CID     string
	CNAME   string
	CPREREQ string
}

func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// The function returns the ID of the cname
// -1 if the course does not exist
func exists(cid string) int {
	cid = strings.ToLower(cid)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	courseid := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "courses" where cid = '%s'`, cid)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		courseid = id
	}
	defer rows.Close()
	return courseid
}

// AddCourse adds a new course to the database
// Returns new Course ID
// -1 if there was an error
func AddCourse(d MSDSCourse) int {
	d.CID = strings.ToLower(d.CID)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	CourseID := exists(d.CID)
	if CourseID != -1 {
		fmt.Println("Course already exists:", Username)
		return -1
	}

	insertStatement := `insert into "courses" ("cid") values ($1)`
	_, err = db.Exec(insertStatement, d.CID)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	CourseID = exists(d.CID)
	if CourseID == -1 {
		return CourseID
	}

	insertStatement = `insert into "coursedata" ("courseid", "cname", "cprereq")
	values ($1, $2, $3)`
	_, err = db.Exec(insertStatement, CourseID, d.CNAME, d.CPREREQ)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return CourseID
}

// DeleteCourse deletes an existing course
func DeleteCourse(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Does the ID exist?
	statement := fmt.Sprintf(`SELECT "cid" FROM "courses" where id = %d`, id)
	rows, err := db.Query(statement)

	var cid string
	for rows.Next() {
		err = rows.Scan(&cid)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if exists(cid) != id {
		return fmt.Errorf("Course with ID %d does not exist", id)
	}

	// Delete from Coursedata
	deleteStatement := `delete from "coursedata" where courseid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	// Delete from Courses
	deleteStatement = `delete from "courses" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	return nil
}

// ListCourses lists all courses in the database
func ListCourses() ([]MSDSCourse, error) {
	Data := []MSDSCourse{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "id","cid","cname","cprereq"
		FROM "courses","coursedata"
		WHERE courses.id = coursedata.courseid`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var cid string
		var cname string
		var cprereq string
		err = rows.Scan(&id, &cid, &cname, &cprereq)
		temp := MSDSCourse{ID: id, CID: cid, CNAME: cname, CPREREQ: cprereq}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// UpdateCourse is for updating an existing course
func UpdateCourse(d MSDSCourse) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	id := exists(d.CID)
	if id == -1 {
		return errors.New("Course does not exist")
	}
	d.ID = id
	updateStatement := `update "coursedata" set "cname"=$1, "cprereq"=$2 where "courseid"=$3`
	_, err = db.Exec(updateStatement, d.CNAME, d.CPREREQ, d.ID)
	if err != nil {
		return err
	}

	return nil
}
