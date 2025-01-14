// autoMigrate.go needs to be executed only when it is required
// to migrate a new relational database schema or to delete
// an existing database
package main

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/pilinux/gorest/config"
	"github.com/pilinux/gorest/database"
	"github.com/pilinux/gorest/database/model"
)

// Load all the models
type auth model.Auth
type twoFA model.TwoFA
type user model.User
type post model.Post
type hobby model.Hobby
type userHobby model.UserHobby

var db *gorm.DB
var errorState int

func main() {
	config.Config()
	configureDB := config.GetConfig().Database.RDBMS
	driver := configureDB.Env.Driver

	/*
	** 0 = default/no error
	** 1 = error
	**/
	errorState = 0

	db = database.InitDB()

	// Auto migration
	/*
		- Automatically migrate schema
		- Only create tables with missing columns and missing indexes
		- Will not change/delete any existing columns and their types
	*/

	// Careful! It will drop all the tables!
	dropAllTables()

	// Automatically migrate all the tables
	migrateTables()

	// Manually set foreign keys for MySQL and PostgreSQL
	// Disable it if foreign keys are already set in the database
	if driver != "sqlite3" {
		setPkFk()
	}

	if errorState == 0 {
		fmt.Println("Auto migration is completed!")
	} else {
		fmt.Println("Auto migration failed!")
	}
}

func dropAllTables() {
	// Careful! It will drop all the tables!
	if err := db.Migrator().DropTable(
		&userHobby{},
		&hobby{},
		&post{},
		&user{},
		&twoFA{},
		&auth{},
	); err != nil {
		errorState = 1
		fmt.Println(err)
	} else {
		fmt.Println("Old tables are deleted!")
	}
}

func migrateTables() {
	configureDB := config.GetConfig().Database.RDBMS
	driver := configureDB.Env.Driver

	if driver == "mysql" {
		// db.Set() --> add table suffix during auto migration
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
			&auth{},
			&twoFA{},
			&user{},
			&post{},
			&hobby{},
		); err != nil {
			errorState = 1
			fmt.Println(err)
		} else {
			fmt.Println("New tables are  migrated successfully!")
		}
	} else {
		if err := db.AutoMigrate(&auth{},
			&user{}, &post{}, &hobby{}); err != nil {
			errorState = 1
			fmt.Println(err)
		} else {
			fmt.Println("New tables are  migrated successfully!")
		}
	}
}

func setPkFk() {
	// Manually set foreign key for MySQL and PostgreSQL
	if err := db.Migrator().CreateConstraint(&user{}, "Posts"); err != nil {
		errorState = 1
		fmt.Println(err)
	}
}
