# watchdogGO
Just a watchdog 


Po stronie Dynamodb przetrzymuje 

id - primary key 
List of services 
NumOfSecCheck - sec between checking is service running 
NumOfSecWait - sec between try to run service   
NumOfAttempts - attempts to run service if it is dead 

TODO 
[] - Updating settings
[] - Unification aws connection
[] - More tests 
[] - Handling errors (e.g. settings with error)