<template>
  <el-scrollbar>
    <div class="sidebar-logo">
      <div v-if="isCollapse" class="logo-mini"><b>A</b>-</div>
      <el-image v-else style="height:45px;margin-top:15px" src="/images/filebrowser.png" />
    </div>
    <div class="user-panel">
      <div class="pull-left image">
        <el-avatar :size="40" :src="avatar" />
      </div>
      <div class="pull-left info sidebar-userinfo">
        <p style="text-align:left">{{username}}</p>
        <a><font-awesome-icon icon="circle" style="color:green" /> Online</a>
      </div>
    </div>
    <div class="user-panel sidenav" style="text-align:left;padding:8px 10px;margin-top:5px;">
      导航
    </div>
    <el-menu
      active-text-color="#ffd04b"
      background-color="#141f29"
      class="el-menu-vertical-demo"
      :default-active="activeIndex"
      text-color="#fff"
      :unique-opened="true"
      :router="true">
      <el-menu-item index="/"><font-awesome-icon icon="home" />首页</el-menu-item>
      <el-menu-item index="/list"><font-awesome-icon icon="list" />文件列表</el-menu-item>
    </el-menu>
  </el-scrollbar>
</template>

<script setup>
import { onBeforeMount, ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
/* 变量定义 */
const store = useStore()
const router = useRouter()
const avatar = computed(() => {
  return "/images/avator.png"
})
const username = computed(() => {
  return store.state.userInfo.username
})
const isCollapse = ref({})
const activeIndex = ref("/")
/* watch */
watch(
  () => router.currentRoute.value,
  (newVal, oldVal) => {
    activeIndex.value = newVal.path
  }
)
/* 生命周期函数 */
onBeforeMount(async () => {
  isCollapse.value = false
})
</script>