import Vue from 'vue';
import Router, { Route } from 'vue-router';
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
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/About.vue'),
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
