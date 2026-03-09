package zhjw

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/W1ndys/easy-qfnu-api-lite/model"
	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"go.uber.org/zap"
)

// FetchGrades 抓取并解析成绩，返回包含统计信息的响应
func FetchGrades(cookie string, term string, courseType string, courseName string, displayType string) (*model.GradeResponse, error) {
	start := time.Now()

	// 课程类型：支持中文名称或ID，统一转换为ID
	courseType = model.GetCourseTypeID(courseType)

	// 使用工厂函数创建 Client (自带检查功能)
	client := NewClient(cookie)

	targetURL := "http://zhjw.qfnu.edu.cn/jsxsd/kscj/cjcx_list"
	formData := map[string]string{
		"kksj": strings.TrimSpace(term),        // 开课时间
		"kcxz": strings.TrimSpace(courseType),  // 课程性质
		"kcmc": strings.TrimSpace(courseName),  // 课程名称
		"xsfs": strings.TrimSpace(displayType), // 显示方式
	}

	// 记录重要的业务行为
	log := logger.L().With(
		zap.String("term", term),
		zap.String("course_name", courseName),
		zap.String("course_type", courseType),
		zap.String("display_type", displayType),
		zap.Int("cookie_len", len(cookie)),
	)
	log.Info("开始抓取成绩")

	// 发起 POST 请求
	resp, err := client.R().
		SetFormData(formData).
		Post(targetURL)

	if err != nil {
		log.Error("抓取成绩失败",
			zap.Error(err),
			zap.Duration("latency", time.Since(start)),
		)
		return nil, err // 遇到错误立刻返回
	}

	// 解析 HTML (调用内部私有函数)
	grades, err := parseGradesHtml(resp.Body())
	if err != nil {
		log.Error("解析成绩失败",
			zap.Error(err),
			zap.Duration("latency", time.Since(start)),
		)
		return nil, err
	}

	// 计算统计信息
	response := calculateStats(grades)
	log.Info("抓取成绩完成",
		zap.Int("record_count", len(grades)),
		zap.Duration("latency", time.Since(start)),
	)
	return response, nil
}

// calculateStats 计算成绩统计信息
func calculateStats(grades []model.Grade) *model.GradeResponse {
	response := &model.GradeResponse{
		Grades:        grades,
		YearStats:     []model.YearStat{},
		SemesterStats: []model.SemesterStat{},
	}

	// 按学期分组
	semesterMap := make(map[string][]model.Grade)
	// 按学年分组
	yearMap := make(map[string][]model.Grade)

	for _, g := range grades {
		semester := g.Semester
		semesterMap[semester] = append(semesterMap[semester], g)

		// 提取学年 (如 "2023-2024-1" -> "2023-2024")
		parts := strings.Split(semester, "-")
		if len(parts) >= 2 {
			year := parts[0] + "-" + parts[1]
			yearMap[year] = append(yearMap[year], g)
		}
	}

	// 计算每学期统计
	var semesters []string
	for s := range semesterMap {
		semesters = append(semesters, s)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(semesters))) // 按学期倒序

	for _, semester := range semesters {
		stat := calculateGradeStat(semesterMap[semester])
		response.SemesterStats = append(response.SemesterStats, model.SemesterStat{
			Semester: semester,
			Stat:     stat,
		})
	}

	// 计算每学年统计
	var years []string
	for y := range yearMap {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years))) // 按学年倒序

	for _, year := range years {
		stat := calculateGradeStat(yearMap[year])
		response.YearStats = append(response.YearStats, model.YearStat{
			Year: year,
			Stat: stat,
		})
	}

	// 计算总体统计
	response.TotalStat = calculateGradeStat(grades)

	return response
}

// calculateGradeStat 计算一组成绩的加权平均绩点和总学分
func calculateGradeStat(grades []model.Grade) model.GradeStat {
	var totalCredits float64
	var weightedSum float64
	var validCourseCount int

	for _, g := range grades {
		credit, err := strconv.ParseFloat(g.Credit, 64)
		if err != nil || credit <= 0 {
			continue
		}
		totalCredits += credit

		gpa, err := strconv.ParseFloat(g.GPA, 64)
		if err != nil || gpa < 0 {
			continue
		}
		weightedSum += gpa * credit
		validCourseCount++
	}

	var weightedGPA float64
	if totalCredits > 0 {
		weightedGPA = weightedSum / totalCredits
	}

	return model.GradeStat{
		WeightedGPA:  round2(weightedGPA),
		TotalCredits: round2(totalCredits),
		CourseCount:  len(grades),
	}
}

// round2 保留两位小数
func round2(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}

// parseGradesHtml 是私有函数(小写p)，只在这个文件内部使用，外部不需要知道解析细节
func parseGradesHtml(htmlBody []byte) ([]model.Grade, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var grades []model.Grade // 使用 model.Grade

	doc.Find("#dataList tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		} // 跳过表头
		tds := s.Find("td")
		if tds.Length() < 10 {
			return
		}

		// 组装数据
		g := model.Grade{
			Semester:   strings.TrimSpace(tds.Eq(1).Text()),
			CourseCode: strings.TrimSpace(tds.Eq(2).Text()),
			CourseName: strings.TrimSpace(tds.Eq(3).Text()),
			Score:      strings.TrimSpace(tds.Eq(5).Text()),
			Credit:     strings.TrimSpace(tds.Eq(7).Text()),
			GPA:        strings.TrimSpace(tds.Eq(9).Text()),
			ExamType:   strings.TrimSpace(tds.Eq(11).Text()),
			CourseProp: strings.TrimSpace(tds.Eq(14).Text()),
		}
		grades = append(grades, g)
	})

	if len(grades) == 0 {
		// 语法点 4: 自定义错误
		return nil, fmt.Errorf("解析结果为空，可能是Cookie失效或页面结构变更")
	}

	return grades, nil
}
