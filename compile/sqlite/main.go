package main

import (
	"bufio"
	"database/sql"
	"embed"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

//go:embed *
var fsdir embed.FS

const DatabaseName = "nya.db"

type State struct {
	Running bool   // if main loops is running
	tempdir string // dir to store db and source code in
	db      *sql.DB
}

func main() {
	var err error
	state := State{
		Running: true,
	}

	state.tempdir, err = os.MkdirTemp("", "*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(state.tempdir)

	CopyEmbedToTemp(fsdir, state.tempdir)

	// Setup db and its table
	state.db, err = sql.Open("sqlite", state.tempdir+"/"+DatabaseName)
	if err != nil {
		log.Panic(err)
	}
	defer state.db.Close()
	if _, err = state.db.Exec(tableCreateQuery); err != nil {
		log.Panic(err)
	}

	// Runs save on force quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		Save(state.tempdir)
		os.Exit(0)
	}()

	fmt.Println("Starting! /exit to quit")

	// Main loop
	scanner := bufio.NewScanner(os.Stdin)
	for state.Running {
		fmt.Print("\033[107;30mPROMPT:\033[0m ") // Prompt start
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		if strings.HasPrefix(input, "/") { // Commands
			f, ok := commands[strings.TrimPrefix(input, "/")]
			if ok {
				f(&state)
			} else {
				fmt.Println("Unknown command.")
			}
		} else {
			var rows *sql.Rows
			rows, err = state.db.Query(input)
			if err != nil {
				PrintError(err)
				continue
			}

			var columnTypes []*sql.ColumnType
			columnTypes, err = rows.ColumnTypes()
			if err != nil {
				PrintError(err)
				continue
			}

			// Pretty prints column names
			if len(columnTypes) > 0 {
				fmt.Print("| ")
				for _, v := range columnTypes {
					fmt.Print(v.Name(), " | ")
				}
				fmt.Println()
			}

			for rows.Next() {
				var columns []string
				columns, err = rows.Columns()
				if err != nil {
					PrintError(err)
					return
				}

				// All the column data should get converted to string! :D
				cols := make([]any, len(columns))
				container := make([]string, len(cols))
				for i := range cols {
					cols[i] = &container[i]
				}

				if err = rows.Scan(cols...); err != nil {
					PrintError(err)
					return
				}

				// Pretty prints rows
				if len(cols) > 0 {
					fmt.Print("| ")
					for _, v := range cols {
						fmt.Print(*v.(*string), " | ")
					}
					fmt.Println()
				}
			}
		}
	}

	Save(state.tempdir)
}

func Save(tempdir string) {
	command := fmt.Sprintf(
		"go build -modfile %s/go.mod -o %s %s/*.go",
		tempdir,
		os.Args[0], // Replaces currently ran binary with the new one
		tempdir,
	)

	if err := exec.Command("sh", "-c", command).Run(); err != nil {
		PrintError(err)
	}
}
