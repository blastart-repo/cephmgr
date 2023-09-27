## User
### RGW
user info
* ``rgw user list`` - get use list
* ``rgw user get <UID>`` - get user info
user create, delete
* ``rgw user delete <UID> `` - delete user
* ``rgw user create [FLAGS]``
    * Required flags
        * ``-u "uid"`` -user ID, 
        * ``-f "FullName"`` -user display name
    * Optional flags
        * ``-e "email"`` - email
        * ``--caps "buckets=read"`` - add single user capability
        * ``--caps "buckets=*;users=read;zone=*"`` - add multiple user capabilities
user info modify
* ``rgw user modify <UID> [FLAGS]``  - edit user fullname and email
    * ``-f "FullName"`` -user display name, optional
    * ``-e "email"`` - email, optional
    user key managment
* ``rgw user keys delete <UID> <AccessKey>`` -delete user Keys
* ``rgw user keys add <UID>`` - add new keys to user
user caps
* ``rgw user caps add <uid> [flag]``
    * ``--caps "buckets=read"`` - add single user capability
    * ``--caps "buckets=*;users=read;zone=*"`` - add multiple user 
* ``rgw user caps remove <uid> [flag]``
    * ``--caps "buckets=read"`` - remove single user capability
    * ``--caps "buckets=*;users=read;zone=*"`` - remove multiple user 

## Bucket
### RGW
Problems:
* rgw ei lase bucketit lisada
* bucketi muutmine?


* ``rgw bucket list`` -get the Bucket list
* ``rgw bucket info <Bucket>`` get bucket info
* ``rgw bucket delete <Bucket> [flags]``- Deletes only empty buckets
    * ``--populated`` deletes populated bucket