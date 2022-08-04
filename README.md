### Introduction into testing
Litte golang api server with broken architecture and no tests.
It'll become better over iterations.

### How to run
```
go run cmd/main.go --conf=conf.yaml
```

#### Chapters
- v0.0.1: added some unit tests, fixed bug in GET /orders and decouple pool&repo from usecase.
- v0.0.2: added intergration tests for gateway-usecase layers. Also added tests with fakes.
