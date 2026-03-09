package zhjw

import (
	"bytes"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/W1ndys/easy-qfnu-api-lite/model"
	"github.com/W1ndys/easy-qfnu-api-lite/pkg/logger"
	"go.uber.org/zap"
)

// FetchSelectionResults 抓取并解析成绩，返回包含统计信息的响应
func FetchSelectionResults(cookie string, term string) ([]model.SelectionResult, error) {
	start := time.Now()

	// 使用工厂函数创建 Client (自带检查功能)
	client := NewClient(cookie)

	targetURL := "http://zhjw.qfnu.edu.cn/jsxsd/xkgl/loadXsxkjgList"
	formData := map[string]string{
		"xnxqid": strings.TrimSpace(term), // 学期id

	}

	// 记录重要的业务行为
	log := logger.L().With(
		zap.String("term", term),
		zap.Int("cookie_len", len(cookie)),
	)
	log.Info("开始抓取选课结果")

	// 发起 POST 请求
	resp, err := client.R().
		SetFormData(formData).
		Post(targetURL)

	if err != nil {
		log.Error("抓取选课结果失败",
			zap.Error(err),
			zap.Duration("latency", time.Since(start)),
		)
		return nil, err // 遇到错误立刻返回
	}

	// 解析 HTML (调用内部私有函数)
	selectionResults, err := parseSelectionResultsHtml(resp.Body())
	if err != nil {
		log.Error("解析选课结果失败",
			zap.Error(err),
			zap.Duration("latency", time.Since(start)),
		)
		return nil, err
	}

	log.Info("抓取选课结果完成",
		zap.Int("record_count", len(selectionResults)),
		zap.Duration("latency", time.Since(start)),
	)
	return selectionResults, nil
}

// parseSelectionResultsHtml 解析选课结果 HTML
func parseSelectionResultsHtml(htmlBody []byte) ([]model.SelectionResult, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var schedules []model.SelectionResult = make([]model.SelectionResult, 0)

	doc.Find(".Nsb_r_list tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return // 跳过表头
		}
		tds := s.Find("td")

		// 检查是否是“未查询到数据”的提示行
		if tds.Length() == 1 && strings.Contains(tds.Text(), "未查询到数据") {
			return
		}

		// 正常数据行应该有 10 列
		if tds.Length() < 10 {
			return
		}

		// 组装数据
		// 序号	课程名称	课程编号	上课老师	总学时	学分	课程属性	课程性质	选课操作人	选课时间
		getText := func(i int) string {
			return strings.TrimSpace(tds.Eq(i).Text())
		}
		es := model.SelectionResult{
			Index:      getText(0), // 序号
			CourseName: getText(1), // 课程名称
			CourseId:   getText(2), // 课程编号
			Teacher:    getText(3), // 上课老师
			Hours:      getText(4), // 总学时
			Credit:     getText(5), // 学分
			CourseAttr: getText(6), // 课程属性
			CourseProp: getText(7), // 课程性质
			Operator:   getText(8), // 选课操作人
			SelectTime: getText(9), // 选课时间
		}
		schedules = append(schedules, es)
	})

	return schedules, nil
}
