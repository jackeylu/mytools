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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("email called")

		checkInput()
		main()
		fmt.Println("email end")
	},
}

func checkInput() {
	if imapUsername == "" {
		imapUsername = viper.GetString("email.username")
	}
	if imapPassword == "" {
		imapPassword = viper.GetString("email.password")
	}
	if imapUsername == "" || imapPassword == "" {
		log.Fatal("username or password is empty")
	}
	if imapHost == "" {
		imapHost = viper.GetString("email.host")
	}
	if imapPort == 0 {
		imapPort = viper.GetInt("email.port")
	}
	if imapHost == "" || imapPort == 0 {
		log.Fatal("host or port is empty")
	}
}

func init() {
	rootCmd.AddCommand(emailCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// emailCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// emailCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	emailCmd.Flags().StringVarP(&imapUsername, "username", "u", "", "imap username")
	emailCmd.Flags().StringVarP(&imapPassword, "password", "p", "", "imap password")
	emailCmd.Flags().StringVarP(&imapHost, "host", "H", "imap.qq.com", "imap host")
	emailCmd.Flags().IntVarP(&imapPort, "port", "P", 993, "imap port")
	emailCmd.Flags().Uint32VarP(&startFetch, "start", "s", 0, "start")
	emailCmd.Flags().Uint32VarP(&endFetch, "end", "e", 0, "end")
	emailCmd.Flags().Uint32VarP(&latestSize, "latest", "l", 50, "latest N email to retreive")
}

func main() {
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
		startFetch = mbox.Messages - latestSize
	}
	imap.CharsetReader = charset.Reader
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(startFetch, endFetch)

	// Get the whole message body
	var section imap.BodySectionName
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
			ans = append(ans, []string{v.Date.Format("2006-01-02 15:04:05"),
				v.From, strings.Join(v.To, ","), v.Subject,
				strings.Join(v.Attachments, ",")})
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

func handleOneMessage(msg *imap.Message, section imap.BodySectionName) (info emailInfo) {
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(&section)
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
