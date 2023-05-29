package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

func getTotalCount(db *sql.DB) (int, error) {
	if !tableExists(db) {
		return 0, nil
	}
	query := "SELECT COUNT(*) FROM system_resources;"
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getCPUUtilizationAboveThreshold(db *sql.DB, threshold int) (int, error) {
	query := "SELECT COUNT(*) FROM system_resources WHERE cpu_percent > ? AND cpu_percent <= ?;"
	var count int
	err := db.QueryRow(query, float64(threshold-10), float64(threshold)).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getRAMUtilizationAboveThreshold(db *sql.DB, threshold int) (int, error) {
	query := "SELECT COUNT(*) FROM system_resources WHERE mem_percent > ? AND mem_percent <= ?;"
	var count int
	err := db.QueryRow(query, float64(threshold-10), float64(threshold)).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func View() {
	db := conn()
	defer db.Close()

	count, err := getTotalCount(db)
	if err != nil {
		log.Fatalf("Failed to get total count: %v", err)
	}

	if count == 0 {
		fmt.Println("No data found..")
		return
	}

	threshold := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 95}
	tRange := make(map[int]rMonitor)

	for _, t := range threshold {

		cpuCount, err := getCPUUtilizationAboveThreshold(db, t)
		if err != nil {
			log.Fatalf("Failed to get CPU utilization above %d%%: %v", t, err)
		}

		memCount, err := getRAMUtilizationAboveThreshold(db, t)
		if err != nil {
			log.Fatalf("Failed to get memory utilization above %d%%: %v", t, err)
		}

		tRange[t] = rMonitor{
			cpuPercent: float64(cpuCount) / float64(count) * 100,
			memPercent: float64(memCount) / float64(count) * 100,
		}

	}

	cpuTable := tablewriter.NewWriter(os.Stdout)
	cpuTable.SetHeader([]string{"Threshold", "CPU Utilization"})
	for _, t := range threshold {
		if tRange[t].cpuPercent > 0 {
			cpuTable.Append([]string{fmt.Sprintf("%d%% ~ %d%%", t-10, t), fmt.Sprintf("%.2f%%", tRange[t].cpuPercent)})
		}
	}
	fmt.Println("CPU Utilization:")
	cpuTable.Render()

	ramTable := tablewriter.NewWriter(os.Stdout)
	ramTable.SetHeader([]string{"Threshold", "RAM Utilization"})
	for _, t := range threshold {
		if tRange[t].memPercent > 0 {
			ramTable.Append([]string{fmt.Sprintf("%d%% ~ %d%%", t-10, t), fmt.Sprintf("%.2f%%", tRange[t].memPercent)})
		}
	}
	fmt.Println("RAM Utilization:")
	ramTable.Render()

}
