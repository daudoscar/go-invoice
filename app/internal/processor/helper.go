package processor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CleanAmount(raw string) int {
	clean := regexp.MustCompile(`[^0-9.]`).ReplaceAllString(raw, "")
	re := regexp.MustCompile(`\d{1,3}(\.\d{3})*`)
	match := re.FindString(clean)
	val := strings.ReplaceAll(match, ".", "")
	n, _ := strconv.Atoi(val)
	return n
}

func ExtractMoney(text, pattern string) int {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(text)
	if len(m) > 1 {
		return CleanAmount(m[1])
	}
	return 0
}

func ParseDate(content string) (string, string) {
	re := regexp.MustCompile(`(\d{2}/\d{2}/\d{2})`)
	match := re.FindString(content)

	// Parsing format DD/MM/YY
	t, err := time.Parse("02/01/06", match)
	if err != nil {
		return "UNKNOWN MONTH", match
	}

	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	bulanIndo := fmt.Sprintf("%s %d", months[t.Month()-1], t.Year())
	tglLengkap := t.Format("02/01/2006")

	return bulanIndo, tglLengkap
}

func GetCleanLocation(content string) string {
	startKey := "Alamat Pembeli:"
	endKey := "No. Handphone"

	idxS := strings.Index(content, startKey)
	if idxS == -1 {
		return "Undetected"
	}

	sub := content[idxS+len(startKey):]
	idxE := strings.Index(sub, endKey)

	rawAddr := ""
	if idxE != -1 {
		rawAddr = strings.ToLower(sub[:idxE])
	}

	mapping := map[string][]string{
		"Head Office": {"graha irama", "rasuna said", "setia budi"},
		"Cibubur":     {"madison", "gunung putri", "ciangsana", "cibubur", "bogor"},
		"Pemuda":      {"pemuda raya", "pulo gadung", "rawamangun"},
		"Surabaya":    {"niaga gapura", "lidah kulon", "lakar santri", "surabaya"},
		"Bandung":     {"sunda", "sumur bandung", "kb. pisang", "bandung"},
		"BSD":         {"griya loka", "bsd", "serpong", "rawa buntu"},
		"Kemang":      {"city view", "kemang timur", "pejaten", "pasar minggu"},
		"Kebayoran":   {"sinabung", "kebayoran baru"},
		"Puri":        {"puri indah", "kembangan"},
		"Alam Sutera": {"weston lane", "alam sutera", "panunggangan", "tangerang"},
	}

	for loc, keys := range mapping {
		for _, k := range keys {
			if strings.Contains(rawAddr, k) {
				return loc
			}
		}
	}

	return "Undetected"
}
