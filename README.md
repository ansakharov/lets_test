### Introduction into testing
Litte golang api server with broken architecture and no tests.
It'll become better over iterations.

### How to run
```
go run cmd/main.go --conf=conf.yaml

```

#### Chapters
iter0: base api routes with some bugs (not in codebase) 
iter1: added some unit tests, fixed bug in GET /orders and decouple pool&repo from usecase
