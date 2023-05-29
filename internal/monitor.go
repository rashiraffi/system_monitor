package internal

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func getCPUUtilization() (float64, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return cpuPercent[0], nil
}

func getMemoryUtilization() (float64, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memory.UsedPercent, nil
}

func insertData(db *sql.DB, cpuPercent, memPercent float64) error {
	query := "INSERT INTO system_resources (cpu_percent, mem_percent) VALUES (?, ?);"
	_, err := db.Exec(query, cpuPercent, memPercent)
	if err != nil {
		return err
	}
	return nil
}

func StartMonitor() {
	db := conn()
	defer db.Close()

	createTable(db)

	for {
		cpuPercent, err := getCPUUtilization()
		if err != nil {
			log.Printf("Failed to get CPU utilization: %v", err)
			continue
		}

		memPercent, err := getMemoryUtilization()
		if err != nil {
			log.Printf("Failed to get memory utilization: %v", err)
			continue
		}

		cpuPercent = float64(int(cpuPercent*100)) / 100
		memPercent = float64(int(memPercent*100)) / 100

		err = insertData(db, cpuPercent, memPercent)
		if err != nil {
			log.Printf("Failed to insert data into the database: %v", err)
		}

		fmt.Printf("%v\tCPU: %.2f%%\tMemory: %.2f%%\n", time.Now().Format("2006-01-02 15:04:05"), cpuPercent, memPercent)

		time.Sleep(2 * time.Second)
	}
}
