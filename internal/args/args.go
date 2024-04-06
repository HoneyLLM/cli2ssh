package args

import (
	"bytes"
	"fmt"
	"net"
	"text/template"

	"github.com/PeronGH/cli2ssh/internal/utils"
	"github.com/charmbracelet/ssh"
)

type Session struct {
	User       string
	Host       string
	Port       string
	Command    string
	Subsystem  string
	PublicKey  string
	RemoteAddr string
}

func NewSession(s ssh.Session) *Session {
	user := s.User()
	host, port, _ := net.SplitHostPort(s.RemoteAddr().String())
	var publicKey string
	if pk := s.PublicKey(); pk != nil {
		publicKey = utils.StringifyPublicKey(pk)
	}
	return &Session{
		User:       user,
		Host:       host,
		Port:       port,
		Command:    s.RawCommand(),
		Subsystem:  s.Subsystem(),
		PublicKey:  publicKey,
		RemoteAddr: s.RemoteAddr().String(),
	}
}

func (s *Session) FormatArg(arg string) string {
	return formatTemplate(arg, s)
}

type ArrayArg []string

func (a *ArrayArg) String() string {
	return fmt.Sprintf("%v", *a)
}

func (a *ArrayArg) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (s *Session) FormatArgs(args []string) []string {
	formatted := make([]string, len(args))
	for i, arg := range args {
		formatted[i] = formatTemplate(arg, s)
	}
	return formatted
}

func formatTemplate(templateStr string, data any) string {
	template, err := template.New("").Parse(templateStr)
	if err != nil {
		return templateStr
	}

	var buf bytes.Buffer
	err = template.Execute(&buf, data)
	if err != nil {
		return templateStr
	}
	return buf.String()
}
