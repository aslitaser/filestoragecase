## Running the Code

To run the File Storage API, follow these steps:

1. Clone the repository:

```
git clone https://github.com/aslitaser/filestoragecase.git
```

2. Change directory to the project folder:

```
cd filestoragecase
```

3. Make sure you have the required dependencies installed:

```
go get -u github.com/gorilla/mux
go get -u github.com/dgrijalva/jwt-go
```

4. Set the JWT secret environment variable:

```
export JWT_SECRET=your_jwt_secret
```

5. Start the MongoDB server and make sure it is running on `mongodb://localhost:27017`.

6. Run the main.go file:

```
go run main.go
```

The server will start listening on `localhost:8080`.

## System Design

The File Storage API is built using a modular and scalable design. The main components of the system are:

- Database: MongoDB is used as the database to store user and file metadata information.
- Web Server: A RESTful API server is built using the Gorilla Mux library for routing and handling HTTP requests.
- Authentication: JSON Web Tokens (JWT) are used for user authentication and authorization.
- Handlers: Handlers are responsible for processing incoming HTTP requests and generating responses.

The system is designed to be easily extensible by adding new routes and handlers as needed.

## Libraries and Frameworks

The following libraries and frameworks are used in the File Storage API:

- [Go](https://golang.org/): The programming language used to build the API.
- [MongoDB](https://www.mongodb.com/): The database used to store user and file metadata information.
- [Gorilla Mux](https://github.com/gorilla/mux): A powerful HTTP router and URL matcher for building Go web servers.
- [JWT-go](https://github.com/dgrijalva/jwt-go): A Go (Golang) implementation of JSON Web Tokens (JWT).

## Design Choices

- **Modular design**: The API is designed to be modular, with each component being responsible for a specific functionality. This makes it easy to extend the API and add new features as needed.
- **Scalable architecture**: The API is built with a scalable architecture, allowing it to handle a high volume of requests efficiently.
- **Secure authentication**: JSON Web Tokens (JWT) are used for secure user authentication and authorization.
- **Error handling**: Proper error handling is implemented throughout the API, ensuring that meaningful error messages are returned to the client in case of failures.
- **File storage**: Files are stored on the local file system, but the design could be easily extended to use cloud storage services like Amazon S3 or Google Cloud Storage.


# File Storage API Documentation

## Overview

The File Storage API allows users to upload, download, and delete files while managing user accounts and permissions. Only authenticated users can upload, download, or delete files. The API is built using Go, MongoDB, Gorilla Mux, and JWT-go.

## Base URL

```
http://localhost:8080
```

## Endpoints

### User Service

#### Register

- **URL**: `/register`
- **Method**: `POST`
- **Description**: Registers a new user with a username and password, and returns a JWT upon successful registration.
- **Request Body**:

```json
{
  "username": "example",
  "password": "password"
}
```

- **Response**:

```json
{
  "token": "your_jwt_token"
}
```

#### Login

- **URL**: `/login`
- **Method**: `POST`
- **Description**: Authenticates a user with a username and password, and returns a JWT upon successful authentication.
- **Request Body**:

```json
{
  "username": "example",
  "password": "password"
}
```

- **Response**:

```json
{
  "token": "your_jwt_token"
}
```

### File Service

#### Upload

- **URL**: `/files/upload`
- **Method**: `POST`
- **Description**: Allows a user to upload a file to the server. Stores file metadata (including `userID`, `filename`, `uploadDate`, `fileSize`) in MongoDB. Requires a valid JWT in the `Authorization` header.
- **Headers**:

```
Authorization: Bearer your_jwt_token
```

- **Request Body**: `multipart/form-data` with a file field named `file`.
- **Response**:

```json
{
  "message": "File uploaded successfully",
  "fileMetadata": {
    "userID": "user_id",
    "filename": "file_name",
    "uploadDate": "date_time",
    "fileSize": "file_size"
  }
}
```

#### Download

- **URL**: `/files/download/{id}`
- **Method**: `GET`
- **Description**: Allows a user to download a file by its `ID`. Requires a valid JWT in the `Authorization` header.
- **Headers**:

```
Authorization: Bearer your_jwt_token
```

- **Response**: The requested file.

#### Delete

- **URL**: `/files/delete/{id}`
- **Method**: `DELETE`
- **Description**: Allows a user to delete a file by its `ID`. Requires a valid JWT in the `Authorization` header.
- **Headers**:

```
Authorization: Bearer your_jwt_token
```

- **Response**:

```json
{
  "message": "File deleted successfully"
}
```

## Error Handling

The API returns appropriate error messages and status codes in case of failures, such as invalid JWT, file not found, or invalid user input.

## Testing

Included meaningful unit tests for each functionality, covering various scenarios.
