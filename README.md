## table service

# RUN


### 1. init mongoDb

create mongo db user

```
1. login your mongo server as admin
mongo --port 27017 -u  mongoadmin  -p secret  --authenticationDatabase admin 

2. db.createUser({
     "user": "admin",
   	 "pwd": "admin",
     "roles": ["readWrite"],
     "mechanisms": ["SCRAM-SHA-1"]
   })
``` 
