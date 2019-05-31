<template>
  <b-container fluid id="home-wrapper">
    <b-row id="column-names">
      <b-col><p>Auto-detected Packages</p></b-col>
      <b-col><p>Unclassified Packages</p></b-col>
      <b-col><p>Manually Classified Packages</p></b-col>
    </b-row>
    <b-row class="main-row">
      <b-col class="package-list-col">
        <b-list-group id="parsed" class="package-list">
          <b-list-group-item
            v-for="p in this.$store.state.home.parsed"
            @click="newTab(p.name, p.version)"
            href="#"
          >
            {{p.name}} @ {{p.version}}
          </b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col class="package-list-col">
        <b-list-group id="not-parsed" class="package-list">
          <b-list-group-item
            v-for="p in this.$store.state.home.notParsed"
            @click="newTab(p.name, p.version)"
            href="#"
          >
            {{p.name}} @ {{p.version}}
          </b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col class="package-list-col">
        <b-list-group id="manual" class="package-list">
          <b-list-group-item
            v-for="p in this.$store.state.home.manual"
            @click="newTab(p.name, p.version)"
            href="#"
          >
            {{p.name}} @ {{p.version}}
          </b-list-group-item>
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
import axios, {AxiosResponse} from 'axios';

const components = { MyButton };

Component.registerHooks([
  'beforeRouteEnter',
  'beforeRouteLeave',
]);

@Component({components})
export default class Home extends Vue {
  public counter: number = 1;

  @Prop()
  public name!: string;

  @Prop()
  public version!: string;

  @Emit()
  public newTab(name: string, version: string) {
    store.commit('newTab', {name, version});
    axios.get(store.state.endpoint_back + '/api/packages/' + name + '@' + version)
      .then((res: AxiosResponse) => {
        store.commit('setPackageData', res.data);
        router.push('/package/' + name + '@' + version);
      });
  }

  public onXHR() {
    axios.get(store.state.endpoint_back + '/api/packages?kind=parsed')
      .then((res: AxiosResponse) => {
        store.commit('getParsed', res.data);
      });
    axios.get(store.state.endpoint_back + '/api/packages?kind=notparsed')
      .then((res: AxiosResponse) => {
        store.commit('getNotParsed', res.data);
      });
    axios.get(store.state.endpoint_back + '/api/packages?kind=manual')
      .then((res: AxiosResponse) => {
        store.commit('getManual', res.data);
      });
  }

  public created() {
    this.onXHR();
  }

  public beforeRouteEnter(to: Route, from: Route, next: (arg?: any) => void) {
    next((component: Home) => {
      component.counter = component.$store.state.home.lastCounter;
    });
  }

  public beforeRouteLeave(to: Route, from: Route, next: (arg?: any) => void) {
    const hs = new HomeState();
    hs.lastCounter = this.counter;
    store.state.home = hs;
    next();
  }
}
</script>
