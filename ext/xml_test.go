package ext

import (
	"log"
	"testing"
)

const (
	EXAMPLE = "<?xml version=\"1.0\"?>\n<d:multistatus xmlns:d=\"DAV:\" xmlns:s=\"http://sabredav.org/ns\" xmlns:oc=\"http://owncloud.org/ns\"><d:response><d:href>/remote.php/dav/files/naudachu/temp/id/</d:href><d:propstat><d:prop><oc:fileid>33225</oc:fileid></d:prop><d:status>HTTP/1.1 200 OK</d:status></d:propstat></d:response></d:multistatus>\n"
)

// [ ] todo normal test...
func TestGetFileID(t *testing.T) {

	str, err := getFileIDFromRespBody([]byte(EXAMPLE))
	log.Print(str, err)

}
