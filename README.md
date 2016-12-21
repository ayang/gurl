# gurl
a commandline tool like curl but very simple

## Usage:
    gurl [-body|-head] <method> <url> [name=value...]

## Examples:
    gurl get wwww.google.com
    echo hello | gurl post http://httpbin.org/post
    gurl post http://httpbin.org/post id=123 "name=jack bower"
    gurl -header a:b,c:d -basic admin:pass put http://httpbin.org/put id=123

## Install:
    go get github.com/ayang/gurl
