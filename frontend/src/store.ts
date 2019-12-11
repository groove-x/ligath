import Vue from 'vue';
import Vuex from 'vuex';
import { HomeState, LicenseState, Package, Tab } from '@/model';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    endpoint_front: window.location.hostname + ':' + window.location.port,
    endpoint_back: 'http://' + window.location.hostname + ':3939',
    fresh: true,
    tabs: new Map<string, Tab>(),
    tabsTick: 0,
    packages: new Map<string, Package>(),
    packagesTick: 0,
    home: new HomeState(),
    license: new LicenseState(),
  },
  mutations: {
    newTab(state, t) {
      // increment tabsTick to make reaction happen
      state.tabsTick += 1;
      state.tabs.set(t.name + t.version + t.kind, new Tab(t.name, t.version, t.kind));
    },
    closeTab(state, t) {
      state.tabsTick += 1;
      state.tabs.delete(t.name + t.version + t.kind);
    },
    setPackageData(state, arg) {
      state.packagesTick += 1;
      const p = new Package(arg.pkg);
      state.packages.set(arg.pkg.name + arg.pkg.version + arg.kind, arg.pkg);
    },
    fresh(state) {
      state.fresh = false;
    },
    getParsed(state, pkgs) {
      state.home.parsed = pkgs;
    },
    getNotParsed(state, pkgs) {
      state.home.notParsed = pkgs;
    },
    getVerified(state, pkgs) {
      state.home.verified = pkgs;
    },
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
