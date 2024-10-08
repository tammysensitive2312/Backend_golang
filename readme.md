# Hướng dẫn deploy backend lên docker bằng Dockerfile

Dockerfile này được thiết kế để xây dựng một ứng dụng Go bằng Docker. Dưới đây là các bước thực hiện trong Dockerfile để tạo ra một Docker image có thể sử dụng để chạy ứng dụng Go.

### Build image backend api 

```dockerfile
sử dụng iamge go chính thức làm base image nên đặt cùng với phiên bản 
GO ở go.mod
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
## Deploy 
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

# API OVERVIEW
## Endpoints
### Lấy danh sách người dùng
- GET: /golang-web/api/projects/all?page=1&pageSize=10
- Tham số query:
    - `page`: số trang (mặc định: 1)
    - `pageSize`: số lượng kết quả trên mỗi trang (mặc định: 10)
- Response:
  ```json
  {
        "projects": [
            {
                "ID": 1,
                "Name": "Project 1",
                "Category": "Category A",
                "ProjectSpend": 1000,
                "ProjectVariance": 100,
                "RevenueRecognised": 0,
                "ProjectStartedAt": "2024-08-16T07:00:00+07:00",
                "ProjectEndedAt": null,
                "CreatedAt": "2024-08-15T14:36:24.129+07:00",
                "UpdatedAt": "2024-08-15T14:36:24.129+07:00",
                "DeletedAt": null,
                "Users": null
            }
        ],
        "total_pages": 1,
        "total_records": 1,
        "current_page": 1
    },

### Lấy access token và refresh token
- POST: golang-web/api/users/login
- Tham số query:
  - `email`
  - `password`
- Response:
``` json
 {
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiZXhwIjoxNzIzOTE0MDU2fQ.ICd-e1VBcyaSZW2WJxCmvYzajKUeiWw83_bd4D_hWRM",
  "refresh_token": "Ql6E2cqPPmaAVFIkejPd7bEi7hx1F-EYdjGtNXlaYg7y_8AJHCjfrQZESTY7vUV0DIcYu8eH0H7SEfj0zNvbJcDNqRbikF1NcJd0yLL1gkYliogEpMtJVF7jepR9Gp_fmRbPB-UfIoTcRBodSq8t2BlBM0XgX-rZkJEEDfMbYPwdRsPmi2W3O4y-YtI46ZI="
 }
```

### Lấy access token dựa trên refresh token
- POST: golang-web/api/refresh
- Tham số query:
  - `refresh_token`: refresh_token được client lưu trữ ở đâu đó để request
- Response:
``` json
 {
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiZXhwIjoxNzIzOTE0MDU2fQ.ICd-e1VBcyaSZW2WJxCmvYzajKUeiWw83_bd4D_hWRM"
 }
```

## Vấn đề còn tồn tại 
- quản lý secret-key ở mã nguồn 