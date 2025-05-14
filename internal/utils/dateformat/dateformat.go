package dateformat

import (
	"fmt"
	"time"
)

func FormatDate(input time.Time) string {
	return input.Format("2006-01-02 15:04:05")
}
func FormatDateSecond(input time.Time) string {
	return input.Format("2006-01-02 15:04:05.999999")
}

func CalculateAge(birthdateStr string, refDate time.Time) (int, error) {
	birthdate, err := time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		return 0, fmt.Errorf("format tanggal salah: %v", err)
	}

	years := refDate.Year() - birthdate.Year()
	if refDate.YearDay() < birthdate.YearDay() {
		years--
	}

	return years, nil
}
func CalculateAgeDetailed(birthdateStr string, refDate time.Time) (string, error) {
	// Konversi string ke time.Time
	birthdate, err := time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		return "", fmt.Errorf("format tanggal salah: %v", err)
	}

	// Hitung selisih tahun
	years := refDate.Year() - birthdate.Year()
	if refDate.YearDay() < birthdate.YearDay() {
		years-- // Belum ulang tahun tahun ini
	}

	// Hitung selisih bulan & hari
	birthdateThisYear := birthdate.AddDate(years, 0, 0)             // Tanggal ulang tahun terakhir
	months := int(refDate.Sub(birthdateThisYear).Hours() / 24 / 30) // Perkiraan bulan
	birthdateWithMonths := birthdateThisYear.AddDate(0, months, 0)  // Tambahkan bulan
	days := int(refDate.Sub(birthdateWithMonths).Hours() / 24)      // Hitung sisa hari

	// Format hasil
	return fmt.Sprintf("%dy %dm %dd", years, months, days), nil
}
