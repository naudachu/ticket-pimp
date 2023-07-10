package helpers

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

func GitNaming(input string) string {
	// Remove leading and trailing whitespace
	input = strings.TrimSpace(input)

	// Replace non-Latin letters with spaces
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	input = strings.TrimSpace(reg.ReplaceAllString(input, " "))

	// Split into words
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	// Join words and return
	return strings.Join(words, "-")
}

type MultistatusObj struct {
	XMLName     xml.Name `xml:"multistatus"`
	Multistatus struct {
		XMLName  xml.Name `xml:"response"`
		Propstat struct {
			XMLName xml.Name `xml:"propstat"`
			Prop    struct {
				XMLName xml.Name `xml:"prop"`
				FileID  struct {
					XMLName xml.Name `xml:"fileid"`
					ID      string   `xml:",chardata"`
				}
			}
		}
	}
}

func GetFileIDFromRespBody(str []byte) int {

	var multi MultistatusObj

	err := xml.Unmarshal(str, &multi)
	if err != nil {
		return 0
	}

	id, err := strconv.Atoi(multi.Multistatus.Propstat.Prop.FileID.ID)
	if err != nil {
		return 0
	}

	return id
}
