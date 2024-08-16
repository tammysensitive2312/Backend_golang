# Hướng dẫn deploy backend lên docker bằng Dockerfile

Dockerfile này được thiết kế để xây dựng một ứng dụng Go bằng Docker. Dưới đây là các bước thực hiện trong Dockerfile để tạo ra một Docker image có thể sử dụng để chạy ứng dụng Go.

## Image nền tảng (Base Image)

```dockerfile
FROM golang:1.21 AS build

Thiết lập thư mục làm việc
WORKDIR /app

sao chép các tệp khai báo phụ thuộc
sao chép các tệp go.mod và go.sum từ máy tính của bạn vào thư mục /app trong container. 
Các tệp này chứa danh sách các thư viện phụ thuộc mà ứng dụng Go của bạn yêu cầu.

COPY go.mod ./
COPY go.sum ./

Câu lệnh này tải về tất cả các thư viện phụ thuộc được liệt kê trong go.mod và go.sum. 
Bằng cách chạy lệnh này riêng biệt.
Docker có thể cache bước này nghĩa là nếu không có sự thay đổi trong go.mod và go.sum
Docker sẽ không cần tải lại các thư viện giúp tăng tốc quá trình build
RUN go mod download

Sao chép mã nguồn ứng dụng
COPY . ./

RUN ls

biên dịch ứng dụng 
RUN go build -o myrepo-test .

```
## Chạy ứng dụng kèm theo db 
``` docker-compose.yml

version: '3.8'

services:
  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: truong
      MYSQL_DATABASE: example_database_golang
      MYSQL_USER: root
      MYSQL_PASSWORD: truong
    ports:
      - "3306:3306"
    networks:
      - backend-network

// build và chạy image server back-end
  backend:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
    environment:
      - DB_HOST=mysql-server
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=truong
      - DB_NAME=example_database_golang
    networks:
      - backend-network

// db và back-end giao tiếp thông qua một kênh chung tên là networks
networks:
  backend-network:
    driver: bridge

```


