package api

import (
	"fmt"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"sort"

)

const dateFormat = "20060102"

// Проверяем чтобы дата была после now
func afterNow(now, date time.Time) bool {
	dateTrunc := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowTrunc := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return dateTrunc.After(nowTrunc)
}

// Вычисляем следующую дату по правилу повторения
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	//Возвращаем ошибку если repeat пустая строка
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	// Получаем date и возвращаем ошибку если дата неверная
	start, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", errors.New("invalid date format in dateStr")
	}

	// Разбиваем repeat на cостовляющие
	parts := strings.Split(repeat, " ")
	ruleType := parts[0]

	switch ruleType {
	// d - переносим на указанное число дней
	case "d":
	if len(parts) != 2 {
		return "", errors.New("invalid format: expected 'd <number>'")
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("invalid number")
	}
	if interval <= 0 || interval > 400 {
		return "", errors.New("the interval should be between 1 and 400")
	}

	// Ищем результат с начальной даты
	date := start
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(now, date) {
			break
		}
	}
	// Преобразуем и возвращаем дату в формате "20060102"
	return date.Format(dateFormat), nil

	case "y":
	// y - выполняем задачу ежегодно
	if len(parts) != 1 {
		return "", errors.New("invalid format: expected just 'y'")
	}

	date := start
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(now, date) {
			break
	}
	}
		return date.Format(dateFormat), nil

	// w - выполняем задачу в указанные дни недели
	case "w":
	if len(parts) != 2 {
		return "", errors.New("invalid format for w rule: expected 'w <days>'")
	}

	// Разделяем дни запятыми
	dayStrs := strings.Split(parts[1], ",")
	targetDays := make([]int, 0)

	for _, s := range dayStrs {
		day, err := strconv.Atoi(s)
		if err != nil || day < 1 || day > 7 {
			return "", errors.New("invalid day: must be 1–7")
	}
	targetDays = append(targetDays, day)
	}

	// Проверяем дни с начальной даты
	date := start
	for {
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7 
	}

	// Ищем подходящую дату
	for _, target := range targetDays {
		if weekday == target && afterNow(now, date) {
			return date.Format(dateFormat), nil
		}
	}

	// Переходим на следующий день
	date = date.AddDate(0, 0, 1)
	}

	//m - выполняем задачу в указанные дни месяца
	case "m":
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("invalid format: expected 'm <days> [months]'")
	}

	// Парсим дни месяца
	dayStrs := strings.Split(parts[1], ",")
	targetDays := make([]int, 0)
	for _, s := range dayStrs {
		day, err := strconv.Atoi(s)
			if err != nil || day < -2 || day > 31 || day == 0 {
				return "", errors.New("invalid day: must be -2, -1 or 1–31")
		}
	targetDays = append(targetDays, day)
	}

	// Парсим месяцы
	var targetMonths []int
	if len(parts) == 3 {
		monthStrs := strings.Split(parts[2], ",")
	for _, s := range monthStrs {
		month, err := strconv.Atoi(s)
		if err != nil || month < 1 || month > 12 {
			return "", errors.New("invalid month: must be 1–12")
	}
	targetMonths = append(targetMonths, month)
	}
	} else {
		targetMonths = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	}

	// Рассчитываем последний день месяца
	lastDayOfMonth := func(year, month int) int {
	nextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)
	last := nextMonth.AddDate(0, 0, -1)
		return last.Day()
	}

	// Начинаем с начальной даты
	date := start
	for {
		year := date.Year()
		month := int(date.Month())

	// Пропускаем месяцы, не входящие в targetMonths
	skip := true
	for _, m := range targetMonths {
			if m == month {
		skip = false
			break
	}
	}
	if skip {
		date = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, date.Location())
			continue
	}

	// Собираем все подходящие даты в текущем месяце
	suitableDate := make([]time.Time, 0)
	for _, day := range targetDays {
	var candidate time.Time
		if day > 0 {
		candidate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, date.Location())
	} else if day == -1 {
		last := lastDayOfMonth(year, month)
		candidate = time.Date(year, time.Month(month), last, 0, 0, 0, 0, date.Location())
	} else if day == -2 {
		last := lastDayOfMonth(year, month)
		candidate = time.Date(year, time.Month(month), last-1, 0, 0, 0, 0, date.Location())
	}

	// Проверяем корректность даты
	if candidate.Month() != time.Month(month) {
		continue
	}
	suitableDate = append(suitableDate, candidate)
	}

	// Сортируем даты по возрастанию
	sort.Slice(suitableDate, func(i, j int) bool {
		return suitableDate[i].Before(suitableDate[j])
	})

	// Ищем первую дату после now
	for _, candidate := range suitableDate {
	if afterNow(now, candidate) {
		return candidate.Format(dateFormat), nil
	}
	}

	// Если в текущем месяце подходящей даты нет, то переходим к следующему
	date = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, date.Location())
	}


	default:
		return "", errors.New("unsupported format")
	}
}

	
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	//Получаем Get-параметры
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	//Проверяем параметр date
	if dateStr == "" {
		http.Error(w, "parameter 'date' is required", http.StatusBadRequest)
			return
	}
	//Проверяем параметр repeat
	if repeat == "" {
		http.Error(w, "parameter 'repeat' is required", http.StatusBadRequest)
			return
	}

	//Если параметр now передан, то парсим его в time.Time
	var now time.Time
	if nowStr != "" {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid 'now' date format: expected YYYYMMDD", http.StatusBadRequest)
			return
	}
	} else {
		now = time.Now()
	}

	//Вызываем функцию nextDate для расчёта следующей даты
	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		log.Printf("error calculating next date: %v (date: %s, repeat: %s)", err, dateStr, repeat)
		http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}
	
	//Выводим ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = fmt.Fprintf(w, "%s", nextDate)
	if err != nil {
		log.Printf("error writing text response: %v", err)
	}
}
