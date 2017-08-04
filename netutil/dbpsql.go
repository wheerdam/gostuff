package netutil

import (
	"fmt"
	"bufio"
	"errors"
	"os"
	"strings"
	"strconv"
	"database/sql"
	_ "github.com/lib/pq"
)

func OpenPostgresDBFromConfig(path string) (*sql.DB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("File does not exist")
	}
	
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())		
		return OpenPostgresDB(tokens)
	}
	return nil, errors.New("Invalid Database Configuration File Format")
}

func OpenPostgresDB(tokens []string) (*sql.DB, error) {
	if len(tokens) != 6 {
			return nil, errors.New("Invalid Database Configuration Format")
	}
	port, err := strconv.Atoi(tokens[3])
	if err != nil {
		return nil, err
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			tokens[2], port, tokens[0], tokens[1], tokens[4])
	db, err := sql.Open("postgres", psqlInfo)		
	if err != nil {
		return nil, err
	}
	return db, nil
}