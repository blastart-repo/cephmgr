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
	Usage       struct {
		RgwMain struct {
			Size           *uint64 `json:"size"`
			SizeActual     *uint64 `json:"size_actual"`
			SizeUtilized   *uint64 `json:"size_utilized"`
			SizeKb         *uint64 `json:"size_kb"`
			SizeKbActual   *uint64 `json:"size_kb_actual"`
			SizeKbUtilized *uint64 `json:"size_kb_utilized"`
			NumObjects     *uint64 `json:"num_objects"`
		} `json:"rgw.main"`
	}
}

type QuotaSpec struct {
	UID        string `json:"user_id" url:"uid"`
	Bucket     string `json:"bucket" url:"bucket"`
	Enabled    *bool  `json:"enabled" url:"enabled"`
	MaxSize    *int64 `json:"max_size" url:"max-size"`
	MaxSizeKb  *int   `json:"max_size_kb" url:"max-size-kb"`
	MaxObjects *int64 `json:"max_objects" url:"max-objects"`
}

type ResponseQuota struct {
	UID        string `json:"user_id" url:"uid"`
	Bucket     string `json:"bucket" url:"bucket"`
	Enabled    *bool  `json:"enabled" url:"enabled"`
	MaxSize    string `json:"max_size" url:"max-size"`
	MaxObjects *int64 `json:"max_objects" url:"max-objects"`
}

type CLIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
type StringSlice struct {
	Buckets []string `json:"bucket-list"`
}

type BucketInfo struct {
	ID     string `json:"id"`
	Bucket string `json:"bucket" url:"bucket"`
	Owner  string `json:"owner"`
}

type BucketInfoUsage struct {
	Bucket     string  `json:"bucket" url:"bucket"`
	Size       string  `json:"size"`
	NumObjects *uint64 `json:"num_objects"`
}
type UserCapsResponse struct {
	UID  string        `json:"user_id" url:"uid"`
	Caps []UserCapSpec `json:"caps"`
}
