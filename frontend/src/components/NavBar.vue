<template>
  <b-nav tabs id="nav">
    <b-nav-item exact class="function-tab" active-class="active" to="/"><font-awesome-icon icon="igloo" />Home</b-nav-item>
    <b-nav-item exact class="function-tab" active-class="active" to="/licenses"><font-awesome-icon icon="file-alt" />Licenses</b-nav-item>
    <b-nav-item
      exact
      active-class="active"
      v-for="[key, tab] in this.$store.getters.tabsAsList"
      :key="key"
      v-bind:to="`/package/${tab.name}@${tab.version}@${tab.kind}`">
      {{tab.kind.charAt(0).toUpperCase() + tab.kind.slice(1)}}: {{tab.name}} @ {{tab.version}}
    </b-nav-item>
  </b-nav>
</template>

<script lang="ts">
  import { Vue } from 'vue-property-decorator';

  export default class NavBar extends Vue {
  }
</script>

<style scoped>
  .function-tab svg {
    margin-right: 6px;
  }
</style>