package cmd

func getUserHelpTemplate() string {
	return `
Usage: cephmgr rgw user get <UID> [FLAGS]
Example: cepfmgr rgw user get user1

Flags:
  -h, --help    help for get
  -j, --json    Return values as json
  -c, --cluster Cluster override
`
}

func listUsersTemplate() string {
	return `
Usage: cephmgr rgw user list [FLAGS]
Example: cepfmgr rgw user list

Flags:
  -h, --help    help for list
  -j, --json    Return values as json
  -c, --cluster Cluster override
`
}
func userCreateTemplate() string {
	return `
Usage: cephmgr rgw user create [FLAGS]
Example: cepfmgr rgw user create -u=user1 -f="User One" --caps "buckets=read" 
  * Required flags
	* -u -user ID
	* -f -user display name
  * Optional flags
	* -e "email" - email
	* --caps "buckets=read" - add single user capability
	* --caps "buckets=*;users=read;zone=*" - add multiple user capabilities
	
Flags:
  -h, --help                help for create
  -u, --user string         Ceph user ID (required)
  -f, --fullname string     Ceph user full name (required)
  -e, --email string        Ceph user e-mail
	  --caps string         User capabilities
	  --config string       Config file (default is $HOME/.cephmgr.yaml)
  -j, --json                Return values as json
  -c, --cluster Cluster override
`
}

func userDeleteTemplate() string {
	return `
Usage: cephmgr rgw user delete <UID> [FLAGS]
Example: cepfmgr rgw user delete user1
	
Flags:
  -h, --help                help for create
  -j, --json                Return values as json
  -c, --cluster Cluster override
`
}
func userGetCapsTemplate() string {
	return `
Usage: cephmgr rgw user caps get <UID> [FLAGS]
Example: cepfmgr rgw user caps get user1

Flags:
  -h, --help    help for get
  -j, --json    Return values as json
  -c, --cluster Cluster override
`
}

func userAddCapsTemplate() string {
	return `
Usage: cephmgr rgw user caps add <UID> [FLAGS]
Example: cepfmgr rgw user caps add user1 --caps "buckets=*;users=read"

Flags:
  -h, --help    help for add
  -j, --json    Return values as json
  -c, --cluster Cluster override
`
}

func userRemoveCapsTemplate() string {
	return `
Usage: cephmgr rgw user caps remove <UID> [FLAGS]
Example: cepfmgr rgw user caps remove user1 --caps "buckets=*;users=read"

Flags:
  -h, --help    help for remove
  -j, --json    Return values as json
  -c, --cluster Cluster override
`
}
func userModifyTemplate() string {
	return `
Usage: cephmgr rgw user modify <UID> [FLAGS]
Example: cepfmgr rgw user modify user1 -f="New Name" -e="newemail@example.com"

Flags:
  -j, --json    Return values as json
  -h, --help                help for modify
  -f, --fullname string     Updated full name for the user
  -e, --email string        Updated email for the user
  -c, --cluster Cluster override
`
}
func userKeysTemplate() string {
	return `
Usage: cephmgr rgw user keys <subcommand> [FLAGS]
Available Subcommands:
  add <UID> - Add new key to user
  remove <UID> <AccessKey> - Remove user keys
  
Flags:
  -j, --json    Return values as json
  -h, --help    help for keys
  -c, --cluster Cluster override
`
}

func userAddKeyTemplate() string {
	return `
Usage: cephmgr rgw user keys add <UID> [FLAGS]
Example: cepfmgr rgw user keys add user1

Flags:
  -j, --json    Return values as json
  -h, --help    help for add
  -c, --cluster Cluster override
`
}

func userRemoveKeyTemplate() string {
	return `
Usage: cephmgr rgw user keys remove <UID> <AccessKey> [FLAGS]
Example: cepfmgr rgw user keys remove user1 ABC123XYZ

Flags:
  -j, --json    Return values as json
  -h, --help    help for remove
  -c, --cluster Cluster override
`
}
func userQuotaSetTemplate() string {
	return `
Usage: cephmgr rgw user quota set <UID> [FLAGS]
Example: cephmgr rgw user quota set user1 --max-objects=100 --max-size=1GB --enabled=true

Flags:
  --max-objects     Max Objects 
  --max-size    Max Size (e.g., 1GB)
  --enabled           Enable or disable quotas

  -j, --json    Return values as json
  -h, --help    help for set
  -c, --cluster Cluster override
`
}
func userQuotaGetTemplate() string {
	return `
Usage: cephmgr rgw user quota get <UID> [FLAGS]
Example: cephmgr rgw user quota get user1

Flags:
  -j, --json    Return values as json
  -h, --help    help for get
  -c, --cluster Cluster override
`
}
func bucketDeleteTemplate() string {
	return `
Usage: cephmgr rgw bucket delete <BUCKET_NAME> [flags]
Example: cephmgr rgw bucket delete bucket1 --populated

Flags:
  --populated   Delete populated buckets
  -j, --json    Return values as json
  -h, --help    help for delete
  -c, --cluster Cluster override
`
}
func bucketListTemplate() string {
	return `
Usage: cephmgr rgw bucket list [FLAGS]
Example: cephmgr rgw bucket list

Flags:
  -j, --json    Return values as json
  -h, --help    help for list
  -c, --cluster Cluster override
`
}

func bucketInfoTemplate() string {
	return `
Usage: cephmgr rgw bucket info <BUCKET_NAME> [FLAGS]
Example: cephmgr rgw bucket info bucket1

Flags:
  -u, --usage   Display bucket usage information
  -q, --quota   Display bucket quota information
  -j, --json    Return values as json
  -h, --help    help for info
  -c, --cluster Cluster override
`
}
func bucketQuotaGetTemplate() string {
	return `
Usage: cephmgr rgw bucket quota get <BUCKET_NAME> [FLAGS]
Example: cephmgr rgw bucket quota get bucket1

Flags:
  -j, --json    Return values as json
  -h, --help    help for get
  -c, --cluster Cluster override
`
}

func bucketQuotaSetTemplate() string {
	return `
Usage: cephmgr rgw bucket quota set <UID> <BUCKET_NAME> [FLAGS]
Example: cephmgr rgw bucket quota set user1 bucket1 --max-objects=1000 --max-size=1GB --enabled

Flags:
  --max-objects=<int>   Max Objects Quota
  --max-size=<size>     Max Size Quota (e.g., 1GB)
  --enabled             Enable or disable quotas
  -j, --json            Return values as json
  -h, --help            help for set
  -c, --cluster         Cluster override
`
}

func clusterListTemplate() string {
	return `
Usage: cephmgr rgw cluster list [FLAGS]
Example: cephmgr rgw cluster list 

Flags:
  -j, --json    Return values as json
  -h, --help    help for list
`
}

func clusterGetActiveTemplate() string {
	return `
Usage: cephmgr rgw cluster get_active [FLAGS]
Example: cephmgr rgw cluster get_active --json

Flags:
  -j, --json    Return values as json
  -h, --help    help for get active
`
}

func clusterSetActiveTemplate() string {
	return `
Usage: cephmgr rgw cluster set_active [FLAGS]
Example: cephmgr rgw cluster set_active -n=cluster1

Flags:
  -n, --name    Cluster name (required)
  -j, --json    Return values as json
  -h, --help    help for set active
`
}

func clusterAddNewTemplate() string {
	return `
Usage: cephmgr rgw cluster add [FLAGS]
Example: cephmgr rgw cluster add -n=cluster4 -k=A574AGHBCNA66 -s=FHJA67820FHASS73HHFA9 -e=https://cluster.url.here --json

Flags:
  -n, --name            Cluster name (required)
  -k, --access_key      Cluster access key (required)
  -s, --access_secret   Cluster access secret (required)
  -e, --endpoint_url    Cluster endpoint URL with scheme (required)
  -j, --json            Return values as json
  -h, --help            help for add
`
}

func clusterRemoveTemplate() string {
	return `
Usage: cephmgr rgw cluster remove [FLAGS]
Example: cephmgr rgw cluster remove -n=cluster2

Flags:
  -n, --name    Cluster name (required)
  -j, --json    Return values as json
  -h, --help    help for remove
`
}
