package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Command line flags
	var (
		dbHost  = flag.String("db_host", "localhost", "Postgres database host")
		dbPort  = flag.Int("db_port", 5432, "Postgres database port")
		dbUser  = flag.String("db_user", "postgres", "Postgres database user")
		dbPass  = flag.String("db_pass", "root123", "Postgres database password")
		dbName  = flag.String("db_name", "registry", "Postgres database name")
		baseDir = flag.String("base_dir", "basedir", "Base directory to scan, for example: /var/lib/registry/docker/registry/v2")
		dryRun  = flag.Bool("dry_run", false, "Whether to skip deleting files")
	)
	flag.Parse()

	// Validate command line arguments
	if *dbUser == "" || *dbName == "" || *baseDir == "" {
		log.Fatalf("db_user, db_name, and base_dir are required")
	}

	// Connect to database
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*dbHost, *dbPort, *dbUser, *dbPass, *dbName))
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer db.Close()

	// Query digest list
	digestMap := make(map[string]bool)
	tobeDeleted := make(map[string]bool)
	rows, err := db.Query("SELECT substr(digest, 8) FROM blob")
	if err != nil {
		log.Fatalf("failed to query database: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var digest string
		if err := rows.Scan(&digest); err != nil {
			log.Fatalf("failed to scan database row: %s", err)
		}
		digestMap[digest] = true
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("failed to iterate database rows: %s", err)
	}

	// avoid delete all files if there is no blob in database
	if len(digestMap) == 0 {
		log.Fatalf("no blob in database")
	}

	// Walk base directory and delete files
	var totalSize, deleteCnt int64
	blobSha256Dir := filepath.Join(*baseDir, "blobs", "sha256")
	if err := filepath.Walk(blobSha256Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == `data` {
			// Get digest from file path
			subPath := strings.TrimPrefix(path, blobSha256Dir)
			parts := strings.Split(subPath, string(os.PathSeparator))
			if len(parts) != 3 {
				log.Printf("invalid file path: %s", subPath)
				return nil
			}
			digest := parts[1]

			// Delete file if digest not in map
			if !digestMap[digest] {
				deleteCnt++
				size := info.Size()
				totalSize += size
				tobeDeleted[digest] = true
				if *dryRun {
					log.Printf("would delete %s (size: %d)", path, size)
				} else {
					if err := os.Remove(path); err != nil {
						log.Printf("failed to delete file: %s", err)
					}
					parent := filepath.Dir(path)
					if err := os.Remove(parent); err != nil {
						log.Printf("failed to delete parent directory: %s", err)
					}
				}
			}
		}
		return nil
	}); err != nil {
		log.Fatalf("failed to walk base directory: %s", err)
	}

	// Walk repositories directory and delete files
	repositoriesDir := filepath.Join(*baseDir, "repositories")
	if err := filepath.Walk(repositoriesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == `link` {
			// Get digest from file path
			subPath := strings.TrimPrefix(path, repositoriesDir)
			parts := strings.Split(subPath, string(os.PathSeparator))
			if len(parts) != 4 {
				log.Printf("invalid file path: %s", subPath)
				return nil
			}
			digest := parts[2]
			if tobeDeleted[digest] {
				if *dryRun {
					log.Printf("would delete link file %s", path)
				} else {
					log.Printf("would delete link file %s", path)
					os.Remove(path)
				}
			}
		}
		return nil
	}); err != nil {
		log.Fatalf("failed to walk repositories directory: %s", err)
	}

	// Print summary
	if *dryRun {
		log.Printf("total files to delete: %d, size %d", deleteCnt, totalSize)
	} else {
		log.Printf("deleted %d files, freed %d bytes", deleteCnt, totalSize)
	}
}
