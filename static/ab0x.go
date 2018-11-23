package static


import (
  "bytes"
  
  "io"
  "log"
  "net/http"
  "os"

  "golang.org/x/net/webdav"
  "golang.org/x/net/context"
)

var ( 
  // CTX is a context for webdav vfs
  CTX = context.Background()

  
  // FS is a virtual memory file system
  FS = webdav.NewMemFS()
  

  // Handler is used to server files through a http handler
  Handler *webdav.Handler

  // HTTP is the http file system
  HTTP http.FileSystem = new(HTTPFS)
)

// HTTPFS implements http.FileSystem
type HTTPFS struct {}



// FileClassifyJSON is a file
var FileClassifyJSON = []byte("\x5b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x61\x6d\x65\x22\x3a\x20\x22\x63\x6f\x6d\x6d\x61\x6e\x64\x65\x72\x2d\x74\x79\x70\x65\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x61\x67\x73\x22\x3a\x20\x5b\x22\x63\x6f\x6d\x6d\x61\x6e\x64\x65\x72\x22\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x67\x6f\x22\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x22\x3a\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x79\x70\x65\x22\x3a\x20\x5b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x53\x75\x6d\x6d\x6f\x6e\x20\x4c\x65\x67\x65\x6e\x64\x7c\x4c\x65\x67\x65\x6e\x64\x61\x72\x79\x20\x43\x72\x65\x61\x74\x75\x72\x65\x29\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x61\x6d\x65\x22\x3a\x20\x22\x63\x6f\x6d\x6d\x61\x6e\x64\x65\x72\x2d\x74\x65\x78\x74\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x61\x67\x73\x22\x3a\x20\x5b\x22\x63\x6f\x6d\x6d\x61\x6e\x64\x65\x72\x22\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x67\x6f\x22\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x22\x3a\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x65\x78\x74\x22\x3a\x20\x5b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x63\x61\x6e\x20\x62\x65\x20\x79\x6f\x75\x72\x20\x63\x6f\x6d\x6d\x61\x6e\x64\x65\x72\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x61\x6d\x65\x22\x3a\x20\x22\x72\x61\x6d\x70\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x61\x67\x73\x22\x3a\x20\x5b\x22\x72\x61\x6d\x70\x22\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x67\x6f\x22\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x22\x3a\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x65\x78\x74\x22\x3a\x20\x5b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x61\x64\x64\x5b\x5e\x2e\x5d\x2b\x7b\x5b\x5e\x7d\x5d\x2b\x7d\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x61\x64\x64\x5b\x5e\x2e\x5d\x2b\x6d\x61\x6e\x61\x22\x2c\x0a\x09\x09\x09\x09\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x70\x6c\x61\x79\x5b\x5e\x2e\x5d\x2b\x61\x64\x64\x69\x74\x69\x6f\x6e\x61\x6c\x5b\x5e\x2e\x5d\x2b\x6c\x61\x6e\x64\x73\x3f\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x70\x75\x74\x5b\x5e\x2e\x5d\x2b\x6c\x61\x6e\x64\x73\x3f\x5b\x5e\x2e\x5d\x2b\x62\x61\x74\x74\x6c\x65\x66\x69\x65\x6c\x64\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x73\x65\x61\x72\x63\x68\x20\x79\x6f\x75\x72\x20\x6c\x69\x62\x72\x61\x72\x79\x20\x66\x6f\x72\x5b\x5e\x2e\x5d\x2b\x28\x66\x6f\x72\x65\x73\x74\x7c\x73\x77\x61\x6d\x70\x7c\x69\x73\x6c\x61\x6e\x64\x7c\x70\x6c\x61\x69\x6e\x73\x7c\x6d\x6f\x75\x6e\x74\x61\x69\x6e\x7c\x6c\x61\x6e\x64\x29\x73\x3f\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x6f\x74\x79\x70\x65\x22\x3a\x20\x5b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x2f\x6c\x61\x6e\x64\x2f\x69\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x61\x6d\x65\x22\x3a\x20\x22\x64\x72\x61\x77\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x67\x6f\x22\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x61\x67\x73\x22\x3a\x20\x5b\x22\x64\x72\x61\x77\x22\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x22\x3a\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x65\x78\x74\x22\x3a\x20\x5b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x63\x61\x72\x64\x73\x3f\x20\x69\x6e\x74\x6f\x20\x79\x6f\x75\x72\x20\x68\x61\x6e\x64\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x64\x72\x61\x77\x20\x28\x61\x7c\x78\x7c\x74\x68\x61\x74\x20\x6d\x61\x6e\x79\x7c\x6f\x6e\x65\x7c\x5c\x5c\x64\x2b\x29\x20\x63\x61\x72\x64\x73\x3f\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x61\x6d\x65\x22\x3a\x20\x22\x62\x6f\x61\x72\x64\x2d\x77\x69\x70\x65\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x67\x6f\x22\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x61\x67\x73\x22\x3a\x20\x5b\x22\x72\x65\x6d\x6f\x76\x61\x6c\x22\x2c\x22\x62\x6f\x61\x72\x64\x2d\x77\x69\x70\x65\x22\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x22\x3a\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x65\x78\x74\x22\x3a\x20\x5b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x09\x09\x22\x28\x3f\x69\x29\x28\x5e\x7c\x20\x29\x28\x73\x61\x63\x72\x69\x66\x69\x63\x65\x7c\x64\x65\x73\x74\x72\x6f\x79\x7c\x65\x78\x69\x6c\x65\x29\x20\x61\x6c\x6c\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x61\x6d\x65\x22\x3a\x20\x22\x74\x72\x69\x62\x65\x2d\x73\x75\x70\x70\x6f\x72\x74\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x61\x67\x73\x22\x3a\x20\x5b\x22\x74\x72\x69\x62\x65\x2d\x73\x75\x70\x70\x6f\x72\x74\x22\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x67\x6f\x22\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x22\x3a\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x65\x78\x74\x22\x3a\x20\x5b\x0a\x09\x09\x09\x09\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x28\x61\x64\x64\x69\x74\x69\x6f\x6e\x7c\x69\x73\x7c\x69\x6e\x73\x74\x61\x6e\x63\x65\x73\x3f\x7c\x68\x6f\x6f\x73\x65\x73\x3f\x7c\x73\x68\x61\x72\x65\x73\x3f\x7c\x61\x6c\x6c\x29\x5b\x5e\x2e\x5d\x2b\x63\x72\x65\x61\x74\x75\x72\x65\x20\x74\x79\x70\x65\x73\x3f\x7c\x63\x72\x65\x61\x74\x75\x72\x65\x5b\x5e\x2e\x5d\x2b\x63\x68\x6f\x69\x63\x65\x29\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x6e\x61\x6d\x65\x22\x3a\x20\x22\x63\x6f\x75\x6e\x74\x65\x72\x73\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x61\x67\x73\x22\x3a\x20\x5b\x22\x63\x6f\x75\x6e\x74\x65\x72\x73\x22\x5d\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x67\x6f\x22\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x72\x78\x22\x3a\x20\x7b\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x74\x65\x78\x74\x22\x3a\x20\x5b\x0a\x09\x09\x09\x09\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x70\x72\x6f\x6c\x69\x66\x65\x72\x61\x74\x65\x73\x3f\x7c\x28\x28\x61\x6e\x79\x7c\x61\x6e\x6f\x74\x68\x65\x72\x29\x20\x6b\x69\x6e\x64\x20\x6f\x66\x29\x3f\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x20\x28\x66\x72\x6f\x6d\x7c\x74\x6f\x7c\x6f\x6e\x7c\x61\x6c\x72\x65\x61\x64\x79\x29\x3f\x29\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x09\x09\x09\x09\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x42\x6f\x6c\x73\x74\x65\x72\x20\x5c\x5c\x64\x2b\x7c\x53\x75\x73\x70\x65\x6e\x64\x20\x5c\x5c\x64\x2b\x7c\x7b\x45\x7d\x29\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x09\x09\x09\x09\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x5b\x2b\x2d\x5d\x5c\x5c\x64\x2b\x2f\x5b\x2b\x2d\x5d\x5c\x5c\x64\x2b\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x41\x67\x65\x7c\x41\x69\x6d\x7c\x41\x72\x72\x6f\x77\x7c\x41\x72\x72\x6f\x77\x68\x65\x61\x64\x7c\x41\x77\x61\x6b\x65\x6e\x69\x6e\x67\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x42\x6c\x61\x7a\x65\x7c\x42\x6c\x6f\x6f\x64\x7c\x42\x6f\x75\x6e\x74\x79\x7c\x42\x72\x69\x62\x65\x72\x79\x7c\x42\x72\x69\x63\x6b\x7c\x43\x61\x72\x72\x69\x6f\x6e\x7c\x43\x68\x61\x72\x67\x65\x7c\x43\x52\x41\x4e\x4b\x21\x7c\x43\x72\x65\x64\x69\x74\x7c\x43\x6f\x72\x70\x73\x65\x7c\x43\x72\x79\x73\x74\x61\x6c\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x43\x75\x62\x65\x7c\x43\x75\x72\x72\x65\x6e\x63\x79\x7c\x44\x65\x61\x74\x68\x7c\x44\x65\x6c\x61\x79\x7c\x44\x65\x70\x6c\x65\x74\x69\x6f\x6e\x7c\x44\x65\x73\x70\x61\x69\x72\x7c\x44\x65\x76\x6f\x74\x69\x6f\x6e\x7c\x44\x69\x76\x69\x6e\x69\x74\x79\x7c\x44\x6f\x6f\x6d\x7c\x44\x72\x65\x61\x6d\x7c\x45\x63\x68\x6f\x7c\x45\x67\x67\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x45\x6c\x69\x78\x69\x72\x7c\x45\x6e\x65\x72\x67\x79\x7c\x45\x6f\x6e\x7c\x45\x78\x70\x65\x72\x69\x65\x6e\x63\x65\x7c\x45\x79\x65\x62\x61\x6c\x6c\x7c\x46\x61\x64\x65\x7c\x46\x61\x74\x65\x7c\x46\x65\x61\x74\x68\x65\x72\x7c\x46\x69\x6c\x69\x62\x75\x73\x74\x65\x72\x7c\x46\x6c\x6f\x6f\x64\x7c\x46\x75\x6e\x67\x75\x73\x7c\x46\x75\x73\x65\x7c\x47\x65\x6d\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x47\x6c\x79\x70\x68\x7c\x47\x6f\x6c\x64\x73\x7c\x47\x72\x6f\x77\x74\x68\x7c\x48\x61\x74\x63\x68\x6c\x69\x6e\x67\x7c\x48\x65\x61\x6c\x69\x6e\x67\x7c\x48\x69\x74\x7c\x48\x6f\x6f\x66\x70\x72\x69\x6e\x74\x7c\x48\x6f\x75\x72\x7c\x48\x6f\x75\x72\x67\x6c\x61\x73\x73\x7c\x48\x75\x6e\x67\x65\x72\x7c\x49\x63\x65\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x49\x6e\x63\x75\x62\x61\x74\x69\x6f\x6e\x7c\x49\x6e\x66\x65\x63\x74\x69\x6f\x6e\x7c\x49\x6e\x74\x65\x72\x76\x65\x6e\x74\x69\x6f\x6e\x7c\x49\x73\x6f\x6c\x61\x74\x69\x6f\x6e\x7c\x4a\x61\x76\x65\x6c\x69\x6e\x7c\x4b\x69\x7c\x4c\x65\x76\x65\x6c\x7c\x4c\x6f\x72\x65\x7c\x4c\x6f\x79\x61\x6c\x74\x79\x7c\x4c\x75\x63\x6b\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x4d\x61\x67\x6e\x65\x74\x7c\x4d\x61\x6e\x69\x66\x65\x73\x74\x61\x74\x69\x6f\x6e\x7c\x4d\x61\x6e\x6e\x65\x71\x75\x69\x6e\x7c\x4d\x61\x73\x6b\x7c\x4d\x61\x74\x72\x69\x78\x7c\x4d\x69\x6e\x65\x7c\x4d\x69\x6e\x69\x6e\x67\x7c\x4d\x69\x72\x65\x7c\x4d\x75\x73\x69\x63\x7c\x4d\x75\x73\x74\x65\x72\x7c\x4e\x65\x74\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x4f\x6d\x65\x6e\x7c\x4f\x72\x65\x7c\x50\x61\x67\x65\x7c\x50\x61\x69\x6e\x7c\x50\x61\x72\x61\x6c\x79\x7a\x61\x74\x69\x6f\x6e\x7c\x50\x65\x74\x61\x6c\x7c\x50\x65\x74\x72\x69\x66\x69\x63\x61\x74\x69\x6f\x6e\x7c\x50\x68\x79\x6c\x61\x63\x74\x65\x72\x79\x7c\x50\x69\x6e\x7c\x50\x6c\x61\x67\x75\x65\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x50\x6f\x69\x73\x6f\x6e\x7c\x50\x6f\x6c\x79\x70\x7c\x50\x72\x65\x73\x73\x75\x72\x65\x7c\x50\x72\x65\x79\x7c\x50\x75\x70\x61\x7c\x51\x75\x65\x73\x74\x7c\x52\x75\x73\x74\x7c\x53\x63\x72\x65\x61\x6d\x7c\x53\x68\x65\x6c\x6c\x7c\x53\x68\x69\x65\x6c\x64\x7c\x53\x69\x6c\x76\x65\x72\x7c\x53\x68\x72\x65\x64\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x53\x6c\x65\x65\x70\x7c\x53\x6c\x65\x69\x67\x68\x74\x7c\x53\x6c\x69\x6d\x65\x7c\x53\x6c\x75\x6d\x62\x65\x72\x7c\x53\x6f\x6f\x74\x7c\x53\x70\x6f\x72\x65\x7c\x53\x74\x6f\x72\x61\x67\x65\x7c\x53\x74\x72\x69\x66\x65\x7c\x53\x74\x75\x64\x79\x7c\x54\x68\x65\x66\x74\x7c\x54\x69\x64\x65\x7c\x54\x69\x6d\x65\x7c\x54\x6f\x77\x65\x72\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x2c\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x22\x28\x3f\x69\x29\x28\x28\x5e\x7c\x20\x29\x28\x54\x72\x61\x69\x6e\x69\x6e\x67\x7c\x54\x72\x61\x70\x7c\x54\x72\x65\x61\x73\x75\x72\x65\x7c\x56\x65\x6c\x6f\x63\x69\x74\x79\x7c\x56\x65\x72\x73\x65\x7c\x56\x69\x74\x61\x6c\x69\x74\x79\x7c\x56\x6f\x6c\x61\x74\x69\x6c\x65\x7c\x57\x61\x67\x65\x7c\x57\x69\x6e\x63\x68\x7c\x57\x69\x6e\x64\x7c\x57\x69\x73\x68\x29\x20\x63\x6f\x75\x6e\x74\x65\x72\x73\x3f\x29\x28\x20\x7c\x5c\x5c\x2e\x7c\x24\x29\x22\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x20\x5d\x0a\x09\x09\x7d\x0a\x20\x20\x20\x20\x20\x20\x20\x20\x7d\x0a\x5d\x0a")



func init() {
  if CTX.Err() != nil {
		log.Fatal(CTX.Err())
	}


var err error

  




  var f webdav.File
  

  
  

  f, err = FS.OpenFile(CTX, ".classify.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  
  _, err = f.Write(FileClassifyJSON)
  if err != nil {
    log.Fatal(err)
  }
  

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
  


  Handler = &webdav.Handler{
    FileSystem: FS,
    LockSystem: webdav.NewMemLS(),
  }
}

// Open a file
func (hfs *HTTPFS) Open(path string) (http.File, error) {
  f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
  if err != nil {
    return nil, err
  }

  return f, nil
}

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
  f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
  if err != nil {
    return nil, err
  }

  buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

  // If the buffer overflows, we will get bytes.ErrTooLarge.
  // Return that as an error. Any other panic remains.
  defer func() {
    e := recover()
    if e == nil {
      return
    }
    if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
      err = panicErr
    } else {
      panic(e)
    }
  }()
  _, err = buf.ReadFrom(f)
  return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func WriteFile(filename string, data []byte, perm os.FileMode) error {
  f, err := FS.OpenFile(CTX, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
  if err != nil {
    return err
  }
  n, err := f.Write(data)
  if err == nil && n < len(data) {
    err = io.ErrShortWrite
  }
  if err1 := f.Close(); err == nil {
    err = err1
  }
  return err
}

// FileNames is a list of files included in this filebox
var FileNames = []string {
  ".classify.json",
  
}
