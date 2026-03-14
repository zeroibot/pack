# TODO

* Make ds.List comparable by default, then add list compare functions as List methods
* Create ds.CmpList for non-comparable types, and add a compare function (e.g. maps.Equal, slices.Equal)
* qb package: QueryBuilder - create Instance that will be used instead of global package variables
* Cache struct type, accepts DB call function it will use for cache misses 
* Query condition is usable for both DB queries and Cache queries 
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