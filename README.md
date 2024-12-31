# Build your own BlockChain

base on https://github.com/tensor-programming/golang-blockchain/tree/master. 

I think my version has better readability

# How to run
1. Currently, the "transaction version" have some bugs, please use the commit `feat: cmd`
by running `git checkout f9e6309`
2. start Redis server(though it's not a ACID database, it's good enough for this demo project)
3. `go run main.go`

If you want to use badgerDB, you can change the option in main.go:
```
db := NewStorage(optionBadgerDB)
```