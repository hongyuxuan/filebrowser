<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>配置</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p v-html="info" />
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" @tab-change="changeTab" class="qingcloud-tab">
      <el-tab-pane name="1" label="用户管理" v-if="role==='admin'">
          <keep-alive>
          <user />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="2" label="S3接入管理">
          <keep-alive>
          <s3 />
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import s3 from './s3.vue'
import user from './user.vue'
/* 变量定义 */
const store = useStore()
const role = computed(() => {
  return store.state.userInfo.role
})
const activeName = ref("1")
const info = ref("配置页面为您提供 Filebrowser 平台设置，包括用户设置、S3接入配置等。部分功能仅管理员有权限设置。如有问题，请联系您的管理员。")
/* methods */
const changeTab = (tabName) => {
  
}
</script>