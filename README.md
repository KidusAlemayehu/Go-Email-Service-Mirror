# Email Service

This project implements an email service using Go. The service supports asynchronous email sending with attachments, retries for failed deliveries, and integration with RabbitMQ for task queuing.

## Features

- **Asynchronous email processing**: Ensures non-blocking API responses by processing emails in the background.
- **Attachment support**: Allows sending emails with file attachments encoded in base64 format.
- **Retry mechanism**: Automatically retries email delivery on failure, up to a configurable limit.
- **RabbitMQ integration**: Uses RabbitMQ for reliable task queuing and delivery.

## Prerequisites

Before running the project, ensure you have the following installed:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/email-service.git
   cd email-service
   ```

2. Create an `.env` file in the root directory with the following content:
   ```env
   SMTP_HOST=smtp.example.com
   SMTP_PORT=587
   SMTP_USER=your-email@example.com
   SMTP_PASSWORD=your-password
   RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   DATABASE_URL=your-database-connection-string
   DB_HOST=your-database-host # eg: db or localhost
   DB_PORT=your-database-port # eg: 5432
   DB_USER=your-database-username # eg: postgres
   DB_NAME=your-database-name # eg: postgres
   DB_PASSWORD=your-database-password 
   SECRET_KEY=your-secret-key
   ```

3. Build and start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```

### API Usage

The following API endpoint is available:

#### Send Email

**Endpoint:** `POST /api/v1/email`

**Request Body:**
```json
{
  "from": "sender@example.com",
  "to": "recipient@example.com",
  "subject": "Test Email",
  "body": "This is a test email.",
  "attachments": [
    {
      "filename": "test.txt",
      "content_type": "text/plain",
      "data": "base64-encoded-content"
    }
  ],
  "cc": ["cc@example.com"],
  "bcc": ["bcc@example.com"],
  "reply_to": "reply@example.com"
}
```

**Response:**
- Success: HTTP 202 Accepted
  ```json
  {
    "message": "Email task successfully enqueued",
    "data": null,
    "status_code":202
  }
  ```
- Failure: HTTP 400/500 with an error message

## Project Structure

```
email-service/
├── cmd/            # App Entry Point
│   ├── api/        # API
│   ├── woker/      # Woker logic
├── internal/
│   ├── http/       # HTTP Handlers
│   ├── dto/        # Data Transfer Objects
│   ├── models/     # Database models
│   ├── services/   # Business logic
├── migrate/        # Migration function
├── utils/          # Utility functions
├── Dockerfile      # Docker build configuration
├── docker-compose.yml # Docker Compose configuration
└── README.md       # Project documentation
```

## Running the Worker

The worker process listens for tasks from RabbitMQ and processes email sending tasks. The worker is automatically started when you run the `docker-compose up` command.

## Docker Configuration

The project uses Docker for containerization. The `docker-compose.yml` file defines two services:

- `api`: The main API service
- `worker`: The background worker for processing email tasks

### Build and Run

To start the application:
```bash
docker-compose up --build
```

To stop the application:
```bash
docker-compose down
```

## Environment Variables

| Variable       | Description                            |
|----------------|----------------------------------------|
| SMTP_HOST      | SMTP server host                      |
| SMTP_PORT      | SMTP server port                      |
| SMTP_USER      | SMTP username                         |
| SMTP_PASSWORD  | SMTP password                         |
| RABBITMQ_URL   | RabbitMQ connection URL               |
| DATABASE_URL   | Database connection string            |
| DB_NAME        | Database Name                         |
| DB_HOST        | Database Host                         |
| DB_USER        | Database User                         |
| DB_PORT        | Database Port                         |
| DB_PASSWORD    | Database password                     |

## License

This project is licensed under the [Apache 2.0 License](./LICENSE).


