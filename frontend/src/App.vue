<template>
  <div id="app">
    <b-nav tabs id="nav">
      <b-nav-item exact active-class="active" to="/"><font-awesome-icon icon="igloo" /></b-nav-item>
      <b-nav-item exact active-class="active" to="/about">Link</b-nav-item>
      <b-nav-item
        exact
        active-class="active"
        v-for="tab in this.$store.state.tabs"
        v-bind:to="'/packages/' + tab.tabName">
        {{tab.tabName}}
      </b-nav-item>
    </b-nav>
    <router-view/>
  </div>
</template>

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
