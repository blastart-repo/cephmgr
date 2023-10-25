## User

flag ```--json``` short ```-j``` prints all responses in json
### RGW
user info
* ```cephmgr rgw user list``` - get use list
* ```cephmgr rgw user get <UID>``` - get user info
user create, delete
* ```cephmgr rgw user delete <UID> ``` - delete user
* ```cephmgr rgw user create [FLAGS]```
    * Required flags
        * ```-u "uid"``` -user ID, 
        * ```-f "FullName"``` -user display name
    * Optional flags
        * ```-e "email"``` - email
        * ```--caps "buckets=read"``` - add single user capability
        * ```--caps "buckets=*;users=read;zone=*"``` - add multiple user capabilities
user modify
* ```cephmgr rgw user modify <UID> [FLAGS]```  - edit user fullname and email
    * ```-f "FullName"``` -user display name, optional
    * ```-e "email"``` - email, optional
    * Kas lisan siia max bucketid?
user key managment
* ```cephmgr rgw user keys delete <UID> <AccessKey>``` -delete user Keys
* ```cephmgr rgw user keys add <UID>``` - add new keys to user
user caps
* ```cephmgr rgw user caps get <UID>``` - get user caps
* ```cephmgr rgw user caps add <uid> [flag]```
    * ```--caps "buckets=read"``` - add single user capability
    * ```--caps "buckets=*;users=read;zone=*"``` - add multiple user 
* ```cephmgr rgw user caps remove <uid> [flag]```
    * ```--caps "buckets=read"``` - remove single user capability
    * ```--caps "buckets=*;users=read;zone=*"``` - remove multiple user 
user quota
* ```cephmgr rgw user quota get <UID>``` - get user quotas
* ```cephmgr rgw user quota set <UID> [flag]``` - get user quotas
    * ```--max-objects=<int>``` user quota max objects
    * ```--max-size=<string>```  user quota max size in bites
        * example ```--max-size=2gb|2tb|459mb``` -units not case sensitive
    * ```--enabled=<bool>``` enable/disable user quotas
## Bucket
### RGW



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


## Cluster
### RGW


* ```rgw cluster add``` Add new cluster
  * ```--name``` or ```-n```
  * ```--access_key``` or ```-k```
  * ```--access_secret``` or ```-s```
  * ```--endpoint_url``` or ```-e```
  * ```--sensitive``` shows access_key and access_secret fields
* ```rgw cluster get_active``` Get default active cluster info
  * ```--sensitive``` shows access_key and access_secret fields
* ```rgw cluster list``` Get a list of clusters
  * ```--sensitive``` shows access_key and access_secret fields
* ```rgw cluster remove``` Removes cluster
  * ```--name``` or ```-n```
  * ```--sensitive``` shows access_key and access_secret fields
* ```rgw cluster set_active``` Set default active cluster
  * ```--name``` or ```-n```
  * ```--sensitive``` shows access_key and access_secret fields