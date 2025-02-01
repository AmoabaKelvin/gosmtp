// email server written from the socket layer. Simplified version of RFC 831, plan
// will be to extend it to 5321.

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type SMTPSession struct {
	conn       net.Conn
	writer     *bufio.Writer
	reader     *bufio.Reader
	serverName string
	mailFrom   string
	mailTo     string
	message    strings.Builder
}

func NewSMTPSession(conn net.Conn) *SMTPSession {
	return &SMTPSession{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (s *SMTPSession) reply(response string) {
	s.writer.WriteString(response)
	s.writer.Flush()
}

func removeAngleBrackets(s string) string {
	cleaned := strings.ReplaceAll(s, "<", "")
	cleaned = strings.ReplaceAll(cleaned, ">", "")
	return cleaned
}

func main() {
	fmt.Println("Starting a new server to listen on 25")
	server, err := net.Listen("tcp", ":25")

	if err != nil {
		os.Exit(1)
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		session := NewSMTPSession(conn)
		if err != nil {
			panic(err)
		}

		go handleConnection(session)
	}
}

func handleConnection(session *SMTPSession) {
	defer session.conn.Close()

	session.reply("220 Connection successful \n")

	for {
		message, err := session.reader.ReadString('\n')
		if err != nil {
			break
		}

		line := strings.TrimSpace(message)
		command := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(command[0])

		switch cmd {
		case "MAIL":
			// process mail commands
			processMailCommand(command[1], session)

		case "RCPT":
			processRcptCommand(command[1], session)

		case "DATA":
			processDataCommand(session)

		case "QUIT":
			processQuitCommand(session)

		default:
			session.reply("550 I did not quite catch that, try again please")
		}
	}

}

func processMailCommand(command string, session *SMTPSession) {
	fromEmail := strings.Split(command, ":")[1]
	session.mailFrom = removeAngleBrackets(fromEmail)

	session.reply("250 OK. \n")
}

func processRcptCommand(command string, session *SMTPSession) {
	rcpt := strings.Split(command, ":")[1]
	cleanedEmail := removeAngleBrackets(rcpt)
	session.mailTo = cleanedEmail
	session.reply("250 OK. \n")
}

func processQuitCommand(session *SMTPSession) {
	fmt.Println("Processing quit command")
	session.conn.Close()
}

func processDataCommand(session *SMTPSession) {
	// check that there is a recipient before trying to write the mail
	if len(session.mailFrom) == 0 {
		session.reply("503 Need RCPT command first")
	}
	session.reply("354 Start mail input; end with <CRLF>.<CRLF> \n")

	for {
		data, err := session.reader.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}

		fmt.Println(data)

		if data == ".\r\n" {
			// marks the end of the email message, we can store the message and then write
			break
		}

		session.message.WriteString(data)
	}

	fileName := fmt.Sprintf("%s.txt", session.mailTo)
	file, err := os.Create(fileName)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	var emailContentBuilder strings.Builder

	fmt.Fprintf(&emailContentBuilder, "From: %s\r\n", session.mailFrom)
	fmt.Fprintf(&emailContentBuilder, "To: %s\r\n", session.mailTo)
	emailContentBuilder.WriteString(session.message.String())

	_, err = file.WriteString(emailContentBuilder.String())
	if err != nil {
		os.Exit(1)
	}
	session.reply("250 OK. \n")
}
