<template>
  <b-container fluid id="home-wrapper">
    <b-row id="column-names">
      <b-col><p>Unclassified Packages</p></b-col>
      <b-col><p>Manually Classified Packages</p></b-col>
      <b-col><p>Auto-detected Packages</p></b-col>
    </b-row>
    <b-row class="main-row">
      <b-col class="package-list-col">
        <b-list-group id="unclassified" class="package-list">
          <b-list-group-item
            v-for="p in this.$store.state.home.unclassified"
            @click="newTab(p.id)"
            href="#"
          >
            {{p.id}}
          </b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col class="package-list-col">
        <b-list-group id="manually-classified" class="package-list">
          <b-list-group-item href="#">Foo</b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col class="package-list-col">
        <b-list-group id="classified" class="package-list">
          <b-list-group-item href="#">Foo</b-list-group-item>
        </b-list-group>
      </b-col>
    </b-row>
  </b-container>
</template>

<style lang="scss">
  #home-wrapper {
    display: flex;
    flex-direction: column;
  }
  #column-names {
    flex-grow: 0;
    flex-shrink: 0;
  }
  #main-row {
    height: 100%;
  }
  .package-list-col {
    overflow: auto;
  }
  .package-list > .list-group-item {
    height: 30px;
    padding: 3px;
    -webkit-font-smoothing: antialiased;
  }
</style>

<script lang="ts">
import { Component, Emit, Prop, Vue } from 'vue-property-decorator';
import MyButton from '@/components/MyButton.vue';
import { HomeState, Package, Tab } from '@/model';
import { Route } from 'vue-router/types/router';
import router from '@/router';
import store from '@/store';
import axios, {AxiosResponse} from "axios";

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
    store.commit('newTab', new Tab(tabName));
    router.push('/package/' + tabName);
  }

  public onXHR() {
    axios.get(store.state.endpoint_back + "/api/packages")
      .then((res: AxiosResponse) => {
        store.commit('getUnclassified', res.data);
        console.log(res.data)
      })
  }

  public created() {
    this.onXHR();
  }

  public beforeRouteEnter(to: Route, from: Route, next: Function) {
    next((component: Home) => {
      component.counter = component.$store.state.home.lastCounter;
    });
  }

  public beforeRouteLeave(to: Route, from: Route, next: Function) {
    const hs = new HomeState();
    hs.lastCounter = this.counter;
    store.state.home = hs;
    next();
  }
}
</script>
