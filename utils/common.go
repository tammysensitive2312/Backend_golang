package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

// thực hiện insert theo batch, mỗi bacth 1000 row, mở tối đa 10 connection trong connenction pool
// insert 10000 row trong khoảng thời gian ngắn

func main() {
	// Kết nối đến database
	db, err := sql.Open("mysql", "root:truong@tcp(localhost:3306)/example_database_golang")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Tạo prepared statement để insert dữ liệu
	stmt, err := db.Prepare("INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Lặp qua từng file CSV trong thư mục "data"
	err = filepath.WalkDir("data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".csv" {
			// Đọc dữ liệu từ file CSV
			// Đọc dữ liệu từ file CSV
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				file.Close()
				return err
			}

			// Bỏ qua dòng đầu tiên nếu đó là tiêu đề
			startIndex := 0
			if len(records) > 0 && records[0][0] == "username" {
				startIndex = 1
			}

			// Thực hiện insert dữ liệu theo batch 1000 dòng
			for i := startIndex; i < len(records); i += 1000 {
				end := i + 1000
				if end > len(records) {
					end = len(records)
				}

				tx, err := db.Begin()
				if err != nil {
					file.Close()
					return err
				}

				for _, record := range records[i:end] {
					// Kiểm tra xem giá trị có phải là định dạng ngày tháng không
					if record[3] != "" && record[3] != "created_at" {
						createdAt, err := time.Parse("2006-01-02", record[3])
						if err != nil {
							tx.Rollback()
							file.Close()
							return err
						}

						updatedAt, err := time.Parse("2006-01-02", record[4])
						if err != nil {
							tx.Rollback()
							file.Close()
							return err
						}
						_, err = tx.Stmt(stmt).Exec(record[0], record[1], record[2], createdAt, updatedAt)
						if err != nil {
							tx.Rollback()
							file.Close()
							return err
						}
					}
				}

				err = tx.Commit()
				if err != nil {
					file.Close()
					return err
				}
			}

			file.Close()

		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Insert completed successfully!")
}
