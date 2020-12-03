package cfg

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var (
	ErrNoEndingBracket = fmt.Errorf("Ending bracket missing")
	ErrNoEqualSign     = fmt.Errorf("Equal sign missing")
)

func (c *Config) Parse(r io.Reader) error {
	var firstErr, err error
	section := "unknown"

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		section, err = c.parseLine(section, scanner.Text())
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if firstErr != nil {
		return firstErr
	}
	return scanner.Err()
}

// TODO: actual line error reports
func (c *Config) parseLine(section, line string) (string, error) {
	line = strings.TrimSpace(line)
	commentIndex := strings.IndexRune(line, ';')
	if commentIndex != -1 {
		line = line[:commentIndex]
	}
	if line == "" {
		return section, nil
	}
	if strings.HasPrefix(line, "[") {
		endingBracket := strings.IndexRune(line, ']')
		if endingBracket == -1 {
			return section, ErrNoEndingBracket
		}
		// TODO: maybe copy this string somehow?
		section = line[1:endingBracket]
		return section, nil
	}
	equalSign := strings.IndexRune(line, '=')
	if equalSign == -1 {
		return section, ErrNoEqualSign
	}
	c.Set(section, strings.TrimSpace(line[:equalSign]), String(strings.TrimSpace(line[equalSign+1:])))
	return section, nil
}
