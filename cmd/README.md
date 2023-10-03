## User
### RGW
user info
* ```rgw user list``` - get use list
* ```rgw user get <UID>``` - get user info
user create, delete
* ```rgw user delete <UID> ``` - delete user
* ```rgw user create [FLAGS]```
    * Required flags
        * ```-u "uid"``` -user ID, 
        * ```-f "FullName"``` -user display name
    * Optional flags
        * ```-e "email"``` - email
        * ```--caps "buckets=read"``` - add single user capability
        * ```--caps "buckets=*;users=read;zone=*"``` - add multiple user capabilities
user info modify
* ```rgw user modify <UID> [FLAGS]```  - edit user fullname and email
    * ```-f "FullName"``` -user display name, optional
    * ```-e "email"``` - email, optional
    * Kas lisan siia max bucketid?
    user key managment
* ```rgw user keys delete <UID> <AccessKey>``` -delete user Keys
* ```rgw user keys add <UID>``` - add new keys to user
user caps
* ```rgw user caps add <uid> [flag]```
    * ```--caps "buckets=read"``` - add single user capability
    * ```--caps "buckets=*;users=read;zone=*"``` - add multiple user 
* ```rgw user caps remove <uid> [flag]```
    * ```--caps "buckets=read"``` - remove single user capability
    * ```--caps "buckets=*;users=read;zone=*"``` - remove multiple user 
user quota
* ```rgw user quota get <UID>``` - get user quotas
* ```rgw user quota set <UID> [flag]``` - get user quotas
    * ```--max-objects=<int>``` user quota max objects
    * ```--max-size=<int>```  user quota max size in bites
    * ```--max-size-kb=<int>``` user quota max size in Kb
    * ```--enabled=<bool>``` enable/disable user quotas
## Bucket
### RGW
Problems:
* rgw ei lase bucketit lisada
* bucketi muutmine?


* ```rgw bucket list``` -get the Bucket list
* ```rgw bucket info <Bucket> [flag]``` get bucket info
    * ```--quota or -q``` get bucket quotas
    * ```--usage or -u``` get bucket usage
* ```rgw bucket delete <Bucket> [flags]```- Deletes only empty buckets
    * ```--populated``` deletes populated bucket
* ```rgw bucket quota get <bucket>``` get bucket quota
* ```rgw bucket quota set <UID> <Bucket> [flags]```
    * ```--max-objects=<int>``` bucket quota max objects
    * ```--max-size=<int>```  bucket quota max size in bites
    * ```--max-size-kb=<int>``` bucket quota max size in Kb
    * ```--enabled=<bool>``` enable/disable bucket quotas