/*
Copyright © 2023 Lyu Lin <lvlin@whu.edu.cn>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	emailFile string
)

// Course is the course settings from configuration file
type Course struct {
	// CourseName is the name of the course
	CourseName string
	// Labs is the list of labs
	Labs []string
	// CourseStudents is the list of students read from the configuration file
	CourseStudents []CourseStudent
}

// emailResult is the result of the email with course information
type emailResult struct {
	// StudentName is the name of the student
	StudentName string
	// StudentID is the ID of the student
	StudentID string
	// Course is the course name
	Course string
	// Lab is the lab name
	Lab string
	// Time is the time of submission
	Time string
	// Email is the email address of the student
	Email string
	// Subject is the subject of the email
	Subject string
	// Attachment is the attachment of the email
	Attachment string
	// Notes is the notes of the processing result
	Notes string
}

// courseCmd represents the course command
var courseCmd = &cobra.Command{
	Use:   "course",
	Short: "为拉取的邮件，从邮件主题和附件命名中寻找可能的学生姓名、学号和课程",
	Long: `为email命令的输出文件，从邮件主题和附件命名中寻找可能的学生姓名、学号和课程。

The "course" command is used to find the student name, student ID, and course from the email.

Example:

$ email course -f emails.xlsx

The output will be like:

姓名	学号	课程   实验名   提交时间  提交人邮件地址 邮件主题  附件名

to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 设置日志文件的格式
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
		// 创建一个 LoggerWriter 对象
		logger := util.NewLoggerWriter("logfile.txt")
		defer logger.Close()
		// 将日志同时输出到终端和日志文件
		log.SetOutput(logger)

		var emails []EmailInfo
		if err := readAttachmentEmailFromFetchedEmailFile(emailFile, &emails); err != nil {
			log.Println(err)
			return err
		}
		//fmt.Println(result)
		var courseInfo []Course
		if err := readCourseFile(&courseInfo); err != nil {
			return err
		}
		var labsMap = buildLabsMap(courseInfo)
		// fmt.Println(courseInfo)
		// 补充完整学生和课程信息
		result := updateEmailResultWithCourseInfo(emails, labsMap)

		// fmt.Println(result)
		saveResult(result)

		return nil
	},
}

func init() {
	emailCmd.AddCommand(courseCmd)
	courseCmd.Flags().StringVarP(&emailFile, "file", "f", "email.xlsx", "the fetched email file by email command")
}

func saveResult(result []emailResult) {
	header := []string{"姓名", "学号", "课程", "实验名", "提交时间", "提交人邮件地址", "邮件主题", "附件名", "备注"}
	columns := make([][]string, len(result))
	for i, v := range result {
		columns[i] = []string{v.StudentName, v.StudentID, v.Course, v.Lab, v.Time, v.Email, v.Subject, v.Attachment, v.Notes}
	}
	util.WriteExcelFile("email_course.xlsx", header, columns)
}

// readAttachmentEmailFromFetchedEmailFile reads the fetched email file, and build the preliminary result
func readAttachmentEmailFromFetchedEmailFile(emailFile string, emails *[]EmailInfo) error {
	util.ReadExcelFile(emailFile, func(row int, columns []string) error {
		if row == 0 {
			if reflect.DeepEqual(columns, ExcelFileHeader()) {
				return nil
			} else {
				return fmt.Errorf("邮件标题不匹配，应当是%v", ExcelFileHeader())
			}
		}
		// ignore email without any attachment
		if len(columns) < 6 {
			return nil
		}
		date := DecodeTime(columns[1])
		num, err := strconv.ParseUint(columns[0], 10, 32)
		if err != nil {
			return err
		}
		// handle the contents
		*emails = append(*emails, EmailInfo{
			SeqNum:      uint32(num),
			Date:        date,
			From:        columns[2],
			To:          strings.Split(columns[3], ","),
			Subject:     columns[4],
			Attachments: DecodeAttachments(columns[5]),
		})
		return nil
	}, false)
	return nil
}

func readCourseFile(courseInfo *[]Course) error {
	dataset := viper.GetStringMap("course")
	for _, value := range dataset {
		// fmt.Println(key, value)
		if reflect.TypeOf(value).Kind() != reflect.Map {
			return fmt.Errorf("课程配置错误，应该是map")
		}
		value := value.(map[string]interface{})
		course := Course{
			CourseName: value["name"].(string),
		}
		labs := value["labs"].([]interface{})
		for _, lab := range labs {
			if reflect.TypeOf(lab).Kind() != reflect.String {
				return fmt.Errorf("实验配置错误，应该是string")
			}
			course.Labs = append(course.Labs, lab.(string))
		}
		classes := value["classes"].([]interface{})
		for _, class := range classes {
			if reflect.TypeOf(class).Kind() != reflect.String {
				return fmt.Errorf("班级配置错误，应该是string")
			}
			course.CourseStudents = append(course.CourseStudents, ReadNameList(class.(string))...)
		}

		*courseInfo = append(*courseInfo, course)
	}

	return nil
}

func updateEmailResultWithCourseInfo(emails []EmailInfo, labsMap map[string]Course) []emailResult {
	var ans []emailResult
	for _, email := range emails {
		results := findAndBuildResults(email, labsMap)
		ans = append(ans, results...)
	}
	return ans
}

func findAndBuildResults(email EmailInfo, labsMap map[string]Course) []emailResult {
	name, id, courseName, labs, err := extractStudentNameAndIDAndLabs(email.Subject, email.Attachments, labsMap)
	if err != nil {
		log.Printf("Failed to find student name and ID from email %v: %v\n", email, err)
		return []emailResult{
			{
				Time:       email.Date.Format("2006-01-02 15:04:05"),
				Email:      email.From,
				Subject:    email.Subject,
				Attachment: EncodeAttachments(email.Attachments),
				Notes:      "Failed",
			},
		}
	}
	var results []emailResult
	for _, lab := range labs {
		// TODO validate the sno and name and labname
		results = append(results, emailResult{
			StudentName: name,
			StudentID:   id,
			Course:      courseName,
			Lab:         lab,
			Time:        email.Date.Format("2006-01-02 15:04:05"),
			Email:       email.From,
			Subject:     email.Subject,
			Attachment:  EncodeAttachments(email.Attachments),
			Notes:       "Success",
		},
		)
	}

	return results
}

func findLab(subjectOrFilename string, labsMap map[string]Course) (fullLabName, removedLabName, courseName string, err error) {
	fullLabName, removedLabName, courseName, err = "", "", "", nil
	longestSubstr := ""
	for key, course := range labsMap {
		if strings.Contains(strings.ToUpper(subjectOrFilename), strings.ToUpper(key)) {
			fullLabName = key
			removedLabName = strings.Replace(strings.ToUpper(subjectOrFilename), strings.ToUpper(key), "", 1)
			courseName = course.CourseName
			return
		} else {
			// try with longest common substring
			substr := util.LongestCommonSubstr(strings.ToUpper(subjectOrFilename), strings.ToUpper(key))
			if substr != "" && len(substr) > 4 && len(substr) > len(longestSubstr) {
				fullLabName = key
				removedLabName = strings.Replace(strings.ToUpper(subjectOrFilename), strings.ToUpper(substr), "", 1)
				longestSubstr = substr
				courseName = course.CourseName
				// return
			}
		}
	}
	if fullLabName == "" {
		err = fmt.Errorf("邮件主题/附件名称中未找到实验名")
	}
	return
}

func extractStudentNameAndIDAndLabName(subjectOrAttachment string, labsMap map[string]Course) (
	name, id, courseName, lab string, err error) {
	name, id, err0 := extractStudentNameAndID(subjectOrAttachment)
	// 根据email中的subject和attachments字段，
	lab, removedLabName, courseName, err := findLab(subjectOrAttachment, labsMap)
	if err != nil {
		return
	}
	if err0 != nil || len(name) > len("四字姓名") {
		// 未找到学生信息，将课程实验名去掉再找找看
		name, id, err = extractStudentNameAndID(removedLabName)
	}
	return
}

func cleanStudentProjectName(projectName string) string {
	ans := projectName
	// 删除可能的中文句号
	ans = strings.Replace(ans, "。", "", -1)
	// 替换n个英文句号中的n-1个为'-'符号，保留最后一个英文句号
	counts := strings.Count(ans, ".")
	if counts > 1 {
		ans = strings.Replace(ans, ".", "-", counts-1)
	}
	// 移除头尾的空格
	ans = strings.TrimSpace(ans)
	// 将所有的下划线替换成'-'
	ans = strings.ReplaceAll(ans, "_", "-")
	// 遍历字符串，将遇到的单个空格或连续多个空格替换成一个'-'
	temp := ans
	ans = ""
	for _, char := range temp {
		if char == ' ' {
			if ans[len(ans)-1] != '-' {
				ans = ans + "-"
			}
		} else {
			ans = ans + string(char)
		}
	}
	// 去掉头尾的'-'
	ans = strings.Trim(ans, "-")
	// 去掉中间的连续多个'-'
	temp = ans
	ans = ""
	for _, char := range temp {
		if char == '-' {
			if ans[len(ans)-1] != '-' {
				ans = ans + string(char)
			}
		} else {
			ans = ans + string(char)
		}
	}
	// 去掉字符串中可能存在的圆括号及圆括号内的内容
	for {
		if strings.Contains(ans, "(") && strings.Contains(ans, ")") {
			// 获取第一个左括号的索引和第一个右括号的索引
			// 如果左括号在右括号之前，则将左括号和右括号之间的内容都去掉
			leftIndex := strings.Index(ans, "(")
			rightIndex := strings.Index(ans, ")")
			if leftIndex < rightIndex {
				ans = strings.Replace(ans, ans[leftIndex:rightIndex+1], "", 1)
			}
		} else {
			break
		}
	}

	return ans
}

func extractStudentNameAndIDAndLabs(subject string, attachments []string, labsMap map[string]Course) (
	name, id, courseName string, labs []string, err error) {
	subject = cleanStudentProjectName(subject)
	for i, v := range attachments {
		attachments[i] = cleanStudentProjectName(v)
	}

	if len(subject) == 0 {
		err = fmt.Errorf("no subject found")
		return
	}
	if len(attachments) == 0 {
		err = fmt.Errorf("no attachments found")
		return
	}
	name, id, courseName, labs, err = "", "", "", []string{}, nil

	name, id, courseName, lab, err := extractStudentNameAndIDAndLabName(subject, labsMap)
	if err != nil {
		// 主题中未找到实验名
		for _, attachment := range attachments {
			name_temp, id_temp, courseName0, lab, err0 := extractStudentNameAndIDAndLabName(attachment, labsMap)
			if name == "" && name_temp != "" {
				name = name_temp
			}
			if id == "" && id_temp != "" {
				id = id_temp
			}
			if courseName == "" && courseName0 != "" {
				courseName = courseName0
			}
			if lab != "" {
				labs = append(labs, lab)
			}
			if err0 != nil && (name == "" || id == "" || lab == "") {
				log.Printf("Failed to find student name/ID and lab name from attachment %s: %v\n", attachment, err0)
			}
		}
		err = nil
		labs = removeDuplicate(labs)
		return
	} else {
		// 主题中找到了课程名、实验名、姓名、学号
		labs = append(labs, lab)
		// 遍历每个实验报告，查找实验名称
		for _, attachment := range attachments {
			if lab, _, _, err = findLab(attachment, labsMap); err == nil {
				labs = append(labs, lab)
			} else {
				log.Printf("Unable finding lab name from attachment %s: %v. But its's ok as it was found in the subject.\n", attachment, err)
				// 有记录就不失败
				err = nil
				return
			}
		}
		labs = removeDuplicate(labs)
		return
	}
}

func removeDuplicate(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// extractStudentNameAndID extracts student name and id from a string removed the labname
func extractStudentNameAndID(s string) (name, id string, err error) {
	s = cleanStudentProjectName(s)
	fields := strings.Split(s, "-")
	name, id, err = "", "", nil
	for i, v := range fields {
		v = strings.TrimSpace(v)
		// TODO 一个个field检查，是否是学号、姓名
		if v == "" {
			fields = append(fields[:i], fields[i+1:]...)
		}
	}
	if len(fields) < 2 {
		// try with regular expression
		if name, id, err = extractStudentNameAndIDRegex(s); err == nil {
			return
		}
		err = fmt.Errorf("illegal input: %s, lack of name and id", s)
		return
	} else {
		name, id = fields[0], fields[1]

		// id 不都是数字组成，name 都是数字组成，则进行交换
		if !util.IsAllCharacterDigit(id) && util.IsAllCharacterDigit(name) {
			name, id = id, name
		}
		if !util.IsAllCharacterDigit(id) || len(id) < 6 {
			name, id = "", ""
			err = fmt.Errorf("illegal input: %s, lack of name and id", s)
		}
		return
	}
}

func extractStudentNameAndIDRegex(tidy string) (name, id string, err error) {
	s, e := -1, -1
	for i, v := range tidy {
		if v >= '0' && v <= '9' {
			if s == -1 {
				s = i
			} else if e == -1 || e+1 == i {
				e = i
			} else {
				err = fmt.Errorf("illegal input: %s, lack of name and id", tidy)
				return
			}
		}
	}
	// 获得的学号长度不够
	if s == -1 || e == -1 || (e-s) < 5 {
		err = fmt.Errorf("illegal input: %s, lack of name and id", tidy)
		return
	}
	id = tidy[s : e+1]
	name = strings.TrimSpace(strings.ReplaceAll(tidy, id, ""))
	return
}

func buildLabsMap(courseInfo []Course) map[string]Course {
	var labsMap = make(map[string]Course)
	for _, course := range courseInfo {
		for _, lab := range course.Labs {
			labsMap[lab] = course
		}
	}
	return labsMap
}
