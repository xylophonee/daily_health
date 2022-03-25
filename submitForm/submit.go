
package submitForm

import (
	"bytes"
	"dailyPic/config"
	"dailyPic/login"
	"dailyPic/mail"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type form struct {
	TBRUID string `json:"TBRUID"`
	XM string `json:"XM"`
	TBRQ string `json:"TBRQ"`
	JRTW string `json:"JRTW"`
	SZSS string `json:"SZSS"`
	JKMZT string `json:"JKMZT"`
	SFYCQK string `json:"SFYCQK"`
	TBZT string `json:"TBZT"`
	SFBT string `json:"SFBT"`
	YMJZQK string `json:"YMJZQK"`
	XCMSFDX string `json:"XCMSFDX"`
	JKMTOKEN JKMTOKEN `json:"JKMTOKEN"`
	XCMTOKEN XCMTOKEN `json:"XCMTOKEN"`
}
type JKMTOKEN struct {
	Token string `json:"token"`
	Ids string `json:"ids"`
}
type XCMTOKEN struct {
	Token string `json:"token"`
	Ids string `json:"ids"`
}


type Pic struct {
	Msg string `json:"msg"`
	Code int `json:"code"`
	Success bool `json:"success"`
	Data struct {
		Wid string `json:"wid"`
		Token string `json:"token"`
		FileName string `json:"fileName"`
		FileURL string `json:"fileUrl"`
		FileFormat string `json:"fileFormat"`
		Suffix string `json:"suffix"`
		URL string `json:"url"`
		Preview bool `json:"preview"`
		PreviewURL interface{} `json:"previewUrl"`
	} `json:"data"`
}

const url_pic = "http://wfw.tjut.edu.cn/orange/app/grmrjktzyd/fileupload/uploadFile.do"
const final_url = "http://wfw.tjut.edu.cn/orange/app/grmrjktzyd/saveData/4eb645ccd2224b039a17c65d548891ff.do"

//上传图片
func uploadPic(l *login.Login,path string)([]string,error){

	buff,ct,err := creatBuff(path)
	if err != nil{
		return nil,err
	}
	req,_ := http.NewRequest("POST", url_pic, &buff)
	for _,cookie := range l.Client.Cookies {
		req.AddCookie(cookie)
	}
	//fmt.Println(req.Cookies())
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57",
		"Accept": "application/json, text/plain, */*",
		"Referer": "http://wfw.tjut.edu.cn/orange/app/grmrjktzyd/index.html",
		"Content-Type":ct,
		"Host": "wfw.tjut.edu.cn",
	}
	for k,v := range headers{
		req.Header.Set(k,v)
	}

	res, err := l.Client.C.Do(req)
	if err != nil {
		return nil, err
	}
	p := Pic{}
	data,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(data))
	err = json.Unmarshal(data,&p)
	if err != nil {
		return nil, err
	}
	if p.Success {
		return []string{p.Data.Wid,p.Data.Token},nil
	}
	return nil,errors.New("上传pic错误")
}

//组合上传图片的数据
func creatBuff(path string)(bytes.Buffer,string,error){
	var buff bytes.Buffer
	// 创建一个Writer
	writer := multipart.NewWriter(&buff)

	// 写入一般的表单字段
	err := writer.WriteField("orderId", "0")
	if err != nil {
		return bytes.Buffer{},"", err
	}
	err = writer.WriteField("uploadType", "temp")
	if err != nil {
		return bytes.Buffer{}, "",err
	}

	// 写入文件字段
		w, err := writer.CreateFormFile("file", path)
		if err != nil {
			return buff,"", err
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return buff, "",err
		}
		// 把文件内容写入cd
	_, err = w.Write(data)
	if err != nil {
		return bytes.Buffer{},"", err
	}
	err = writer.Close()
	if err != nil {
		return bytes.Buffer{},"", err
	}

	return buff,writer.FormDataContentType(),nil
}

// SubmitForm 提交表单
func SubmitForm(l *login.Login,user config.Users) error{
	fmt.Println("上传健康码&行程码")
	u1 ,err := uploadPic(l,"./pic/1.jpg")
	if err != nil {
		return err
	}
	u2 ,err:= uploadPic(l,"./pic/2.jpg")
	if err != nil {
		return err
	}
	fmt.Println("健康码&行程码上传成功！")
	data := form{
		TBRUID:user.Username,
		XM: "",
		TBRQ: time.Now().Format("2006-01-02"),
		JRTW: "36.5",
		SZSS: user.Address,
		JKMZT: "1",
		SFYCQK: "0",
		TBZT: "0",
		SFBT: "0",
		YMJZQK: "06",
		XCMSFDX:"0",
		JKMTOKEN: JKMTOKEN{
			Token: u1[0],
			Ids: u1[1],
		},
		XCMTOKEN: XCMTOKEN{
			Token: u2[0],
			Ids: u2[1],
		},
	}

	j ,err := json.Marshal(data)
	if err != nil{
		return err
	}

	buff := bytes.Buffer{}
	buff.WriteString("sfbt=0&formData=")
	buff.Write(j)
	headers := map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36 Edg/89.0.774.57",
			"Accept": "application/json, text/plain, */*",
			"Referer": "http://wfw.tjut.edu.cn/orange/app/grmrjktzyd/index.html",
			"Content-Type": "application/x-www-form-urlencoded",
			"Host": "wfw.tjut.edu.cn",
	}

	req,_ := http.NewRequest("POST", final_url, &buff)
	for k,v := range headers{
		req.Header.Set(k,v)
	}
	for _,cookie := range l.Client.Cookies {
		req.AddCookie(cookie)
	}
	fmt.Println("提交表单...")
	response, err := l.Client.C.Do(req)
	if err != nil {
		return err
	}
	d, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(d))
	str := string(d)
	err = mail.SendMail("每日打卡", user.Email,str , str, "html")
	if err != nil {
		return err
	}
	return nil
}