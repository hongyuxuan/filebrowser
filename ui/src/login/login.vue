<template>
  <el-container>
    <el-aside class="login-aside" width="40%">
      <el-image style="height:300px;margin-top:calc(50vh - 120px)" src="/images/filebrowser_banner.png" />
    </el-aside>
    <el-container>
      <el-main style="padding:0">
        <el-row justify="center" style="flex-direction:column;align-items:center;height:calc(100vh)">
          <div style="color:#141f29;font-size:25px;">WELCOME TO FILEBROWSER</div>
          <el-form
            ref="login"
            label-position="top"
            label-width="auto"
            :model="form"
            :rules="rules"
            style="width:400px;margin-top:30px">
            <el-form-item label="Username" prop="username">
              <el-input v-model="form.username" size="large" />
            </el-form-item>
            <el-form-item label="Password" prop="password">
              <el-input v-model="form.password" type="password" show-password size="large" @keyup.enter="submit(login)" />
            </el-form-item>
            <el-form-item>
              <el-button @click="submit(login)" style="width:100%" type="primary" size="large">Login</el-button>
            </el-form-item>
          </el-form>
        </el-row>
      </el-main>
    </el-container>
  </el-container>
</template>
<script setup>
import { onBeforeMount, ref, reactive } from 'vue'
import { axios } from '/src/assets/util/axios'
import { ElMessage } from 'element-plus'
const form = ref({})
const rules = reactive({
  username: [{required: true, message: '请输入用户名'}],
  password: [{required: true, message: '请输入密码'}],
})
const login = ref(null)
/* methods */
const submit = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let response = await axios.post(`/filebrowser/auth/login`, form.value)
      localStorage.access_token = response.access_token
      window.location.href = "/"
    }
  })
}
</script>
<style>
.login-aside {
  background-color: #141f29;
  height: calc(100vh);
  text-align:center;
}
.login-panel {
  text-align:center;
}
</style>