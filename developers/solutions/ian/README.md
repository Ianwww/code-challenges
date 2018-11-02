# ByteCubed Code challenge


## Solution overview
The solution attempts to seat all people at tables in a round robin fashion starting with the largest group size to the smallest group size. The tables are checked for space in order from laregest to smallest table.

## Running solution
1. download and configure golang 
2. download this repository in your go workspace
3. cd into the wp directory
4. go build -o wp
5. ./wp
6. to run test cases run `go test -v wp_test.go wp.go` (without backticks)