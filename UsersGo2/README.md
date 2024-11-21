# Code Anatomy
# Stack
### Golang
description

### Postgres
description

## imports
(note: tried to use a little packages as possible)

- main.go
    - fmt
        - provides functions for formatted I/O, inc printing to the console
    - log
        - provides logging functionality, specifically for error messages/debugging
- api.go
    - encoding/json
        - dfsdfgdss.
    - net/http
    - github.com/gorilla/mux
        - ghjgh
- storage.go
    - database/sql
        - Provides the core database interface in Go. It defines sql.DB for database connections and methods for executing querie
    - github.com/lib/pq
        - The PostgreSQL driver for Go. The _ import means its init functions are executed to register the driver with database/sql.
- types.go
    - math/rand
    - time

## File Structure
### makefile
- Contains instructions used for automating specific tasks
    - build - compiles the source code in the current directory into an executable binary, to the sepcified output file (-o bin/usersgo2). note: "@" prevents the command from being echoed in the terminal (cleaner approach). after running an executable is created in the bin directory
    - run: build - dependant on build target, when make run is invoked it will first execute the build target, and then runs the compiled binary. after running the program is executed.
    - test - [placeholder], eventually will run unit tests.

# main.go
#### Execution Flow:
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
newAPIServer method called to create a new API server instance. ":3000" specifies the port in use and store is passing the database connection, allowing it to interact with the database
``server.Run()``
called the Run() method on the server object to start the API Server

# storage.go

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
#NewPostgressStore
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
