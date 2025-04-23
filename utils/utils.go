package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yafyx/baak-api/models"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/time/rate"
)

// Initialize the random number generator with a unique seed
func init() {
	rand.Seed(time.Now().UnixNano())

	// Initialize cookie jar for session persistence
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create cookie jar: %v", err))
	}
	httpClient.Jar = jar
}

const (
	BaseURL = "https://baak.gunadarma.ac.id"
	BaseIP  = "103.23.40.57" // IP address for direct access
)

// Session cookies and visited pages for more authentic requests
var (
	visitedPages = []string{
		BaseURL,
		BaseURL + "/jadwal",
		BaseURL + "/kalender",
	}
	clientMutex = &sync.RWMutex{}
)

var (
	httpClient = &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     60 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableCompression:  true,
			ForceAttemptHTTP2:   false,
			DisableKeepAlives:   false,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip verification for testing
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
)

var (
	Limiter = rate.NewLimiter(rate.Limit(5), 10)
)

// List of common user agents to rotate through
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1",
}

// Additional realistic accept-language values
var acceptLanguages = []string{
	"en-US,en;q=0.9",
	"en-GB,en;q=0.8,en-US;q=0.9",
	"en-CA,en-US;q=0.9,en;q=0.8",
	"en-AU,en;q=0.9,en-GB;q=0.8",
	"id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7",
}

// Simulate human-like delays
func humanDelay() {
	// Random delay between 1-3 seconds to simulate human interaction
	delay := 1000 + rand.Intn(2000)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// getClient returns an HTTP client, potentially with a different proxy
func getClient() *http.Client {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	// If we have proxies configured, try to use them
	proxyManager := GetProxyManager()
	if proxyManager.HasProxies() {
		client, err := proxyManager.GetClient()
		if err == nil {
			return client
		}
		// If error, fall back to default client
	}

	return httpClient
}

// directIPRequest makes a request directly to the server's IP address
func directIPRequest() error {
	// Attempt direct connection via IP
	directURL := fmt.Sprintf("http://%s", BaseIP)
	req, err := http.NewRequest("GET", directURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create direct IP request: %v", err)
	}

	// Set headers to appear as a genuine browser request
	req.Host = "baak.gunadarma.ac.id" // Set the Host header to the domain name
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Execute the request
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("direct IP request failed: %v", err)
	}
	defer res.Body.Close()

	// Read and discard the body to ensure connection reuse
	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Check cookies
	baseUrl, _ := url.Parse(BaseURL)
	clientMutex.RLock()
	hasCookies := len(httpClient.Jar.Cookies(baseUrl)) > 0
	clientMutex.RUnlock()

	if !hasCookies {
		return fmt.Errorf("no cookies established with direct IP request")
	}

	return nil
}

// simpleRequest makes a very basic request to the given URL
func simpleRequest(targetURL string) error {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()

	// Read and discard the body to ensure connection reuse
	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	return nil
}

// Warm up the session by visiting the homepage first
func ensureSession() error {
	// Visit the homepage first to establish cookies if we haven't done so already
	baseUrl, err := url.Parse(BaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Check if we already have cookies for this domain
	clientMutex.RLock()
	hasCookies := len(httpClient.Jar.Cookies(baseUrl)) > 0
	clientMutex.RUnlock()

	if !hasCookies {
		fmt.Println("[DEBUG] No cookies found, trying to establish session")

		// Try HTTP first (some sites redirect HTTP to HTTPS)
		err := simpleRequest("http://baak.gunadarma.ac.id")
		if err == nil {
			// Check if we got cookies
			clientMutex.RLock()
			hasCookies := len(httpClient.Jar.Cookies(baseUrl)) > 0
			clientMutex.RUnlock()

			if hasCookies {
				fmt.Println("[DEBUG] Session established using HTTP request")
				return nil
			}
		} else {
			fmt.Printf("[DEBUG] simpleRequest(HTTP: %s) failed: %v\n", "http://baak.gunadarma.ac.id", err)
		}

		// Try HTTPS
		err = simpleRequest(BaseURL)
		if err == nil {
			// Check if we got cookies
			clientMutex.RLock()
			hasCookies := len(httpClient.Jar.Cookies(baseUrl)) > 0
			clientMutex.RUnlock()

			if hasCookies {
				fmt.Println("[DEBUG] Session established using HTTPS request")
				return nil
			}
		} else {
			fmt.Printf("[DEBUG] simpleRequest(HTTPS: %s) failed: %v\n", BaseURL, err)
		}

		// Try direct IP access as a last resort
		// fmt.Println("[DEBUG] Trying direct IP access method")
		// err = directIPRequest()
		// if err != nil {
		// 	 return fmt.Errorf("failed to establish session: %v", err)
		// } else {
		// 	 fmt.Println("[DEBUG] Session established using direct IP method")
		// }

		// If all methods failed without cookies
		clientMutex.RLock()
		hasCookiesAfterAttempts := len(httpClient.Jar.Cookies(baseUrl)) > 0
		clientMutex.RUnlock()
		if !hasCookiesAfterAttempts {
			return fmt.Errorf("failed to establish session after trying HTTP, HTTPS (check logs for details)")
		}

	} else {
		fmt.Println("[DEBUG] Session already established")
	}

	return nil
}

// Fetch a document with proper referrer and headers
func FetchDocumentWithRetry(url string, referrer string, maxRetries int) (*goquery.Document, error) {
	backoffFactor := 2.0
	initialBackoff := 1 * time.Second
	var lastErr error

	// Use a default referrer if none provided
	if referrer == "" {
		if len(visitedPages) > 0 {
			referrer = visitedPages[rand.Intn(len(visitedPages))]
		} else {
			referrer = BaseURL
		}
	}

	// Get a client (may have a different proxy)
	client := getClient()

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Add human-like delay between attempts
		if attempt > 0 {
			humanDelay()

			// For retry attempts, try to get a fresh client with potentially different proxy
			if attempt > 1 {
				client = getClient()
			}
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Create a new request
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Randomize User-Agent and other headers
		userAgent := userAgents[rand.Intn(len(userAgents))]
		acceptLang := acceptLanguages[rand.Intn(len(acceptLanguages))]

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", acceptLang)
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Referer", referrer)
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Cache-Control", "max-age=0")

		// Add a pseudo-random request ID to make each request unique
		req.Header.Set("X-Request-ID", fmt.Sprintf("%d", time.Now().UnixNano()))

		// Execute the request
		res, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to fetch URL: %v", err)
			backoffTime := time.Duration(float64(initialBackoff) * (backoffFactor * float64(attempt)))
			time.Sleep(backoffTime)
			continue
		}
		defer res.Body.Close()

		// Handle response based on status code
		if res.StatusCode != http.StatusOK {
			if res.StatusCode == http.StatusForbidden {
				lastErr = fmt.Errorf("access forbidden (403): the server might be restricting access or detecting automated requests")
				// For 403 errors, use a longer backoff with random jitter
				jitter := float64(1.0 + (rand.Float64() * 0.5)) // 1.0-1.5 jitter factor
				backoffTime := time.Duration(float64(initialBackoff*3) * (backoffFactor * float64(attempt) * jitter))
				time.Sleep(backoffTime)
				continue
			}

			lastErr = fmt.Errorf("unexpected status code: %d %s", res.StatusCode, res.Status)
			if attempt < maxRetries-1 {
				backoffTime := time.Duration(float64(initialBackoff) * (backoffFactor * float64(attempt)))
				time.Sleep(backoffTime)
				continue
			}
			return nil, lastErr
		}

		// Store current URL as visited page for future referrers
		if len(visitedPages) > 5 {
			// Keep the list to a reasonable size
			visitedPages = visitedPages[1:]
		}
		visitedPages = append(visitedPages, url)

		// Successfully got a 200 OK response, parse the document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML: %v", err)
		}

		return doc, nil
	}

	// If we got here, all attempts failed
	return nil, fmt.Errorf("all retry attempts failed: %v", lastErr)
}

func FetchDocument(url string) (*goquery.Document, error) {
	// Ensure we have an active session
	if err := ensureSession(); err != nil {
		return nil, err
	}

	// Add slight random delay to mimic human behavior
	humanDelay()

	return FetchDocumentWithRetry(url, "", 5) // Increase max retries to 5
}

// GetCSRFToken fetches a page and extracts the CSRF token from a hidden input field.
func GetCSRFToken(url string) (string, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch document for CSRF token: %w", err)
	}

	token, exists := doc.Find("input[name=\"_token\"]").First().Attr("value")
	if !exists {
		// Optionally log the HTML body here for debugging if token is not found
		// html, _ := doc.Html()
		// fmt.Println("DEBUG: HTML body:\n", html)
		return "", errors.New("CSRF token input field not found on page: " + url)
	}

	return token, nil
}

func GetJadwal(url string) (models.Jadwal, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return models.Jadwal{}, err
	}

	jadwal := models.Jadwal{}
	hariMap := map[string]*[]models.MataKuliah{
		"Senin":  &jadwal.Senin,
		"Selasa": &jadwal.Selasa,
		"Rabu":   &jadwal.Rabu,
		"Kamis":  &jadwal.Kamis,
		"Jum'at": &jadwal.Jumat,
		"Sabtu":  &jadwal.Sabtu,
	}

	timeStampLUT, err := GetTimeStampLUT()
	if err != nil {
		return models.Jadwal{}, err
	}

	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 5 {
			return
		}

		hari := strings.TrimSpace(cells.Eq(1).Text())
		waktu := strings.TrimSpace(cells.Eq(3).Text())
		jam := convertWaktuToJam(waktu, timeStampLUT)

		mataKuliah := models.MataKuliah{
			Nama:  strings.TrimSpace(cells.Eq(2).Text()),
			Waktu: waktu,
			Jam:   jam,
			Ruang: strings.TrimSpace(cells.Eq(4).Text()),
			Dosen: strings.TrimSpace(cells.Eq(5).Text()),
		}

		if hariSlice, ok := hariMap[hari]; ok {
			*hariSlice = append(*hariSlice, mataKuliah)
		}
	})

	return jadwal, nil
}

func GetTimeStampLUT() ([][]string, error) {
	doc, err := FetchDocument(BaseURL + "/kuliahUjian/6")
	if err != nil {
		return nil, err
	}

	var result [][]string
	doc.Find("table.cell-xs-6 tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() >= 2 {
			timeRange := strings.TrimSpace(cells.Eq(1).Text())
			timeRange = strings.ReplaceAll(timeRange, " ", "")
			timeRange = strings.ReplaceAll(timeRange, ".", ":")
			times := strings.Split(timeRange, "-")
			if len(times) == 2 {
				result = append(result, times)
			}
		}
	})

	return result, nil
}

func convertWaktuToJam(waktu string, timeStampLUT [][]string) string {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllString(waktu, -1)

	if len(matches) < 2 || len(matches) > 3 || len(timeStampLUT) == 0 {
		return ""
	}

	start, _ := strconv.Atoi(matches[0])
	end, _ := strconv.Atoi(matches[len(matches)-1])

	if start < 1 || start > len(timeStampLUT) || end < 1 || end > len(timeStampLUT) {
		return ""
	}

	return timeStampLUT[start-1][0] + " - " + timeStampLUT[end-1][1]
}

func GetKegiatan(url string) ([]models.Kegiatan, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}

	var kegiatanList []models.Kegiatan
	var parentKegiatan string

	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() == 2 {
			kegiatanText := strings.TrimSpace(cells.Eq(0).Text())
			tanggalText := strings.TrimSpace(cells.Eq(1).Text())

			if tanggalText == "" {
				parentKegiatan = kegiatanText
				return
			}

			start, end := parseTanggal(tanggalText)

			fullKegiatan := kegiatanText
			if parentKegiatan != "" && isSubItem(kegiatanText) {
				fullKegiatan = parentKegiatan + " " + kegiatanText
			} else {
				parentKegiatan = ""
			}

			kegiatan := models.Kegiatan{
				Kegiatan: fullKegiatan,
				Tanggal:  tanggalText,
				Start:    start,
				End:      end,
			}
			kegiatanList = append(kegiatanList, kegiatan)
		} else {
			parentKegiatan = ""
		}
	})

	return kegiatanList, nil
}

func isSubItem(text string) bool {
	return regexp.MustCompile(`^[a-z]\..+`).MatchString(text)
}

func parseTanggal(tanggal string) (start, end string) {
	parts := strings.Split(tanggal, "-")
	if len(parts) == 2 {
		start = strings.TrimSpace(parts[0])
		end = strings.TrimSpace(parts[1])
	} else if len(parts) == 1 {
		start = strings.TrimSpace(parts[0])
		end = start
	}

	return start, end
}

func GetKelasbaru(baseURL string) ([]models.KelasBaru, error) {
	var kelasBaru []models.KelasBaru
	page := 1

	for {
		url := fmt.Sprintf("%s&page=%d", baseURL, page)
		doc, err := FetchDocument(url)
		if err != nil {
			return nil, err
		}

		doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() == 5 {
				mhs := models.KelasBaru{
					NPM:       strings.TrimSpace(cells.Eq(1).Text()),
					Nama:      strings.TrimSpace(cells.Eq(2).Text()),
					KelasLama: strings.TrimSpace(cells.Eq(3).Text()),
					KelasBaru: strings.TrimSpace(cells.Eq(4).Text()),
				}
				kelasBaru = append(kelasBaru, mhs)
			}
		})

		if doc.Find(`a[rel="next"]`).Length() == 0 {
			break
		}

		page++
	}

	return kelasBaru, nil
}

func GetMahasiswaBaru(url string) ([]models.MahasiswaBaru, error) {
	var mahasiswaBaru []models.MahasiswaBaru
	page := 1

	for {
		pageURL := fmt.Sprintf("%s&page=%d", url, page)
		doc, err := FetchDocument(pageURL)
		if err != nil {
			return nil, err
		}

		doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() == 6 {
				mhs := models.MahasiswaBaru{
					NoPend:     strings.TrimSpace(cells.Eq(1).Text()),
					Nama:       strings.TrimSpace(cells.Eq(2).Text()),
					NPM:        strings.TrimSpace(cells.Eq(3).Text()),
					Kelas:      strings.TrimSpace(cells.Eq(4).Text()),
					Keterangan: strings.TrimSpace(cells.Eq(5).Text()),
				}
				mahasiswaBaru = append(mahasiswaBaru, mhs)
			}
		})

		if doc.Find(`a[rel="next"]`).Length() == 0 {
			break
		}

		page++
	}

	return mahasiswaBaru, nil
}

func GetUTS(url string) ([]models.UTS, error) {
	doc, err := FetchDocument(url)
	if err != nil {
		return nil, err
	}

	var utsList []models.UTS
	doc.Find("table").First().Find("tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() == 5 {
			uts := models.UTS{
				Nama:  strings.TrimSpace(cells.Eq(1).Text()),
				Waktu: strings.TrimSpace(cells.Eq(2).Text()),
				Ruang: strings.TrimSpace(cells.Eq(3).Text()),
				Dosen: strings.TrimSpace(cells.Eq(4).Text()),
			}
			utsList = append(utsList, uts)
		}
	})

	return utsList, nil
}

// EnsureSessionPublic is a public wrapper around ensureSession
func EnsureSessionPublic() error {
	return ensureSession()
}
