# TODO

* qb package: QueryBuilder - create Instance that will be used instead of global package variables
* Cache struct type, accepts DB call function it will use for cache misses 
* Query condition is usable for both DB queries and Cache queries 
* Pipeline-able data fetch functions: db -> cache -> decorator / stripper 
* Experiment: EndIf channel ← error, that stops execution to reduce if err != nil
