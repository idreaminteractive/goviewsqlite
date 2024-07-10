package commands

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "modernc.org/sqlite"
)

func RunSeed(w io.Writer, dbPath, schemaPath string) error {

	// Open the SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Fprintf(w, "Error %v\n", err)
		return err
	}
	defer db.Close()

	// Read the SQL DDL file
	ddlBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Fprintf(w, "Error %v\n", err)
		return err
	}
	ddl := string(ddlBytes)

	// Execute the SQL statements
	_, err = db.Exec(ddl)
	if err != nil {
		fmt.Fprintf(w, "Error %v\n", err)
		return err
	}

	fmt.Fprint(w, "Database schema created successfully!\n")
	return nil
}
