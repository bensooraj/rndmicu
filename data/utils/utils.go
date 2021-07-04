package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bensooraj/rndmicu/data/database"
	"github.com/bensooraj/rndmicu/data/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
}

func main() {
	var err error
	if len(os.Args) != 2 {
		log.Println("utils needs exactly one argument")
		return
	}

	dbCfg := database.Config{
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_PATH"),
	}
	log.Println("dbCfg: ", dbCfg)

	db, err := database.Open(dbCfg)
	if err != nil {
		log.Fatalln("Error establishing db connection:", err)
	}
	defer db.Close()

	switch os.Args[1] {
	case "create":
		err = ExecuteSqlFile(db, "data/schema/sql", "schema.sql")
		if err != nil {
			log.Println("Error creating the tables:", err)
			return
		}
	case "dropall":
		err = ExecuteSqlFile(db, "data/utils/sql", "drop.sql")
		if err != nil {
			log.Println("Error dropping the tables:", err)
			return
		}
	case "seed":
		// Delete all existing data
		// err = ExecuteSqlFile(db, "data/utils/sql", "delete_all.sql")
		// if err != nil {
		// 	log.Println("Error deleting date from the tables:", err)
		// 	return
		// }
		_, err = SeedCreators(db)
		if err != nil {
			log.Println("Error seeding the creator table:", err)
			return
		}
		err = SeedAudioShorts(db)
		if err != nil {
			log.Println("Error seeding the audio_shorts table:", err)
			return
		}

	default:
		log.Println("Command doesn't exist")

	}
}

// ExecuteSqlFile executes a given SQL file against a database.
// The queries are run in a transaction and rolled back if any fail.
func ExecuteSqlFile(db *sqlx.DB, dirName, sqlFilename string) error {
	// Load SQL
	path := filepath.Join(dirName, sqlFilename)
	sqlByteArray, ioErr := ioutil.ReadFile(path)
	if ioErr != nil {
		log.Println("Error reading the schema.sql file:", ioErr)
	}
	sql := string(sqlByteArray)
	log.Println("sql: ", sql)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(sql); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// SeedCreators creates creators using the generated models. The queries are run in a
// transaction and rolled back if any fail.
func SeedCreators(db *sqlx.DB) ([]models.Creator, error) {
	args := []models.CreateCreatorParams{}
	for _, i := range makeRange(1, 4) {
		args = append(args, models.CreateCreatorParams{
			Name:  fmt.Sprintf("ben_%d", i),
			Email: fmt.Sprintf("ben_%d@gmail.com", i),
		})
	}
	creators := []models.Creator{}

	q := models.New(db)
	for _, arg := range args {
		c, err := q.CreateCreator(context.Background(), arg)
		if err != nil {
			return nil, err
		}
		creators = append(creators, c)
		log.Println("[CREATOR]", c.ID, c.Name, c.Email)
	}

	return creators, nil
}

// SeedAudioShorts creates creators using the generated models. The queries are run in a
// transaction and rolled back if any fail.
func SeedAudioShorts(db *sqlx.DB) error {
	var err error
	args := []models.CreateAudioShortParams{}
	for _, i := range makeRange(1, 4) {
		args = append(args, models.CreateAudioShortParams{
			ID:           uuid.New(),
			CreatorID:    i,
			Title:        fmt.Sprintf("Title %d", i),
			Description:  fmt.Sprintf("Description %d", i),
			Category:     fmt.Sprintf("Category %d", i),
			AudioFileUrl: fmt.Sprintf("audio_file_url_%d.mp3", i),
			DateCreated:  time.Now(),
			DateUpdated:  time.Now(),
		})
	}

	q := models.New(db)
	for _, arg := range args {
		c, err := q.CreateAudioShort(context.Background(), arg)
		if err != nil {
			return err
		}
		log.Println("[AUDIO SHORT]", c.ID, c.Title, c.CreatorID, c.AudioFileUrl)
	}

	return err
}

// Helpers funcs
func makeRange(min, max int32) []int32 {
	a := make([]int32, max-min+1)
	for i := range a {
		a[i] = min + int32(i)
	}
	return a
}
