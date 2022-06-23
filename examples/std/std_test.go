package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func importAllStdLibCmds(shouldRun bool) string {
	pkgListStr := `archive/tar
	archive/zip
	bufio
	bytes
	compress/bzip2
	compress/flate
	compress/gzip
	compress/lzw
	compress/zlib
	container/heap
	container/list
	container/ring
	context
	crypto
	crypto/aes
	crypto/cipher
	crypto/des
	crypto/dsa
	crypto/ecdsa
	crypto/elliptic
	crypto/hmac
	crypto/md5
	crypto/rand
	crypto/rc4
	crypto/rsa
	crypto/sha1
	crypto/sha256
	crypto/sha512
	crypto/subtle
	crypto/tls
	crypto/x509
	crypto/x509/pkix
	database/sql
	database/sql/driver
	debug/dwarf
	debug/elf
	debug/gosym
	debug/macho
	debug/pe
	debug/plan9obj
	encoding
	encoding/ascii85
	encoding/asn1
	encoding/base32
	encoding/base64
	encoding/binary
	encoding/csv
	encoding/gob
	encoding/hex
	encoding/json
	encoding/pem
	encoding/xml
	errors
	expvar
	flag
	fmt
	hash
	hash/adler32
	hash/crc32
	hash/crc64
	hash/fnv
	html
	html/template
	image
	image/color
	image/color/palette
	image/draw
	image/gif
	image/jpeg
	image/png
	index/suffixarray
	io
	io/ioutil
	log
	log/syslog
	math
	math/big
	math/cmplx
	math/rand
	mime
	mime/multipart
	mime/quotedprintable
	net
	net/http
	net/http/cgi
	net/http/cookiejar
	net/http/fcgi
	net/http/httptest
	net/http/httptrace
	net/http/httputil
	net/http/pprof
	net/mail
	net/rpc
	net/rpc/jsonrpc
	net/smtp
	net/textproto
	net/url
	os
	os/exec
	os/signal
	os/user
	path
	path/filepath
	plugin
	reflect
	regexp
	regexp/syntax
	runtime
	sort
	strconv
	strings
	sync
	sync/atomic
	text/scanner
	text/tabwriter
	text/template
	text/template/parse
	time
	unicode
	unicode/utf16
	unicode/utf8
	`

	pkgListStr = strings.ReplaceAll(pkgListStr, "\r", "")

	pkgList := strings.Split(pkgListStr, "\n")

	res := ""

	for _, inputPkg := range pkgList {
		inputPkg = strings.TrimSpace(inputPkg)
		if inputPkg != "" {
			outputPath := strings.ReplaceAll(inputPkg, "/", ".")
			outputPath = "js/" + outputPath
			outputGo := outputPath + ".go"
			outputJs := outputPath + ".ts"

			outputPkg := "js"

			args := []string{"goja_imports", "-p", outputPkg, "-i", inputPkg, "-go", outputGo, "-ja", outputJs}
			cmd := strings.Join(args, " ")
			res += cmd + "\r\n"

			if shouldRun {
				c := exec.Command(args[0], args[1:]...)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				err := c.Run()
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return res
}

func TestImportAllStdLibCmds(t *testing.T) {
	cmd := importAllStdLibCmds(false)
	fmt.Println(cmd)
}

func TestImportAllStdLibCmdsRun(t *testing.T) {
	os.RemoveAll("./js")
	importAllStdLibCmds(true)
	//fmt.Println(cmd)
}

func TestRun(t *testing.T) {
	//json.Marshal()
}
