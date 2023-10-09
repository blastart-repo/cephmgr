package cmd

func getUserHelpTemplate() string {
	return `
Usage: cephmgr rgw user get <UID> [FLAGS]
Example: cepfmgr rgw user get user1

Flags:
  -h, --help    help for get
  -j, --json    Return values as json
`
}

func listUsersHelpTemplate() string {
	return `
Usage: cephmgr rgw user list [FLAGS]
Example: cepfmgr rgw user list

Flags:
  -h, --help    help for list
  -j, --json    Return values as json
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
	  
`
}

func userDeleteTemplate() string {
	return `
	Usage: cephmgr rgw user delete <UID> [FLAGS]
	Example: cepfmgr rgw user delete user1
	 	
	Flags:
	  -h, --help                help for create
	  -j, --json                Return values as json

`
}
func userGetCapsHelpTemplate() string {
	return `
Usage: cephmgr rgw user caps get <UID> [FLAGS]
Example: cepfmgr rgw user caps get user1

Flags:
  -h, --help    help for get
  -j, --json    Return values as json
`
}

func userAddCapsHelpTemplate() string {
	return `
Usage: cephmgr rgw user caps add <UID> [FLAGS]
Example: cepfmgr rgw user caps add user1 --caps "buckets=*;users=read"

Flags:
  -h, --help    help for add
  -j, --json    Return values as json
`
}

func userRemoveCapsHelpTemplate() string {
	return `
Usage: cephmgr rgw user caps remove <UID> [FLAGS]
Example: cepfmgr rgw user caps remove user1 --caps "buckets=*;users=read"

Flags:
  -h, --help    help for remove
  -j, --json    Return values as json
`
}
func userModifyTemplate() string {
	return `
	Usage: cephmgr rgw user modify <UID> [FLAGS]
	Example: cepfmgr rgw user modify user1 -f="New Name" -e="newemail@example.com"
	
	Flags:
	  -h, --help                help for modify
	  -f, --fullname string     Updated full name for the user
	  -e, --email string        Updated email for the user
`
}
