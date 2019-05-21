import pathlib
import re
import subprocess
from os import PathLike
from typing import Dict, List, Iterator, Tuple

from debian import copyright as cr


def gather_apt_licenses() -> Tuple[Dict[str, List[cr.ParagraphTypes]], List[str]]:
    out: bytes = subprocess.check_output("dpkg --get-selections", shell=True)
    packages: List[str] = [
        s.split("\t")[0].split(":")[0] for s in out.decode().split("\n")
    ]

    lovot: List[str] = []
    exist: List[Tuple[str, PathLike]] = []
    not_exist: List[str] = []

    for p in packages:
        if p.startswith("lovot"):
            lovot.append(p)

        path = pathlib.Path(f"/usr/share/doc/{p}/copyright")
        if path.exists():
            exist.append((p, path))
        else:
            print(f"{p}: copyright file was not found")
            not_exist.append(path)

    not_parsed = not_exist
    parsed: Dict[str, List[cr.ParagraphTypes]] = dict()

    for package, path in exist:
        try:
            with open(path) as f:
                c = cr.Copyright(f, strict=False)
        except (AttributeError, cr.NotMachineReadableError):
            not_parsed.append(package)
            continue
        except ValueError as e:
            if len(e.args) > 0 and e.args[0] == "value must not have blank lines":
                # will be raised in debian.deb822:1149
                not_parsed.append(package)
                continue
            raise RuntimeError(f"unknown ValueError: {e}")

        paragraphs: Iterator[cr.AllParagraphTypes] = c.all_paragraphs()
        file_par = [p for p in paragraphs if isinstance(p, cr.FilesParagraph)]
        lice_par = [
            p
            for p in paragraphs
            if isinstance(p, cr.LicenseParagraph) and p.known_format()
        ]

        parsed[package] = file_par + lice_par

    print(f"Summary: parsed({len(parsed.items())}), parse error({len(not_parsed)})")
    return parsed, not_parsed


gather_apt_licenses()
