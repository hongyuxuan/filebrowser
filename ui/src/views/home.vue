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
    <el-button-group v-if="role==='admin'" class="pull-right">
      <el-button @click="currentS3Name=k;show.create=true" type="primary">+ 创建桶</el-button>
    </el-button-group>
  </template>
  <el-space wrap alignment="start" :size="15">
    <el-card class="pointer hover-container" style="position:relative;width:300px;"  v-for="(item,i) in v.buckets" :key="i" @click="gotoObjectList(k, item)">
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
<el-dialog v-model="show.create" title="创建桶" width="600px">
  <el-form :model="form" :rules="rules" ref="bucketRef" label-width="100" style="margin-top:15px">
    <el-form-item label="桶名" prop="bucket_name">
      <el-input v-model="form.bucket_name" size="large" />
    </el-form-item>
    <el-form-item label="启用版本控制" prop="versioning">
      <el-checkbox v-model="form.versioning" />
    </el-form-item>
    <el-form-item label="启用对象锁定" prop="object_locking">
      <el-checkbox v-model="form.object_locking" />
    </el-form-item>
  </el-form>
  <template #footer>
    <div class="dialog-footer">
      <el-button @click="show.create=false">取消</el-button>
      <el-button type="primary" @click="createBucket(bucketRef)">提交</el-button>
    </div>
  </template>
</el-dialog>
</template>
<script setup>
import { ref, onMounted, reactive, computed } from "vue"
import { ArrowRight, DeleteFilled } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import {axios} from '../assets/util/axios'
const router = useRouter()
const store = useStore()
const role = computed(() => {
  return store.state.userInfo?.role
})
const drivers = ref([])
const s3list = ref({})
const show = ref({
  create: false
})
const form = ref({})
const currentS3Name = ref("")
const rules = reactive({
  bucket_name: [{required: true, message: '请填写桶名称'}],
})
const bucketRef = ref(null)
/* 生命周期函数 */
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
const createBucket = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      await axios.post(`/filebrowser/s3/createbucket?name=${currentS3Name.value}`, form.value)
      getS3()
      show.value.create = false
    }
  })
}
</script>