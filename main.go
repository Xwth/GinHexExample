package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"

	_ "github.com/jinzhu/gorm/dialects/postgres" // Necessary for the postgres gorm.Open()
)

// All fields regardless of tags are used by gorm
// You can't ignore marshall or unmarshall like a
// One way mirror, unless you implement the
// Unmarshal/Marshal interface funcs
type user struct {
	gorm.Model

	email     string `gorm:"unique"`
	firstname string
	lastname  string
	creds
}

type creds struct {
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

// Extends the marshal method and
// strips the password from the marshalling
func (c user) MarshalJSON() ([]byte, error) {
	type u user // prevents recursion
	x := u(c)
	x.creds.Password = ""
	return json.Marshal(x)
}

type database struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
}

func (db *database) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", db.host, db.port, db.user, db.dbName, db.password)
}

func main() {
	conn := &database{
		host:     "localhost",
		port:     5432,
		user:     "userapp",
		password: "password",
		dbName:   "postgres",
	}

	// Opens the connection to the db
	db, err := gorm.Open("postgres", conn.String())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Migrates the DB models listed, go to the function
	DBMigrate(db)

	// passes the wired dependencies to the engine
	handlersUser := HandlersUserApi(db)
	engine := router(handlersUser)

	// It is a shortcut for http.ListenAndServe(addr, router)
	// Blank uses default :8080
	engine.Run()
}

// injectar deps? para hacer los handlers
func router(uAPI userAPI) *gin.Engine {
	// Defaults returns an Engine instance
	// with the Logger and Recovery middleware already attached.
	// gin.New() returns a new blank Engine instance without any middleware attached.
	r := gin.Default()

	// Adds middlewares
	r.Use(CORS(), Common())

	// Set the /login endpoint using the
	// userAPI.login function as handler
	r.POST("/login", uAPI.login)
	return r
}

/*------Middlewares----*/

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pre-flght
		switch c.Request.Method {
		case "OPTIONS":
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE")
			c.Header("Access-Control-Max-Age", "86400")
			c.Status(http.StatusOK)
		}

		c.Next()
	}
}

func Common() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CommonHeaders
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")

		c.Next()
	}
}

/*------Middlewares----*/

/*------User API------*/
type userAPI struct {
	service userService
}

func newUserAPI(u userService) userAPI {
	return userAPI{service: u}
}

func (u *userAPI) login(c *gin.Context) {
	creds := &creds{}
	c.BindJSON(creds)

	// TODO: Add auth service? or handle the login in the
	// 		auth service, separate from the user handlers
	// Needs a way a way to check the user credentials
	// comparing the password hash

	c.JSON(http.StatusOK, fmt.Sprintf("Logged in %s", creds.Username))
}

/*------User API------*/

/*----User Service---*/
type userService struct {
	repo userRepository
}

func newUserService(r userRepository) userService {
	return userService{repo: r}
}

func (s *userService) getUser(id uint) {
	s.repo.getUser(id)
}

/*----User Service---*/

/*--User Repository--*/
type userRepository struct {
	db *gorm.DB
}

func newUserRepository(db *gorm.DB) userRepository {
	return userRepository{db: db}
}

func (s *userRepository) getUser(id uint) *user {
	u := &user{}
	s.db.First(&u, id)
	return u
}

// Models is a list of all the models to migrate
var Models = []interface{}{
	&user{},
}

// DBMigrate uses the `...` notation at
// the end to unpack the slice
func DBMigrate(db *gorm.DB) {
	db.AutoMigrate(Models...)
}

/*--User Repository--*/
