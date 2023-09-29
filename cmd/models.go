package cmd

type User struct {
	ID          string        `json:"user_id" url:"uid"`
	DisplayName string        `json:"display_name" url:"display-name"`
	Email       string        `json:"email" url:"email"`
	Keys        []UserKeySpec `json:"keys"`
	Caps        []UserCapSpec `json:"caps"`
	UserCaps    string        `url:"user-caps"`
}

type UserKeySpec struct {
	User      string `json:"user"`
	AccessKey string `json:"access_key" url:"access-key"`
	SecretKey string `json:"secret_key" url:"secret-key"`
	// Request fields
	UID         string `url:"uid"`          // The user ID to receive the new key
	KeyType     string `url:"key-type"`     // s3 or swift
	GenerateKey *bool  `url:"generate-key"` // Generate a new key pair and add to the existing keyring
}

type UserCapSpec struct {
	Type string `json:"type"`
	Perm string `json:"perm"`
}

type Bucket struct {
	ID          string    `json:"id"`
	Bucket      string    `json:"bucket" url:"bucket"`
	Owner       string    `json:"owner"`
	BucketQuota QuotaSpec `json:"bucket_quota"`
}

type QuotaSpec struct {
	UID        string `json:"user_id" url:"uid"`
	Bucket     string `json:"bucket" url:"bucket"`
	QuotaType  string `url:"quota-type"`
	Enabled    *bool  `json:"enabled" url:"enabled"`
	MaxSize    *int64 `json:"max_size" url:"max-size"`
	MaxSizeKb  *int   `json:"max_size_kb" url:"max-size-kb"`
	MaxObjects *int64 `json:"max_objects" url:"max-objects"`
}
