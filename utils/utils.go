package utils

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yafyx/baak-api/models"
	"golang.org/x/time/rate"
)

// Initialize the random number generator with a unique seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	BaseURL = "https://baak.gunadarma.ac.id"
)

var (
	httpClient = &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     20,
			IdleConnTimeout:     20 * time.Second,
			DisableCompression:  false,
		},
	}
)

var (
	Limiter = rate.NewLimiter(rate.Limit(5), 10)
)

// List of common user agents to rotate through :)
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1",
}

func FetchDocument(url string) (*goquery.Document, error) {
	maxRetries := 3
	backoffFactor := 2.0
	initialBackoff := 1 * time.Second

	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create a new request
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set headers to appear more like a browser
		req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Cache-Control", "max-age=0")

		// Execute the request
		res, err := httpClient.Do(req)
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
				// For 403 errors, use a longer backoff
				backoffTime := time.Duration(float64(initialBackoff*2) * (backoffFactor * float64(attempt)))
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
