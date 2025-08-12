package entity

type File struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	FileSize      int64      `json:"file_size"`
	TosURI        string     `json:"tos_uri"`
	Status        FileStatus `json:"status"`
	Comment       string     `json:"comment"`
	Source        FileSource `json:"source"`
	CreatorID     string     `json:"creator_id"`
	CozeAccountID int64      `json:"coze_account_id"`
	ContentType   string     `json:"content_type"`
	CreatedAt     int64      `json:"created_at"`
	UpdatedAt     int64      `json:"updated_at"`
}

type FileStatus int32

const (
	FileStatusInvalid FileStatus = 0
	FileStatusValid   FileStatus = 1
)

type FileSource int32

const (
	FileSourceAPI FileSource = 1
)
