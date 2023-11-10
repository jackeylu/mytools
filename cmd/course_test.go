package cmd

import (
	"testing"
)

var labMaps = make(map[string]Course)

func init() {
	var phpCourse = Course{
		CourseName: "PHP程序设计",
		Labs: []string{
			"Lab1-PHP开发环境搭建",
			"Lab2-PHP基础知识",
		},
		CourseStudents: []CourseStudent{
			{
				Sno:  "220301093",
				Name: "易思敏",
			},
		},
	}
	var pythonCourse = Course{
		CourseName: "Python程序设计",
		Labs: []string{
			"Lab0-基础语法",
			"Homework1-Python字符串",
		},
		CourseStudents: []CourseStudent{
			{
				Sno:  "230301004",
				Name: "李星雨",
			},
		},
	}
	labMaps["Lab1-PHP开发环境搭建"] = phpCourse
	labMaps["Lab2-PHP基础知识"] = phpCourse
	labMaps["Lab0-基础语法"] = pythonCourse
	labMaps["Homework1-Python字符串"] = pythonCourse
}

type result struct {
	name       string
	id         string
	courseName string
	labName    string
	err        error
}

type multiResult struct {
	name       string
	id         string
	courseName string
	labs       []string
	err        error
}

func TestCleanStudentProjectName(t *testing.T) {
	testCases := []struct {
		desc     string
		given    string
		expected string
	}{
		{
			desc:     "能够处理姓名和实验名之间没有间隔的场景",
			given:    "220301093易思敏Lab1-PHP开发环境搭建",
			expected: "220301093易思敏Lab1-PHP开发环境搭建",
		},
		{
			desc:     "包含中文的句号",
			given:    "220301093易思敏.Lab1-PHP开发环境搭建。doc",
			expected: "220301093易思敏-Lab1-PHP开发环境搭建.doc",
		},
		{
			desc:     "有多个英文.",
			given:    "220301093.易思敏.Lab1-PHP开发环境搭建.doc",
			expected: "220301093-易思敏-Lab1-PHP开发环境搭建.doc",
		},
		{
			desc:     "头尾有空格",
			given:    " 220301093-易思敏-Lab1-PHP开发环境搭建.doc ",
			expected: "220301093-易思敏-Lab1-PHP开发环境搭建.doc",
		},
		{
			desc:     "中间有空格",
			given:    "220301093 易思敏 Lab1-PHP开发环境搭建.doc",
			expected: "220301093-易思敏-Lab1-PHP开发环境搭建.doc",
		},
		{
			desc:     "中间有连续多个空格",
			given:    "220301093 易思敏  Lab1-PHP开发环境搭建.doc",
			expected: "220301093-易思敏-Lab1-PHP开发环境搭建.doc",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := cleanStudentProjectName(tC.given)
			if actual != tC.expected {
				t.Errorf("Expected %v, but got %v", tC.expected, actual)
			}
		})
	}
}

func TestExtractStudentNameAndLabName(t *testing.T) {

	testCases := []struct {
		desc  string
		given string

		expected result
	}{
		{
			desc:  "能够处理姓名和实验名之间没有间隔的场景",
			given: "220301093易思敏Lab1-PHP开发环境搭建",
			expected: result{
				"易思敏",
				"220301093",
				"PHP程序设计",
				"Lab1-PHP开发环境搭建",
				nil,
			},
		},
		{
			desc:  "能够处理实验名中的大小写差异",
			given: "230301004-李星雨-Homework1-python字符串",
			expected: result{
				"李星雨",
				"230301004",
				"Python程序设计",
				"Homework1-Python字符串",
				nil,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			name, id, coureName, labName, err := extractStudentNameAndIDAndLabName(tC.given, labMaps)
			if err != nil {
				t.Errorf("Expected error %v, but got %v", nil, err)
			} else {
				if name != tC.expected.name {
					t.Errorf("Expected name %s, but got %s", tC.expected.name, name)
				}

				if id != tC.expected.id {
					t.Errorf("Expected id %s, but go %s", tC.expected.id, id)
				}

				if coureName != tC.expected.courseName {
					t.Errorf("Expected courseName %s, but got %s", tC.expected.courseName, coureName)
				}

				if labName != tC.expected.labName {
					t.Errorf("Expected labName %s, but got %s", tC.expected.labName, labName)
				}
			}

			t.Logf("name: %s, id: %s, courseName: %s, labName: %s", name, id, coureName, labName)
		})
	}
}

func TestExtractStudentNameAndIDAndLabs(t *testing.T) {
	testCases := []struct {
		desc             string
		givenSubject     string
		givenAttachments []string
		expected         multiResult
	}{
		{
			desc:         "附件中没有个人信息，但在主题中有",
			givenSubject: "230301004-李星雨-Homework1-python字符串",
			givenAttachments: []string{
				"李星雨.zip",
			},
			expected: multiResult{
				"李星雨",
				"230301004",
				"Python程序设计",
				[]string{"Homework1-Python字符串"},
				nil,
			},
		},
		{
			desc:         "附件中只有实验名称，主题中有个人信息",
			givenSubject: "220301053孙焦",
			givenAttachments: []string{
				"Lab1-PHP开发环境搭建.doc",
			},
			expected: multiResult{
				"孙焦",
				"220301053",
				"PHP程序设计",
				[]string{"Lab1-PHP开发环境搭建"},
				nil,
			},
		},
		{
			desc:         "附件或主题中的公共子串在多门课出现：基础",
			givenSubject: "220301101项升杰-PHP基础知识",
			givenAttachments: []string{
				"220301101项升杰-PHP基础知识.doc",
			},
			expected: multiResult{
				"项升杰",
				"220301101",
				"PHP程序设计",
				[]string{"Lab2-PHP基础知识"},
				nil,
			},
		},
		{
			desc:         "姓名学号+Lab数字的场景",
			givenSubject: "朱曼桢220301089-Lab1-PHP开发环境搭建",
			givenAttachments: []string{
				"朱曼桢220301089-Lab1-PHP开发环境搭建.doc",
			},
			expected: multiResult{
				"朱曼桢",
				"220301089",
				"PHP程序设计",
				[]string{"Lab1-PHP开发环境搭建"},
				nil,
			},
		},
		{
			desc:         "只有实验名的场景",
			givenSubject: "Lab1-PHP开发环境搭建",
			givenAttachments: []string{
				"Lab1-PHP开发环境搭建.doc",
			},
			expected: multiResult{
				"",
				"",
				"PHP程序设计",
				[]string{"Lab1-PHP开发环境搭建"},
				nil,
			},
		},
		{
			desc:         "加入了年级、班级干扰",
			givenSubject: "22级信安一班 220301013 陈晓艳",
			givenAttachments: []string{
				"220301013陈晓艳Lab1-PHP开发环境搭建.doc",
			},
			expected: multiResult{
				"",
				"",
				"PHP程序设计",
				[]string{"Lab1-PHP开发环境搭建"},
				nil,
			},
		},
		{
			desc:         "名字前面加入了班级",
			givenSubject: "信安一班张芬220301011",
			givenAttachments: []string{
				"220301011张芬-Lab1-PHP开发环境搭建.doc",
			},
			expected: multiResult{
				"",
				"",
				"PHP程序设计",
				[]string{"Lab1-PHP开发环境搭建"},
				nil,
			},
		},
		{
			desc:         "有多个英文半角点号",
			givenSubject: "信安一班.220301006.刘子萱",
			givenAttachments: []string{
				"刘子萱 220301006-Lab1-PHP开发环境搭建.doc",
			},
			expected: multiResult{
				"",
				"",
				"PHP程序设计",
				[]string{"Lab1-PHP开发环境搭建"},
				nil,
			},
		},
		{
			desc:         "有中文句号的场景",
			givenSubject: "220301104   王凯lab2-php基础知识。",
			givenAttachments: []string{
				"220301104   王凯lab2-php基础知识。.doc",
			},
			expected: multiResult{
				"",
				"",
				"PHP程序设计",
				[]string{"Lab1-PHP开发环境搭建"},
				nil,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			name, id, courseName, labs, err := extractStudentNameAndIDAndLabs(tC.givenSubject, tC.givenAttachments, labMaps)
			if err != nil {
				t.Errorf("Expected error %v, but got %v", nil, err)
			} else {
				if name != tC.expected.name {
					t.Errorf("Expected name =[%s], but got = [%s]", tC.expected.name, name)
				}

				if id != tC.expected.id {
					t.Errorf("Expected id = [%s], but got =[%s]", tC.expected.id, id)
				}
				if courseName != tC.expected.courseName {
					t.Errorf("Expected courseName = [%s], but got =[%s]", tC.expected.courseName, courseName)
				}
				if len(labs) != len(tC.expected.labs) {
					t.Errorf("Expected labs %v, but got %v", tC.expected.labs, labs)
				} else {
					for i, lab := range labs {
						if lab != tC.expected.labs[i] {
							t.Errorf("Expected lab =[%s], but got =[%s]", tC.expected.labs[i], lab)
						}
					}
				}
			}
		})
	}
}

func TestExtractStudentNameAndIDRegex(t *testing.T) {
	tests := []struct {
		name         string
		tidy         string
		expectedID   string
		expectedName string
		expectedErr  error
	}{
		{
			name:         "Input with valid characters",
			tidy:         "李四 123456",
			expectedID:   "123456",
			expectedName: "李四",
			expectedErr:  nil,
		},
		{
			name:         "Input with no name and no ID",
			tidy:         "220301033刘徐明",
			expectedID:   "220301033",
			expectedName: "刘徐明",
			expectedErr:  nil,
		},
	}

	for _, test := range tests {
		name, id, err := extractStudentNameAndIDRegex(test.tidy)

		if err != nil && err != test.expectedErr {
			t.Errorf("Expected error %v, but got %v", test.expectedErr, err)
		} else {
			if id != test.expectedID {
				t.Errorf("Expected ID %s, but got %s", test.expectedID, id)
			}

			if name != test.expectedName {
				t.Errorf("Expected name %s, but got %s", test.expectedName, name)
			}
		}
	}
}
