# watchdogGO
Watchdog written in golang, using dynamodb and sns AWS services.

# Setup 

Enviromental variables for aws
```
export AWS_ACCESS_KEY_ID=<YOUR_AWS_ACCESS_KEY>
export AWS_SECRET_ACCESS_KEY=<YOUR_AWS_SECRET>
export AWS_REGION=<YOUR_AWS_REGION>

```
or put these above in .aws folder into config/credentials file

--

Enviromental variables for script 

```
export SNSWATCHDOG=arn:aws:sns:region-2:444455556666:MyTopic
export TABLEWATCHDOG=dynamoDB table name 
```

# Run  


``` 
go run . -id <num of settings id from dynamodb>

OR

go run main.go -id <num of settings id from dynamodb>

OR 

go build && ./watchdogGO
```

log files will be putted into file watchdog.log and if you subscribe for sns topic you can get notifications as well.

# Stop 
To stop proces you need to send SIGTERM to process. 
