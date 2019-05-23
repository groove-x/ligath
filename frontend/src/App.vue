<template>
  <div id="app">
    <NavBar></NavBar>
    <router-view/>
  </div>
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
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin: 10px 10px 10px;
}
#nav {
  margin-bottom: 10px;
}
</style>
