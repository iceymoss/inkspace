<template>
  <el-card v-loading="loading">
    <template #header>知识库全局配额</template>
    <el-alert title="0 表示不限制。配额仅在新增时校验，调整上限不会清理存量内容。" type="info" :closable="false" />
    <el-form label-width="220px" style="margin-top:24px">
      <el-form-item v-for="[key,label] in fields" :key="key" :label="label"><el-input-number v-model="form[key]" :min="0" :max="2147483647" :precision="0" /></el-form-item>
      <el-form-item><el-button type="primary" :loading="saving" @click="save">保存配额</el-button></el-form-item>
    </el-form>
  </el-card>
</template>
<script setup>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/utils/adminApi'
const form=reactive({}),loading=ref(false),saving=ref(false)
const fields=[['knowledge.max_workspaces_per_user','每用户最大工作区数'],['knowledge.max_docs_per_workspace','每工作区最大文档数'],['knowledge.max_storage_per_user_mb','每用户最大存储（MB）']]
onMounted(async()=>{loading.value=true;try{Object.assign(form,(await api.get('/admin/knowledge/quota')).data)}finally{loading.value=false}})
async function save(){if(fields.some(([key])=>!Number.isInteger(form[key])||form[key]<0)){ElMessage.error('请输入非负整数');return};saving.value=true;try{await api.put('/admin/knowledge/quota',form);ElMessage.success('配额已保存')}finally{saving.value=false}}
</script>
