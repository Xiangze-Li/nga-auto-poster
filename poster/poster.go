package poster

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Xiangze-Li/nga-auto-poster/config"
	"github.com/Xiangze-Li/nga-auto-poster/utils"
	cn "golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const NgaApi = "https://bbs.nga.cn/post.php"

func PostReply(cfg config.Config) error {
	log.Println("准备发帖")

	content := utils.ReadAndPop(cfg.ContentFile, cfg.Split)
	if len(content) == 0 {
		log.Println("没得发了! 别鸽了! 快去写文!")
		return nil
	}

	log.Print("帖子内容\n" + content)

	content, _ = cn.GB18030.NewEncoder().String(content)

	form := url.Values{
		"action":       []string{"reply"},
		"fid":          []string{strconv.Itoa(cfg.Fid)},
		"tid":          []string{strconv.Itoa(cfg.Tid)},
		"nojump":       []string{"1"},
		"step":         []string{"2"},
		"lite":         []string{"xml"},
		"post_content": []string{content},
	}
	formStr := form.Encode()

	req, err := http.NewRequest(http.MethodPost, NgaApi, strings.NewReader(formStr))
	if err != nil {
		return err
	}
	for k, v := range cfg.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formStr)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP返回值: %d", resp.StatusCode)
		log.Println("发帖失败")
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	ngaResp := decodeResp(body)

	if strings.Contains(string(ngaResp), "发贴完毕 ...") {
		log.Println("发帖成功")
	} else {
		log.Println(string(ngaResp))
		log.Println("发帖失败")
	}

	return nil
}

func decodeResp(body []byte) string {
	type ngaXmlResp struct {
		XMLName xml.Name `xml:"root"`
		Message struct {
			XMLName xml.Name `xml:"__MESSAGE"`
			Items   []string `xml:"item"`
		} `xml:"__MESSAGE"`
	}

	var resp ngaXmlResp
	d := xml.NewDecoder(bytes.NewReader(body))
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return transform.NewReader(input, cn.GB18030.NewDecoder()), nil
	}

	utils.ExitOnError(d.Decode(&resp))
	return strings.Join(resp.Message.Items, " ")
}
