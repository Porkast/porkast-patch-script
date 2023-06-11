package feed

import "github.com/gogf/gf/v2/text/gstr"


func formatFeedAuthor(author string) (formatAuthor string) {

	if author != "" && gstr.HasSuffix(author, "|") {
		formatAuthor = author[:len(author)-1]
    } else {
		formatAuthor = author
	}

	return
}
