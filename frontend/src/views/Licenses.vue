<template>
  <b-container fluid id="wrapper">
    <b-row id="toolbar" class="mb-2">
      <b-button-toolbar>
        <b-button-group>
          <b-button @click="editPackage" variant="primary">Edit Package</b-button>
        </b-button-group>
      </b-button-toolbar>
    </b-row>
    <b-row id="main-row">
      <b-col class="sub-col col-4">
        <p class="input-title">Licenses</p>
        <b-row class="sub-row">
          <b-col class="rounded sub-col">
            <b-list-group class="items-list">
              <b-list-group-item
                v-for="l in this.$store.state.license.licenses"
                @click="fetchPackages(l.name)"
                href="#"
              >
                {{l.name}}
              </b-list-group-item>
            </b-list-group>
          </b-col>
        </b-row>
      </b-col>
      <b-col class="sub-col col-4">
        <p class="input-title">Filtered Packages</p>
        <b-row class="sub-row">
          <b-col class="rounded sub-col">
            <b-list-group class="items-list">
              <b-list-group-item
                v-for="p in this.$store.state.license.filteredPackages"
                @click="newTab(p.name, p.version, 'verified')"
                href="#"
              >
                {{p.name}} @ {{p.version}}
              </b-list-group-item>
            </b-list-group>
          </b-col>
        </b-row>
      </b-col>
      <b-col class="sub-col col-4">
        <p class="input-title">Lorem ipsum</p>
        <b-row class="sub-row">
          <b-col class="rounded sub-col">
            <b-list-group class="items-list">
              <b-list-group-item href="#">Lorem ipsum dolor sit amet</b-list-group-item>
            </b-list-group>
          </b-col>
        </b-row>
      </b-col>
    </b-row>
  </b-container>
</template>

<style lang="scss">
  #wrapper {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  #toolbar {
    flex-grow: 0;
    flex-shrink: 0;
  }

  #main-row {
    height: 100%;
    flex-grow: 1;
  }

  #copyright-list {
    padding: 0;
    background-color: #eee;
  }

  #raw-copyright {
    font-family: monospace;
  }

  #copyright, #range {
    height: 100px;
  }

  #license-body {
    height: 300px;
  }

  .sub-row {
    height: 100%;
    margin: 0;
  }

  .sub-col {
    display: flex;
    flex-direction: column;
    padding: 0 10px 0;
    height: 100%;

    &:nth-child(1) {
      padding-left: 0;
    }

    &:last-child {
      padding-right: 0;
    }

    .input-title {
      &:nth-child(n+2) {
        margin-top: 10px;
      }
      margin-bottom: 5px;
      padding: 2px 0 2px;
      color: #fff;
      background-color: #6c757d;
      flex-grow: 0;
      flex-shrink: 0;
      border-radius: .25rem;
    }

    .input-body {
      overflow: auto;
      height: 100%;
    }
  }

  .form-group {
    margin-bottom: 5px;
  }

  .package-list-col {
    overflow: auto;
  }

  .items-list {
    overflow: auto;
    .list-group-item {
      height: 30px;
      padding-top: 4px;
      font-size: 14px;
      -webkit-font-smoothing: antialiased;
    }
  }

  .col-form-label {
    display: inline-block;
    padding: 0;
    margin-top: 10px;
  }

  .parsed-buttons {
    padding: 0;
    float: right;
  }

  input {
    display: inline-block;
  }
</style>

<script lang="ts">
import {Component, Emit, Vue} from "vue-property-decorator";
import { Route } from 'vue-router/types/router';
import router from '@/router';
import store from '@/store';
import {Package, Package as PackageObj, FunctionTabEnum} from "@/model";
import axios, { AxiosResponse } from 'axios';

const components = {};

Component.registerHooks([
  'beforeRouteEnter',
]);

@Component({components})
export default class Licenses extends Vue {
  public package: PackageObj = new PackageObj(null);

  @Emit()
  public newTab(name: string, version: string, kind: string) {
    store.commit('newTab', {name, version, kind});
    axios.get(`${store.state.endpoint_back}/api/packages/${name}@${version}?kind=${kind}`)
      .then((res: AxiosResponse) => {
        store.commit('setPackageData', {pkg: new Package(res.data), kind: kind});
        router.push(`/package/${name}@${version}@${kind}`);
      });
  }

  public created() {
    this.fetchLicenses()
  }

  public beforeRouteEnter(to: Route, from: Route, next: (arg?: any) => void) {
    store.commit('enterFunctionTab', FunctionTabEnum.LicenseTab);
    next((component: Licenses) => {
      // component.resetCopyright(to, from, component);
    });
  }

  public fetchLicenses() {
    axios.get(
      store.state.endpoint_back
        + `/api/licenses`
    )
      .then((res: AxiosResponse) => {
        store.state.license.licenses = res.data
      })
  }

  public fetchPackages(license: string) {
    axios.get(
      store.state.endpoint_back
        + `/api/packages?license=${encodeURIComponent(license)}`
    )
      .then((res: AxiosResponse) => {
        store.state.license.filteredPackages = res.data
      })
      .catch((res: AxiosResponse) => {

      })
  }
}
</script>
