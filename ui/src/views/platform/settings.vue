<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-collapse v-model="activeNames">
      <el-collapse-item name="1" v-if="role==='admin'">
        <template #title><h4><b>文件预览大小限制</b></h4></template>
        <el-row>
          <el-col :span="15">浏览本地文件或S3上的文件时，文件大小的最大限制</el-col>
          <el-col :span="9" >
            <el-input v-model="settings.preview_file_size.setting_value" size="large" clearable style="width:400px" placeholder="按回车提交" @change="setValue('preview_file_size')" />
          </el-col>
        </el-row>
      </el-collapse-item>
    </el-collapse>
  </div>
</div>
</template>
<script setup>
import { onBeforeMount, ref, computed } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { useStore } from 'vuex'
import axios from 'axios'
/* 变量定义 */
const store = useStore()
const role = computed(() => {
  return store.state.userInfo.role
})
const activeNames = ref(["1"])
const settings = ref({
  preview_file_size: {}
})
/* 生命周期函数 */
onBeforeMount(async () => {
  getSettings()
})
 /* methods */
const getSettings = async () => {
  let response = await axios.get(`/filebrowser/db/settings?size=1000`)
  for(let x of response.results) {
    if(x.setting_value === 'true' || x.setting_value === 'false')
      x.setting_value = JSON.parse(x.setting_value)
    settings.value[x.setting_key] = x
  }
}
const setValue = async (setting_key) => {
  await axios.put(`/filebrowser/db/settings/${settings.value[setting_key].id}`, {
    setting_value: settings.value[setting_key].setting_value
  })
}
</script>