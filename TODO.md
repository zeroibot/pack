# TODO

* Add assert package for testing - replace repetitive testing code so far 
  * Case2, Case3, Case4, etc. that will be used instead of defining testCase structs inside each test function 
  * assert.Eq, assert.NotEq, assert.EqList, assert.EqMap
* Secrets/Configs package: Load .env file
* Make ds.List comparable by default, then add list compare functions as List methods
* Create ds.CmpList for non-comparable types, and add a compare function (e.g. maps.Equal, slices.Equal)
* Cache struct type, accepts DB call function it will use for cache misses 
* Pipeline-able data fetch functions: db -> cache -> decorator / stripper 
* Experiment: EndIf channel ← error, that stops execution to reduce if err != nil
* Add an RDB testing at app initialization for all registered schemas/types to ensure that all column mapping is correct

## Test Scripts

### test.sh
```bash
go test -cover ./dyn
```

### testcover.sh
```bash
go test -coverprofile=coverage.out ./dyn
go tool cover -html=coverage.out
```