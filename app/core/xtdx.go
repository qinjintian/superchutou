package core

import (
	"encoding/json"
	"github.com/qinjintian/superchutou/pkg/net"
	"github.com/qinjintian/superchutou/pkg/net/http"
	"io/ioutil"
	pkgHttp "net/http"
	"strings"
)

type XTDXService struct {
	ip        string
	cookie    *pkgHttp.Cookie
	cookieStr string
}

const (
	UA = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36"
)

func NewXTDXService() (*XTDXService, error) {
	ip, err := net.GetExtranetIP()
	if err != nil {
		return nil, err
	}

	return &XTDXService{
		ip: ip,
	}, nil
}

// BindStudentLoginByCardNumber 通过卡号绑定学生登录
func (s *XTDXService) BindStudentLoginByCardNumber(card, pwd string) ([]byte, error) {
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = UA

	params := make(map[string]interface{}, 0)
	params["card"] = card
	params["password"] = pwd
	params["IsOauth"] = 0
	params["IsEncryptPasword"] = 1
	params["Specialty_ID"] = ""
	params["notVerifyPhone"] = true
	params["CardNumber"] = card
	params["Password"] = pwd

	b, _ := json.Marshal(params)

	_, resp, err := http.Request("POST", "https://xtdx.web2.superchutou.com/service/eduSuper/Student/BindStudentLoginByCardNumber", strings.NewReader(string(b)), headers)
	if err != nil {
		return nil, err
	}

	s.cookie = resp.Cookies()[0]

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetStudentDetailRegisterSet 获取学生登录资料
func (s *XTDXService) GetStudentDetailRegisterSet(cookie string) ([]byte, error) {
	headers := make(map[string]string, 0)
	//headers["Content-Type"] = "application/json"
	headers["User-Agent"] = UA
	headers["Cookie"] = cookie
	s.cookieStr = cookie
	params := make(map[string]interface{}, 0)
	b, _ := json.Marshal(params)
	_, resp, err := http.Request("GET", "https://xtdx.web2.superchutou.com/service/eduSuper/StudentinfoDetail/GetStudentDetailRegisterSet", strings.NewReader(string(b)), headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetStuSpecialtyCurriculumList 获取学生专业课程列表
func (s *XTDXService) GetStuSpecialtyCurriculumList(url string) ([]byte, error) {
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = UA
	headers["Cookie"] = s.cookieStr

	_, resp, err := http.Request("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GetCourseChaptersNodeList 获取课程章节节点列表
func (s *XTDXService) GetCourseChaptersNodeList(url string) ([]byte, error) {
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = UA
	headers["Cookie"] = s.cookieStr

	_, resp, err := http.Request("GET", url, nil, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// SaveCourseLook 更新视频进度
func (s *XTDXService) SaveCourseLook(courseId uint64) ([]byte, error) {
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json"
	headers["User-Agent"] = UA
	headers["Cookie"] = s.cookieStr

	params := make(map[string]interface{}, 0)
	params["CourseChapters_ID"] = courseId
	params["IP"] = s.ip
	params["LookTime"] = 60
	params["LookType"] = 0

	b, _ := json.Marshal(params)

	_, resp, err := http.Request("POST", "https://xtdx.web2.superchutou.com/service/datastore/WebCourse/SaveCourse_Look", strings.NewReader(string(b)), headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
