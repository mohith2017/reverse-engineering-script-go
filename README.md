The challenge located at https://ciphersprint.pulley.com/

A Go script to help reverse engineer and find the flag

Main file for execution:
- `pulley.go`

How to execute locally:
- `go install <dependency_name>`
- `go run pulley.go`

Dependencies used:
- "encoding/base64"
- "encoding/json"
- "fmt"
- "io/ioutil"
- "net/http"
- "regexp"
- "strconv"
- "strings"
- "github.com/vmihailenco/msgpack/v5"
