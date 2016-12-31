// Main HTTP server package for Squad Up.
package main

// Import deps.
import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"

    "github.com/Noah-Huppert/squad-up/server/models"
    tables "github.com/Noah-Huppert/squad-up/server/models/db"
	"github.com/Noah-Huppert/squad-up/server/handlers"
)

// Main entry point of program.
func main() {
    // Connect to DB
    db, err := gorm.Open("postgres", "host=localhost user=username password=password dbname=squad-up sslmode=disable")
    defer db.Close()
    if err != nil {
        fmt.Println("Error connecting to database: " + err.Error())
        return
    } else {
        fmt.Println("Connected to database on :5432")
    }

    // Setup DB
    db.AutoMigrate(&tables.User{})

    // Create App Context
    config := models.Config{
        GAPIClientId: "432144215744-2n6fha955i4f2en9jubvelfhmdsh1jcv.apps.googleusercontent.com",
        JWTServerURI: "squad-up@server/api/v1",
        JWTHMACKey: "abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrst",
    }

    ctx := models.AppContext{config, db}

	// New HTTP router.
	mux := http.NewServeMux()

	// Attach handlers
    handlerLoader := handlers.NewLoader(mux, &ctx)
    handlerLoader.Load()

	// Start listening on any host, port 5000.
	fmt.Println("Listening on :5000")

	err = http.ListenAndServe(":5000", mux)
	if err != nil { // If err print to console.
		fmt.Println("Error starting HTTP server on :5000: " + err.Error())
	}
}
