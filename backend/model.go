package main

import "sort"

type License struct {
	Name                string `json:"name"`
	MachineReadableName string `json:"machine_readable_name,omitempty"`
	Body                string `json:"body"`
}

func SortLicenses(in *[]License) {
	l := Licenses{
		l: in,
	}
	sort.Sort(l)
}

type Licenses struct {
	l *[]License
}

func (l Licenses) Slice() []License {
	return *l.l
}

func (l Licenses) Len() int {
	return len(*l.l)
}

func (l Licenses) Swap(i, j int) {
	(*l.l)[i], (*l.l)[j] = (*l.l)[j], (*l.l)[i]
}

func (l Licenses) Less(i, j int) bool {
	return (*l.l)[i].Name < (*l.l)[j].Name
}

type Copyright struct {
	Copyright string   `json:"copyright"`
	FileRange []string `json:"file_range"`
	License   License  `json:"license"`
}

type Package struct {
	Name         string      `json:"name"`
	Version      string      `json:"version"`
	Copyrights   []Copyright `json:"copyrights"`
	RawCopyright string      `json:"raw_copyright"`
}

type AptOutput struct {
	Host      string    `json:"host"`
	Parsed    []Package `json:"parsed"`
	NotParsed []Package `json:"not_parsed"`
}

type MigratedPackage struct {
	Name         string
	Version      string
	Copyrights   []Copyright
	RawCopyright string
}

type PackageListItem struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func GPLv2() License {
	return License{
		Name:                "GNU General Public License version 2",
		MachineReadableName: "GPL-v2",
		Body:                "yoyo this is gplv2",
	}
}
