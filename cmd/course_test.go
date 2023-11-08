package cmd

import (
	"fmt"
	"testing"
)

var labMaps = make(map[string]Course)

func init() {
	var phpCourse = Course{
		CourseName: "PHP程序设计",
		Labs: []string{
			"Lab1-PHP开发环境搭建",
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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			name, id, courseName, labs, err := extractStudentNameAndIDAndLabs(tC.givenSubject, tC.givenAttachments, labMaps)
			if err != nil {
				t.Errorf("Expected error %v, but got %v", nil, err)
			} else {
				if name != tC.expected.name {
					t.Errorf("Expected name %s, but got %s", tC.expected.name, name)
				}

				if id != tC.expected.id {
					t.Errorf("Expected id %s, but go %s", tC.expected.id, id)
				}
				if courseName != tC.expected.courseName {
					t.Errorf("Expected courseName %s, but got %s", tC.expected.courseName, courseName)
				}
				if len(labs) != len(tC.expected.labs) {
					t.Errorf("Expected labs %v, but got %v", tC.expected.labs, labs)
				} else {
					for i, lab := range labs {
						if lab != tC.expected.labs[i] {
							t.Errorf("Expected lab %s, but got %s", tC.expected.labs[i], lab)
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
			name:         "Valid input",
			tidy:         "John Doe 12345",
			expectedID:   "12345",
			expectedName: "John Doe",
			expectedErr:  nil,
		},
		{
			name:         "Input with no name",
			tidy:         "12345",
			expectedID:   "12345",
			expectedName: "",
			expectedErr:  fmt.Errorf("illegal input: 12345, lack of name and id"),
		},
		{
			name:         "Input with no name and no ID",
			tidy:         "220301033刘徐明",
			expectedID:   "220301033",
			expectedName: "刘徐明",
			expectedErr:  nil,
		},
		{
			name:         "Input with valid characters",
			tidy:         "李四 12345",
			expectedID:   "12345",
			expectedName: "李四",
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
