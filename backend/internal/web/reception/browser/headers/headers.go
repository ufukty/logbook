package headers

const (
	CacheControl    = "Cache-Control"    // (Request/Response) Directives for caching mechanisms. Eg. `no-cache, no-store, must-revalidate`
	Connection      = "Connection"       // (Request/Response) Controls options for the connection. Eg. `keep-alive`
	ContentEncoding = "Content-Encoding" // (Request/Response) Compression applied to the body. Eg. `gzip`
	ContentLanguage = "Content-Language" // (Request/Response) Natural language of the body. Eg. `en`
	ContentLength   = "Content-Length"   // (Request/Response) Size of the body in bytes. Eg. `348`
	ContentLocation = "Content-Location" // (Request/Response) Alternate location for returned data. Eg. `/index.htm`
	ContentRange    = "Content-Range"    // (Request/Response) Part of a full body returned. Eg. `bytes 200-1000/67589`
	ContentType     = "Content-Type"     // (Request/Response) Media type of the body. Eg. `application/json; charset=utf-8`
	Date            = "Date"             // (Request/Response) Date/time of message origination. Eg. `Tue, 15 Nov 1994 08:12:31 GMT`
	Pragma          = "Pragma"           // (Request/Response) Implementation-specific directives (deprecated, replaced by Cache-Control). Eg. `no-cache`
	Trailer         = "Trailer"          // (Request/Response) Headers present after the body in chunked transfer. Eg. `Expires`
)

const (
	Accept                      = "Accept"                         // (Request) Media types the client can handle. Eg. `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`
	AcceptCharset               = "Accept-Charset"                 // (Request) Character sets accepted. Eg. `utf-8, iso-8859-1;q=0.5`
	AcceptDatetime              = "Accept-Datetime"                // (Request) (Experimental) Acceptable version-date of resource. Eg. `Thu, 31 May 2007 20:35:00 GMT`
	AcceptEncoding              = "Accept-Encoding"                // (Request) Content codings the client can handle. Eg. `gzip, deflate, br`
	AcceptLanguage              = "Accept-Language"                // (Request) Preferred natural languages. Eg. `en-US,en;q=0.5`
	AccessControlRequestHeaders = "Access-Control-Request-Headers" // (Request) Used in CORS preflight to indicate custom headers. Eg. `Content-Type, Authorization`
	AccessControlRequestMethod  = "Access-Control-Request-Method"  // (Request) Used in CORS preflight to indicate method. Eg. `POST`
	AIM                         = "A-IM"                           // (Request) Instance-manipulations the client supports. Eg. `feed`
	Authorization               = "Authorization"                  // (Request) Credentials for authentication. Eg. `Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==`
	Cookie                      = "Cookie"                         // (Request) Client cookies. Eg. `sessionId=abc123; theme=light`
	Expect                      = "Expect"                         // (Request) Indicates expectations. Commonly `100-continue`. Eg. `100-continue`
	Forwarded                   = "Forwarded"                      // (Request) Proxy information. Eg. `for=192.0.2.60;proto=http;by=203.0.113.43`
	From                        = "From"                           // (Request) Email address of user making request. Eg. `ufuktan@ufukty.com`
	Host                        = "Host"                           // (Request) Hostname of server. Required in HTTP/1.1. Eg. `ufukty.com`
	IfMatch                     = "If-Match"                       // (Request) Conditional request: proceed if ETag matches. Eg. `"737060cd8c284d8af7ad3082f209582d"`
	IfModifiedSince             = "If-Modified-Since"              // (Request) Conditional GET: send if newer. Eg. `Sat, 29 Oct 1994 19:43:31 GMT`
	IfNoneMatch                 = "If-None-Match"                  // (Request) Conditional GET: send if ETag doesn’t match. Eg. `"737060cd8c284d8af7ad3082f209582d"`
	IfRange                     = "If-Range"                       // (Request) Conditional range request based on ETag or date. Eg. `"737060cd8c284d8af7ad3082f209582d"`
	IfUnmodifiedSince           = "If-Unmodified-Since"            // (Request) Conditional: proceed only if resource not modified since. Eg. `Sat, 29 Oct 1994 19:43:31 GMT`
	MaxForwards                 = "Max-Forwards"                   // (Request) Limits proxy/forwarding hops. Eg. `10`
	Origin                      = "Origin"                         // (Request) Origin of the request for CORS. Eg. `https://ufukty.com`
	ProxyAuthorization          = "Proxy-Authorization"            // (Request) Credentials for proxy authentication. Eg. `Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==`
	Range                       = "Range"                          // (Request) Requests partial resource. Eg. `bytes=200-1000`
	Referer                     = "Referer"                        // (Request) Page that linked to resource. Eg. `http://www.ufukty.com/start.html`
	SecFetchDest                = "Sec-Fetch-Dest"                 // (Request) Fetch metadata: request destination. Eg. `document` | `image` | `script` | `style` | `iframe`
	SecFetchMode                = "Sec-Fetch-Mode"                 // (Request) Fetch metadata: mode. Eg. `cors` | `no-cors` | `same-origin` | `navigate`
	SecFetchSite                = "Sec-Fetch-Site"                 // (Request) Fetch metadata: relationship of origin. Eg. `same-origin` | `same-site` | `cross-site` | `none`
	SecFetchUser                = "Sec-Fetch-User"                 // (Request) Fetch metadata: user activation. Eg. `?1`
	TE                          = "TE"                             // (Request) Transfer encodings accepted. Eg. `trailers, deflate`
	Upgrade                     = "Upgrade"                        // (Request) Protocol upgrade (e.g. WebSocket). Eg. `websocket`
	UserAgent                   = "User-Agent"                     // (Request) Client software identifier. Eg. `Mozilla/5.0 (Windows NT 10.0; Win64; x64)`
)

const (
	AcceptRanges                  = "Accept-Ranges"                    // (Response) Indicates if server supports partial requests. Eg. `bytes`
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials" // (Response) Whether response can expose credentials (CORS). Eg. `true`
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"     // (Response) Headers permitted in CORS request. Eg. `Content-Type, Authorization, X-Custom-Header`
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"     // (Response) Methods permitted in CORS request. Eg. `GET, POST, PUT, DELETE`
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"      // (Response) Allowed origin(s) for CORS. Eg. `*` OR `https://ufukty.com`
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"    // (Response) Headers accessible to scripts (CORS). Eg. `Content-Length, X-Kuma-Revision`
	AccessControlMaxAge           = "Access-Control-Max-Age"           // (Response) How long preflight response can be cached. Eg. `600`
	Age                           = "Age"                              // (Response) Time since response was generated (seconds). Eg. `3600`
	Allow                         = "Allow"                            // (Response) Valid methods for the resource. Eg. `GET, POST, HEAD`
	ETag                          = "ETag"                             // (Response) Identifier for specific resource version. Eg. `"33a64df551425fcc55e4d42a148795d9f25f89d4"`
	Expires                       = "Expires"                          // (Response) Date/time after which response is stale. Eg. `Wed, 21 Oct 2015 07:28:00 GMT`
	LastModified                  = "Last-Modified"                    // (Response) Timestamp of last modification. Eg. `Tue, 15 Nov 1994 12:45:26 GMT`
	Location                      = "Location"                         // (Response) Redirect target or new resource URI. Eg. `http://www.ufukty.com/newpage.html`
	ProxyAuthenticate             = "Proxy-Authenticate"               // (Response) Authentication method required by proxy. Eg. `Basic realm="Access to the staging site"`
	RetryAfter                    = "Retry-After"                      // (Response) When client can retry request. Eg. `120`
	Server                        = "Server"                           // (Response) Server software details. Eg. `nginx/1.14.1`
	SetCookie                     = "Set-Cookie"                       // (Response) Send cookies from server. Eg. `sessionId=abc123; Path=/; HttpOnly`
	TransferEncoding              = "Transfer-Encoding"                // (Response) Encoding form of body transfer. Eg. `chunked`
	Vary                          = "Vary"                             // (Response) Headers that affect response selection. Eg. `Accept-Encoding`
	Via                           = "Via"                              // (Response) Intermediate proxies information. Eg. `1.0 fred, 1.1 ufukty.com (Apache/1.1)`
	Warning                       = "Warning"                          // (Response) Additional info about status/transformations. Eg. `199 Miscellaneous warning`
	WWWAuthenticate               = "WWW-Authenticate"                 // (Response) Authentication method required by server. Eg. `Basic realm="Access to the site"`
)
