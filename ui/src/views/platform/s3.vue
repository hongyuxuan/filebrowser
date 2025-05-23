<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-row>
      <el-col :span="12">
        <el-button-group>
          <el-button icon="refresh" size="large" style="margin-right:5px" @click="getList(1)" />
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:300px;" />
        </el-button-group>
      </el-col>
      <el-col :span="12">
        <el-button class="pull-right" size="large" type="primary" @click="show=true;edit=false;form={tenant}">+ 新建S3连接</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list"
      class="line-height40" 
      style="width:100%;margin-top:10px">
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="s3_endpoint" label="S3地址" min-width="300" />
      <el-table-column prop="update_at" label="更新时间" width="170">
        <template #default="scope">
        {{ moment(scope.row.update_at).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
      </el-table-column>
      <el-table-column prop="Option" label="操作" width="140">
        <template #default="scope">
          <el-button :icon="EditPen" circle @click="editOne(scope.row)" />
          <el-tooltip effect="dark" content="复制" placement="top">
            <el-button :icon="CopyDocument" circle @click="copyOne(scope.row)" />
          </el-tooltip>
          <el-popconfirm title="确认删除？" @confirm="deleteOne(scope.row)">
            <template #reference>
              <el-button :icon="Delete" circle />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination 
      class="pull-right"
      background 
      v-model:page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper" 
      :total="pageTotal"
      @size-change="handleSizeChange"
      @current-change="getList"
      v-model:current-page="current" />
  </div>
</div>
<el-drawer v-model="show" direction="rtl" size="600px">
  <template #header>
    <h4 v-if="edit===false">新建S3连接</h4>
    <h4 v-if="edit===true">编辑S3连接</h4>
  </template>
  <template #default>
    <el-form ref="repo" :model="form" :rules="rules" label-width="110px">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" size="large" clearable />
      </el-form-item>
      <el-form-item label="S3地址" prop="s3_endpoint">
        <el-input v-model="form.s3_endpoint" size="large" clearable />
      </el-form-item>
      <el-form-item label="S3 Region" prop="s3_region">
        <el-input v-model="form.s3_region" size="large" clearable />
      </el-form-item>
      <el-form-item label="S3 AccessKey" prop="s3_access_key">
        <el-input v-model="form.s3_access_key" size="large" clearable />
      </el-form-item>
      <el-form-item label="S3 SecretKey" prop="s3_secret_key">
        <el-input v-model="form.s3_secret_key" type="password" size="large" clearable />
      </el-form-item>
      <el-form-item label="开启Secure" prop="use_secure">
        <el-checkbox v-model="form.use_secure" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(repo)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { Search,EditPen,CopyDocument,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 变量定义 */
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const show = ref(false)
const form = ref({})
const rules = reactive({
  repo_url: [{required: true, message: '请填写仓库地址'}],
  repo_type: [{required: true, message: '请选择仓库类型', trigger: 'change'}],
})
const edit = ref(false)
const repo = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList(1)
})
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey.value != "") url += `&search=repo_url==${searchKey.value}`
  let response = await axios.get(`/filebrowser/db/s3_repository?${url}`)
  list.value = response.results.map(x => {
    if(x.use_secure === 1) 
      x.use_secure = true
    else
      x.use_secure = false
    return x
  })
  pageTotal.value = response.total
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.update_at = moment()
      if(edit.value === false) {
        await axios.post(`/filebrowser/db/s3_repository`, params)
        getList(1)
        current.value = 1
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/filebrowser/db/s3_repository/${id}`, params)
        getList(current.value)
      }
      show.value = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const editOne = async (row) => {
  form.value = Object.assign({}, row)
  edit.value = true
  show.value = true
}
const deleteOne = async (row) => {
  await axios.delete(`/filebrowser/db/s3_repository/${row.id}`)
  getList(current.value)
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  edit.value = false
  show.value = true
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
</script>