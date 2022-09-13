package main

type ClipDto struct {
	Id            string        `json:"id"`
	SubId         string        `json:"subId"`
	FileSize      int64         `json:"fileSize"`
	CanDelete     bool          `json:"canDelete"`
	Created       string        `json:"created"`
	EncryptedText string        `json:"encryptedText,omitempty"`
	ValidUntil    string        `json:"validUntil"`
	IsClipOwner   bool          `json:"isClipOwner"`
	MaxTransfers  int           `json:"maxTransfers"`
	Transfers     int           `json:"transfers"`
	Files         []ClipFileDto `json:"files"`
	Version       int           `json:"version,omitempty"`
	IVData        QrcIVData
}

type ClipFileDto struct {
	Name       string `json:"name"`
	Index      int    `json:"index"`
	Size       int64  `json:"size"`
	ChunkCount int    `json:"chunkCount"`
}

type UpdateClipDto struct {
	EncryptedText    string `json:"encryptedText"`
	FileSize         int64  `json:"fileSize,omitempty"`
	ExpiresInMinutes int    `json:"expiresInMinutes,omitempty"`
	MaxTransfers     int    `json:"maxTransfers,omitempty"`
	AllowDelete      bool   `json:"allowDelete,omitempty"`
	FirstChunkSize   int64  `json:"firstChunkSize,omitempty"`
	Version          int    `json:"version,omitempty"`
	Storage          string `json:"storage,omitempty"`
}

type UpdateClipResponseDto struct {
	Ok             bool            `json:"ok"`
	Info           string          `json:"info"`
	ExpirationDate string          `json:"expirationDate,omitempty"`
	PreSignedPost  S3PreSignedPost `json:"preSignedPost,omitempty"`
	PreSignedPut   string          `json:"preSignedPut,omitempty"`
}

type CreateClipDto struct {
	ReceivingMode bool `json:"receivingMode"`
}

type S3PreSignedPost struct {
	Url    string                `json:"url"`
	Fields S3PreSignedPostFields `json:"fields"`
}

type S3PreSignedPostFields struct {
	Key            string `json:"key"`
	Bucket         string `json:"bucket"`
	XAmzAlgorithm  string `json:"X-Amz-Algorithm"`
	XAmzCredential string `json:"X-Amz-Credential"`
	XAmzDate       string `json:"X-Amz-Date"`
	Policy         string `json:"Policy"`
	XAmzSignature  string `json:"X-Amz-Signature"`
}

type GetFileChunkUploadLink struct {
	FileIndex  int   `json:"fileIndex"`
	ChunkIndex int   `json:"chunkIndex"`
	Size       int64 `json:"size"`
}

type FileChunkUploadLinkResponse struct {
	FileIndex     int             `json:"fileIndex"`
	ChunkIndex    int             `json:"chunkIndex"`
	PreSignedPost S3PreSignedPost `json:"preSignedPost,omitempty"`
	PreSignedPut  string          `json:"preSignedPut,omitempty"`
}

type FileUploadFinishedDto struct {
	Files []FileUploadFinishedFileDto `json:"files"`
}

type FileUploadFinishedFileDto struct {
	Name       string `json:"name"`
	Index      int    `json:"index"`
	Size       int64  `json:"size"`
	ChunkCount int    `json:"chunkCount"`
}

type FileUploadFinishedResponseDto struct {
	Ok             bool   `json:"ok"`
	ExpirationDate string `json:"expirationDate,omitempty"`
}

type LoginApprovalDto struct {
	Id         string `json:"id"`
	Key        string `json:"key"`
	Expiration string `json:"expiration"`
}

type LogInResponseDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Error        string `json:"error,omitempty"`
}

type LogInDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ClipLimitsDto struct {
	Text             int `json:"text"`
	FileMb           int `json:"fileMb"`
	FileNumber       int `json:"fileNumber"`
	ExpiresInMinutes int `json:"expiresInMinutes"`
	MaxTransfers     int `json:"maxTransfers"`
}

type RefreshTokenRequestDto struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

type LogInRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetFileDownloadTicketResponseDto struct {
	Url     string  `json:"url"`
	Id      string  `json:"id"`
	Key     string  `json:"key"`
	Credits float64 `json:"credits"`
	Error   int     `json:"error"`
}

type GetFileDownloadChunkResponseDto struct {
	Url   string `json:"url"`
	Chunk int    `json:"chunk"`
	Error string `json:"error"`
}

type QRClipConfigDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Storage      string `json:"storage"`
}

type QRCStorageLocation struct {
	Index int
	Code  string
	Name  string
}
