<template lang="pug">
  v-app#inspire
    v-navigation-drawer(v-model="drawer" app)
      v-list
        v-list-item(link)
          v-list-item-action
            v-icon mdi-home
          v-list-item-content
            v-list-item-title Home

    v-app-bar(app)
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      v-toolbar-title Application

    v-main
      v-container.fill-height
        v-row
          v-col.text-center
            v-btn
              v-icon(large) mdi-code-tags
    v-footer(app)
      span &copy; {{ new Date().getFullYear() }}
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'

@Component({
  components: {},
})
export default class VuetifyTestPage extends Vue {
  drawer = null
}
</script>
