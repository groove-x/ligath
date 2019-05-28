<template>
  <b-container fluid id="package-wrapper">
    <b-row id="toolbar" class="mb-2">
      <b-button-toolbar>
        <b-button-group class="mr-1">
          <b-button variant="primary">Save</b-button>
          <b-button @click="onClose" variant="secondary">Save & Close</b-button>
        </b-button-group>
        <b-button-group>
          <b-button variant="danger">Delete</b-button>
        </b-button-group>
      </b-button-toolbar>
    </b-row>
    <b-row id="main-row">
      <b-col class="sub-col col-6">
        <p class="input-title">Raw copyright</p>
        <b-form-textarea id="copyright" class="input-body" v-model="this.package.rawCopyright"></b-form-textarea>
      </b-col>
      <b-col class="sub-col col-6">
        <p class="input-title">Parsed copyright</p>
        <b-row class="sub-row">
          <b-col class="sub-col col-4">
            <b-list-group id="unclassified" class="package-list">
              <b-list-group-item
                v-for="p in this.$store.state.packages.get(this.$route.params.name).copyrights"
                href="#"
              >
                {{p.fileRange}}, {{p.license.name}}
              </b-list-group-item>
              <b-list-group-item>+ Add new notice</b-list-group-item>
            </b-list-group>
          </b-col>
          <form class="col-8">
            <div class="form-group">
              <label for="range" class="col col-form-label">Range</label>
              <input id="range" class="form-control"/>
            </div>
            <div class="form-group">
              <label for="owner" class="col col-form-label">Owner</label>
              <input id="owner" class="form-control"/>
            </div>
            <div class="form-group">
              <label for="license" class="col col-form-label">License</label>
              <select id="license" class="form-control">
                <option></option>
              </select>
            </div>
          </form>
          <div class="col text-right parsed-buttons">
            <b-button class="btn mr-1">Save</b-button>
            <b-button class="btn-danger">Delete</b-button>
          </div>
        </b-row>
      </b-col>
    </b-row>
  </b-container>
</template>

<style lang="scss">
  #package-wrapper {
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
  .sub-row {
    margin: 10px 0 0;
  }
  .sub-col {
    display: flex;
    flex-direction: column;
    padding: 0 5px 0;
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
      margin-bottom: 0;
      padding: 2px 0 2px;
      color: #fff;
      background-color: #6c757d;
      flex-grow: 0;
      flex-shrink: 0;
    }

    .input-body {
      overflow: auto;
      height: 100%;
    }
  }
  .package-list-col {
    overflow: auto;
  }
  .package-list > .list-group-item {
    height: 30px;
    padding: 3px;
    -webkit-font-smoothing: antialiased;
  }

  #copyright {
    font-family: monospace;
  }

  .form-group {
    display: flex;
  }

  label.col-form-label {
    display: inline-block;
    padding-left: 0;
  }

  input {
    display: inline-block;
  }

  .parsed-buttons {
    float: right;
  }
</style>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Route } from 'vue-router/types/router';
import store from '@/store';

const components = {};

Component.registerHooks([
  'beforeRouteEnter',
  'beforeRouteLeave',
]);

@Component({components})
export default class Package extends Vue {
  public package: any = {};

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

  public created() {
    const pkg = store.state.packages.get(this.$route.params.name);
    if (pkg) {
      this.package = pkg;
    }
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
