<template>
  <div>
    <b-button @click="clicked" variant="primary">New Tab</b-button>
  </div>
</template>

<script lang="ts">
import {Component, Emit, Vue} from 'vue-property-decorator';

@Component
export default class MyButton extends Vue {
  @Emit()
  public clicked() {
  }
}
</script>

<style scoped lang="scss">
</style>