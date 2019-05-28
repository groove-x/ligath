class Tab {
  public tabName: string;

  constructor(tabName: string) {
    this.tabName = tabName;
  }
}

class License {
  public name: string;
  public machineReadableNames: string[];
  public body: string;

  constructor(values: any) {
    this.name = values.name;
    this.machineReadableNames = values.machineReadableNames;
    this.body = values.body;
  }
}

class Copyright {
  public notice: string;
  public fileRange: string;
  public license: License;

  constructor(values: any) {
    this.notice = values.notice;
    this.fileRange = values.file_range;
    this.license = new License(values.license);
  }
}

class Package {
  public id: string;
  public version: string;
  public copyrights: Copyright[];
  public rawCopyright: string;

  constructor(values: any) {
    this.id = values.id;
    this.version = values.version;
    this.copyrights = new Array<Copyright>();
    values.copyrights.forEach((c: any) => {
      this.copyrights.push(new Copyright(c));
    });
    this.rawCopyright = values.raw_copyright;
  }
}

class HomeState {
  public lastCounter: number;
  public unclassified: Package[];
  public manualClassified: Package[];
  public autoClassified: Package[];

  constructor() {
    this.lastCounter = 0;
    this.unclassified = new Array<Package>();
    this.manualClassified = new Array<Package>();
    this.autoClassified = new Array<Package>();
  }
}

export { Tab, Package, HomeState };
