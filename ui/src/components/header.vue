<template>
<el-menu
  mode="horizontal"
  :ellipsis="false"
  @select="handleSelect">
  <el-menu-item index="0">
    <span class="portal-title">File Browser</span>
  </el-menu-item>
  <div class="flex-grow" />
  <el-menu-item index="1">
    <div>
      <el-avatar :size="40" :src="avatar" style="vertical-align:middle" />
      <span style="margin-left:15px">{{username}}</span>
    </div>
  </el-menu-item>
  <el-sub-menu index="2">
    <template #title>
      <font-awesome-icon icon="gears" style="font-size:20px" />
    </template>
    <el-menu-item index="2-1">配置</el-menu-item>
    <el-menu-item index="2-2">修改密码</el-menu-item>
    <el-menu-item index="2-3">注销</el-menu-item>
  </el-sub-menu>
</el-menu>
<el-dialog v-model="show.modify" title="修改密码" width="600px">
  <el-form ref="refModify" :model="form" :rules="rules" label-width="140px">
    <el-form-item label="请输入旧密码" prop="oldPassword">
      <el-input v-model="form.oldPassword" size="large" type="password" show-password />
    </el-form-item>
    <el-form-item label="请输入新密码" prop="oldPassword">
      <el-input v-model="form.newPassword" size="large" type="password" show-password />
    </el-form-item>
    <el-form-item label="请再次输入新密码" prop="oldPassword">
      <el-input v-model="form.confirmPassword" size="large" type="password" show-password />
    </el-form-item>
    <el-form-item>
      <el-button @click="show.modify=false">取消</el-button>
      <el-button @click="submit(refModify)" type="primary">确定</el-button>
    </el-form-item>
  </el-form>
</el-dialog>
</template>

<script setup>
import { computed, ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const router = useRouter()
const store = useStore()
const avatar = computed(() => {
  return "/images/avator.png"
})
const username = computed(() => {
  return store.state.userInfo.username
})
const show = ref({
  modify: false
})
const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})
const rules = reactive({
  oldPassword: [{required: true, message: '请输入旧密码'}],
  newPassword: [{required: true, message: '请输入新密码'}],
  confirmPassword: [{required: true, message: '请再次输入新密码'}],
})
const refModify = ref(null)
/* methods */
const handleSelect = (index) => {
  switch(index) {
    case '1': router.push('/profile');break
    case '2-1': router.push('/platform/settings');break
    case '2-2': show.value.modify = true;break
    case '2-3': logout();break
  }
}
const logout = async () => {
  localStorage.clear()
  window.location.href = "/login/index.html"
}
const submit = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      if(params.newPassword !== params.confirmPassword) {
        ElMessage.error('两次输入密码不一致')
        return
      }
      params.username = localStorage.username
      delete params.confirmPassword
      await axios.post(`/filebrowser/auth/chpasswd`, params)
      window.location.href = "/login/index.html"
    } else {
      ElMessage.warning('必填项未填完')
    }
  })
}
</script>

<style>
.el-header {
  padding: 0 !important;
}
.sidebar-toggle {
  float: left;
  padding: 0 15px;
  color: inherit;
}
.portal-title {
  float:left;
  padding-left: 10px;
  font-size: 16px;
}
.flex-grow {
  flex-grow: 1;
}
.el-header .el-menu--horizontal .el-sub-menu__icon-arrow {
  display: none;
}
.el-header .el-menu--horizontal>.el-menu-item.is-active,
.el-header .el-menu--horizontal>.el-submenu.is-active>.el-submenu__title {
  border-bottom: 0px solid #ffffff !important;
}
</style>