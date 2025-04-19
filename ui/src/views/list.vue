<template>
<el-alert v-show="show.donwload_status" 
  type="success" 
  :closable="true"
  @close="show.donwload_status=false"
  class="top-fixed-alert">
  <div style="display: flex; align-items: center;">
    <span style="margin-right: 10px;">下载进度：</span>
    <el-progress v-if="downloadProgress<100" :percentage="downloadProgress" style="flex:1;width:380px;" />
    <el-progress v-else :percentage="downloadProgress" status="success" style="flex:1;width:380px;" />
  </div>
</el-alert>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>文件列表</el-breadcrumb-item>
</el-breadcrumb>
<el-backtop :right="80" :bottom="80" />
<el-card>
  <template #header>
    <span v-if="route.query.source!=='s3'">文件列表</span>
    <span v-else>对象存储：{{ s3Info.name }} ( {{ s3Info.s3_endpoint }} )</span>
  </template>
  <el-breadcrumb :separator-icon="ArrowRight" style="background-color:#f5f5f5;padding:15px">
    <el-breadcrumb-item v-for="(item,i) in paths" :key="i">
      <el-link :underline="false" @click="clickBreadcrumb(item,i)">{{ item.Name }}</el-link>
    </el-breadcrumb-item>
  </el-breadcrumb>
  <el-row>
    <el-col :span="16">
      <el-button-group>
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(currentpath)" />
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="search()" clearable style="width:300px;" />
      </el-button-group>
    </el-col>
    <el-col :span="8">
      <el-button-group class="pull-right">
        <el-button size="large" type="primary" @click="show.upload=true" style="margin-right:5px"><el-icon><Upload /></el-icon> 上传文件</el-button>
        <el-button size="large" @click="mkdir" >+ 创建目录</el-button>
      </el-button-group>
    </el-col>
  </el-row>
  <el-table :data="list" class="line-height40" v-loading="loading" element-loading-text="奋力加载中...">
    <el-table-column width=30>
      <template #default="scope">
        <font-awesome-icon icon="folder" v-if="scope.row.IsDir" />
      </template>
    </el-table-column>
    <el-table-column label="文件/目录" prop="Name" sortable min-width="300">
      <template #default="scope">
        <el-link v-if="scope.row.IsDir" :underline="false" @click="gotoDir(scope.row)">{{ scope.row.Name }}</el-link>
        <span v-else>{{ scope.row.Name }}</span>
      </template>
    </el-table-column>
    <el-table-column label="修改时间" prop="LastUpdate" sortable width="170">
      <template #default="scope">
        {{ scope.row.LastUpdate ? moment(scope.row.LastUpdate).format('YYYY-MM-DD HH:mm:ss') : "" }}
      </template>
    </el-table-column>
    <el-table-column label="大小" prop="Size" sortable width="100">
      <template #default="scope">{{ scope.row.IsDir ? "" : filesize(scope.row.Size) }}</template>
    </el-table-column>
    <el-table-column label="操作" width="60">
      <template #default="scope">
        <el-button circle v-if="!scope.row.IsDir" :icon="Download" @click="download(scope.row)" />
      </template>
    </el-table-column>
  </el-table>
  <el-pagination 
      class="pull-right"
      background 
      layout="total" 
      :total="pageTotal" />
</el-card>
<el-drawer v-model="show.upload" direction="rtl" size="600px">
  <template #header>
    <h4>上传文件</h4>
  </template>
  <el-form ref="uploadForm" :model="form" label-width="80px">
    <el-form-item label="文件属主" v-if="route.query.source!=='s3'">
      <el-input v-model="form.uid" size="large" clearable placeholder="linux请填写uid，windows留空" />
    </el-form-item>
    <el-form-item label="文件属组" v-if="route.query.source!=='s3'">
      <el-input v-model="form.gid" size="large" clearable placeholder="linux请填写gid，windows留空" />
    </el-form-item>
    <el-form-item label="选择文件">
      <el-upload 
        ref="uploadRef" 
        drag 
        :data="formData" 
        :action="`${route.query.source==='s3'?'/filebrowser/s3/upload':'/filebrowser/upload'}`"
        :headers="{
          'Authorization': `Bearer ${accessToken}`
        }"
        show-file-list
        multiple
        :auto-upload="false"
        :before-upload="beforeUpload"
        :on-success="uploadSuccess"
        :on-error="uploadFailed"
        style="width:100%">
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          拖拽文件到此或者<em>点我选择文件</em>
        </div>
      </el-upload>
    </el-form-item>
  </el-form>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.upload=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(uploadForm)">上传</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { ref, onMounted, computed } from "vue";
import { ArrowRight, Download, Refresh, Search, Upload, Share  } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '../assets/util/axios'
import moment from 'moment'
import { filesize } from "filesize"
/* 变量定义 */
const accessToken = computed(() => {
  return localStorage.access_token
})
const route = useRoute()
const all = ref([])
const list = ref([])
const paths = ref([]) 
const rootpath = ref("")
const currentpath = ref("")
const searchKey = ref("")
const show = ref({
  upload: false,
  donwload_status: false
})
const uploadRef = ref(null)
const uploadForm = ref(null)
const loading = ref(false)
const pageTotal = ref(0)
const form = ref({})
const formData = computed(() => {
  let data = {
    path: currentpath.value
  }
  if(form.value.uid) data.uid = form.value.uid
  if(form.value.gid) data.gid = form.value.gid
  if(route.query.source === 's3') data.name = route.query.name
  return data
})
const origin = ref("")
const downloadProgress = ref(0)
const s3Info = ref({})
onMounted(async () => {
  origin.value = window.location.origin
  rootpath.value = route.query.path
  if(!route.query.path) {
    rootpath.value = await axios.get(`/filebrowser/rootpath`)
  }
  paths.value.push({
    Path: rootpath.value,
    Name: rootpath.value
  })
  if(route.query.source === 's3') getS3()
  getList(rootpath.value)
})
/* methods */
const getS3 = async () => {
  let response = await axios.get(`/filebrowser/db/s3_repository?filter=name==${route.query.name}`)
  if(response.total > 0) {
    s3Info.value = response.results[0]
  }
}
const getList = async (path) => {
  currentpath.value = path
  loading.value = true
  var response
  if(route.query.source === 's3' && route.query.name)
    response = await axios.get(`/filebrowser/s3/listobjects?name=${route.query.name}&path=${path}`)
  else
    response = await axios.get(`/filebrowser/listfile?path=${path}`)
  response.sort((a, b) => {
    // 规则1：目录优先
    if (a.IsDir !== b.IsDir) {
      return b.IsDir - a.IsDir;
    }
    
    // 规则2：英文优先
    const isAChinese = /[\u4e00-\u9fa5]/.test(a.Name);
    const isBChinese = /[\u4e00-\u9fa5]/.test(b.Name);
    
    if (isAChinese !== isBChinese) {
      return isAChinese - isBChinese; // 修改这里，使英文排在前面
    }
    
    // 规则3：字符串顺序
    return a.Name.localeCompare(b.Name, isAChinese ? 'zh' : 'en');
  })
  all.value = response.map(x => {
    if(x.LastUpdate === "0001-01-01T00:00:00Z")
      delete x.LastUpdate
    return x
  })
  search()
}
const search = () => {
  if(!searchKey.value) {
    list.value = all.value
  } else {
    list.value = all.value.filter(n => n.Name.includes(searchKey.value))
  }
  pageTotal.value = list.value.length
  loading.value = false
}
const gotoDir = (row) => {
  paths.value.push(row)
  getList(row.Path)
}
const clickBreadcrumb = (item, index) => {
  paths.value.splice(index+1, paths.value.length-index)
  getList(item.Path)
}
const download = async (row) => {
  if(route.query.source === 's3') 
    downloadFile(`/filebrowser/s3/downloadobject?name=${route.query.name}&path=${row.Path}`, {}, row.Name)
  else
    downloadFile(`/filebrowser/download?file=${row.Path}`, {}, row.Name)
}
const beforeUpload = () => {
  ElMessage("上传任务已提交，请等待")
  console.log(formData.value)
}
const uploadSuccess = (res) => {
  ElMessage.success(res.message)
  getList(currentpath.value)
}
const uploadFailed = (err) => {
  ElMessage.error(`上传失败: ${err}`)
}
const confirmClick = async () => {
  uploadRef.value.submit()
}
const downloadFile = async (url, params, filename = 'file') => {
  try {
    show.value.donwload_status = true
    const response = await fetch(url, {
      method: 'POST', // 或 'GET'
      body: JSON.stringify(params),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.access_token}`,
      }
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || '下载失败');
    }
    const contentLength = parseInt(response.headers.get('content-length'));
    const reader = response.body.getReader();
    const chunks = [];
    let receivedLength = 0;
    while(true) {
      const {done, value} = await reader.read();
      if (done) break;
      chunks.push(value);
      receivedLength += value.length;
      downloadProgress.value = Math.ceil((receivedLength / contentLength) * 100)
    }
    // 合并chunks
    const blob = new Blob(chunks);
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = filename;
    link.click();
  } catch (error) {
    ElMessage.error({message: "下载文件失败：" + error})
  } finally {
    show.value.donwload_status = false
  }
}
const mkdir = async () => {
  if(route.query.source === 's3') {
    ElMessageBox.prompt(`当前目录是：${currentpath.value}`, '创建目录', {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      inputPattern:
        /^(?!\/).*\/$/,
      inputErrorMessage: '目录不能以 / 开头，必须以 / 结尾',
    }).then(async ({value}) => {
      let curpath = currentpath.value
      if(!curpath.endsWith("/")) curpath += "/"
      await axios.post(`/filebrowser/s3/mkdir?path=${curpath + value}&name=${route.query.name}`)
      getList(currentpath.value)
    }).catch(() => {})
  }
  else {
    ElMessageBox.prompt(`当前目录是：${currentpath.value}`, '创建目录', {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      inputPattern:
        /\/.*/,
      inputErrorMessage: '目录必须以 / 开头',
    }).then(async ({value}) => {
      let curpath = currentpath.value
      if(!curpath.endsWith("/")) curpath += "/"
      await axios.post(`/filebrowser/mkdir?path=${curpath + value}`)
      getList(currentpath.value)
    }).catch(() => {})
  }
}
</script>