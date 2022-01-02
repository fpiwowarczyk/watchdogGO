# watchdogGO
Watchdog written in golang, using dynamodb and sns AWS services.

# Run 
You need have correct settings and config in .aws file in home directory of profile that is running script. 

Also you need correct config.json file with 

``` json 
{
    "tables":{
        "watchdog" : "table-name"
    },
    "sns":{
        "watchdog":"arn:aws:sns:us-east-2:444455556666:MyTopic"
    }
}

```


``` 
go run . -id <num of settings id from dynamodb>

OR

go run main.go -id <num of settings id from dynamodb>
```

To stop proces you need to send SIGTERM to process. 
