import Vue from 'vue';
import Vuex from 'vuex';
import { HomeState } from '@/model';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    endpoint_front: window.location.hostname + ":" + window.location.port,
    endpoint_back: "http://" + window.location.hostname + ":3939",
    fresh: true,
    tabs: new Map<string, any>(),
    tabsTick: 0,
    home: new HomeState(),
  },
  mutations: {
    newTab(state, tab) {
      // increment tabsTick to make reaction happen
      state.tabsTick += 1;
      state.tabs.set(tab.tabName, tab);
    },
    closeTab(state, name) {
      state.tabsTick += 1;
      state.tabs.delete(name);
    },
    fresh(state) {
      state.fresh = false;
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
