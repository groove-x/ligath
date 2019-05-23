import Vue from 'vue';
import Router, { Route } from 'vue-router';
import About from '@/views/About.vue';
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
      path: '/about',
      name: 'about',
      component: About,
    },
    {
      path: '/package/:name',
      name: 'package',
      component: Package,
    },
  ],
});

r.beforeEach((to: Route, from: Route, next: Function) => {
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
