package main

type License struct {
	Name                 string   `json:"name"`
	MachineReadableNames []string `json:"machine_readable_name,omitempty"`
	Body                 string   `json:"body"`
}

type Copyright struct {
	Notice    string  `json:"notice"`
	FileRange string  `json:"file_range"`
	License   License `json:"license"`
}

type Package struct {
	ID         string      `json:"id"`
	Version    string      `json:"version"`
	Copyrights []Copyright `json:"copyrights"`
}

func GPLv2() License {
	return License{
		Name: "GNU General Public License version 2",
		MachineReadableNames: []string{
			"GPL-v2",
		},
		Body: "yoyo this is gplv2",
	}
}
