<template>
  <b-container fluid id="app">
      <NavBar></NavBar>
      <router-view/>
  </b-container>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';
  import NavBar from '@/components/NavBar.vue'

  const components = {NavBar} ;

  @Component({components})
  export default class App extends Vue {
  }
</script>

<style lang="scss">
html, body {
  height: 100%;
}
body {
  padding: 10px;
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 0;
}
#nav {
  width: 100%;
  margin-bottom: 10px;
  flex-grow: 0;
  flex-shrink: 0;
}
</style>
