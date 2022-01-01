# watchdogGO
Watchdog written in golang, using dynamodb and sns AWS services.

# Run 
``` 
go run . -id <num of settings id from dynamodb>

OR

go run main.go -id <num of settings id from dynamodb>
```

# TODO

- [ ] Add tests
- [ ] Clean code
- [ ] Re run watchdogs only if settings differ
- [ ] Update README