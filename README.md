# user-actions-api

### **Running the Application and Tests Using the `Makefile`**

#### 1. **To Run the Application Locally**:
   You can use the following command to run the API server locally:
   ```bash
   make run
   ```
   This will use `go run main.go` to start the application and start the API server on `localhost:8080` by default.

#### 2. **To Run Tests**:
   The `make test` command will run all the tests in your Go project:
   ```bash
   make test
   ```
   This will run `go test ./...` to execute all tests in the project.

---

### Endpoints:

---

### 1. **`GET /users/:id`**  
   **Description**:  
   Retrieves a user by their unique `id`.

   - **Success (StatusOK)**: Returns the user data, including `id`, `name`, and `createdAt` fields.  
     Example response:
     ```json
     {
       "id": 2,
       "name": "Alice",
       "createdAt": "2021-07-04T12:47:09.888Z"
     }
     ```

   - **Error (StatusNotFound)**: If the user with the provided `id` does not exist.  
   
   - **Error (StatusBadRequest)**: If the `id` is invalid or missing in the request.

---

### 2. **`GET /users/:id/actions/count`**  
   **Description**:  
   Retrieves the count of actions performed by a user with the specified `id`.

   - **Success (StatusOK)**: Returns the number of actions taken by the user.  
     Example response:
     ```json
     {
       "count": 5
     }
     ```

   - **Error (StatusBadRequest)**: If the `id` is invalid or missing in the request.

---

### 3. **`GET /actions/:type/next-probability`**  
   **Description**:  
   Retrieves the probability of the next action of a specified type.

   - **Success (StatusOK)**: Returns the probability of the next action of the specified type.  
     Example response:
     ```json
     {
       "CONNECT_CRM": 0.5,
       "VIEW_CONTACTS": 0.5
     }
     ```
   
   - **Error (StatusBadRequest)**: If the `type` is invalid or missing in the request.

   - **Error (StatusNotFound)**: If no data is available for the given action type.

---

### 4. **`GET /users/referal-index`**  
   **Description**:  
   Retrieves the referral index for users.

   - **Success (StatusOK)**: Returns the referral index data.
     Example response:
     ```json
     {
       "1": 3,
       "3": 7
     }
     ```
   
---