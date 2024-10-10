# Build your own BlockChain

base on https://github.com/tensor-programming/golang-blockchain/tree/master. 

I think my version has better readability

# How to run
```
go run main.go
```
\
By default, for a better cross-platform experience, 
it will use redis as the blockchain storage, 
if you want to use badgerDB(better run it in Linux), you can change the option in main.go:
```
db := NewStorage(optionBadgerDB)
```