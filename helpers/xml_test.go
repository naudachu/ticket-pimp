package helpers

import (
	"encoding/xml"
	"fmt"
	"log"
	"testing"
)

/*

<?xml version=\"1.0\"?>
<d:multistatus xmlns:d=\"DAV:\" xmlns:s=\"http://sabredav.org/ns\" xmlns:oc=\"http://owncloud.org/ns\">
    <d:response>
        <d:href>/remote.php/dav/files/naudachu/temp/id/</d:href>
        <d:propstat>
            <d:prop>
                <oc:fileid>33225</oc:fileid>
            </d:prop>
            <d:status>HTTP/1.1 200 OK</d:status>
        </d:propstat>
    </d:response>
</d:multistatus>

*/
/*
type MultistatusObj struct {
	XMLName     xml.Name `xml:"multistatus"`
	Multistatus struct {
		XMLName xml.Name `xml:"response"`
		Other   string   `xml:",innerxml"`
	}
}*/

type MultistatusObj struct {
	XMLName     xml.Name `xml:"multistatus"`
	Multistatus struct {
		XMLName  xml.Name `xml:"response"`
		Propstat struct {
			XMLName xml.Name `xml:"propstat"`
			Prop    struct {
				XMLName xml.Name `xml:"prop"`
				Other   string   `xml:",innerxml"`
			}
		}
	}
}

const (
	EXAMPLE = "<?xml version=\"1.0\"?>\n<d:multistatus xmlns:d=\"DAV:\" xmlns:s=\"http://sabredav.org/ns\" xmlns:oc=\"http://owncloud.org/ns\"><d:response><d:href>/remote.php/dav/files/naudachu/temp/id/</d:href><d:propstat><d:prop><oc:fileid>33225</oc:fileid></d:prop><d:status>HTTP/1.1 200 OK</d:status></d:propstat></d:response></d:multistatus>\n"
)

func GetFileID(str string) string {

	var multi MultistatusObj
	err := xml.Unmarshal([]byte(str), &multi)
	if err != nil {
		fmt.Print(err)
	}
	return multi.Multistatus.Propstat.Prop.Other
}

func TestGetFileID(t *testing.T) {
	str := GetFileID(EXAMPLE)
	log.Print(str)

}
