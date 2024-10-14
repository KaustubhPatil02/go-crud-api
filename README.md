# Step 1: Install Dependencies
echo "Installing Go packages..."
```bash
go install github.com/gin-gonic/gin@latest
go install github.com/lib/pq@latest
echo "Dependencies installed."
```

# Step 2: Set Up PostgreSQL Database
echo "Setting up PostgreSQL database..."
```bash
DB_NAME="go_crud_db"
DB_USER="postgres"
DB_PASSWORD="<YOUR_PASSWORD>"
```
## Instructions
Replace <YOUR_PASSWORD> with your actual PostgreSQL password in the script.

# Create the database
psql -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;"

# Create the items table
psql -U "$DB_USER" -d "$DB_NAME" -c "
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);"

echo "Database and table created."

# Step 3: Configure the connection string in main.go
CONN_STR="postgres://$DB_USER:$DB_PASSWORD@localhost/$DB_NAME?sslmode=disable"
sed -i "s|connStr := .*|connStr := \"$CONN_STR\"|" main.go
echo "Connection string configured in main.go."

# Step 4: Run the application
echo "Starting the Go application..."
```bash
go run main.go
```


# Testing the API using Postman
1. POST /items - Create an Item
To create an item, use the following command:
    ```
    curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name": "Item1"}'
    ```

2. GET /items - Get All Items
To fetch all items, use:
    ```
    curl -X GET http://localhost:8080/items
    ```

3. GET /items/:id - Get a Specific Item
To get a single item by ID, use:
    ```
    curl -X GET http://localhost:8080/items/1
    ```
4. PUT /items/:id - Update an Item
To update an existing item by ID, use:
    ```
    curl -X PUT http://localhost:8080/items/1 -H "Content-Type:     application/json" -d '{"name": "Updated Item"}' 
    ```
5. DELETE /items/:id - Delete an Item
To delete an item by ID, use:
    ```
    {
    "message": "Item deleted"
    }
    ```