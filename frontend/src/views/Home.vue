<template>
  <div class="home">
    <MyButton :tabname="tabName" @clicked="onButtonClicked"></MyButton>
  </div>
</template>

<script lang="ts">
import { Component, Emit, Prop, Vue } from 'vue-property-decorator';
import MyButton from '@/components/MyButton.vue';
import { HomeState, Tab } from '@/model';
import { Route } from 'vue-router/types/router';

const components = { MyButton };

Component.registerHooks([
  'beforeRouteEnter',
  'beforeRouteLeave',
]);

@Component({components})
export default class Home extends Vue {
  public counter: number = 1;

  @Prop()
  public tabName!: string;

  @Emit()
  public newTab(tabName: string) {
    this.$store.commit('newTab', new Tab(this.counter.toString()));
  }

  public onButtonClicked() {
    this.newTab(this.counter.toString());
    this.counter += 1;
  }

  public beforeRouteEnter(to: Route, from: Route, next: Function) {
    next((component: Home) => {
      component.counter = component.$store.state.home.lastCounter;
    });
  }

  public beforeRouteLeave(to: Route, from: Route, next: Function) {
    const hs = new HomeState();
    hs.lastCounter = this.counter;
    this.$store.state.home = hs;
    next();
  }
}
</script>