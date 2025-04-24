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
    <el-button-group v-if="role==='admin'&&route.query.source==='s3'" class="pull-right">
      <el-button @click="deleteBucket()" type="danger" :icon="Delete">删除桶</el-button>
    </el-button-group>
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
        <el-link v-else @click="preview(scope.row)" :underline="false">{{ scope.row.Name }}</el-link>
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
    <el-table-column label="操作" width="100">
      <template #default="scope">
        <el-button circle :disabled="scope.row.IsDir" :icon="Download" @click="download(scope.row)" />
        <el-popconfirm title="确定删除?" confirm-button-text="确认" cancel-button-text="取消" @confirm="deleteOne(scope.row)">
          <template #reference>
            <el-button icon="Delete" circle :disabled="role!=='admin'" />
          </template>
        </el-popconfirm>
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
<el-image-viewer v-if="show.preview_image" 
  :url-list="imageList.map(x => x.blob)" 
  show-progress 
  :initial-index="imageIndex"
  @close="show.preview_image = false">
  <template #toolbar="{ actions, prev, next, reset, activeIndex, setActiveItem }">
    <el-icon @click="prev"><Back /></el-icon>
    <el-icon @click="next"><Right /></el-icon>
    <el-icon @click="actions('zoomOut')"><ZoomOut /></el-icon>
    <el-icon @click="actions('zoomIn')"><ZoomIn /></el-icon>
    <el-icon @click="actions('clockwise')"><RefreshRight /></el-icon>
    <el-icon @click="actions('anticlockwise')"><RefreshLeft /></el-icon>
    <el-icon @click="download(imageList[activeIndex])"><Download /></el-icon>
  </template>
</el-image-viewer>
<el-dialog v-model="show.preview_plain" :show-close="false" :fullscreen="fullscreen.plain" width="90%" top="50px">
  <template #header="{ close }">
    <span class="el-dialog__title">{{ currentRow.Name }}</span>
    <el-button-group class="pull-right">
      <el-button circle @click="expandScreen('plain')">
        <font-awesome-icon icon="expand" v-if="!fullscreen.plain" />
        <font-awesome-icon icon="minimize" v-else />
      </el-button>
      <el-button circle :icon="Download" @click="download(currentRow)" />
      <el-button circle :icon="Close" @click="close" />
    </el-button-group>
    <el-switch v-model="wrapLine" inactive-text="强制折行" class="pull-right" style="margin-right:10px" />
  </template>
  <v-ace-editor
    id="plain"
    v-model:value="currentRow.filecontent"
    lang="text"
    theme="terminal"
    style="width:100%"
    class="minimize-screen-height"
    :options="{
      wrap: wrapLine,
      enableBasicAutocompletion: true,
      enableSnippets: true,
      enableLiveAutocompletion: true,
      tabSize: 2,
      showPrintMargin: false,
      fontSize: 14,
      minLines: 10,
  }" />
</el-dialog>
<el-dialog v-model="show.preview_pdf" :show-close="false" :fullscreen="fullscreen.pdf" width="90%" top="50px">
  <template #header="{ close }">
    <span class="el-dialog__title">{{ currentRow.Name }}</span>
    <el-button-group class="pull-right">
      <el-button circle @click="expandScreen('pdf')">
        <font-awesome-icon icon="expand" v-if="!fullscreen.pdf" />
        <font-awesome-icon icon="minimize" v-else />
      </el-button>
      <el-button circle :icon="Close" @click="close" />
    </el-button-group>
  </template>
  <embed id="pdf" :src="pdfSrc" width="100%" class="minimize-screen-height" />
</el-dialog>
</template>
<script setup>
import { ref, onMounted, onUnmounted, computed } from "vue";
import { ArrowRight, Download, Refresh, Search, Upload, Close, Delete } from '@element-plus/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '../assets/util/axios'
import moment from 'moment'
import { filesize } from "filesize"
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-text'
import terminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
ace.config.setModuleUrl('ace/theme/terminal', terminalUrl)
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const accessToken = computed(() => {
  return localStorage.access_token
})
const route = useRoute()
const router = useRouter()
const store = useStore()
const role = computed(() => {
  return store.state.userInfo?.role
})
const settings = computed(() => {
  return store.state.settings
})
const all = ref([])
const list = ref([])
const paths = ref([]) 
const rootpath = ref("")
const currentpath = ref("")
const searchKey = ref("")
const show = ref({
  upload: false,
  donwload_status: false,
  preview_image: false
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
const imageList = ref([])
const imageIndex = ref(0)
const currentRow = ref({})
const wrapLine = ref(false)
const fullscreen = ref({
  plain: false,
  pdf: false
})
const pdfSrc = ref(null)
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
onUnmounted(() => {
  for(let x of imageList.value) {
    if(x.blob?.startsWith('blob:')) {
      URL.revokeObjectURL(x)
    }
  }
  if(pdfSrc.value?.startsWith('blob:')) {
    URL.revokeObjectURL(pdfSrc.value)
  }
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
  // getPreviewList()
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
  if(route.query.source === 's3') {
    // downloadFile(`/filebrowser/s3/downloadobject?name=${route.query.name}&path=${row.Path}`, {}, row.Name)
    window.location = `/filebrowser/s3/downloadobject?name=${route.query.name}&path=${row.Path}`
  }
  else {
    // downloadFile(`/filebrowser/download?file=${row.Path}`, {}, row.Name)
    window.location = `/filebrowser/download?file=${row.Path}`
  }
}
const uploadSuccess = (res) => {
  ElMessage.success(res.message)
  getList(currentpath.value)
  uploadRef.value.clearFiles()
  show.value.upload = false
}
const uploadFailed = (err) => {
  ElMessage.error(`上传失败: ${err}`)
}
const confirmClick = async () => {
  uploadRef.value.submit()
}
// fetch + blob 方式会造成浏览器内存撑爆
// streamsaver.js必须联网从CDN下载mitm.html，不适用于内网环境
// const downloadFile = async (url, params, filename = 'file') => {
//   try {
//     show.value.donwload_status = true
//     const response = await fetch(url, {
//       method: 'POST', // 或 'GET'
//       body: JSON.stringify(params),
//       headers: {
//         'Content-Type': 'application/json',
//         'Authorization': `Bearer ${localStorage.access_token}`,
//       }
//     });
//     if (!response.ok) {
//       const error = await response.json();
//       throw new Error(error.message || '下载失败');
//     }
//     const contentLength = parseInt(response.headers.get('content-length'));
//     const reader = response.body.getReader();
//     const chunks = [];
//     let receivedLength = 0;
//     while(true) {
//       const {done, value} = await reader.read();
//       if (done) break;
//       chunks.push(value);
//       receivedLength += value.length;
//       downloadProgress.value = Math.ceil((receivedLength / contentLength) * 100)
//     }
//     // 合并chunks
//     const blob = new Blob(chunks);
//     const downloadUrl = window.URL.createObjectURL(blob)
//     const link = document.createElement('a');
//     link.href = downloadUrl;
//     link.download = filename;
//     link.click();
//     window.URL.revokeObjectURL(downloadUrl)
//   } catch (error) {
//     console.error(error)
//     ElMessage.error({message: "下载文件失败：" + error})
//   } finally {
//     show.value.donwload_status = false
//   }
// }
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
const preview = async (row) => {
  currentRow.value = row
  const sizeLimit = parseInt(settings.value.preview_file_size)
  if(row.Size > sizeLimit) {
    ElMessage.error(`大于 ${filesize(sizeLimit)} 的文件不提供预览，请下载文件后查看`)
    return
  }
  let url = `/filebrowser/preview?file=${row.Path}`
  if(route.query.source === 's3')
    url = `/filebrowser/s3/preview?name=${route.query.name}&path=${row.Path}`
  // 获取图片预览
  const ext = row.Name.split('.').pop()
  if(['jpg','jpeg','png','gif','bmp'].includes(ext)) {
    if(!imageList.value.find(n => n.Path === row.Path)) {
      let response = await axios.get(url, {
        responseType: 'blob'
      })
      row.blob = URL.createObjectURL(new Blob([response]))
      imageList.value.push(row)
    }
    imageIndex.value = imageList.value.findIndex(n => n.Path === row.Path)
    show.value.preview_image = true
  } else if (ext === 'svg') {
    if(!imageList.value.find(n => n.Path === row.Path)) {
      let response = await axios.get(url)
      row.blob = 'data:image/svg+xml,' + encodeURIComponent(response)
      imageList.value.push(row)
    }
    imageIndex.value = imageList.value.findIndex(n => n.Path === row.Path)
    show.value.preview_image = true
  } else if(ext === 'pdf') {
    let response = await axios.get(url, {
      responseType: 'arraybuffer'
    })
    const blob = new Blob([response], { type: 'application/pdf' })
    pdfSrc.value = URL.createObjectURL(blob)
    show.value.preview_pdf = true
  } else {
    // 文本文件预览
    
    currentRow.value.filecontent = await axios.get(url)
    show.value.preview_plain = true
  }
}
const expandScreen = (name) => {
  let dom = document.getElementById(name)
  if(!fullscreen.value[name]) {
    dom.classList.remove('minimize-screen-height')
    dom.classList.add('full-screen-height')
  } else {
    dom.classList.remove('full-screen-height')
    dom.classList.add('minimize-screen-height')
  }
  fullscreen.value[name] = !fullscreen.value[name]
}
const deleteOne = async (row) => {
  if(route.query.source === 's3') {
    if(row.IsDir && !row.Path.endsWith("/")) {
      row.Path += "/"
    } 
    await axios.delete(`/filebrowser/s3/delete?name=${route.query.name}&path=${row.Path}`)
  } else {
    await axios.delete(`/filebrowser/delete?file=${row.Path}`)
  }
  getList(currentpath.value)
}
const deleteBucket = async () => {
  ElMessageBox.confirm(
    '确认删除?',
    'Warning',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    await axios.delete(`/filebrowser/s3/deletebucket?name=${route.query.name}`, {data: {bucket_name: rootpath.value.replace(/\/$/, "")}})
    router.push({
      path: '/'
    })
  }).catch(() => {})
}
</script>