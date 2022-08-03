package main

var tableCreateQuery = `
	CREATE TABLE IF NOT EXISTS Data (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    input TEXT NOT NULL,
	    date TEXT NOT NULL
	);
`

var insertDataQuery = `
INSERT INTO Data (input, date) VALUES (?, ?)
`
