package x_http_client

const (
	TempFilePrefix = "x-go-temp-" // Temp file prefix
	TempFileSuffix = ".temp"      // Temp file suffix
	Version        = "v1.0.0"
)

type AuthVersionType string

const (
	// AuthV1 v1
	AuthV1 AuthVersionType = "v1"
	// AuthV2 v2
	AuthV2 AuthVersionType = "v2"
)

// HTTPMethod HTTP request method
type HTTPMethod string

const (
	// HTTPGet HTTP GET
	HTTPGet HTTPMethod = "GET"

	// HTTPPut HTTP PUT
	HTTPPut HTTPMethod = "PUT"

	// HTTPHead HTTP HEAD
	HTTPHead HTTPMethod = "HEAD"

	// HTTPPost HTTP POST
	HTTPPost HTTPMethod = "POST"

	// HTTPDelete HTTP DELETE
	HTTPDelete HTTPMethod = "DELETE"
)

// HTTP headers
const (
	HTTPHeaderAcceptEncoding string = "Accept-Encoding"
	HTTPHeaderAuthorization         = "Authorization"
	// HTTPHeaderCacheControl              = "Cache-Control"
	// HTTPHeaderContentDisposition        = "Content-Disposition"
	// HTTPHeaderContentEncoding           = "Content-Encoding"
	HTTPHeaderContentLength   = "Content-Length"
	HTTPHeaderContentMD5      = "Content-MD5"
	HTTPHeaderContentType     = "Content-Type"
	HTTPHeaderContentLanguage = "Content-Language"
	HTTPHeaderDate            = "Date"
	// HTTPHeaderEtag                      = "ETag"
	// HTTPHeaderExpires                   = "Expires"
	HTTPHeaderHost = "Host"
	// HTTPHeaderLastModified              = "Last-Modified"
	// HTTPHeaderRange                     = "Range"
	// HTTPHeaderLocation                  = "Location"
	// HTTPHeaderOrigin                    = "Origin"
	// HTTPHeaderServer                    = "Server"
	HTTPHeaderUserAgent = "User-Agent"
	// HTTPHeaderIfModifiedSince           = "If-Modified-Since"
	// HTTPHeaderIfUnmodifiedSince         = "If-Unmodified-Since"
	// HTTPHeaderIfMatch                   = "If-Match"
	// HTTPHeaderIfNoneMatch               = "If-None-Match"
	// HTTPHeaderACReqMethod               = "Access-Control-Request-Method"
	// HTTPHeaderACReqHeaders              = "Access-Control-Request-Headers"

	HTTPHeaderSecurityToken = "X-Security-Token"
	// HTTPHeaderOssServerSideEncryption        = "X-Oss-Server-Side-Encryption"
	// HTTPHeaderOssServerSideEncryptionKeyID   = "X-Oss-Server-Side-Encryption-Key-Id"
	// HTTPHeaderOssServerSideDataEncryption    = "X-Oss-Server-Side-Data-Encryption"
	// HTTPHeaderSSECAlgorithm                  = "X-Oss-Server-Side-Encryption-Customer-Algorithm"
	// HTTPHeaderSSECKey                        = "X-Oss-Server-Side-Encryption-Customer-Key"
	// HTTPHeaderSSECKeyMd5                     = "X-Oss-Server-Side-Encryption-Customer-Key-MD5"
	// HTTPHeaderOssCopySource                  = "X-Oss-Copy-Source"
	// HTTPHeaderOssCopySourceRange             = "X-Oss-Copy-Source-Range"
	// HTTPHeaderOssCopySourceIfMatch           = "X-Oss-Copy-Source-If-Match"
	// HTTPHeaderOssCopySourceIfNoneMatch       = "X-Oss-Copy-Source-If-None-Match"
	// HTTPHeaderOssCopySourceIfModifiedSince   = "X-Oss-Copy-Source-If-Modified-Since"
	// HTTPHeaderOssCopySourceIfUnmodifiedSince = "X-Oss-Copy-Source-If-Unmodified-Since"
	// HTTPHeaderOssMetadataDirective           = "X-Oss-Metadata-Directive"
	// HTTPHeaderOssNextAppendPosition          = "X-Oss-Next-Append-Position"
	HTTPHeaderTrackID   = "X-Track-Id"
	HTTPHeaderRequestID = "X-Request-Id"
	// HTTPHeaderOssCRC64                       = "X-Oss-Hash-Crc64ecma"
	// HTTPHeaderOssSymlinkTarget               = "X-Oss-Symlink-Target"
	// HTTPHeaderOssStorageClass                = "X-Oss-Storage-Class"
	// HTTPHeaderOssCallback                    = "X-Oss-Callback"
	// HTTPHeaderOssCallbackVar                 = "X-Oss-Callback-Var"
	// HTTPHeaderOssRequester                   = "X-Oss-Request-Payer"
	// HTTPHeaderOssTagging                     = "X-Oss-Tagging"
	// HTTPHeaderOssTaggingDirective            = "X-Oss-Tagging-Directive"
	// HTTPHeaderOssTrafficLimit                = "X-Oss-Traffic-Limit"
	// HTTPHeaderOssForbidOverWrite             = "X-Oss-Forbid-Overwrite"
	// HTTPHeaderOssRangeBehavior               = "X-Oss-Range-Behavior"
	// HTTPHeaderOssTaskID                      = "X-Oss-Task-Id"
)
