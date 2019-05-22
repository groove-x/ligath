import Vue from 'vue';
import Vuex from 'vuex';
import { HomeState } from '@/model';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    tabs: new Array(0),
    home: new HomeState(),
  },
  mutations: {
    newTab(state, tab) {
      state.tabs.push(tab);
    },
  },
  actions: {

  },
});
