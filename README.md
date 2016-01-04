# go-idcf
Golang library for IDCF Cloud services

## Example

### DNS

```go
package main

import (
	"fmt"

	"github.com/atsaki/go-idcf/dns"
)

func main() {

	APIKey := "IDCF_API_KEY"
	SecretKey := "IDCF_SECRET_KEY"

	c, _ := dns.NewClient(APIKey, SecretKey)

	zs, _ := c.Zones()
	for _, z := range zs {
		fmt.Printf("Zone: %s\n", z.Name)
		rs, _ := c.Records(z.UUID)
		for _, r := range rs {
			fmt.Printf("%s\t%s\t%v\n", r.Name, r.Type, r.Content)
		}
		fmt.Println("")
	}
}
```
