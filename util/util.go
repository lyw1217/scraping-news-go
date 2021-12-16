package util

import (
	"log"

	// "github.com/djimenez/iconv-go" // https://pkg.go.dev/github.com/djimenez/iconv-go
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

/*
// euc-kr 문자열을 utf-8 문자열로 변환 (iconv-go)
func ConvEuckrToUtf8(input string) string {
	output, err := iconv.ConvertString(input, "euc-kr", "utf-8")
	if err != nil {
		log.Println(err)
		return ""
	}
	return output
}
*/
// euc-kr 문자열을 utf-8 문자열로 변환 (korean, transform)
func TransEuckrToUtf8(input string) string {
	euckrDec := korean.EUCKR.NewDecoder()

	output, _, err := transform.String(euckrDec, input)
	if err != nil {
		log.Println(err)
		return ""
	}
	return output
}
