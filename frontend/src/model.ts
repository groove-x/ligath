class Tab {
  public name: string;
  public version: string;
  public kind: string;

  constructor(name: string, version: string, kind: string) {
    this.name = name;
    this.version = version;
    this.kind = kind;
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
  public license: License;

  private fileRange: string[];

  private get range(): string {
    if (this.fileRange === undefined) {
      return '';
    } else {
      return this.fileRange.join('\n');
    }
  }

  private set range(newVal: string) {
    this.fileRange = newVal.split('\n');
  }

  constructor(values: any) {
    this.copyright = values.copyright;
    this.fileRange = values.file_range || values.fileRange;
    this.license = new License(values.license);
  }

  public jsonCompatible(): any {
    return {
      copyright: this.copyright,
      file_range: this.fileRange,
      license: this.license,
    };
  }
}

class Package {
  public name: string;
  public version: string;
  public copyrights: Copyright[];
  public rawCopyright: string;

  constructor(values: any) {
    this.copyrights = new Array<Copyright>();

    if (values === null) {
      this.name = '';
      this.version = '';
      this.rawCopyright = '';
    } else {
      this.name = values.name;
      this.version = values.version;
      values.copyrights.forEach((c: any) => {
        this.copyrights.push(new Copyright(c));
      });
      this.rawCopyright = values.raw_copyright;
    }
  }

  public jsonCompatible(): any {
    const copyrightsj = new Array<Copyright>();
    this.copyrights.forEach((c: Copyright) => {
      copyrightsj.push(c.jsonCompatible());
    });
    return {
      name: this.name,
      version: this.version,
      copyrights: copyrightsj,
      raw_copyright: this.rawCopyright,
    };
  }
}

class PackageListItem {
  public name: string;
  public version: string;
  public verified: boolean;

  constructor(name: string, version: string, verified: boolean) {
    this.name = name;
    this.version = version;
    this.verified = verified;
  }
}

class HomeState {
  public lastCounter: number;
  public parsed: PackageListItem[];
  public notParsed: PackageListItem[];
  public verified: PackageListItem[];

  constructor() {
    this.lastCounter = 0;
    this.parsed = new Array<PackageListItem>();
    this.notParsed = new Array<PackageListItem>();
    this.verified = new Array<PackageListItem>();
  }
}

class LicenseState {
  public licenses: string[];
  public filteredPackages: PackageListItem[];

  constructor() {
    this.licenses = new Array<string>();
    this.filteredPackages = new Array<PackageListItem>();
  }
}

enum FunctionTabEnum {
  HomeTab = 0,
  LicenseTab = 1,
  BulkRenameTab = 2,
}

export { Tab, License, Copyright, Package, HomeState, LicenseState, FunctionTabEnum };
