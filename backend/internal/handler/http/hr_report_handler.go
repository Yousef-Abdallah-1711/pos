package http

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pos-saas/internal/middleware"
)

type HRReportHandler struct {
	db *sql.DB
}

func NewHRReportHandler(db *sql.DB) *HRReportHandler {
	return &HRReportHandler{db: db}
}

func (h *HRReportHandler) GetHeadcountReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Database is not available")
		return
	}

	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	department := strings.TrimSpace(r.URL.Query().Get("department"))
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	employmentType := strings.TrimSpace(r.URL.Query().Get("employment_type"))
	if status != "" {
		switch status {
		case "active", "inactive", "terminated", "on_leave":
		default:
			respondError(w, http.StatusBadRequest, "Invalid status filter")
			return
		}
	}
	if employmentType != "" {
		switch employmentType {
		case "full_time", "part_time", "contract", "intern":
		default:
			respondError(w, http.StatusBadRequest, "Invalid employment_type filter")
			return
		}
	}

	conditions := []string{
		"tenant_id = $1",
		"restaurant_id = $2",
		"is_active = true",
	}
	args := []interface{}{tenantID, restaurantID}

	if department != "" {
		args = append(args, department)
		conditions = append(conditions, "department = $"+strconv.Itoa(len(args)))
	}
	if status != "" {
		args = append(args, status)
		conditions = append(conditions, "employment_status = $"+strconv.Itoa(len(args)))
	}
	if employmentType != "" {
		args = append(args, employmentType)
		conditions = append(conditions, "employment_type = $"+strconv.Itoa(len(args)))
	}

	where := strings.Join(conditions, " AND ")

	var total int64
	if err := h.db.QueryRow("SELECT COUNT(*) FROM employees WHERE "+where, args...).Scan(&total); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
		return
	}

	byDepartment := make([]map[string]interface{}, 0)
	rows, err := h.db.Query(`
		SELECT COALESCE(NULLIF(TRIM(department), ''), 'Unassigned') AS department, COUNT(*) AS count
		FROM employees
		WHERE `+where+`
		GROUP BY 1
		ORDER BY 2 DESC, 1 ASC
	`, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
		return
	}
	for rows.Next() {
		var dept string
		var count int64
		if scanErr := rows.Scan(&dept, &count); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
			return
		}
		byDepartment = append(byDepartment, map[string]interface{}{
			"department": dept,
			"count":      count,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
		return
	}

	byStatus := make([]map[string]interface{}, 0)
	rows, err = h.db.Query(`
		SELECT employment_status, COUNT(*) AS count
		FROM employees
		WHERE `+where+`
		GROUP BY 1
		ORDER BY 2 DESC, 1 ASC
	`, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
		return
	}
	for rows.Next() {
		var itemStatus string
		var count int64
		if scanErr := rows.Scan(&itemStatus, &count); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
			return
		}
		byStatus = append(byStatus, map[string]interface{}{
			"status": itemStatus,
			"count":  count,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
		return
	}

	byType := make([]map[string]interface{}, 0)
	rows, err = h.db.Query(`
		SELECT employment_type, COUNT(*) AS count
		FROM employees
		WHERE `+where+`
		GROUP BY 1
		ORDER BY 2 DESC, 1 ASC
	`, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
		return
	}
	for rows.Next() {
		var itemType string
		var count int64
		if scanErr := rows.Scan(&itemType, &count); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
			return
		}
		byType = append(byType, map[string]interface{}{
			"employment_type": itemType,
			"count":           count,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate headcount report")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"filters": map[string]interface{}{
				"department":      department,
				"status":          status,
				"employment_type": employmentType,
			},
			"summary": map[string]interface{}{
				"total_employees": total,
			},
			"by_department":      byDepartment,
			"by_status":          byStatus,
			"by_employment_type": byType,
			"generated_at":       time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *HRReportHandler) GetTurnoverReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Database is not available")
		return
	}

	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	from, to, err := parseHRDateRange(r.URL.Query().Get("date_from"), r.URL.Query().Get("date_to"), 365)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var hires, departures int64
	if err = h.db.QueryRow(`
		SELECT COUNT(*)
		FROM employees
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND hire_date BETWEEN $3 AND $4
	`, tenantID, restaurantID, from, to).Scan(&hires); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate turnover report")
		return
	}
	if err = h.db.QueryRow(`
		SELECT COUNT(*)
		FROM employees
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND termination_date IS NOT NULL
		  AND termination_date BETWEEN $3 AND $4
	`, tenantID, restaurantID, from, to).Scan(&departures); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate turnover report")
		return
	}

	var headcountStart, headcountEnd int64
	if err = h.db.QueryRow(`
		SELECT COUNT(*)
		FROM employees
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND hire_date <= $3
		  AND (termination_date IS NULL OR termination_date > $3)
	`, tenantID, restaurantID, from).Scan(&headcountStart); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate turnover report")
		return
	}
	if err = h.db.QueryRow(`
		SELECT COUNT(*)
		FROM employees
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND hire_date <= $3
		  AND (termination_date IS NULL OR termination_date > $3)
	`, tenantID, restaurantID, to).Scan(&headcountEnd); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate turnover report")
		return
	}

	averageHeadcount := float64(headcountStart+headcountEnd) / 2.0
	turnoverRate := 0.0
	if averageHeadcount > 0 {
		turnoverRate = (float64(departures) / averageHeadcount) * 100
	}

	monthly := make([]map[string]interface{}, 0)
	rows, err := h.db.Query(`
		WITH months AS (
			SELECT generate_series(
				date_trunc('month', $3::date),
				date_trunc('month', $4::date),
				interval '1 month'
			) AS month_start
		),
		hires AS (
			SELECT date_trunc('month', hire_date)::date AS month_start, COUNT(*) AS hires
			FROM employees
			WHERE tenant_id = $1
			  AND restaurant_id = $2
			  AND hire_date BETWEEN $3 AND $4
			GROUP BY 1
		),
		departures AS (
			SELECT date_trunc('month', termination_date)::date AS month_start, COUNT(*) AS departures
			FROM employees
			WHERE tenant_id = $1
			  AND restaurant_id = $2
			  AND termination_date IS NOT NULL
			  AND termination_date BETWEEN $3 AND $4
			GROUP BY 1
		)
		SELECT
			to_char(m.month_start, 'YYYY-MM') AS month,
			COALESCE(h.hires, 0) AS hires,
			COALESCE(d.departures, 0) AS departures
		FROM months m
		LEFT JOIN hires h ON h.month_start = m.month_start::date
		LEFT JOIN departures d ON d.month_start = m.month_start::date
		ORDER BY m.month_start ASC
	`, tenantID, restaurantID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate turnover report")
		return
	}
	for rows.Next() {
		var month string
		var monthlyHires, monthlyDepartures int64
		if scanErr := rows.Scan(&month, &monthlyHires, &monthlyDepartures); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate turnover report")
			return
		}
		monthly = append(monthly, map[string]interface{}{
			"month":      month,
			"hires":      monthlyHires,
			"departures": monthlyDepartures,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate turnover report")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"filters": map[string]interface{}{
				"date_from": from.Format("2006-01-02"),
				"date_to":   to.Format("2006-01-02"),
			},
			"summary": map[string]interface{}{
				"hires":             hires,
				"departures":        departures,
				"average_headcount": averageHeadcount,
				"turnover_rate":     turnoverRate,
			},
			"trend":        monthly,
			"generated_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *HRReportHandler) GetAbsenteeismReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Database is not available")
		return
	}

	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	from, to, err := parseHRDateRange(r.URL.Query().Get("date_from"), r.URL.Query().Get("date_to"), 30)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var totalRecords, absentRecords int64
	if err = h.db.QueryRow(`
		SELECT COUNT(*)
		FROM attendance
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND attendance_date BETWEEN $3 AND $4
	`, tenantID, restaurantID, from, to).Scan(&totalRecords); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
		return
	}
	if err = h.db.QueryRow(`
		SELECT COUNT(*)
		FROM attendance
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND attendance_date BETWEEN $3 AND $4
		  AND status = 'absent'
	`, tenantID, restaurantID, from, to).Scan(&absentRecords); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
		return
	}

	rate := 0.0
	if totalRecords > 0 {
		rate = (float64(absentRecords) / float64(totalRecords)) * 100
	}

	patterns := make([]map[string]interface{}, 0)
	rows, err := h.db.Query(`
		SELECT
			EXTRACT(DOW FROM attendance_date)::INT AS day_of_week,
			COUNT(*) AS absent_days
		FROM attendance
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND attendance_date BETWEEN $3 AND $4
		  AND status = 'absent'
		GROUP BY 1
		ORDER BY 2 DESC, 1 ASC
	`, tenantID, restaurantID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
		return
	}
	for rows.Next() {
		var dayOfWeek int
		var absentDays int64
		if scanErr := rows.Scan(&dayOfWeek, &absentDays); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
			return
		}
		patterns = append(patterns, map[string]interface{}{
			"day_of_week": dayOfWeek,
			"day_name":    weekdayName(dayOfWeek),
			"absent_days": absentDays,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
		return
	}

	topAbsentees := make([]map[string]interface{}, 0)
	rows, err = h.db.Query(`
		SELECT
			a.employee_id,
			TRIM(COALESCE(e.first_name, '') || ' ' || COALESCE(e.last_name, '')) AS employee_name,
			COUNT(*) AS absent_days
		FROM attendance a
		LEFT JOIN employees e
			ON e.id = a.employee_id
			AND e.tenant_id = a.tenant_id
			AND e.restaurant_id = a.restaurant_id
		WHERE a.tenant_id = $1
		  AND a.restaurant_id = $2
		  AND a.attendance_date BETWEEN $3 AND $4
		  AND a.status = 'absent'
		GROUP BY a.employee_id, employee_name
		ORDER BY absent_days DESC, employee_name ASC
		LIMIT 10
	`, tenantID, restaurantID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
		return
	}
	for rows.Next() {
		var employeeID int64
		var employeeName string
		var absentDays int64
		if scanErr := rows.Scan(&employeeID, &employeeName, &absentDays); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
			return
		}
		topAbsentees = append(topAbsentees, map[string]interface{}{
			"employee_id":   employeeID,
			"employee_name": strings.TrimSpace(employeeName),
			"absent_days":   absentDays,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate absenteeism report")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"filters": map[string]interface{}{
				"date_from": from.Format("2006-01-02"),
				"date_to":   to.Format("2006-01-02"),
			},
			"summary": map[string]interface{}{
				"total_records":    totalRecords,
				"absent_records":   absentRecords,
				"absenteeism_rate": rate,
			},
			"patterns":      patterns,
			"top_absentees": topAbsentees,
			"generated_at":  time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *HRReportHandler) GetOvertimeReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Database is not available")
		return
	}

	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	from, to, err := parseHRDateRange(r.URL.Query().Get("date_from"), r.URL.Query().Get("date_to"), 30)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	page, pageSize := parsePageAndSize(r, 1, 20, 100)
	offset := (page - 1) * pageSize

	var totalHours sql.NullFloat64
	var employeesWithOvertime int64
	if err = h.db.QueryRow(`
		SELECT COALESCE(SUM(overtime_hours), 0), COUNT(DISTINCT employee_id)
		FROM attendance
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND attendance_date BETWEEN $3 AND $4
		  AND overtime_hours > 0
	`, tenantID, restaurantID, from, to).Scan(&totalHours, &employeesWithOvertime); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate overtime report")
		return
	}

	totalOvertimeHours := 0.0
	if totalHours.Valid {
		totalOvertimeHours = totalHours.Float64
	}

	items := make([]map[string]interface{}, 0)
	rows, err := h.db.Query(`
		SELECT
			a.employee_id,
			TRIM(COALESCE(e.first_name, '') || ' ' || COALESCE(e.last_name, '')) AS employee_name,
			COALESCE(NULLIF(TRIM(e.department), ''), 'Unassigned') AS department,
			COALESCE(SUM(a.overtime_hours), 0) AS overtime_hours,
			COALESCE(e.base_salary, 0) AS base_salary,
			COALESCE(e.working_hours_per_week, 40) AS working_hours_per_week
		FROM attendance a
		JOIN employees e
			ON e.id = a.employee_id
			AND e.tenant_id = a.tenant_id
			AND e.restaurant_id = a.restaurant_id
		WHERE a.tenant_id = $1
		  AND a.restaurant_id = $2
		  AND a.attendance_date BETWEEN $3 AND $4
		  AND a.overtime_hours > 0
		GROUP BY a.employee_id, employee_name, department, e.base_salary, e.working_hours_per_week
		ORDER BY overtime_hours DESC, employee_name ASC
		LIMIT $5 OFFSET $6
	`, tenantID, restaurantID, from, to, pageSize, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate overtime report")
		return
	}
	totalOvertimeCost := 0.0
	for rows.Next() {
		var employeeID int64
		var employeeName, department string
		var overtimeHours, baseSalary, workingHoursPerWeek float64
		if scanErr := rows.Scan(&employeeID, &employeeName, &department, &overtimeHours, &baseSalary, &workingHoursPerWeek); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate overtime report")
			return
		}

		monthlyHours := workingHoursPerWeek * 4.33
		hourlyRate := 0.0
		if monthlyHours > 0 {
			hourlyRate = baseSalary / monthlyHours
		}
		overtimeCost := overtimeHours * hourlyRate * 1.5
		totalOvertimeCost += overtimeCost

		items = append(items, map[string]interface{}{
			"employee_id":    employeeID,
			"employee_name":  strings.TrimSpace(employeeName),
			"department":     department,
			"overtime_hours": overtimeHours,
			"hourly_rate":    hourlyRate,
			"overtime_cost":  overtimeCost,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate overtime report")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"filters": map[string]interface{}{
				"date_from": from.Format("2006-01-02"),
				"date_to":   to.Format("2006-01-02"),
			},
			"summary": map[string]interface{}{
				"total_overtime_hours": totalOvertimeHours,
				"employees_count":      employeesWithOvertime,
				"total_overtime_cost":  totalOvertimeCost,
			},
			"pagination": map[string]interface{}{
				"page":      page,
				"page_size": pageSize,
			},
			"items":        items,
			"generated_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *HRReportHandler) GetLeaveBalancesReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Database is not available")
		return
	}

	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	year := time.Now().UTC().Year()
	if yearRaw := strings.TrimSpace(r.URL.Query().Get("year")); yearRaw != "" {
		parsedYear, err := strconv.Atoi(yearRaw)
		if err != nil || parsedYear < 2000 || parsedYear > 3000 {
			respondError(w, http.StatusBadRequest, "Invalid year")
			return
		}
		year = parsedYear
	}
	page, pageSize := parsePageAndSize(r, 1, 50, 200)
	offset := (page - 1) * pageSize

	const annualAllowance = 21.0
	items := make([]map[string]interface{}, 0)
	rows, err := h.db.Query(`
		SELECT
			e.id AS employee_id,
			TRIM(COALESCE(e.first_name, '') || ' ' || COALESCE(e.last_name, '')) AS employee_name,
			COALESCE(NULLIF(TRIM(e.department), ''), 'Unassigned') AS department,
			COALESCE(SUM(CASE WHEN l.status = 'approved' THEN l.total_days ELSE 0 END), 0) AS approved_leave_days,
			COALESCE(SUM(CASE WHEN l.status = 'approved' AND l.leave_type = 'annual' THEN l.total_days ELSE 0 END), 0) AS annual_leave_used
		FROM employees e
		LEFT JOIN leaves l
			ON l.employee_id = e.id
			AND l.tenant_id = e.tenant_id
			AND l.restaurant_id = e.restaurant_id
			AND EXTRACT(YEAR FROM l.start_date) = $3
		WHERE e.tenant_id = $1
		  AND e.restaurant_id = $2
		  AND e.is_active = true
		  AND e.employment_status = 'active'
		GROUP BY e.id, employee_name, department
		ORDER BY employee_name ASC
		LIMIT $4 OFFSET $5
	`, tenantID, restaurantID, year, pageSize, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate leave balances report")
		return
	}

	totalRemainingAnnual := 0.0
	for rows.Next() {
		var employeeID int64
		var employeeName, department string
		var approvedLeaveDays, annualLeaveUsed float64
		if scanErr := rows.Scan(&employeeID, &employeeName, &department, &approvedLeaveDays, &annualLeaveUsed); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate leave balances report")
			return
		}
		remaining := annualAllowance - annualLeaveUsed
		if remaining < 0 {
			remaining = 0
		}
		totalRemainingAnnual += remaining

		items = append(items, map[string]interface{}{
			"employee_id":            employeeID,
			"employee_name":          strings.TrimSpace(employeeName),
			"department":             department,
			"approved_leave_days":    approvedLeaveDays,
			"annual_leave_used":      annualLeaveUsed,
			"annual_leave_allowance": annualAllowance,
			"annual_leave_remaining": remaining,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate leave balances report")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"filters": map[string]interface{}{
				"year": year,
			},
			"summary": map[string]interface{}{
				"employees_count":        len(items),
				"annual_leave_allowance": annualAllowance,
				"total_remaining_annual": totalRemainingAnnual,
			},
			"pagination": map[string]interface{}{
				"page":      page,
				"page_size": pageSize,
			},
			"items":        items,
			"generated_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (h *HRReportHandler) GetTrainingComplianceReport(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		respondError(w, http.StatusServiceUnavailable, "Database is not available")
		return
	}

	tenantID := int(middleware.GetTenantID(r))
	restaurantID := int(middleware.GetRestaurantID(r))
	if tenantID == 0 || restaurantID == 0 {
		respondError(w, http.StatusBadRequest, "Missing tenant context")
		return
	}

	today := time.Now().UTC().Format("2006-01-02")

	var activeEmployees, certifiedEmployees int64
	if err := h.db.QueryRow(`
		SELECT COUNT(*)
		FROM employees
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND is_active = true
		  AND employment_status = 'active'
	`, tenantID, restaurantID).Scan(&activeEmployees); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
		return
	}

	if err := h.db.QueryRow(`
		SELECT COUNT(DISTINCT tr.employee_id)
		FROM training_records tr
		JOIN employees e
			ON e.id = tr.employee_id
			AND e.tenant_id = tr.tenant_id
			AND e.restaurant_id = tr.restaurant_id
		WHERE tr.tenant_id = $1
		  AND tr.restaurant_id = $2
		  AND e.is_active = true
		  AND e.employment_status = 'active'
		  AND tr.status IN ('completed', 'certified')
		  AND (tr.certificate_expiry IS NULL OR tr.certificate_expiry >= $3::date)
	`, tenantID, restaurantID, today).Scan(&certifiedEmployees); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
		return
	}

	complianceRate := 0.0
	if activeEmployees > 0 {
		complianceRate = (float64(certifiedEmployees) / float64(activeEmployees)) * 100
	}

	byDepartment := make([]map[string]interface{}, 0)
	rows, err := h.db.Query(`
		WITH active_employees AS (
			SELECT id, COALESCE(NULLIF(TRIM(department), ''), 'Unassigned') AS department
			FROM employees
			WHERE tenant_id = $1
			  AND restaurant_id = $2
			  AND is_active = true
			  AND employment_status = 'active'
		),
		certified AS (
			SELECT DISTINCT employee_id
			FROM training_records
			WHERE tenant_id = $1
			  AND restaurant_id = $2
			  AND status IN ('completed', 'certified')
			  AND (certificate_expiry IS NULL OR certificate_expiry >= $3::date)
		)
		SELECT
			ae.department,
			COUNT(*) AS total_employees,
			COUNT(c.employee_id) AS certified_employees
		FROM active_employees ae
		LEFT JOIN certified c ON c.employee_id = ae.id
		GROUP BY ae.department
		ORDER BY ae.department ASC
	`, tenantID, restaurantID, today)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
		return
	}
	for rows.Next() {
		var department string
		var totalEmployeesByDept, certifiedByDept int64
		if scanErr := rows.Scan(&department, &totalEmployeesByDept, &certifiedByDept); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
			return
		}

		deptRate := 0.0
		if totalEmployeesByDept > 0 {
			deptRate = (float64(certifiedByDept) / float64(totalEmployeesByDept)) * 100
		}

		byDepartment = append(byDepartment, map[string]interface{}{
			"department":          department,
			"total_employees":     totalEmployeesByDept,
			"certified_employees": certifiedByDept,
			"compliance_rate":     deptRate,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
		return
	}

	byType := make([]map[string]interface{}, 0)
	rows, err = h.db.Query(`
		SELECT
			COALESCE(NULLIF(TRIM(training_type), ''), 'Unspecified') AS training_type,
			COUNT(*) AS total_records,
			SUM(CASE WHEN status IN ('completed', 'certified') THEN 1 ELSE 0 END) AS completed_records
		FROM training_records
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		GROUP BY 1
		ORDER BY total_records DESC, training_type ASC
	`, tenantID, restaurantID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
		return
	}
	for rows.Next() {
		var trainingType string
		var totalRecords, completedRecords int64
		if scanErr := rows.Scan(&trainingType, &totalRecords, &completedRecords); scanErr != nil {
			rows.Close()
			respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
			return
		}
		completionRate := 0.0
		if totalRecords > 0 {
			completionRate = (float64(completedRecords) / float64(totalRecords)) * 100
		}
		byType = append(byType, map[string]interface{}{
			"training_type":     trainingType,
			"total_records":     totalRecords,
			"completed_records": completedRecords,
			"completion_rate":   completionRate,
		})
	}
	if closeErr := rows.Close(); closeErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate training compliance report")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"summary": map[string]interface{}{
				"active_employees":    activeEmployees,
				"certified_employees": certifiedEmployees,
				"compliance_rate":     complianceRate,
			},
			"by_department":    byDepartment,
			"by_training_type": byType,
			"generated_at":     time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func parseHRDateRange(fromRaw, toRaw string, defaultDays int) (time.Time, time.Time, error) {
	if defaultDays <= 0 {
		defaultDays = 30
	}

	now := time.Now().UTC()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	from := to.AddDate(0, 0, -(defaultDays - 1))

	var err error
	if strings.TrimSpace(fromRaw) != "" {
		from, err = time.Parse("2006-01-02", fromRaw)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	if strings.TrimSpace(toRaw) != "" {
		to, err = time.Parse("2006-01-02", toRaw)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	return from.UTC(), to.UTC(), nil
}

func weekdayName(day int) string {
	switch day {
	case 0:
		return "Sunday"
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	default:
		return "Unknown"
	}
}

func parsePageAndSize(r *http.Request, defaultPage, defaultPageSize, maxPageSize int) (int, int) {
	page := defaultPage
	if rawPage := strings.TrimSpace(r.URL.Query().Get("page")); rawPage != "" {
		if parsed, err := strconv.Atoi(rawPage); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := defaultPageSize
	if rawSize := strings.TrimSpace(r.URL.Query().Get("page_size")); rawSize != "" {
		if parsed, err := strconv.Atoi(rawSize); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
