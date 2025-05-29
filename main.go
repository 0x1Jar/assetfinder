package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	var subsOnly bool
	var sourcesList string
	flag.BoolVar(&subsOnly, "subs-only", false, "Only include subdomains of search domain")
	flag.StringVar(&sourcesList, "sources", "all", "Comma-separated list of sources to use (e.g., certspotter,crtsh)")
	flag.Parse()

	var domains io.Reader
	domains = os.Stdin

	domain := flag.Arg(0)
	if domain != "" {
		domains = strings.NewReader(domain)
	}

	allSources := map[string]fetchFn{
		"certspotter":     fetchCertSpotter,
		"hackertarget":    fetchHackerTarget,
		"threatcrowd":     fetchThreatCrowd,
		"crtsh":           fetchCrtSh,
		"facebook":        fetchFacebook,
		"wayback":         fetchWayback,
		"virustotal":      fetchVirusTotal,
		"findsubdomains":  fetchFindSubDomains,
		"urlscan":         fetchUrlscan,
		"bufferoverrun":   fetchBufferOverrun,
	}

	selectedSources := []fetchFn{}
	if sourcesList == "all" {
		for _, fn := range allSources {
			selectedSources = append(selectedSources, fn)
		}
	} else {
		sourceNames := strings.Split(sourcesList, ",")
		for _, name := range sourceNames {
			name = strings.TrimSpace(name)
			if fn, ok := allSources[name]; ok {
				selectedSources = append(selectedSources, fn)
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Unknown source '%s' specified.\n", name)
			}
		}
	}

	if len(selectedSources) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No valid sources selected. Exiting.")
		os.Exit(1)
	}

	out := make(chan string)
	var wg sync.WaitGroup

	sc := bufio.NewScanner(domains)
	rl := newRateLimiter(time.Second)

	for sc.Scan() {
		domain := strings.ToLower(sc.Text())

		// call each of the source workers in a goroutine
		for _, source := range selectedSources {
			wg.Add(1)
			fn := source

			go func() {
				defer wg.Done()

				rl.Block(fmt.Sprintf("%#v", fn))
				names, err := fn(domain)

				if err != nil {
					//fmt.Fprintf(os.Stderr, "err: %s\n", err)
					return
				}

				for _, n := range names {
					n = cleanDomain(n)
					if subsOnly && !strings.HasSuffix(n, domain) {
						continue
					}
					out <- n
				}
			}()
		}
	}

	// close the output channel when all the workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// track what we've already printed to avoid duplicates
	printed := make(map[string]bool)

	for n := range out {
		if _, ok := printed[n]; ok {
			continue
		}
		printed[n] = true

		fmt.Println(n)
	}
}

type fetchFn func(string) ([]string, error)

func httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	raw, err := io.ReadAll(res.Body)

	res.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return raw, nil
}

func cleanDomain(d string) string {
	d = strings.ToLower(d)

	// no idea what this is, but we can't clean it ¯\_(ツ)_/¯
	if len(d) < 2 {
		return d
	}

	if d[0] == '*' || d[0] == '%' {
		d = d[1:]
	}

	if d[0] == '.' {
		d = d[1:]
	}

	return d

}

func fetchJSON(url string, wrapper interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	return dec.Decode(wrapper)
}
