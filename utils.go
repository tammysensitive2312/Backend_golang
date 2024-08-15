package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func connectDB() {
	db, err := sql.Open("mysql", "user:truong@tcp(localhost:3306)/example_database_golang")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Gọi hàm xuất dữ liệu ra CSV
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

	// Ghi tiêu đề của các cột vào tệp CSV
	headers := []string{"id", "name", "email", "created_at"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Truy vấn dữ liệu từ cơ sở dữ liệu
	rows, err := db.Query("SELECT id, name, category, project_spend, project_variance, revenue_recognised FROM projects")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Lấy các giá trị từ các hàng và ghi chúng vào tệp CSV
	for rows.Next() {
		var id int
		var name, category string
		var projectSpend, project_variance, revenue_recognised int

		err := rows.Scan(&id, &name, &category, &projectSpend, &projectSpend)
		if err != nil {
			return err
		}

		record := []string{strconv.Itoa(id), name, email, createdAt}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	// Kiểm tra lỗi khi quét các hàng
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func main() {

}
