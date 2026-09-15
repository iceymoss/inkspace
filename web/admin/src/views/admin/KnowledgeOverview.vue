<template><div v-loading="loading"><h2>知识库运营概览</h2><div class="cards"><el-card v-for="[key,label] in cards" :key="key"><el-statistic :title="label" :value="data[key]||0" /></el-card></div><el-card><template #header>近 7 日新增文档</template><div v-for="item in data.trend" :key="item.day" class="trend"><span>{{item.day}}</span><el-progress :percentage="maxCount?item.count/maxCount*100:0" :format="()=>String(item.count)" /></div></el-card><el-card><template #header>公开 Wiki 浏览量 Top 10</template><el-table :data="data.top_docs||[]"><el-table-column prop="title" label="文档"/><el-table-column prop="workspace_id" label="工作区"/><el-table-column prop="view_count" label="浏览量"/></el-table></el-card><p>附件按当前文档引用统计；发布态附件不单独计入，多次引用重复计数。</p></div></template>
<script setup>
import {ref,computed,onMounted} from 'vue'
import api from '@/utils/adminApi'
const data=ref({}),loading=ref(false)
const cards=[['workspaces','工作区总数'],['public_workspaces','公开工作区'],['private_workspaces','私有工作区'],['docs','文档总数'],['draft_docs','草稿'],['published_docs','已发布'],['blocked_docs','已下架'],['shares','分享链接数'],['storage_bytes','附件存储（字节）']]
const maxCount=computed(()=>Math.max(0,...(data.value.trend||[]).map(x=>x.count)))
onMounted(async()=>{loading.value=true;try{data.value=(await api.get('/admin/knowledge/overview')).data}finally{loading.value=false}})
</script>
<style scoped>.cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:16px;margin-bottom:20px}.el-card{margin-bottom:16px}.trend{display:flex;gap:20px;margin:12px}.el-progress{flex:1}p{color:#909399;font-size:13px}</style>
