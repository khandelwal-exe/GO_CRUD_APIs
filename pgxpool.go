package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func main() {
	var (
		username string
		password string
	)

	// Prompt the user for the database username and password
	fmt.Print("Enter the database username: ")
	_, _ = fmt.Scan(&username)
	fmt.Print("Enter the database password: ")
	_, _ = fmt.Scan(&password)

	// Replace these with your actual database connection details
	host := "localhost"
	port := 5432
	dbname := "OfficeManagement"

	// Create a database connection
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Println("Unable to connect to the database:", err)
		return
	}
	defer conn.Close(context.Background())

	// Query the "users" table
	rows, err := conn.Query(context.Background(), "SELECT id, username, email FROM users")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	// Process the query results
	for rows.Next() {
		var id int
		var username, email string
		err := rows.Scan(&id, &username, &email)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Printf("User ID: %d, Username: %s, Email: %s\n", id, username, email)
	}
}
