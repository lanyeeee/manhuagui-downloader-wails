<script setup lang="ts">
import {onMounted, ref} from "vue"
import {useNotification} from "naive-ui"
import {useDownloaderStore} from "../../stores/downloader"
import {GetCpuNum} from "../../../wailsjs/go/api/UtilsApi"
import {ChooseDirectory} from "../../../wailsjs/go/api/SettingsApi"

const store = useDownloaderStore()
const notification = useNotification()

const maxExportConcurrentCount = ref<number>(0)

onMounted(async () => {
  maxExportConcurrentCount.value = await GetCpuNum()
})

async function onChooseCacheDirectory() {
  const response = await ChooseDirectory(store.cacheDirectory)
  if (response.code != 0) {
    notification.create({type: "error", title: "选择缓存目录失败", meta: response.msg,})
  } else if (response.data !== "") {
    store.cacheDirectory = response.data as string
    store.downloadTreeOptions = []
  }
}

async function onChooseExportDirectory() {
  const response = await ChooseDirectory(store.exportDirectory)
  if (response.code != 0) {
    notification.create({type: "error", title: "选择导出目录失败", meta: response.msg,})
  } else if (response.data !== "") {
    store.exportDirectory = response.data as string
  }
}

</script>

<template>
  <n-flex vertical justify="space-around" style="height: 90vh">
    <n-popover trigger="hover">
      <template #trigger>
        <n-input v-model:value="store.proxyUrl" placeholder="">
          <template #prefix>代理地址：</template>
        </n-input>
      </template>
      <span>如果不使用代理，请清空这个输入框</span>
    </n-popover>
    <n-popover trigger="hover">
      <template #trigger>
        <n-input v-model:value="store.cacheDirectory" readonly placeholder=""
                 @click="onChooseCacheDirectory">
          <template #prefix>缓存目录：</template>
        </n-input>
      </template>
      <span>{{ store.cacheDirectory }}</span>
    </n-popover>
    <n-popover trigger="hover">
      <template #trigger>
        <n-input v-model:value="store.exportDirectory" placeholder="" readonly
                 @click="onChooseExportDirectory">
          <template #prefix>导出目录：</template>
        </n-input>
      </template>
      <span>{{ store.exportDirectory }}</span>
    </n-popover>
    <n-popover trigger="hover">
      <template #trigger>
        <n-input-number v-model:value="store.downloadInterval" placeholder="0" :min="0" :max="600">
          <template #prefix>
            下载间隔：
          </template>
          <template #suffix>
            秒
          </template>
        </n-input-number>
      </template>
      <span>每章漫画下载完成后暂停的时间，不设置间隔或间隔太短则容易被ban IP</span>
    </n-popover>
    <n-popover trigger="hover">
      <template #trigger>
        <n-input-number v-model:value="store.downloadConcurrentCount" :min="1" :max="5">
          <template #prefix>
            下载并发数：
          </template>
        </n-input-number>
      </template>
      <span>下载某个章节时同时下载的图片数量。章节不支持并发下载(并发下载章节容易被ban IP)</span>
    </n-popover>
    <n-popover trigger="hover">
      <template #trigger>
        <n-input-number v-model:value="store.exportConcurrentCount" :min="1" :max="maxExportConcurrentCount">
          <template #prefix>
            导出并发数：
          </template>
        </n-input-number>
      </template>
      <span>生成PDF时的并发数</span>
    </n-popover>
    <n-popover trigger="hover">
      <template #trigger>
        <n-input-number v-model:value="store.exportTreeMaxDepth" :min="0">
          <template #prefix>
            导出树最大深度：
          </template>
        </n-input-number>
      </template>
      <span>此参数用于限制导出树最大深度，防止扫描缓存目录的时间过长</span>
    </n-popover>
  </n-flex>
</template>


