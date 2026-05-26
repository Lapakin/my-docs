# my-docs
This is a microservices-based file storage system designed for comprehensive file and folder management.

## Description
My-docs is a Go-based file storage service that handles user management, file/folder operations, and sharing capabilities. It features an API Gateway (KrakenD) for routing, a PostgreSQL database for metadata, and MinIO for object storage, all presented through an integrated web UI.

## Setup

### 1. Clone the Repository
First, download the project files by running the following command:

```bash
git clone https://github.com/Lapakin/my-docs.git

```

### 2. Configure Environment Variables

Before launching the application, create a `.env` file in the project root to store your database and object storage credentials. Use the following format as an example:

```env
DB_USER=dev
DB_PASSWORD=12345

MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=admin123

```

### 3. Launch the Application

Start the services using Docker and Docker Compose via the provided Makefile:

```bash
make up

```

## Services and Endpoints

| Service | Endpoint |
| --- | --- |
| Web UI / API Gateway | `http://localhost:8080` |

## Unit Testing

To execute the full suite of unit tests, use the following command:

```bash
make go-unit-tests

```
