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

// InitSMTPDialer 会初始化一个私有的 dialer,
// 以后可以使用 Send 方法时会使用这个内部的私有 dialer 对象，
// 该方法适用于一个程序中只需要一个 SMTP dialer 的情况，如果
// 需要多个 dialer 对象，请使用 Init 和 SendWithDialer 方法
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

// Init 同 gomail.NewDialer
func Init(host, username, password string, port int) *ge.Dialer {
	return ge.NewDialer(host, port, username, password)
}

func getSMTPDialer() SMTPDialer {
	if dialer == nil || dialer.dialer == nil {
		panic("SMTP not connected！, please do InitSMTPDialer")
	}
	return *dialer
}

// GetDialer 在使用 InitSMTPDialer 初始化内部私有 dialer 后，可以使用该函数获取此 dialer 的指针
func GetDialer() *ge.Dialer {
	return getSMTPDialer().dialer
}

func formatAddressList(l []Role) []string {
	res := make([]string, len(l))
	m := ge.NewMessage()
	for i, v := range l {
		res[i] = m.FormatAddress(v.Address, v.Name)
	}
	return res
}

// SendWithDialer 使用自定义的 dialer 发送邮件
func SendWithDialer(dia *ge.Dialer, c *Context) (err error) {
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
	err = dia.DialAndSend(m)
	return err
}

// Send 配合 InitSMTPDialer 使用
func Send(c *Context) (err error) {
	dia := getSMTPDialer()
	return SendWithDialer(dia.dialer, c)
}
