<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    设备和驱动器
  </template>
  <el-space wrap alignment="start" :size="15">
    <el-card class="pointer" style="width:300px;"  v-for="(item,i) in drivers" :key="i" @click="gotoFileList(item)">
      <table>
        <tbody>
        <tr>
          <td width="45"><font-awesome-icon icon="inbox" style="font-size:30px" /></td>
          <td width="230">
            <div style="margin-bottom:8px"><span v-if="item.VolumeName">{{ item.VolumeName }}</span> ( {{ item.DriverName }} )</div>
            <div v-if="item.TotalSpace" style="margin-bottom:8px"><el-progress :percentage="item.UsedSpace/item.TotalSpace*100" :show-text="false" /></div>
            <div v-if="item.TotalSpace">{{ item.FreeSpaceGB}} 可用，共 {{ item.TotalSpaceGB }}</div>
          </td>
        </tr>
        </tbody>
      </table>
    </el-card>
  </el-space>
</el-card>
<el-card v-for="(v,k,i) in s3list" :key="i">
  <template #header>
    对象存储：{{ k }} ( {{ v.s3_endpoint }} )
  </template>
  <el-space wrap alignment="start" :size="15">
    <el-card class="pointer" style="width:300px;"  v-for="(item,i) in v.buckets" :key="i" @click="gotoObjectList(k, item)">
      <table>
        <tbody>
        <tr>
          <td width="45"><font-awesome-icon icon="database" style="font-size:30px" /></td>
          <td width="230">
            <div style="margin-bottom:8px">{{ item.name }}</div>
          </td>
        </tr>
        </tbody>
      </table>
    </el-card>
  </el-space>
</el-card>
</template>
<script setup>
import { ref, onMounted } from "vue"
import { ArrowRight } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import {axios} from '../assets/util/axios'
const router = useRouter()
const drivers = ref([])
const s3list = ref({})
onMounted(async () => {
  getDrivers()
  getS3()
})
/* methods */
const getDrivers = async () => {
  drivers.value = await axios.get(`/filebrowser/drivers`)
}
const gotoFileList = (item) => {
  router.push({
    path: 'list',
    query: {
      path: item.DriverName
    }
  })
}
const gotoObjectList = (name, item) => {
  router.push({
    path: 'list',
    query: {
      path: item.name + '/',
      source: 's3',
      name
    }
  })
}
const getS3 = async () => {
  let response = await axios.get(`/filebrowser/s3/listconnections`)
  for(let [k,v] of Object.entries(response)) {
    s3list.value[k] = {}
    let r = await axios.get(`/filebrowser/s3/listbuckets?name=${k}`)
    s3list.value[k] = {
      s3_endpoint: v.s3_endpoint,
      buckets: r
    }
  }
}
</script>