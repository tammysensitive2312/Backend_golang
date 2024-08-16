package main

import (
	"database/sql"
	"encoding/csv"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
)

func connectDB() {
	db, err := sql.Open("mysql", "root:truong@tcp(localhost:3306)/example_database_golang")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = exportToCSV(db, "output.csv")
	if err != nil {
		log.Fatal("Error exporting to CSV:", err)
	}
}

func exportToCSV(db *sql.DB, filePath string) error {
	// Mở tệp CSV để ghi
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"id", "name", "project_spend", "project_variance", "revenue_recognised"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	rows, err := db.Query("SELECT id, name, category, project_spend, project_variance, revenue_recognised FROM projects")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, category string
		var projectSpend, projectVariance, revenueRecognised int

		err := rows.Scan(&id, &name, &category, &projectSpend, &projectVariance, &revenueRecognised)
		if err != nil {
			return err
		}

		record := []string{strconv.Itoa(id), name, category, strconv.Itoa(projectSpend), strconv.Itoa(projectVariance), strconv.Itoa(revenueRecognised)}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	connectDB()
}
