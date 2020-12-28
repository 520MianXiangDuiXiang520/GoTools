package email

import (
	junePath "github.com/520MianXiangDuiXiang520/GinTools/path"
	ge "gopkg.in/gomail.v2"
	"path"
	"runtime"
	"sync"
)

type SMTPDialer struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	dialer   *ge.Dialer
}

type Role struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Context struct {
	ToList  []Role `json:"to_list"`  // 收件人列表
	CcList  []Role `json:"cc_list"`  // 抄送列表
	BccList []Role `json:"bcc_list"` // 密送列表
	Subject string `json:"subject"`  // 邮件主题
	Body    string `json:"body"`     // 邮件正文
	Path    string `json:"path"`     // 附件路径
}

var dialer *SMTPDialer
var once sync.Once

func InitSMTPDialer(host, username, password string, port int) {
	if dialer == nil || dialer.dialer == nil {
		once.Do(func() {
			dialer = &SMTPDialer{
				dialer:   ge.NewDialer(host, port, username, password),
				Username: username,
				Password: password,
				Host:     host,
				Port:     port,
			}
		})
	}
}

func getSMTPDialer() SMTPDialer {
	if dialer == nil || dialer.dialer == nil {
		panic("SMTP not connected！, please do InitSMTPDialer")
	}
	return *dialer
}

func formatAddressList(l []Role) []string {
	res := make([]string, len(l))
	m := ge.NewMessage()
	for i, v := range l {
		res[i] = m.FormatAddress(v.Address, v.Name)
	}
	return res
}

func Send(c *Context) (err error) {
	dia := getSMTPDialer()
	m := ge.NewMessage()
	m.SetHeader("From", dia.Username)
	m.SetHeader("To", formatAddressList(c.ToList)...)
	if len(c.CcList) > 0 {
		m.SetHeader("Cc", formatAddressList(c.CcList)...)
	}
	if len(c.BccList) > 0 {
		m.SetHeader("Bcc", formatAddressList(c.BccList)...)
	}
	if len(c.Path) > 0 {
		if !junePath.IsAbs(c.Path) {
			_, currently, _, _ := runtime.Caller(1)
			filename := path.Join(path.Dir(currently), c.Path)
			m.Attach(filename)
		} else {
			m.Attach(c.Path)
		}

	}
	m.SetHeader("Subject", c.Subject)
	m.SetBody("text/html", c.Body)
	err = dia.dialer.DialAndSend(m)
	return err
}
