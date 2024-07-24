package migrations

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func Generate(filename string) {
	if filename == "" {
		fmt.Println("please provide a filename for the migration file.")
		return
	}

	path := "./backend/migrations"
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("migrations directory", path, "does not exist")
			return
		}
	}

	now := time.Now()
	formattedTime := now.Format("060102150405")
	file := fmt.Sprintf("%s/%s_%s.go", path, formattedTime, filename)

	if err := createFile(file, content(formattedTime)); err != nil {
		fmt.Println(err)
		return
	}
}

func createFile(filename, content string) error {
	if _, err := os.Stat(filename); err == nil {
		return errors.New("file already exists")
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func content(fileid string) string {
	template := `package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "%s",
		Migrate: func(tx *gorm.DB) error {
			return nil
		},
	}
	MigrationHelper.Register(m)
}`
	return fmt.Sprintf(template, fileid)
}
