package engine

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
)

var validHosts = []string{
	"github.com",
	"raw.githubusercontent.com",
	"gitlab.com",
}

func validateRemoteHost(maybeUrl string) bool {
	u, err := url.Parse(maybeUrl)
	return err == nil && slices.Contains(validHosts, u.Host)
}

func isUrl(maybeUrl string) bool {
	u, err := url.Parse(maybeUrl)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func downloadFile(url string, filepath string) error {
	if !validateRemoteHost(url) {
		return fmt.Errorf("unsupported remote host; { %s } are supported", strings.Join(validHosts, ", "))
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
