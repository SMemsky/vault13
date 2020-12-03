package cfg

import (
	"fmt"
	"bufio"
	"io"
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

func (c *Config) parseLine(section, line string) (string, error) {
	fmt.Println(line)
	return section, nil
}
