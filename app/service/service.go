package service

import (
	"bufio"
	"fmt"
	"github.com/qinjintian/superchutou/app/core"
	"github.com/qinjintian/superchutou/pkg/utils"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Service struct {
	cfg         *Config
	xdSvc       *core.XTDXService
	stuId       string                 // 学生ID
	stuDetailId string                 // 学生详情ID
	curriculums map[uint64]*curriculum // 课程信息
	wg          sync.WaitGroup
}

// curriculum 课程信息
type curriculum struct {
	courseId           uint64
	curriculumId       uint64
	cuName             string // 课程名
	courseChapters     uint64 // 课程集数
	courseReadChapters uint64 // 课程已看集数
}

func NewService(cfg *Config, xdSvc *core.XTDXService) (*Service, error) {
	return &Service{
		cfg:         cfg,
		xdSvc:       xdSvc,
		curriculums: make(map[uint64]*curriculum, 0),
	}, nil
}

func (s *Service) Run() {
	result := s.Login() // 登录

	log.Println(fmt.Sprintf("登录成功，^_^欢迎%s，开始获取专业课程列表~~~~", result.Get("Data.Name")))

	time.Sleep(2 * time.Second)

	// 获取学生专业课程列表
	result = s.GetCurriculums(result.Get("Data.0.StuDetail_ID").String(), result.Get("Data.0.StuID").String())
	var (
		arrs  = result.Get("Data.list").Array()
		count = result.Get("TotalCount").Uint()
	)

	for key, arr := range arrs {
		var (
			courseReadChapters = arr.Get("CourseReadChapters").Uint()
			courseChapters     = arr.Get("CourseChapters").Uint()
		)

		s.curriculums[arr.Get("Curriculum_ID").Uint()] = &curriculum{
			courseId:           arr.Get("Course_ID").Uint(),
			curriculumId:       arr.Get("Curriculum_ID").Uint(),
			cuName:             arr.Get("CuName").String(),
			courseChapters:     arr.Get("CourseChapters").Uint(),
			courseReadChapters: arr.Get("CourseReadChapters").Uint(),
		}

		// 跳过还未可以看的课程
		if courseChapters == 0 {
			log.Println(fmt.Sprintf("[%d/%d] 学期： 第%d学期 | 课程ID： %d | 课程名称： %s | 进度： 已观看数: 0 总数: 0 进度: 0%% | 状态： 未开放", key+1, count, arr.Get("StudyYear").Uint(), arr.Get("Curriculum_ID").Uint(), arr.Get("CuName").String()))
			continue
		}

		rate := utils.Decimal(float64(courseReadChapters) / float64(courseChapters) * 100)

		log.Println(fmt.Sprintf("[%d/%d] 学期： 第%d学期 | 课程ID： %d | 课程名称： %s | 进度： 已观看数: %d 总数: %d 进度: %v%% | 状态： 已开放", key+1, count, arr.Get("StudyYear").Uint(), arr.Get("Curriculum_ID").Uint(), arr.Get("CuName").String(), courseReadChapters, courseChapters, rate))
	}

	scanner := bufio.NewScanner(os.Stdin)

	log.Println("请输入要观看的课程ID，多个ID请用隔空隔开，格式：1100 1122 1133")

ReEnter:
	curriculumIds := make([]string, 0)

	for {
		scanner.Scan()
		curriculumIdStr := scanner.Text()

		if curriculumIdStr == "" {
			log.Println("（。・＿・。）课程ID不能为空，请重新输入~")
			continue
		}

		curriculumIds = strings.Split(curriculumIdStr, " ")
		for key, curriculumId := range curriculumIds {
			id, err := strconv.ParseUint(curriculumId, 10, 64)
			if err != nil {
				log.Println(fmt.Sprintf("（。・＿・。）您输入的第%d个课程ID不正确，请重新输入~", key+1))
				goto ReEnter
			}

			if _, ok := s.curriculums[id]; !ok {
				log.Println(fmt.Sprintf("（。・＿・。）您输入的第%d个课程ID不正确，请重新输入~", key+1))
				goto ReEnter
			}

			break
		}

		break
	}

	for _, curriculumId := range curriculumIds {
		id, _ := strconv.ParseUint(curriculumId, 10, 64)
		c := s.curriculums[id]

		s.wg.Add(1)

		go s.HandleCourseChapters(c)
	}

	s.wg.Wait()

	log.Println("所选课程已经全部播放完毕，请确认进度(*^__^*)")
}

// Login 登录
func (s *Service) Login() gjson.Result {
	data, err := s.xdSvc.GetStudentDetailRegisterSet(s.cfg.Authenticate.Cookie)
	if err != nil {
		log.Fatalln(err)
	}
	result := gjson.ParseBytes(data)
	if !result.Get("SuccessResponse").Bool() {
		log.Fatalln(result.Get("Message").String())
	}
	s.stuId = result.Get("Data.0.StuID").String()
	s.stuDetailId = result.Get("Data.0.StuDetail_ID").String()

	return result
}

// GetCurriculums 获取学生专业课程列表
func (s *Service) GetCurriculums(stuDetailId, stuId string) gjson.Result {
	data, err := s.xdSvc.GetStuSpecialtyCurriculumList(fmt.Sprintf("https://xtdx.web2.superchutou.com/service/eduSuper/Specialty/GetStuSpecialtyCurriculumList?StuDetail_ID=%s&IsStudyYear=1&StuID=%s", stuDetailId, stuId))
	if err != nil {
		log.Fatalln(err)
	}

	result := gjson.ParseBytes(data)
	if !result.Get("SuccessResponse").Bool() {
		log.Fatalln(result.Get("Message").String())
	}

	return result
}

// GetCourseChapters 获取课程章节节点列表
func (s *Service) GetCourseChapters(courseId, curriculumId uint64, stuId, stuDetailId string) gjson.Result {
	data, err := s.xdSvc.GetCourseChaptersNodeList(fmt.Sprintf("https://xtdx.web2.superchutou.com/service/eduSuper/Question/GetCourse_ChaptersNodeList?Valid=1&Course_ID=%d&StuID=%s&Curriculum_ID=%d&Examination_ID=0&StuDetail_ID=%s", courseId, stuId, curriculumId, stuDetailId))
	if err != nil {
		log.Fatalln(err)
	}

	result := gjson.ParseBytes(data)
	if !result.Get("SuccessResponse").Bool() {
		log.Fatalln(result.Get("Message").String())
	}

	return result
}

// HandleCourseChapters 处理课程章节
func (s *Service) HandleCourseChapters(c *curriculum) {
	defer s.wg.Done()

	if c.courseReadChapters == c.courseChapters {
		log.Println(fmt.Sprintf("%s 观看进度已达100%%，已跳过该课程", c.cuName))
		return
	}

	var (
		loops          = 1                        // 轮询更新进度次数
		notLookCourses = make(map[uint64]bool, 0) // 未观看的课程
	)

	for {
		if loops > 60 {
			break
		}

		var (
			isLookCount uint64 = 0
			chapters           = s.GetCourseChapters(c.courseId, c.curriculumId, s.stuId, s.stuDetailId).Get("Data").Array()
		)

		if loops == 1 {
			for _, chapter := range chapters {
				courses := chapter.Get("ChildNodeList").Array()
				for _, course := range courses {
					// 未读的课程
					if course.Get("IsLook").Uint() == 0 {
						notLookCourses[course.Get("ID").Uint()] = true
					}
				}
			}
		}

		for _, chapter := range chapters {
			courses := chapter.Get("ChildNodeList").Array()
			for _, course := range courses {
				if _, ok := notLookCourses[course.Get("ID").Uint()]; ok {
					if course.Get("IsLook").Uint() == 1 {
						isLookCount++
						log.Println(fmt.Sprintf("课程 >> %d %s | %s | %s 已播放完毕，请在个人中心查看课程观看进度", course.Get("ID").Uint(), c.cuName, chapter.Get("Name").String(), course.Get("Name").String()))
						continue
					}
				}

				if course.Get("IsLook").Uint() == 1 {
					isLookCount++
					continue
				}

				if loops == 1 {
					log.Println(fmt.Sprintf("正在以守护进程方式播放 %s > %s > %s，请勿关闭终端窗口", c.cuName, chapter.Get("Name").String(), course.Get("Name").String()))
				}

				_, err := s.xdSvc.SaveCourseLook(course.Get("ID").Uint())
				log.Println(fmt.Sprintf("更新视频进度成功"))
				if err != nil {
					log.Println(fmt.Sprintf("%s | %s | %s | %s", c.cuName, chapter.Get("Name").String(), course.Get("Name").String(), err.Error()))
					continue
				}
			}
		}

		if isLookCount == c.courseChapters {
			log.Println(fmt.Sprintf("课程 %s 已播放完毕，请查看进度", c.cuName))
			break
		}

		time.Sleep(time.Second * 60)

		loops++
	}
}
