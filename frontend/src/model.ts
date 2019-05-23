class Tab {
  public tabName: string;

  constructor(tabName: string) {
    this.tabName = tabName;
  }
}

class Package {
  public id: string;

  constructor(values: any) {
    this.id = values.id;
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
