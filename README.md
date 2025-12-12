# HRMS System (Attendance & Student Management) API

A minimal Human Resource Management System (HRMS) API for managing students and their attendance. This project is a backend system built with Golang and the Gin framework, following a Model-View-ViewModel-Controller (MVVC) architecture.

## Key Features

- **Student Management**: Full CRUD (Create, Read, Update, Delete) functionality for student records.
- **Attendance Management**: Mark and view student attendance.
- **Automated Reporting**: Cron jobs to generate and display weekly and monthly attendance reports to the console.

## Technology Stack

- **Language**: Go
- **Framework**: Gin
- **Database**: MySQL
- **Cron Scheduler**: robfig/cron (v3)
- **Containerization**: Docker

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.25.1 or newer)
- Docker and Docker Compose

### Installation & Setup

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/your-repo-name.git
    cd your-repo-name
    ```

2.  **Start the database:**
    Use Docker Compose to start the MySQL database container. This command will also initialize the database schema using the `init.sql` file.
    ```sh
    docker-compose up -d
    ```

3.  **Install dependencies:**
    Download the required Go modules.
    ```sh
    go mod tidy
    ```

4.  **Run the application:**
    ```sh
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

## API Endpoints

### Student Management

- `POST /students`
  - **Description**: Creates a new student.
  - **Body**: `{"name": "John Doe", "email": "john.doe@example.com"}`

- `GET /students`
  - **Description**: Retrieves a list of all students.

- `GET /students/:id`
  - **Description**: Retrieves a single student by their ID.

- `PUT /students/:id`
  - **Description**: Updates an existing student's details.
  - **Body**: `{"name": "Johnathan Doe", "email": "john.doe.new@example.com"}`

- `DELETE /students/:id`
  - **Description**: Deletes a student by their ID.

### Attendance Management

- `POST /attendance/mark`
  - **Description**: Marks attendance for a student on a specific date.
  - **Body**: `{"student_id": 1, "date": "2025-12-12T10:00:00Z", "status": "Present"}`

- `GET /attendance/:student_id`
  - **Description**: Retrieves all attendance records for a specific student.

## API Documentation (Swagger)

This project uses Swagger (OpenAPI) for interactive API documentation.

1.  **Generate Swagger Documentation (if not already generated):**
    ```sh
    go run github.com/swaggo/swag/cmd/swag init
    ```
2.  **Run the application:**
    ```sh
    go run main.go
    ```
3.  **Access Documentation:**
    Once the server is running, open your web browser and navigate to:
    [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Running Tests

To run the suite of unit tests, execute the following command from the root directory:

```sh
go test ./...
```

## Contribution

Contributions are welcome! Please feel free to submit a pull request.

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request
