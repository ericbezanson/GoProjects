
# Code Anatomy


# Stack
### Golang
Golang, or Go, is a modern, open-source programming language developed by Google. It is designed for simplicity, efficiency, and scalability, making it ideal for building reliable, high-performance applications.

### Postgres
descriPostgreSQL (often referred to as Postgres) is a powerful, open-source relational database management system. It adheres to the SQL standard while also offering additional features

## imports
(note: tried to use a little packages as possible)

| **File**       | **Import**                                   | **Description**                                                                                                                                       |
|-----------------|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| **main.go**     | `fmt`                                       | Provides functions for formatted I/O, including printing to the console.                                                                              |
|                 | `log`                                       | Provides logging functionality, specifically for error messages and debugging.                                                                        |
| **api.go**      | [`encoding/json`](https://pkg.go.dev/encoding/json) | Implements encoding and decoding of JSON as defined in RFC 7159. Includes `Marshal` and `Unmarshal` functions for mapping JSON and Go values.         |
|                 | [`net/http`](https://pkg.go.dev/net/http)   | Provides HTTP client and server implementations.                                                                                                      |
|                 | [`github.com/gorilla/mux`](https://github.com/gorilla/mux) | An HTTP request multiplexer. Matches incoming requests to registered routes and executes the corresponding handler.                                   |
| **storage.go**  | `database/sql`                              | Provides the core database interface in Go. Defines `sql.DB` for database connections and methods for executing queries.                              |
|                 | `github.com/lib/pq`                        | PostgreSQL driver for Go. The `_` import ensures its `init` functions are executed, registering the driver with `database/sql`.                       |
| **types.go**    | `math/rand`                                 | Provides pseudo-random number generation.                                                                                                             |
|                 | `time`                                      | Offers functionality for measuring and displaying time, as well as sleep and time-based operations.                                                   |


## File Structure
### makefile
- Contains instructions used for automating specific tasks
    - build - compiles the source code in the current directory into an executable binary, to the sepcified output file (-o bin/usersgo2). note: "@" prevents the command from being echoed in the terminal (cleaner approach). after running an executable is created in the bin directory
    - run: build - dependant on build target, when make run is invoked it will first execute the build target, and then runs the compiled binary. after running the program is executed.
    - test - [placeholder], eventually will run unit tests.

# main.go
#### Summary
1. connect to PostgreSQL
2. Initialize the Database
3. Create the API Server
4. Start the API Server

#### Code Explaination
`` main()`` 
entry point to every go application, starts the programs execution
[``store, err := NewPostgressStore()``](#NewPostgressStore) 
...._(function built in storage.go)_
creates a new instance of a PostgresStore object
establishes a connection to a PostgreSQL database 
`` if err != nil ``
error handling, if error occurs when connecting to db, logs error and terminates immediately using log.Fatal(err)
``sotre.Init()`` 
is a method called on the PostgresStore object, returning an error if init fails and log.Fatal's it
``fmt.Printf("%+v\n", store)``
prints store object, "%+v\n" formates the object with field names (if it's a struct) (more detailed output)
``server := newAPIServer(":3000", store)`` 
...._(function built in api.go)_
[newAPIServer](#newapiserver) method called to create a new API server instance. ":3000" specifies the port in use and store is passing the database connection, allowing it to interact with the database
[``server.Run()``](#serverrun)
called the Run() method on the server object to start the API Server

# storage.go
#### Summary
1. Defines a Storage interface for account CRUD operations.
2. Implements the interface with a PostgresStore that connects to a PostgreSQL database.
3. Provides methods to initialize the database, create accounts, and retrieve all accounts.
4. Has placeholder methods for updating, deleting, and retrieving accounts by ID.
5. Includes good practices like parameterized queries to prevent SQL injection and error handling for database operations.

``type Storage interface{}``
_interface: A set of method signatures The purpose of defining the Storage interface is to describe what operations can be performed on account-related data without being tied to a specific implementation._

defines the Storage interface***, acting as an abstraction for account-related database opperations. whenever implemented you must define the following methods:
```sh
	CreateAccount(*Account) error - args: pointer to Account*** struct return: error
	DeleteAccount(int) error - args: int return error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)args slice of pointer to Account, error
	GetByAccountID(int) (*Account, error)
```
*** [interface](https://go.dev/tour/methods/9), [pointer](https://go.dev/tour/moretypes/1), [slice](https://go.dev/tour/moretypes/7)

``` 
type PostgresStore struct {
    db *sql.DB
}
```
PostgressStore is a [struct](https://go.dev/tour/moretypes/2) that implements the Storage interface
_note: Since PostgresStore has all the methods that the Storage interface specifies, it automatically implements the Storage interface. There’s no need for any explicit relationship or reference to the Storage interface._

#### Why doesn't PostgresStore need to "call" or "reference" the interface?

Go's interface implementation is structural, not declarative. This means:
The Go compiler checks whether a type has the methods defined in the interface.
if a type has all the methods, it is automatically considered to "implement" the interface.
There’s no need for explicit syntax to declare this relationship.

#### Key Benefits of this Approach
1. Decoupling: You can write code that depends on the `Storage` interface rather than a specific implementation like `PostgresStore`
2. Flexibility: To switch the storage mechanism, you only need to implement the `Storage` interface for the new type. The rest of the code using the interface remains unchanged.
3. Polymorphism: You can use different implementations of `Storage` interchangeably.

--
<h4 id="newPostgresStore"></h4>
```
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=R3dsp@ce sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
```
- Establishes a connection to the PostgreSQL database and returns a new PostgresStore
- connStr specifies the connection details for PostfreSQL (user,dbname,password, sslmode)
- [sql.Open](https://pkg.go.dev/database/sql#Open), Opens the database handle using the PostgreSQL driver. (it does not verify the connection at this point; it just prepares the database handle)
- db.Ping is a test query to ensure the database is reachable
###### note: 
_A database driver is a software component that allows applications to interact with a specific DBMS (Database Management System). It acts as a bridge between an application and a database, making communication and data exchange standardized and easier._
###### note: 
_& returns the memory address of the following variable. * returns the value of the following variable_

-- 

###### init database
```
func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}
```
Calls the CreateAccountTable method to ensure the required database schema exists.

--

###### creating the account table
```
func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}
```
- query: an SQL statement to create the account table
- s.db.Exec: execustes the query, returns an error if it fails
###### note:
_In Go, := is for declaration + assignment, whereas = is for assignment only._

--

###### creating an account
```
func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account
		(first_name, last_name, number, balance, created_at )
		values
		($1, $2, $3, $4, $5)`
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	fmt.Println("%+v\n", resp)
	return nil
}
```
- Inserts a new account into the database
- SQL Query: uses placeholders ($1, $2 etc) for parameterized queries to prevent SQL injection
    - inserts values from the Account object (acc)
- s.db.Query:
    - executes the querty and returns result set 

###### Get all Accounts
```
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

```
- Retrieves all rows from the account table and maps them to Account objects.
- rows.Scan: maps each column in a row to corresponding fields in the Account struct


# api.go
### summary 
- Step 1: Server Init
    - Server Init - newAPIServer is called with the listen address and storage
    - Creates an APIserver instance with the specified configurations
- Step 2: Running the Server
    - Run is called to setup routes and start the server
- Step 3: Request Handling
    - When a request is recieved:
        - Mux matches the route (e.g. ``/account`` or ``/account/{id}``)
        - the associated handler is executed (eg. handleAccount)
        - Operation that takes place is based on method (GET, POST, DELETE, etc)
- Step 4: Response Handling
    -  the handler (or makeHTTPHandleFunc) sends a JSON response using the ``WriteJSON`` helper
    -  
### struct definitions
``APIserver``
- represents the API server. listenAddr = the address (eg ":3000"), store = Am instance of a type implementing the ``Storage`` interface (eg PostgresStore) to handle db opperations.

``ApiError`` 
- used for formatting error responses as JSON objects.

### functions
###### new api server
<h4 id="newapiserver"></h4>
```
func newAPIServer(listenAddr string, store Storage) *APIserver {
	return &APIserver{
		listenAddr: listenAddr,
		store:      store,
	}
}

```
- a constructor used to create an APIserver instance
- returns a pointer to a new APIserver

###### api server Run
<h4 id="serverrun"></h4>
```
func (s *APIserver) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	log.Println("JSON API Server running on port:", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```
- Starts the APIserver
    - create a new router using the Gorilla Mux Library
    - Registers routes and their corresponding handler functions
        - "/account" => handleAccount
        - "/account/{id}" => handleGetAccountByID
    - calls ``http.ListenAndServer`` to bind the server to listenAddr and start handling requests

###### note: 
_routes are registered with ``makeHTTPHandleFunc``, a helper described below, to wrap handlers (apiFunc) in error handling logic. if ``http.ListenAndServe`` fails, it logs the error and terminates the application_
###### note:
_``http.ListenAndServe``: ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections. Accepted connections are configured to enable TCP keep-alives._
###### note: 
_``makeHTTPHandleFunc: A helper that wraps apiFunc handlers to handle errors gracefully. step 1: call the given apiFunc, step 2: if the function returns an error, sends a JSON response with the error message and a ``400 Bad Request`` status_

###### handleAccount
```
func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("mothod not supported", r.Method)
}
```
- handles all opperations on the "/account" endpoint
- chooses logic based on HTTP method
    - GET: fetches all accounts using handleGetAccount
    - POST: creates a new account using handleCreateAccount
    - DELETE: Deletes an account using handleDeleteAccount
- returns error if method is unsupported
- 

###### handle get account
```
func (s *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}
```
- fetches all accounts from storage (s.store.GetAccounts)
- Encodes the accounts as JSON using WriteJSON
###### note
_WriteJSON: A helper that Encodes a Go value(v) as JSON and write it to the response, adds the ``Content-Type: application/json`` header and sets the HTTP status code (status)_
```
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
```
###### handle get account by ID
```
func (s *APIserver) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)
	fmt.Println(id)
	// account := NewAccount("Eric", "Bezanson")
	return WriteJSON(w, http.StatusOK, &Account{})
}
```
- Handles the "/account/{id}" endpoint to fetch a specific account by its ID
- uses ``mux.Vars(r)`` to extract the ``id`` fromt he route parameters
- currently returns a placeholder Account (&Account{})

###### handle Create Account
```
func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreatAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}
```
- Creates a new account based on JSON data from the request body:
    - Decodes the request body into a ``CreateAccountRequest`` object
    - Uses ``NewAccount`` to create an ``Account`` instance
    - persists the account via ``s.store.CreateAccount``
    - sends the created account as a JSON response

###### example
Request POST: ``/account``
```
{
  "firstName": "John",
  "lastName": "Doe"
}
```
Response (200 OK)
```
{
  "id": 1,
  "firstName": "John",
  "lastName": "Doe",
  "number": 1234,
  "balance": 0,
  "createdAt": "2024-11-21T15:30:00Z"
}

```

### Key Concepts

    1. Gorilla Mux: Handles request routing and extracting route parameters.
    2. JSON Serialization: Uses encoding/json for marshaling and unmarshaling.
    3. Error Handling: Centralized error handling via makeHTTPHandleFunc.
    4. Storage Interface: Interacts with the database through the Storage abstraction.
    5. RESTful API Design: Implements HTTP methods for different operations (GET, POST, etc.).