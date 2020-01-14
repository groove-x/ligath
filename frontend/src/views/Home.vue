<template>
  <b-container fluid id="home-wrapper">
    <b-row id="column-names">
      <b-col><p>Parsed Packages  ({{this.$store.state.home.parsed.length}})</p></b-col>
      <b-col><p>Not-parsed Packages  ({{this.$store.state.home.notParsed.length}})</p></b-col>
      <b-col><p>Verified Packages  ({{this.$store.state.home.verified.length}})</p></b-col>
    </b-row>
    <b-row class="main-row">
      <b-col class="package-list-col">
        <b-list-group id="parsed" class="package-list">
          <b-list-group-item
            v-for="p in this.$store.state.home.parsed"
            v-bind:class="{ verified: p.verified }"
            @click="newTab(p.name, p.version, 'parsed')"
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
            v-bind:class="{ verified: p.verified }"
            @click="newTab(p.name, p.version, 'notparsed')"
            href="#"
          >
            {{p.name}} @ {{p.version}}
          </b-list-group-item>
        </b-list-group>
      </b-col>
      <b-col class="package-list-col">
        <b-list-group id="manual" class="package-list">
          <b-list-group-item
            v-for="p in this.$store.state.home.verified"
            @click="newTab(p.name, p.version, 'verified')"
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

    &.verified {
      background-color: #DFD;
    }
  }
</style>

<script lang="ts">
import { Component, Emit, Prop, Vue } from 'vue-property-decorator';
import MyButton from '@/components/MyButton.vue';
import {FunctionTabEnum, HomeState, Package, Tab} from "@/model";
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
  public newTab(name: string, version: string, kind: string) {
    store.commit('newTab', {name, version, kind});
    axios.get(`${store.state.endpoint_back}/api/packages/${name}@${version}?kind=${kind}`)
      .then((res: AxiosResponse) => {
        store.commit('setPackageData', {pkg: new Package(res.data), kind: kind});
        router.push(`/package/${name}@${version}@${kind}`);
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
    axios.get(store.state.endpoint_back + '/api/packages?kind=verified')
      .then((res: AxiosResponse) => {
        store.commit('getVerified', res.data);
      });
  }

  public created() {
    this.onXHR();
  }

  public beforeRouteEnter(to: Route, from: Route, next: (arg?: any) => void) {
    store.commit('enterFunctionTab', FunctionTabEnum.HomeTab);
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
