package email

import (
	"testing"
)

func init() {
	InitSMTPDialer("smtp.163.com", "153649xxxxx@163.com", "xxxxxxx", 25)
}
func TestSend(t *testing.T) {
	e := Send(&Context{
		ToList: []Role{
			{Address: "17xxxx5643@qq.com", Name: "QQ"},
		},
		CcList: []Role{
			{Address: "17xxxx4650@st.xxx.edu.cn", Name: "NUC"},
			{Address: "zhxxxxxo23@outlook.com", Name: "Outlook"},
			{Address: "30xxxx3164@qq.com", Name: "30xxxx3164"},
		},
		BccList: []Role{
			{Address: "31xxxx0040@qq.com", Name: "BCC"},
		},
		Subject: "每日一题",
		Body: "<h1> Strive </h1> \n 今天执行下面的代码, 期望结果为 1， 但实际输出 0，" +
			" 原因在于当函数执行到 `defer` 时，编译器会拷贝延迟函数用到的外部变量，" +
			"而此时 `ab` 都是初始值 0，为了解决这个问题，我们可以使用匿名函数：",
		Path: "./test.txt",
	})
	if e != nil {
		t.Error(e)
	}

}
