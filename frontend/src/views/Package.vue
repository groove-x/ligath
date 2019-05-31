<template>
  <b-container fluid id="wrapper">
    <b-row id="toolbar" class="mb-2">
      <b-button-toolbar>
        <b-button-group class="mr-1">
          <b-button variant="primary">Save</b-button>
          <b-button @click="saveAndClose" variant="secondary">Save & Close</b-button>
        </b-button-group>
        <b-button-group>
          <b-button variant="danger">Delete</b-button>
        </b-button-group>
      </b-button-toolbar>
    </b-row>
    <b-row id="main-row">
      <b-col class="sub-col col-6">
        <p class="input-title">Raw copyright</p>
        <b-form-textarea id="raw-copyright" class="input-body" v-model="this.package.rawCopyright"></b-form-textarea>
      </b-col>
      <b-col class="sub-col col-6">
        <p class="input-title">Parsed copyright</p>
        <b-row class="sub-row">
          <b-col id="copyright-list" class="rounded sub-col col-4">
            <b-list-group class="package-list">
              <b-list-group-item
                v-for="(p, i) in this.$store.state.packages.get(this.$route.params.name+this.$route.params.version).copyrights"
                @click="editCopyright(i)"
                href="#"
              >
                {{p.license.name}}
              </b-list-group-item>
              <b-list-group-item @click="addCopyright" href="#">+ Add new notice</b-list-group-item>
            </b-list-group>
          </b-col>
          <form class="col-8 sub-col">
            <div class="form-group">
              <label for="range" class="col col-form-label">Range</label>
              <textarea id="range" class="form-control" v-model="this.editingCopyright.range"></textarea>
            </div>
            <div class="form-group">
              <label for="copyright" class="col col-form-label">Copyright</label>
              <textarea id="copyright" class="form-control" v-model="this.editingCopyright.copyright"></textarea>
            </div>
            <div class="form-group">
              <label for="license" class="col col-form-label">License</label>
              <select id="license" class="form-control" v-model="this.editingCopyright.license">
                <option></option>
              </select>
            </div>
            <div class="col text-right parsed-buttons">
              <b-button class="btn mr-1">Save</b-button>
              <b-button class="btn-danger">Delete</b-button>
            </div>
          </form>
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
    height: 200px;
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

  .package-list {
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
import { Component, Vue } from 'vue-property-decorator';
import { Route } from 'vue-router/types/router';
import store from '@/store';
import { Copyright } from '@/model';

const components = {};

Component.registerHooks([
  'beforeRouteEnter',
  'beforeRouteLeave',
]);

@Component({components})
export default class Package extends Vue {
  public package: any = {};
  public editingCopyright: any = {};
  public editingCopyrightIndex: number = 0;

  public saveAndClose() {
    let index: number = 0;
    const items: any[][] = Array.from(this.$store.state.tabs);
    items.some((item: any, i) => {
      if (item[1].name === this.$route.params.name
          && item[1].version === this.$route.params.version) {
        index = i;
        return true;
      }
      return false;
    });
    this.$store.commit('closeTab', items[index][1]);

    if (items.length === 1) {
      this.$router.push({name: 'home'});
      return;
    } else if (index === items.length - 1) {
      index -= 1;
    } else {
      index += 1;
    }
    this.$router.push({path: '/package/' + items[index][1].name + '@' + items[index][1].version});
  }

  public addCopyright() {
    const newlen: number = this.package.copyrights.push(new Copyright({
      notice: '',
      fileRange: '',
      license: '',
    }));

    this.editingCopyright = this.package.copyrights[newlen - 1];
    this.editingCopyrightIndex = newlen - 1;
  }

  public editCopyright(i: number) {
    this.editingCopyright = this.package.copyrights[i];
    this.editingCopyrightIndex = i;
  }

  public created() {
    const pkg = store.state.packages.get(this.$route.params.name+this.$route.params.version);
    if (pkg) {
      this.package = pkg;
    }
  }

  public beforeRouteEnter(to: Route, from: Route, next: (arg?: any) => void) {
    next((component: Package) => {});
  }

  public beforeRouteLeave(to: Route, from: Route, next: (arg?: any) => void) {
    next();
  }
}
</script>
