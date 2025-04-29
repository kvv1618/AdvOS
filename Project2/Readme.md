# Instructions to build/run
- To start file server, run `go run file_server.go <file_path> <config_file_path>`
- To start main program, run `go run main.go <N> <C> <file_path> <config_file_path>`
- To start workers, run `go run worker.go <M> <config_file_path>`

_Note <file_path> and <config_file_path> are mandatory arguments. M, N and C are optional arguments. If not provided, default values will be used._
