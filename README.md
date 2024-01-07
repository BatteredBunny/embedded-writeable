# embedded-writeable
Example programs that modify itself while being ran

# [Simple binary modification](https://github.com/BatteredBunny/embedded-writeable/tree/master/simple)
[![asciicast](https://asciinema.org/a/YdSUIE0m5ROlu5GnZSRdfBdR2.svg)](https://asciinema.org/a/YdSUIE0m5ROlu5GnZSRdfBdR2)
Simply adds the data at the end of program binary seperated by a signifier

# Compiling examples
Both of the below programs save data in the binary by recompiling itself, meaning you need [go installed on your system](https://go.dev/doc/install)
## [Simple compile example](https://github.com/BatteredBunny/embedded-writeable/tree/master/compile/simple)
### Program that shows the last time you opened it
Recompiles itself with new information as the program closes

## [Sqlite](https://github.com/BatteredBunny/embedded-writeable/tree/master/compile/sqlite)
### Sqlite database with simple REPL you can play around with
Stores the database in virtual FS and moves it out of that into a temp folder
