<script setup lang="ts">
import {onMounted, ref} from "vue";
import {useDownloaderStore} from "../../stores/downloader";
import {GetCpuNum} from "../../../wailsjs/go/api/UtilsApi";
import CacheDirectoryInput from "./CacheDirectoryInput.vue";
import ExportDirectoryInput from "./ExportDirectoryInput.vue";

const store = useDownloaderStore();

const maxExportConcurrentCount = ref<number>(0);

onMounted(async () => {
  maxExportConcurrentCount.value = await GetCpuNum();
});

//TODO: 优化设置界面
</script>

<template>
  <div class="flex flex-col gap-y-10 p-2">
    <n-popover trigger="hover">
      <template #trigger>
        <n-input class="text-align-left flex" v-model:value="store.proxyUrl" placeholder="">
          <template #prefix>代理地址：</template>
        </n-input>
      </template>
      <span>如果不使用代理，请清空这个输入框</span>
    </n-popover>
    <cache-directory-input/>
    <export-directory-input/>
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
      <span><b>(如果你完全看不懂这个参数的描述，请使用默认值 3)</b><br/>此参数用于限制导出页面中文件树的最大深度，防止扫描缓存目录的时间过长<br/>例如选择C盘根目录作为缓存目录，如果不对深度加以限制，则会扫描整个C盘下的所有文件</span>
    </n-popover>
  </div>
</template>


