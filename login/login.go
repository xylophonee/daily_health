package login

import (
	"dailyPic/untils"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const LoginUrl = "http://authserver.tjut.edu.cn/authserver/login?service=http%3A%2F%2Fwfw.tjut.edu.cn%2Forange%2Fapp%2Fgrmrjktzyd%2Findex.html"
const UserInfoUrl = "http://wfw.tjut.edu.cn/orange/app/grmrjktzyd/getUserInfo.do"
type firstGetParm struct {
	lt                    string
	dllt                  string
	execution             string
	_eventId              string
	rmShown               string
	pwdDefaultEncryptSalt string
}

type Client struct {
	C       *http.Client
	Cookies map[string]*http.Cookie
	nextUrl string
}

type Login struct {
	user string
	pwd string
	param  *firstGetParm
	Client *Client
}

func New(user,pwd string)  *Login{
	return &Login{
		user: user,
		pwd: pwd,
		param: &firstGetParm{},
		Client: &Client{
			C: &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}},
			Cookies: make(map[string]*http.Cookie,0),
		},
	}
}



func (l* Login) Login() error {
	err := l.firstVisitForParam()
	if err != nil {
		return err
	}
	err = l.login()
	if err != nil {
		return err
	}
	err = l.finalVisit()
	if err != nil {
		return err
	}
	return nil
}

func (l* Login) firstVisitForParam() error {
	fmt.Println("开始登录...")
	res, err := l.Client.C.Get(LoginUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200{
		return errors.New("occur error in get params! ")
	}
	l.Client.ProcessCookies(res)
	err = parseParam(res.Body,l.param)
	if err != nil {
		return err
	}
	return nil
}

func (l* Login) login()error{
	fmt.Println("登录中...")
	EnPwd := untils.RandStr(64) + l.pwd
	//加密密码
	enPwd := base64.StdEncoding.EncodeToString(untils.AESCBCEncrypt([]byte(EnPwd), []byte(l.param.pwdDefaultEncryptSalt)))
	enPwd = url.QueryEscape(enPwd)

	strTmp := "captchaResponse=&"+"username="+l.user+"&password="+enPwd+"&lt="+l.param.lt+
		"&dllt="+l.param.dllt+"&execution="+l.param.execution+"&rmShown="+l.param.rmShown+"&_eventId="+l.param._eventId
	bodyP := strings.NewReader(strTmp)
	req,_ := http.NewRequest("POST", LoginUrl, bodyP)
	//添加cookies
	for _,cookie := range l.Client.Cookies {
		req.AddCookie(cookie)
	}
	headerHelp(req)
	//发送登录请求
	res,err := l.Client.C.Do(req)
	if  err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 302 {
		data,_:= ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		panic("登录失败，请重试...")
	}
	l.Client.nextUrl = res.Header.Get("Location")
	l.Client.ProcessCookies(res)
	return nil
}

func (l* Login) finalVisit() error {
	req ,_ := http.NewRequest("GET", l.Client.nextUrl, nil)
	headerHelp(req)
	res, err := l.Client.C.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	l.Client.ProcessCookies(res)
	if res.StatusCode != 302 {
		panic("获取cookies失败！")
	}
	l.Client.nextUrl = res.Header.Get("Location")
	FinalReq ,_ := http.NewRequest("GET", l.Client.nextUrl, nil)
	headerHelp(FinalReq)
	//添加cookies
	for _,cookie := range l.Client.Cookies {
		FinalReq.AddCookie(cookie)
	}
	finalRes, err := l.Client.C.Do(FinalReq)
	if err != nil {
		return err
	}
	defer finalRes.Body.Close()
	if finalRes.StatusCode == 200 {
		dom, err := goquery.NewDocumentFromReader(finalRes.Body)
		if err != nil {
			return err
		}
		if "每日健康状况管理" != dom.Find("title").Text(){
			panic("登录失败...")
		}
		fmt.Println("登录成功！！！")
	}else {
		panic("登录失败...")
	}
	return nil
}

func (c* Client)GetUserInfo() error {

	req, _ := http.NewRequest("POST",UserInfoUrl,nil)
	//添加cookies
	for _,cookie := range c.Cookies {
		req.AddCookie(cookie)
	}
	headerHelp(req)
	res,err := c.C.Do(req)
	defer res.Body.Close()
	if  err != nil {
		return err
	}
	//data ,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(data))
	c.ProcessCookies(res)
	return nil
}

func (c* Client) ProcessCookies(res *http.Response)  {
	for _,v := range res.Cookies(){
		c.Cookies[v.Name] = v
	}
}

func (c* Client)printCookies()  {
	for _,v := range c.Cookies {
		fmt.Println(v.Name+"---"+v.Value)
	}
}

func parseParam(body io.Reader,p *firstGetParm) error {
	fmt.Println("解析参数...")
	b := body
	doc, err := goquery.NewDocumentFromReader(b)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cryKey,exist := doc.Find("input#pwdDefaultEncryptSalt").Attr("value")
	if !exist  {
		fmt.Println("不存在？")
		return errors.New("参数解析错误")
	}
	lt,_ := doc.Find("input[name=lt]").Attr("value")
	dllt,_ := doc.Find("input[name=dllt]").Attr("value")
	execution,_ := doc.Find("input[name=execution]").Attr("value")
	rmShown,_ := doc.Find("input[name=rmShown]").Attr("value")
	_eventId,_ := doc.Find("input[name=_eventId]").Attr("value")
	p.pwdDefaultEncryptSalt = cryKey
	p.lt = lt
	p.rmShown = rmShown
	p._eventId = _eventId
	p.dllt = dllt
	p.execution = execution
	return nil
}


//用于构造请求头
func headerHelp(request *http.Request){
	Headers := map[string]string {
		"Host": "authserver.tjut.edu.cn",
		"Origin": "http://authserver.tjut.edu.cn",
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57",
		"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Referer": "http://authserver.tjut.edu.cn/authserver/login?service=http%3A%2F%2Fwfw.tjut.edu.cn%2Forange%2Fapp%2Fgrmrjktzyd%2Findex.html",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	for hKey,hValue := range Headers{
		request.Header.Set(hKey,hValue)
	}
}