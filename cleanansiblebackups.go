package main

////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// CLEANANSIBLEBACKUPS
// ===================
// Small utlility written in Go to cleanup Ansible backup files created over time when using
// "backup: yes" when copying files.
//
// -----------------------------------------------------------------------------------------------------
// v1.0  MJO 2022/01/25 Initial release
//
// IDEAS
// =====
// - To move the backup files to a different directory (test existance and accessibility)
// 
////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
        "flag"
        "fmt"
        "os"
        "path/filepath"
        "time"

        "github.com/itroot/keysort"
)

const (
        myVersion = "v1.0(5)"
)

type fileWithAge struct {
        name     string
        ageHours int64
        ageDays  int64
}

func main() {
        var path string
        var mincount int
        var age int

        var myFiles []fileWithAge
        var aFile fileWithAge

        fmt.Printf("cleanansiblebackups %s\n\n", myVersion)

        // Define the flags we are accepting
        flag.StringVar(&path, "path", "", "path and base filename where to clean up the Ansible backups")
        flag.IntVar(&mincount, "mincount", 0, "minimum number of backups to preserve regardless of age")
        flag.IntVar(&age, "age", 0, "number of days to keep a backup")
        flag.Parse()

        // Check if we have all mandatory flags
        if len(path) == 0 || mincount == 0 || age == 0 {
                fmt.Println("Usage: cleanansiblebackups -path <path and base filename> -mincount <x> -age <y>")
                flag.PrintDefaults()
                os.Exit(1)
        }

        // Check if the given path and baseline filename exists
        _, err := os.Stat(path)

        if err != nil {
                fmt.Println("Cannot locate the given path and base filename?")
                flag.PrintDefaults()
                os.Exit(2)
        }

        // Prepare a search pattern
        search := fmt.Sprintf("%s.*~", path)

        fmt.Printf("Scanning: %s\n", search)

        matches, err := filepath.Glob(search)

        if err != nil {
                fmt.Println("Error while scanning for given path and base filename?")
                fmt.Println(err)
                os.Exit(3)
        }

        // Fetch file information for each file located
        for _, match := range matches {
                file, err := os.Stat(match)

                if err != nil {
                        fmt.Printf("Could not stat file - skipping: %s\n", match)
                } else {
                        // Store the information we've found
                        now := time.Now()
                        diff := now.Sub(file.ModTime())

                        aFile.name = match
                        aFile.ageHours = int64(diff.Hours())
                        aFile.ageDays = int64(aFile.ageHours / 24)

                        myFiles = append(myFiles, aFile)
                }
        }

        fmt.Printf("Matches found: %d\n", len(myFiles))

        // Sort the information by age
        keysort.Sort(myFiles, func(i int) keysort.Sortable {
                file := myFiles[i]
                return keysort.Sequence{file.ageHours}
        })

        // Determine which files to delete based on the provisioned flag information
        for i, output := range myFiles {
                // Match against conditions
                if output.ageDays < int64(age) {
                        continue
                }
                if i < mincount {
                        continue
                }

                // Report that we remove this file
                fmt.Printf("Removing %s, age %d days\n", output.name, output.ageDays)

                err := os.Remove(output.name)

                if err != nil {
                        fmt.Printf("! Failed to remove this file?\n%s\n", err)
                }
        }
}
