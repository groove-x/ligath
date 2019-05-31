package main

type License struct {
	Name                string `json:"name"`
	MachineReadableName string `json:"machine_readable_name,omitempty"`
	Body                string `json:"body"`
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
