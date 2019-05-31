from __future__ import annotations

import pathlib
import subprocess
from dataclasses import dataclass
from datetime import datetime
from os import PathLike
from socket import gethostname
from typing import Dict, List, Optional, Set, Tuple

from dataclasses_json import dataclass_json
from debian import copyright as cr


@dataclass_json
@dataclass
class License:
    name: str
    machine_readable_name: str
    body: str


@dataclass_json
@dataclass
class Copyright:
    copyright: str
    file_range: List[str]
    license: License


@dataclass_json
@dataclass
class Package:
    name: str
    version: str
    copyrights: List[Copyright]
    raw_copyright: str


@dataclass_json
@dataclass
class Output:
    host: str
    parsed: List[Package]
    not_parsed: List[Package]


known_licenses: Dict[str, License] = dict()


def parse_copyright(
    package: str,
    path: PathLike,
    version: str,
    parsed: List[Package],
    not_parsed: List[Package],
):
    raw_copyright = ""

    try:
        with open(path) as f:
            raw_copyright = f.read()
            f.seek(0)
            c = cr.Copyright(f, strict=False)
    except (AttributeError, cr.NotMachineReadableError):
        not_parsed.append(Package(package, version, [], raw_copyright))
        return
    except ValueError as e:
        if len(e.args) > 0 and e.args[0] == "value must not have blank lines":
            # will be raised in debian.deb822:1149
            not_parsed.append(Package(package, version, [], raw_copyright))
            return
        raise RuntimeError(f"unknown ValueError: {e}")

    file_par = [p for p in c.all_paragraphs() if isinstance(p, cr.FilesParagraph)]
    lice_par = [p for p in c.all_paragraphs() if isinstance(p, cr.LicenseParagraph)]

    copyrights = []

    for lp in lice_par:
        try:
            syn = lp.license.synopsis
        except cr.MachineReadableFormatError as e:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print(f"{package}: {e}")
            return

        if syn in known_licenses:
            li = known_licenses[lp.license.synopsis]
            li.machine_readable_name = lp.license.synopsis
        else:
            li = License(lp.license.synopsis, lp.license.synopsis, lp.license.text)
            known_licenses[lp.license.synopsis] = li

    for fp in file_par:
        try:
            syn = fp.license.synopsis
        except cr.MachineReadableFormatError as e:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print(f"{package}: {e}")
            return
        except AttributeError as e:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print(f"{package}: {e}")
            return

        lic = known_licenses.get(syn, None)
        if fp.license.synopsis is None:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print(f"Unknown License!!: {fp.license.synopsis}")
            return

        cleansed = "\n".join(l.lstrip() for l in fp.copyright.split("\n"))
        copyrights.append(Copyright(cleansed, fp.files, lic))

    parsed.append(Package(package, version, copyrights, raw_copyright))


def gather_apt_licenses():
    out: bytes = subprocess.check_output("dpkg -l | grep '^ii'", shell=True)

    packages: List[Tuple[str, str]] = []
    for l in out.decode().split("\n")[:-1]:
        tokens = [s for s in l.split(" ") if s != ""]
        packages.append((tokens[1].split(":")[0], tokens[2]))

    lovot: List[str] = []
    exist: List[Tuple[str, str, PathLike]] = []
    not_exist: List[Tuple[str, str]] = []

    for p, v in packages:
        if p.startswith("lovot"):
            lovot.append(p)

        path = pathlib.Path(f"/usr/share/doc/{p}/copyright")
        if path.exists():
            exist.append((p, v, path))
        else:
            print(f"{p}: copyright file was not found")
            not_exist.append((p, v))

    parsed: List[Package] = []
    not_parsed: List[Package] = [Package(p, v, [], "") for p, v in not_exist]

    for package, version, path in exist:
        parse_copyright(package, path, version, parsed, not_parsed)

    print(f"Summary: parsed({len(parsed)}), parse error({len(not_parsed)})")

    out: str = Output(gethostname(), parsed, not_parsed).to_json(
        ensure_ascii=False, indent=" "
    )
    with open(datetime.now().strftime("%y%m%d") + ".json", "w") as f:
        f.write(out)


gather_apt_licenses()
