<script lang="ts" setup>
import DownloaderPane from "./components/DownloaderPane.vue"
import SettingsPane from "./components/settings/SettingsPane.vue"
import {useDownloaderStore} from "./stores/downloader"
import {onMounted} from "vue"

const store = useDownloaderStore()
onMounted(async () => {
  await store.init()
})

</script>

<template>
  <n-notification-provider placement="bottom-right">
    <n-layout has-sider sider-placement="right">
      <n-layout-content style="height: 100vh">
        <downloader-pane/>
      </n-layout-content>
      <n-layout-sider
          collapse-mode="transform"
          :collapsed-width="15"
          :width="300"
          :native-scrollbar="false"
          show-trigger="arrow-circle"
          content-style="padding: 24px;"
          bordered
          default-collapsed
          style="text-align: left;height: 100vh"
      >
        <n-flex>
          <settings-pane/>
        </n-flex>
      </n-layout-sider>
    </n-layout>
  </n-notification-provider>
</template>

