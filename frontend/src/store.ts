import Vue from 'vue';
import Vuex from 'vuex';
import { HomeState, Package, Tab } from '@/model';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    endpoint_front: window.location.hostname + ":" + window.location.port,
    endpoint_back: "http://" + window.location.hostname + ":3939",
    fresh: true,
    tabs: new Map<string, Tab>(),
    tabsTick: 0,
    packages: new Map<string, Package>(),
    packagesTick: 0,
    home: new HomeState(),
  },
  mutations: {
    newTab(state, tabName) {
      // increment tabsTick to make reaction happen
      state.tabsTick += 1;
      state.tabs.set(tabName, new Tab(tabName));
    },
    closeTab(state, name) {
      state.tabsTick += 1;
      state.tabs.delete(name);
    },
    setPackageData(state, pkg: any) {
      state.packagesTick += 1;
      const p = new Package(pkg);
      state.packages.set(p.id, p);
    },
    fresh(state) {
      state.fresh = false;
    },
    getUnclassified(state, pkgs) {
      state.home.unclassified = pkgs;
    }
  },
  getters: {
    tabsAsList: (state) => {
      // relate tabsTick with tabs
      return state.tabsTick && Array.from(state.tabs);
    },
  },
  actions: {

  },
});
