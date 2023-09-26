## User
### RGW
* ``rgw user list`` - get use list
* ``rgw user get <UID>`` - get user info
* ``rgw user delete <UID> `` - delete user
* ``rgw user create [FLAGS]``
    * Required flags
        * ``-u "uid"`` -user ID, 
        * ``-f "FullName"`` -user display name
    * Optional flags
        * ``-e "email"`` - email
        * ``--caps "buckets=read"`` - add single user capability
        * ``--caps "buckets=*;users=read;zone=*"`` - add multiple user capabilities
* ``rgw user modify <UID> [FLAGS]``  - edit user fullname and email
    * ``-f "FullName"`` -user display name, optional
    * ``-e "email"`` - email, optional
* ``rgw user keys delete <UID> <AccessKey>`` -delete user Keys
* ``rgw user keys add <UID>`` - add new keys to user
* ``rgw user caps add <uid> [flag]``
    * ``--caps "buckets=read"`` - add single user capability
    * ``--caps "buckets=*;users=read;zone=*"`` - add multiple user 

## Bucket
### RGW