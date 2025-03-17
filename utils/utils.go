package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/yafyx/baak-api/models"
	"golang.org/x/time/rate"
)

const (
	BaseURL = "https://baak.gunadarma.ac.id"
)

var (
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
)

var (
	Limiter = rate.NewLimiter(rate.Limit(5), 10)
)

func FetchDocument(url string) (*goquery.Document, error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	return doc, nil
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
