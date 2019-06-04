<template>
  <b-nav tabs id="nav">
    <b-nav-item exact active-class="active" to="/"><font-awesome-icon icon="igloo" /></b-nav-item>
    <b-nav-item
      exact
      active-class="active"
      v-for="[key, tab] in this.$store.getters.tabsAsList"
      :key="key"
      v-bind:to="`/package/${tab.name}@${tab.version}@${tab.kind}`">
      {{tab.kind}}: {{tab.name}} @ {{tab.version}}
    </b-nav-item>
  </b-nav>
</template>

<script lang="ts">
  import { Vue } from 'vue-property-decorator';

  export default class NavBar extends Vue {
  }
</script>

<style scoped>
</style>