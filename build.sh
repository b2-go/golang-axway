go get github.com/gin-gonic/gin
go get github.com/stretchr/testify
go get gopkg.in/mgo.v2
go test
go build
docker build . -t axwaymws
