package cfg

import (
	"sort"
)

type Config struct {
	Name     string
	Sections []*Section
}

type Section struct {
	Name    string
	Entries []*Entry
}

func New() *Config {
	c := new(Config)
	c.Sections = []*Section{}
	return c
}

func (c *Config) GetSection(name string) *Section {
	idx := sort.Search(len(c.Sections), func(i int) bool{
		return c.Sections[i].Name >= name
	})
	if idx >= len(c.Sections) {
		return nil
	}
	return c.Sections[idx]
}

func (c *Config) Get(section, entry string) *Entry {
	s := c.GetSection(section)
	if s != nil {
		return s.Get(entry)
	}
	return nil
}

func (c *Config) Set(section, entry string, value *Entry) {
	value.Name = entry
	if s := c.GetSection(section); s != nil {
		s.Set(entry, value)
	} else {
		c.Sections = append(c.Sections, &Section{
			Name:    section,
			Entries: []*Entry{value},
		})
		sort.Slice(c.Sections, func(i, j int) bool{
			return c.Sections[i].Name < c.Sections[j].Name
		})
	}
}

func (s *Section) Get(name string) *Entry {
	idx := sort.Search(len(s.Entries), func(i int) bool{
		return s.Entries[i].Name >= name
	})
	if idx >= len(s.Entries) {
		return nil
	}
	return s.Entries[idx]
}

func (s *Section) Set(name string, entry *Entry) {
	entry.Name = name
	idx := sort.Search(len(s.Entries), func(i int) bool{
		return s.Entries[i].Name >= name
	})
	if idx < len(s.Entries) {
		s.Entries[idx] = entry
	} else {
		s.Entries = append(s.Entries, entry)
		sort.Slice(s.Entries, func(i, j int) bool{
			return s.Entries[i].Name < s.Entries[j].Name
		})
	}
}
