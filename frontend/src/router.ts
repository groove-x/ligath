import Vue from 'vue';
import Router, { Route } from 'vue-router';
import Licenses from '@/views/Licenses.vue';
import Home from '@/views/Home.vue';
import Package from '@/views/Package.vue';
import store from '@/store';

Vue.use(Router);

const r = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/licenses',
      name: 'licenses',
      component: Licenses,
    },
    {
      path: '/package/:name@:version@:kind',
      name: 'package',
      component: Package,
    },
  ],
});

r.beforeEach((to: Route, from: Route, next: (arg?: any) => void) => {
  if (!store.state.fresh) {
    next();
    return;
  }

  if (to.name === 'package' || to.name === 'about') {
    next({name: 'home'});
  } else if (to.name === 'home') {
    store.state.fresh = false;
    next();
  } else {
    next();
  }
});

export default r;
