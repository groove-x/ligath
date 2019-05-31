class Tab {
  public name: string;
  public version: string;

  constructor(name: string, version: string) {
    this.name = name;
    this.version = version;
  }
}

class License {
  public name: string;
  public machineReadableName: string;
  public body: string;

  constructor(values: any) {
    this.name = values.name;
    this.machineReadableName = values.machineReadableName;
    this.body = values.body;
  }
}

class Copyright {
  public copyright: string;
  public fileRange: string[];
  public license: License;

  constructor(values: any) {
    this.copyright = values.copyright;
    this.fileRange = values.file_range;
    this.license = new License(values.license);
  }
}

class Package {
  public name: string;
  public version: string;
  public copyrights: Copyright[];
  public rawCopyright: string;

  constructor(values: any) {
    this.name = values.name;
    this.version = values.version;
    this.copyrights = new Array<Copyright>();
    values.copyrights.forEach((c: any) => {
      this.copyrights.push(new Copyright(c));
    });
    this.rawCopyright = values.raw_copyright;
  }
}

class PackageListItem {
  public name: string;
  public version: string;

  constructor(name: string, version: string) {
    this.name = name;
    this.version = version;
  }
}

class HomeState {
  public lastCounter: number;
  public parsed: PackageListItem[];
  public notParsed: PackageListItem[];
  public manual: PackageListItem[];

  constructor() {
    this.lastCounter = 0;
    this.parsed = new Array<PackageListItem>();
    this.notParsed = new Array<PackageListItem>();
    this.manual = new Array<PackageListItem>();
  }
}

export { Tab, License, Copyright, Package, HomeState };
