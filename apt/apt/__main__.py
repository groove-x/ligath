import json
import pathlib
import subprocess
from datetime import datetime
from socket import gethostname
from typing import List, Tuple

from debian import copyright as cr


class License:
    name = ""
    machine_readable_name = ""
    body = ""

    def __init__(self, name="", machine_readable_name="", body=""):
        self.name = name
        self.machine_readable_name = machine_readable_name
        self.body = body

    def asdict(self):
        return {
            "name": self.name,
            "machine_readable_name": self.machine_readable_name,
            "body": self.body,
        }


class Copyright:
    copyright = ""
    file_range = None
    license = None

    def __init__(self, copyright="", file_range=None, lic=None):
        self.copyright = copyright
        self.file_range = [] if file_range is None else file_range
        self.license = License() if lic is None else lic

    def asdict(self):
        return {
            "copyright": self.copyright,
            "file_range": self.file_range,
            "license": self.license.asdict(),
        }


class Package:
    name = ""
    version = ""
    copyrights = None
    raw_copyright = ""

    def __init__(self, name="", version="", copyrights=None, raw_copyright=""):
        self.name = name
        self.version = version
        self.copyrights = [] if copyrights is None else copyrights
        self.raw_copyright = raw_copyright

    def asdict(self):
        return {
            "name": self.name,
            "version": self.version,
            "copyrights": [p.asdict() for p in self.copyrights],
            "raw_copyright": self.raw_copyright,
        }


class Output:
    host = ""
    parsed = None
    not_parsed = None

    def __init__(self, host="", parsed=None, not_parsed=None):
        self.host = host
        self.parsed = [] if parsed is None else parsed
        self.not_parsed = [] if not_parsed is None else not_parsed

    def asdict(self):
        return {
            "host": self.host,
            "parsed": [p.asdict() for p in self.parsed],
            "not_parsed": [np.asdict() for np in self.not_parsed],
        }


def parse_copyright(
    package: str,
    path: pathlib.Path,
    version: str,
    parsed: List[Package],
    not_parsed: List[Package],
):
    raw_copyright = ""

    try:
        with open(str(path), encoding="utf-8") as f:
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
        raise RuntimeError("unknown ValueError: {}".format(e))

    file_par = [p for p in c.all_paragraphs() if isinstance(p, cr.FilesParagraph)]
    lice_par = [p for p in c.all_paragraphs() if isinstance(p, cr.LicenseParagraph)]

    copyrights = []
    licenses = dict()

    for lp in lice_par:
        try:
            syn = lp.license.synopsis
        except cr.MachineReadableFormatError as e:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print("{}: {}".format(package, e))
            return

        if syn not in licenses:
            li = License(lp.license.synopsis, lp.license.synopsis, lp.license.text)
            licenses[lp.license.synopsis] = li

    for fp in file_par:
        try:
            syn = fp.license.synopsis
        except cr.MachineReadableFormatError as e:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print("{}: {}".format(package, e))
            return
        except AttributeError as e:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print("{}: {}".format(package, e))
            return

        if syn is None:
            not_parsed.append(Package(package, version, [], raw_copyright))
            print("Unknown License!! {}".format(fp.license.synopsis))
            return

        lic = licenses.get(syn, None)
        if fp.copyright is not None:
            cleansed = "\n".join(l.lstrip() for l in fp.copyright.split("\n"))
        else:
            cleansed = ""
        copyrights.append(Copyright(cleansed, fp.files, lic))

    parsed.append(Package(package, version, copyrights, raw_copyright))


def gather_apt_licenses():
    out = subprocess.check_output("dpkg -l | grep '^ii'", shell=True)

    packages = []
    for l in out.decode().split("\n")[:-1]:
        tokens = [s for s in l.split(" ") if s != ""]
        packages.append((tokens[1].split(":")[0], tokens[2]))

    exist = []
    not_exist = []

    for p, v in packages:
        path = pathlib.Path("/usr/share/doc/{}/copyright".format(p))
        if path.exists():
            exist.append((p, v, path))
        else:
            print("{}: copyright file was not found".format(p))
            not_exist.append((p, v))

    parsed = []
    not_parsed = [Package(p, v, [], "") for p, v in not_exist]

    for package, version, path in exist:
        parse_copyright(package, path, version, parsed, not_parsed)

    print("Summary: parsed({}), parse error({})".format(len(parsed), len(not_parsed)))

    with open(datetime.now().strftime("%y%m%d") + ".json", "w", encoding="utf-8") as f:
        json.dump(
            Output(gethostname(), parsed, not_parsed).asdict(),
            f,
            ensure_ascii=False,
            indent=" ",
        )


gather_apt_licenses()
