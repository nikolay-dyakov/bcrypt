# Small script to generate bcrypt hash based on password and cost

* `go build -o bcrypt main.go`

```bash
Usage of bcrypt:
  -c int
    	Cost (default 10)
  -cc
    	Copy bcrypt hash to clipboard
  -p string
    	Password to hash (default "bcrypt")
```
