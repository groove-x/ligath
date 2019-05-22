class Tab {
  public tabName: string;

  constructor(tabName: string) {
    this.tabName = tabName;
  }
}

class HomeState {
  public lastCounter: number;

  constructor() {
    this.lastCounter = 0;
  }
}

export { Tab, HomeState };
