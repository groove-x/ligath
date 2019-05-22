<template>
  <div class="package">
    <p>foo</p>
    <b-button @click="onClose">Close</b-button>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Route } from 'vue-router/types/router';

const components = {};

Component.registerHooks([
  'beforeRouteEnter',
  'beforeRouteLeave',
]);

@Component({components})
export default class Package extends Vue {
  public counter: number = 1;

  public onClose() {
    let index: number = 0;
    const items: any[][] = Array.from(this.$store.state.tabs);
    items.some((item: any, i) => {
      if (item[0] === this.$route.params.name) {
        index = i;
        return true;
      }
      return false;
    });
    this.$store.commit('closeTab', this.$route.params.name);

    if (items.length === 1) {
      this.$router.push({name: 'home'});
      return;
    } else if (index === items.length - 1) {
      index -= 1;
    } else {
      index += 1;
    }
    this.$router.push({path: '/package/' + items[index][1].tabName});
  }

  public beforeRouteEnter(to: Route, from: Route, next: Function) {
    next((component: Package) => {
    });
  }

  public beforeRouteLeave(to: Route, from: Route, next: Function) {
    next();
  }
}
</script>
