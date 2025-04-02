# Media Collection API

A RESTful API for managing a personal media collection.

## Configuration

The application uses environment variables for configuration. You can set these variables directly in your environment or use a `.env` file.

### Setting up the .env file

1. Copy the `.env.example` file to `.env`:
   ```
   cp .env.example .env
   ```

2. Edit the `.env` file with your specific configuration values.

### Environment Variables

#### Database Configuration

You can configure the database connection in two ways:

1. **Individual parameters**:
   - `DB_HOST`: Database host (default: localhost)
   - `DB_PORT`: Database port (default: 5432)
   - `DB_USER`: Database user (default: postgres)
   - `DB_PASSWORD`: Database password (required)
   - `DB_NAME`: Database name (default: media_collection)

2. **Complete connection string**:
   - `DATABASE_URL`: Complete PostgreSQL connection string (overrides individual parameters)

#### API Configuration

- `API_PORT`: Port for the API server (default: 8080)

#### JWT Configuration

- `JWT_SECRET`: Secret key for JWT token generation (required)
- `JWT_EXPIRY_HOURS`: JWT token expiry time in hours (default: 24)

## Running the Application

1. Make sure you have Go 1.21 or later installed.

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run .
   ```

The server will start on the configured port (default: 8080).

## API Endpoints

### Public Endpoints

- `POST /api/register`: Register a new user
- `POST /api/login`: Login with username and password

### Protected Endpoints (requires authentication)

- `GET /api/user/profile`: Get user profile
- `PUT /api/user/profile`: Update user profile
- `POST /api/user/password`: Change password
- `POST /api/user/2fa/setup`: Set up 2FA
- `POST /api/user/2fa/enable`: Enable 2FA
- `POST /api/user/2fa/disable`: Disable 2FA
- `GET /api/media`: List media items
- `POST /api/media`: Create a media item
- `GET /api/media/:id`: Get a media item
- `PUT /api/media/:id`: Update a media item
- `DELETE /api/media/:id`: Delete a media item 