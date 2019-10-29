// client.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	client := http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:3000/articles/add", strings.NewReader("message"))
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now()

	req.Header.Add("x-req-id", t.Format("15:04:05"))
	req.Header.Write(os.Stdout)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

}
