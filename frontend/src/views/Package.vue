<template>
  <b-container fluid id="wrapper">
    <b-row id="toolbar" class="mb-2">
      <b-button-toolbar>
        <b-button-group>
          <b-button @click="save" variant="secondary">Save</b-button>
          <b-button @click="close" variant="secondary">Close</b-button>
          <b-button @click="saveAndClose" variant="primary">Save & Close</b-button>
        </b-button-group>
        <b-button-group class="ml-1">
          <b-button variant="danger">Delete</b-button>
        </b-button-group>
        <b-button-group class="ml-2">
          <div
            id="state-saving"
            class="btn rounded-pill btn-outline-dark"
            v-if="$data.state === 1"
          >
            Saving
          </div>
          <div
            id="state-saved"
            class="btn rounded-pill btn-outline-primary"
            v-if="$data.state === 2"
          >
            Saved successfully
          </div>
          <div
            id="state-fail"
            class="btn rounded-pill btn-outline-danger"
            v-if="$data.state === 3"
          >
            Failed
          </div>
        </b-button-group>
      </b-button-toolbar>
    </b-row>
    <b-row id="main-row">
      <b-col class="sub-col col-6">
        <p class="input-title">Raw copyright</p>
        <b-form-textarea id="raw-copyright" class="input-body" v-model="$data.package.rawCopyright"></b-form-textarea>
      </b-col>
      <b-col class="sub-col col-6">
        <p class="input-title">Parsed copyright</p>
        <b-row class="sub-row">
          <b-col id="copyright-list" class="rounded sub-col col-3">
            <b-list-group class="package-list">
              <b-list-group-item
                v-for="(p, i) in this.$store.state.packages.get(this.$route.params.name+this.$route.params.version+this.$route.params.kind).copyrights"
                @click="editCopyright(i)"
                href="#"
              >
                {{p.license.name.trim() === "" ? "(No Name)" : p.license.name.trim()}}
              </b-list-group-item>
              <b-list-group-item @click="addCopyright" href="#">+ Add new notice</b-list-group-item>
            </b-list-group>
          </b-col>
          <b-col class="col-9 sub-col">
            <form>
              <div class="form-group">
                <label for="range" class="col col-form-label">File Range</label>
                <textarea id="range" class="form-control" v-model="$data.editingCopyright.range" v-bind:disabled="!isEditingCopyright"></textarea>
                <label for="copyright" class="col col-form-label">Copyright</label>
                <textarea id="copyright" class="form-control" v-model="$data.editingCopyright.copyright" v-bind:disabled="!isEditingCopyright"></textarea>
                <label for="license-name" class="col col-form-label">License Identifier</label>
                <input id="license-name" class="form-control" v-model="$data.editingCopyright.license.name" v-bind:disabled="!isEditingCopyright"/>
                <label for="license-body" class="col col-form-label">License Body</label>
                <textarea id="license-body" class="form-control" v-model="$data.editingCopyright.license.body" v-bind:disabled="!isEditingCopyright"></textarea>
              </div>
            </form>
            <b-row class="sub-row">
              <b-col class="col-12 text-right parsed-buttons">
                <b-button class="btn-danger">Delete</b-button>
              </b-col>
            </b-row>
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

  #import {
    margin-bottom: 10px;
    flex-grow: 0;
    flex-shrink: 0;
  }

  #state-saving {
    &:hover {
      color: #007bff;
      background-color: unset;
    }
  }

  #state-ok {
    &:hover {
      color: #007bff;
      background-color: unset;
    }
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
import { Component, Vue, Watch } from 'vue-property-decorator';
import { Route } from 'vue-router/types/router';
import store from '@/store';
import {Copyright, FunctionTabEnum, License, Package as PackageObj} from "@/model";
import axios, { AxiosResponse } from 'axios';

const components = {};

enum State {
  Nothing = 0,
  Saving,
  Succeeded,
  Failed,
}

Component.registerHooks([
  'beforeRouteEnter',
]);

@Component({components})
export default class Package extends Vue {
  public package: PackageObj = new PackageObj(null);
  public editingCopyright: any = {};
  public editingCopyrightIndex: number = 0;
  public isEditingCopyright: boolean = false;
  public state: State = State.Nothing;

  public save() {
    this.state = State.Saving;
    axios.put(
      store.state.endpoint_back
        + `/api/packages/${this.$route.params.name}@${this.$route.params.version}?kind=${this.$route.params.kind}`,
      this.package.jsonCompatible(),
    )
      .then((res: AxiosResponse) => {
        if (res.status == 200) {
          this.state = State.Succeeded;
        } else {
          this.state = State.Failed;
        }
      })
      .catch((res: AxiosResponse) => {
        this.state = State.Failed;
        console.error(res.toString());
      })
  }

  public close() {
    let index: number = 0;
    const items: any[][] = Array.from(this.$store.state.tabs);
    items.some((item: any, i) => {
      if (item[1].name === this.$route.params.name
          && item[1].version === this.$route.params.version
          && item[1].kind === this.$route.params.kind) {
        index = i;
        return true;
      }
      return false;
    });
    this.$store.commit('closeTab', items[index][1]);

    if (items.length === 1) {
      switch (this.$store.state.lastFunctionTab) {
      case FunctionTabEnum.HomeTab:
        this.$router.push({name: 'home'});
        break;
      case FunctionTabEnum.LicenseTab:
        this.$router.push({name: 'licenses'});
        break;
      case FunctionTabEnum.BulkRenameTab:
        this.$router.push({name: 'bulk'});
        break;
      }
      return;
    } else if (index === items.length - 1) {
      index -= 1;
    } else {
      index += 1;
    }
    this.$router.push({path: `/package/${items[index][1].name}@${items[index][1].version}@${items[index][1].kind}`});
  }

  public saveAndClose() {
    this.save();
    this.close();
  }

  public addCopyright() {
    const newlen: number = this.package.copyrights.push(this.createDummyCopyright());

    this.editingCopyright = this.package.copyrights[newlen - 1];
    this.editingCopyrightIndex = newlen - 1;
  }

  public editCopyright(i: number) {
    this.editingCopyright = this.package.copyrights[i];
    this.editingCopyrightIndex = i;
    this.isEditingCopyright = true;
  }

  public created() {
    console.log(State.Succeeded);
    this.editingCopyright = this.createDummyCopyright();
  }

  public beforeRouteEnter(to: Route, from: Route, next: (arg?: any) => void) {
    next((component: Package) => {
      component.resetCopyright(to, from, component);
    });
  }

  @Watch('$route')
  public updateRoute(to: Route, from: Route) {
    this.resetCopyright(to, from, this);
    this.editingCopyright = this.createDummyCopyright();
  }

  private resetCopyright(to: Route, from: Route, component: Package) {
    const pkg = store.state.packages.get(to.params.name + to.params.version + to.params.kind);
    if (pkg) {
      component.package = pkg;
    }
  }

  private createDummyCopyright(): Copyright {
    return new Copyright({
      copyright: '',
      fileRange: new Array<string>(),
      license: new License({
        name: '',
        machineReadableName: '',
        body: '',
      }),
    });
  }
}
</script>
