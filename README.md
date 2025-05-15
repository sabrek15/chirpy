# Chirpy API Documentation

## **Health Check**

### **GET /api/healthz**

- **Description**: Checks the readiness of the server.
- **Response**:
  - **200 OK**: Returns "OK".

---

## **Metrics**

### **GET /admin/metrics**

- **Description**: Retrieves server metrics.
- **Response**:
  - **200 OK**: Returns metrics in plain text.

---

## **Reset Users**

### **POST /admin/reset**

- **Description**: Deletes all users. Only available in the development environment.
- **Response**:
  - **200 OK**: Users deleted successfully.
  - **403 Forbidden**: Access denied in non-development environments.
  - **500 Internal Server Error**: Failed to delete users.

---

## **Create User**

### **POST /api/users**

- **Description**: Creates a new user.
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response**:
  - **201 Created**: Returns the created user.
  - **400 Bad Request**: Invalid request body.
  - **404 Not Found**: Failed to create the user.
  - **409 Conflict**: Password hashing failed.

---

## **Update User**

### **PUT /api/users**

- **Description**: Updates user credentials.
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response**:
  - **200 OK**: Returns updated user details.
  - **400 Bad Request**: Invalid request body or token.
  - **401 Unauthorized**: Invalid or missing token.

---

## **Login**

### **POST /api/login**

- **Description**: Logs in a user and generates access and refresh tokens.
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response**:
  - **200 OK**: Returns user details, access token, and refresh token.
  - **400 Bad Request**: Invalid request body or token generation failed.
  - **401 Unauthorized**: Incorrect password.
  - **404 Not Found**: User not found.

---

## **Refresh Token**

### **POST /api/refresh**

- **Description**: Refreshes the access token using a valid refresh token.
- **Request Headers**:
  - `Authorization: Bearer <refresh_token>`
- **Response**:
  - **200 OK**: Returns a new access token.
  - **401 Unauthorized**: Invalid, revoked, or expired refresh token.

---

## **Revoke Refresh Token**

### **POST /api/revoke**

- **Description**: Revokes a refresh token.
- **Request Headers**:
  - `Authorization: Bearer <refresh_token>`
- **Response**:
  - **204 No Content**: Token revoked successfully.
  - **401 Unauthorized**: Invalid or missing token.

---

## **Post Chirp**

### **POST /api/chirps**

- **Description**: Creates a new chirp.
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Request Body**:
  ```json
  {
    "body": "string"
  }
  ```
- **Response**:
  - **201 Created**: Returns the created chirp.
  - **400 Bad Request**: Invalid request body or chirp too long.
  - **401 Unauthorized**: Invalid or missing token.
  - **404 Not Found**: Failed to create chirp.

---

## **Get Chirps**

### **GET /api/chirps**

- **Description**: Retrieves all chirps or chirps by a specific author.
- **Query Parameters**:
  - `author_id` (optional): UUID of the author.
- **Response**:
  - **302 Found**: Returns the list of chirps.
  - **400 Bad Request**: Invalid author ID.
  - **404 Not Found**: Chirps not found.

---

## **Get Chirp by ID**

### **GET /api/chirps/{chirpid}**

- **Description**: Retrieves a chirp by its ID.
- **Path Parameters**:
  - `chirpid`: UUID of the chirp.
- **Response**:
  - **201 Created**: Returns the chirp.
  - **404 Not Found**: Chirp not found or invalid ID.

---

## **Delete Chirp by ID**

### **DELETE /api/chirps/{chirpid}**

- **Description**: Deletes a chirp by its ID.
- **Request Headers**:
  - `Authorization: Bearer <token>`
- **Path Parameters**:
  - `chirpid`: UUID of the chirp.
- **Response**:
  - **200 OK**: Chirp deleted successfully.
  - **403 Forbidden**: User is not the owner of the chirp.
  - **404 Not Found**: Chirp not found or invalid ID.
  - **401 Unauthorized**: Invalid or missing token.

---

## **Polka Webhook**

### **POST /api/polka/webhooks**

- **Description**: Handles Polka webhook events.
- **Request Headers**:
  - `Authorization: ApiKey <key>`
- **Request Body**:
  ```json
  {
    "event": "string",
    "data": {
      "user_id": "string"
    }
  }
  ```
- **Response**:
  - **204 No Content**: Event handled successfully.
  - **400 Bad Request**: Invalid request body.
  - **401 Unauthorized**: Invalid or missing API key.
  - **404 Not Found**: User ID does not exist.

---

## Contributing

We welcome contributions to Chirpy! To contribute, follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature-or-bugfix-name
   ```
3. Make your changes and commit them with clear and concise messages:
   ```bash
   git commit -m "Description of changes"
   ```
4. Push your changes to your fork:
   ```bash
   git push origin feature-or-bugfix-name
   ```
5. Open a pull request to the main repository.

### Guidelines

- Ensure your code follows the project's coding standards.
- Write tests for any new functionality.
- Update the documentation if necessary.

Thank you for contributing to Chirpy!
