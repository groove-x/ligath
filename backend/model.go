package main

import "encoding/json"

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

type MigratedPackage struct {
	Name         string
	Version      string
	Copyrights   []Copyright
	RawCopyright string
}

type PackageListItem struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Verified bool   `json:"verified"`
}

func GPLv2() License {
	return License{
		Name:                "GNU General Public License version 2",
		MachineReadableName: "GPL-v2",
		Body:                "yoyo this is gplv2",
	}
}

func CreateSamplePackageInfo() string {
	j := []Package{
		{
			Name:    "libfoo",
			Version: "1.2.3-4",
			Copyrights: []Copyright{
				{
					Copyright: "Copyright (C) 2015-2019 Foobar, Inc.",
					FileRange: []string{"include/*", "src/*"},
					License: License{
						Name:                "GPL-2",
						MachineReadableName: "",
						Body:                "This software is free software; you can redistribute ...",
					},
				},
				{
					Copyright: "Copyright (C) 1995-2019 Hogepiyo, Inc.",
					FileRange: []string{"lib/*"},
					License: License{
						Name:                "GPL-3",
						MachineReadableName: "",
						Body:                "This software is free software; you can redistribute ...",
					},
				},
			},
			RawCopyright: "libfoo\n\ninclude/*, src/*\nCopyright (C) 2015-2019 Foobar, Inc.\n\nlib/*\nCopyright (C) 1995-2019 Hogepiyo, Inc.\n\nThis software is free software; you can redistribute ...",
		},
	}
	out, err := json.MarshalIndent(&j, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(out)
}
