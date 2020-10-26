package modem

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jorgefuertes/mister-modemu/internal/console"
)

// finds an AT route by its name
func (p *parser) findByPath(path string) (*route, error) {
	for _, r := range p.routes {
		if r.path == path {
			return &r, nil
		}
	}
	return nil, errors.New("Route not found")
}

// AT - Add an AT route
func (p *parser) AT(path string, cb func(s *Status)) error {
	if _, err := p.findByPath(path); err == nil {
		return errors.New("Route already exists")
	}

	needle := strings.ReplaceAll(path, `*`, `.*`)
	needle = strings.ReplaceAll(needle, `+`, `\+`)
	needle = strings.ReplaceAll(needle, `?`, `\?`)
	reg := regexp.MustCompile(`\A` + needle + `\z`)
	p.routes = append(p.routes, route{path, reg, cb})

	return nil
}

// Parse - parse an AT command
func (p *parser) Parse(m *Status, cmd string) {
	p.Cmd = cmd
	p.Err = nil
	for _, r := range p.routes {
		if r.e.MatchString(cmd) {
			console.Debug(`AT/PARSER`, "Matched route: ", r.path)
			r.cb(m)
			if p.Err != nil {
				console.Error(`AT/PARSER`, p.Err.Error())
			}
			return
		}
	}

	console.Debug(`AT/PARSER`, "No matching route, writing OK by default")
	m.OK()
}

// Error - set an error to be logged
func (p *parser) Error(err interface{}) {
	switch e := err.(type) {
	case error:
		p.Err = e
	case string:
		p.Err = errors.New(e)
	default:
		panic("Parser error: Only error or string allowed as error")
	}
}

// GetArg - one arg line, even if it has colon sep args
func (p *parser) GetArg() string {
	r := regexp.MustCompile(`^[A-Z_]+\=\"*(?P<Arg>.*)\"*$`)
	m := r.FindStringSubmatch(p.Cmd)
	console.Debug(`AT/PARSER/ARG`, m)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// GetArgs - slice of args from colon sep args
func (p *parser) GetArgs() []string {
	args := strings.Split(p.GetArg(), ",")
	for i, a := range args {
		args[i] = strings.Trim(a, `"`)
		console.Debug(`AT/PARSER/ARGS`, args[i])
	}
	return args
}
