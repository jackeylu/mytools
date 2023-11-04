/*
Copyright © 2023 Lyu Lin <lvlin@whu.edu.cn>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"math"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/jackeylu/mytools/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	imapUsername string
	imapPassword string
	imapHost     string
	imapPort     int
	startFetch   uint32
	endFetch     uint32
	latestSize   uint32
)

// emailCmd represents the email command
var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "从邮箱拉取一批邮件的基础信息和附件名称，并存储在email.xlsx中.",
	Long: `该程序用来从指定的邮箱中拉取一批邮件的基础信息和附件名称. 拉取的信息将会存储在当前目录的email.xlsx文件中.

使用方法:

mytools email -u <username> -p <password> -H [host] -P [port] -s [startFetch] -e [endFetch] -l [latestSize]

参数说明:

-u, --username <username>  邮箱用户名
-p, --password <password>  邮箱密码
-H, --host <host>          邮箱主机地址: 默认是 imap.qq.com
-P, --port <port>          邮箱端口号: 默认是 993
-s, --startFetch <startFetch>  起始邮件序号
-e, --endFetch <endFetch>    结束邮件序号
-l, --latestSize <latestSize>  拉取的最新邮件数量

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 设置日志文件的格式
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
		// 创建一个 LoggerWriter 对象
		logger := util.NewLoggerWriter("logfile.txt")
		defer logger.Close()
		// 将日志同时输出到终端和日志文件
		log.SetOutput(logger)

		if err := checkInput(); err != nil {
			return err
		}

		fetchAndSaveEmails()
		return nil
	},
}

func checkInput() error {
	if imapUsername == "" {
		imapUsername = viper.GetString("email.username")
	}
	if imapPassword == "" {
		imapPassword = viper.GetString("email.password")
	}
	if imapUsername == "" || imapPassword == "" {
		return fmt.Errorf("username or password is empty")
	}
	if imapHost == "" {
		imapHost = viper.GetString("email.host")
	}
	if imapPort == 0 {
		imapPort = viper.GetInt("email.port")
	}
	if imapHost == "" || imapPort == 0 {
		return fmt.Errorf("host or port is empty")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(emailCmd)
	emailCmd.Flags().StringVarP(&imapUsername, "username", "u", "", "imap username")
	emailCmd.Flags().StringVarP(&imapPassword, "password", "p", "", "imap password")
	emailCmd.Flags().StringVarP(&imapHost, "host", "H", "imap.qq.com", "imap host")
	emailCmd.Flags().IntVarP(&imapPort, "port", "P", 993, "imap port")
	emailCmd.Flags().Uint32VarP(&startFetch, "start", "s", 0, "start")
	emailCmd.Flags().Uint32VarP(&endFetch, "end", "e", 0, "end")
	emailCmd.Flags().Uint32VarP(&latestSize, "latest", "l", 50, "latest N email to retreive")
}

func fetchAndSaveEmails() {
	// Connect to the IMAP server
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", imapHost, imapPort), nil)
	if err != nil {
		log.Fatalf("failed to dial IMAP server: %v", err)
	}
	defer c.Close()

	if err := c.Login(imapUsername, imapPassword); err != nil {
		log.Fatalf("failed to login: %v", err)
	}

	// 获取邮箱列表
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("邮箱列表:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// 选择收件箱
	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// Get the last message
	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	} else {
		log.Println("Messages:", mbox.Messages)
	}
	if endFetch == 0 {
		endFetch = mbox.Messages
	}
	if latestSize != 0 {
		startFetch = uint32(math.Max(0, float64(mbox.Messages-latestSize-1)))
	}
	imap.CharsetReader = charset.Reader
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(startFetch, endFetch)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	// msg := <-messages
	var count int = 1
	result := make([]emailInfo, 0)
	for m := range messages {
		log.Printf("Message %d\n", count)
		info := handleOneMessage(m, section)
		if info.Date.IsZero() {
			log.Fatal("Date is zero")
		}
		result = append(result, info)
		count++
	}

	if err := c.Logout(); err != nil {
		log.Fatal(err)
	}

	util.WriteExcelFileByFunction("email.xlsx", []string{"Date", "From", "To", "Subject", "Attachments"}, func() [][]string {
		var ans [][]string
		for _, v := range result {
			ans = append(ans, []string{v.Date.Format("2006-01-02T15:04:05 +080000"),
				v.From, strings.Join(v.To, ","), v.Subject,
				strings.Join(v.Attachments, "\r\n")})
		}
		return ans
	})
}

type emailInfo struct {
	Date        time.Time
	From        string
	To          []string
	Subject     string
	Attachments []string
}

func handleOneMessage(msg *imap.Message, section *imap.BodySectionName) (info emailInfo) {
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal("On creating message reader: ", err)
	}

	header := mr.Header
	if date, err := header.Date(); err == nil {
		log.Println("Date:", date)
		info.Date = date
	}
	if from, err := header.AddressList("From"); err == nil {
		log.Println("From:", from[0].Address)
		info.From = from[0].Address
	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to[0].Address)
		info.To = make([]string, len(to))
		for i, addr := range to {
			info.To[i] = addr.Address
		}
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
		info.Subject = subject
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("On reading next part: ", err)
		}

		switch h := p.Header.(type) {

		case *mail.AttachmentHeader:
			filename, _ := h.Filename()

			if strings.HasSuffix(filename, ".doc") || strings.HasSuffix(filename, ".docx") ||
				strings.HasSuffix(filename, ".zip") || strings.HasSuffix(filename, ".rar") {
				info.Attachments = append(info.Attachments, filename)
			} else {
				log.Printf("Ignore attachment: %s\n", filename)
			}
		}
	}

	return
}
