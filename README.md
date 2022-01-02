# watchdogGO
Watchdog written in golang, using dynamodb and sns AWS services.

# Run 
You need have correct settings and config in .aws file in home directory of profile that is running script. 

``` 
go run . -id <num of settings id from dynamodb>

OR

go run main.go -id <num of settings id from dynamodb>
```
