package cronJob

import (
	"fmt"
	"hrms_backend/internal/services"
	"log"
)

type AttendanceCron struct {
	service services.AttendanceService
}

func NewAttendanceCron(service services.AttendanceService) *AttendanceCron {
	return &AttendanceCron{service: service}
}

// RunWeeklyReport generates and prints the attendance stats
func (j *AttendanceCron) RunWeeklyReport() {
	log.Println("------ 🗓️ Starting Weekly Attendance Report ------")

	records, err := j.service.GetWeeklyAttendance()
	if err != nil {
		log.Printf("❌ Error fetching weekly attendance: %v\n", err)
		return
	}

	if len(records) == 0 {
		log.Println("No attendance records found for the last 7 days.")
		return
	}

	// Aggregation Logic: Map[StudentID] -> Stats
	type Stats struct {
		Name    string
		Present int
		Total   int
	}
	report := make(map[uint]*Stats)

	for _, rec := range records {
		if _, exists := report[rec.StudentID]; !exists {
			report[rec.StudentID] = &Stats{Name: rec.StudentName, Present: 0, Total: 0}
		}

		stats := report[rec.StudentID]
		stats.Total++
		if rec.Status == "present" {
			stats.Present++
		}
	}

	// Print Summary
	for id, stats := range report {
		// e.g. "Student Alice (ID: 1) was present 4/5 days"
		fmt.Printf("🎓 Student %s (ID: %d) was present %d/%d days this week.\n",
			stats.Name, id, stats.Present, stats.Total)
	}

	log.Println("------ ✅ Weekly Report Completed ------")
}
