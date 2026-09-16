<template>
  <el-card v-loading="loading">
    <template #header><strong>{{ title }}</strong></template>
    <el-form inline @submit.prevent="search">
      <el-form-item v-for="field in fields" :key="field.key" :label="field.label">
        <el-select v-if="field.options" v-model="filters[field.key]" clearable style="width:140px"><el-option v-for="o in field.options" :key="o[0]" :value="o[0]" :label="o[1]" /></el-select>
        <el-input v-else v-model="filters[field.key]" clearable style="width:150px" />
      </el-form-item>
      <el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button>
    </el-form>
    <el-table :data="rows" :row-class-name="({row})=>row.over_quota?'over-quota':''">
      <el-table-column v-for="column in columns" :key="column[0]" :prop="column[0]" :label="column[1]" min-width="110" show-overflow-tooltip>
        <template #default="{row}"><template v-if="column[0]==='owner_id'">{{ row.owner_id }} · {{ row.owner_name || '—' }}</template><template v-else>{{ display(row,column[0]) }}</template></template>
      </el-table-column>
      <el-table-column v-if="editable" label="操作" width="230" fixed="right"><template #default="{row}">
        <el-button v-if="kind!=='shares'" link type="primary" @click="detail(row)">详情</el-button>
        <el-button v-if="kind!=='shares'" link type="warning" @click="audit(row)">{{ row.audit_status ? '恢复':'下架' }}</el-button>
        <el-button link type="danger" @click="remove(row)">{{kind==='shares'?'强制删除':'删除'}}</el-button>
      </template></el-table-column>
    </el-table>
    <el-empty v-if="!loading&&!rows.length" description="暂无数据" />
    <el-pagination v-model:current-page="page" v-model:page-size="size" :total="total" :page-sizes="[10,20,50,100]" layout="total,sizes,prev,pager,next" @current-change="load" @size-change="search" />
    <p v-if="kind==='usage'" class="footnote">附件按文档当前 AttachmentID 统计；发布态附件不单独计入，同一附件被多个文档引用时重复计数。红色行表示超出全局配额。</p>
    <el-dialog v-model="detailVisible" title="元数据与公开内容" width="75%">
      <el-descriptions :column="2" border><el-descriptions-item v-for="(value,key) in metadata" :key="key" :label="key">{{value}}</el-descriptions-item></el-descriptions>
      <el-alert v-if="kind==='docs'&&selected.content===undefined" title="私有内容，仅展示元数据（已下架内容也不开放预览）" type="info" :closable="false" />
      <pre v-if="selected.content!==undefined" class="preview">{{selected.content}}</pre>
    </el-dialog>
  </el-card>
</template>
<script setup>
import { computed, ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/adminApi'
const props=defineProps({kind:String,title:String})
const rows=ref([]),page=ref(1),size=ref(10),total=ref(0),loading=ref(false),filters=reactive({}),selected=ref({}),detailVisible=ref(false)
const editable=computed(()=>['workspaces','docs','shares'].includes(props.kind))
const auditOptions=[[0,'正常'],[1,'已下架']]
const fields=computed(()=>({workspaces:[{key:'keyword',label:'关键词'},{key:'owner_id',label:'所有者ID'},{key:'is_public',label:'可见性',options:[['true','公开'],['false','私有']]},{key:'audit_status',label:'审核',options:auditOptions}],docs:[{key:'keyword',label:'关键词'},{key:'owner_id',label:'所有者ID'},{key:'workspace_id',label:'工作区ID'},{key:'status',label:'状态',options:[[0,'草稿'],[1,'已发布']]},{key:'audit_status',label:'审核',options:auditOptions},{key:'kind',label:'类型',options:['markdown','text','code','file'].map(k=>[k,k])}],shares:[{key:'owner_id',label:'所有者ID'},{key:'doc_id',label:'文档ID'},{key:'state',label:'状态',options:[['active','生效中'],['disabled','已禁用'],['expired','已过期']]}],usage:[{key:'sort',label:'排序',options:[['storage_bytes','存储占用'],['doc_count','文档数'],['workspace_count','工作区数']]}],'audit-logs':[{key:'admin_id',label:'管理员ID'},{key:'target_type',label:'对象',options:[['workspace','工作区'],['doc','文档'],['share_link','分享链接']]},{key:'action',label:'操作',options:[['block','下架'],['unblock','恢复'],['delete','删除']]},{key:'start',label:'开始日期 YYYY-MM-DD'},{key:'end',label:'结束日期 YYYY-MM-DD'}]})[props.kind]||[])
const columns=computed(()=>({workspaces:[['id','ID'],['name','名称'],['owner_id','所有者'],['is_public','可见性'],['doc_count','文档数'],['audit_status','审核'],['audit_reason','处置理由'],['updated_at','更新时间']],docs:[['id','ID'],['title','标题'],['workspace_id','工作区'],['owner_id','所有者'],['kind','类型'],['status','状态'],['audit_status','审核'],['audit_reason','处置理由']],shares:[['id','ID'],['token','分享标识'],['doc_id','文档'],['owner_id','所有者'],['state','状态'],['expires_at','到期时间'],['view_count','浏览量']],usage:[['id','用户ID'],['username','用户'],['workspace_count','工作区数'],['doc_count','文档数'],['storage_bytes','存储 MB'],['max_workspace_docs','单工作区最多文档']], 'audit-logs':[['id','ID'],['admin_username','管理员'],['action','操作'],['target_type','对象类型'],['target_id','对象ID'],['target_title','标题'],['reason','理由'],['ip','IP'],['created_at','时间']]})[props.kind]||[])
const metadata=computed(()=>Object.fromEntries(Object.entries(selected.value).filter(([k])=>!['content','content_html'].includes(k))))
function display(row,key){if(key==='audit_status')return row[key]?'已下架':'正常';if(key==='status')return row[key]?'已发布':'草稿';if(key==='is_public')return row[key]?'公开':'私有';if(key==='storage_bytes')return (row[key]/1048576).toFixed(2);if(key==='state')return !row.enabled?'已禁用':row.expires_at&&new Date(row.expires_at)<=new Date()?'已过期':'生效中';return row[key]??'—'}
async function load(){loading.value=true;try{const params={page:page.value,page_size:size.value,...Object.fromEntries(Object.entries(filters).filter(([,v])=>v!==''&&v!=null))};const r=await api.get(`/admin/knowledge/${props.kind}`,{params});rows.value=r.data.list||[];total.value=r.data.total}finally{loading.value=false}}
function search(){page.value=1;load()}function reset(){Object.keys(filters).forEach(k=>delete filters[k]);search()}
async function reason(title){const r=await ElMessageBox.prompt('必填：2–200 字处置理由',title,{inputType:'textarea',inputValidator:v=>v&&v.trim().length>=2&&v.trim().length<=200||'请输入 2–200 字理由'});return r.value.trim()}
async function audit(row){try{const text=await reason(row.audit_status?'恢复内容':'下架内容');await api.put(`/admin/knowledge/${props.kind}/${row.id}/audit`,{audit_status:row.audit_status?0:1,reason:text});ElMessage.success('处置成功');await load()}catch(e){if(e!=='cancel'&&e!=='close')console.error(e)}}
async function remove(row){try{if(props.kind!=='shares'){const title=row.name||row.title;await ElMessageBox.prompt(`请输入标题「${title}」以确认删除`,'删除确认',{inputValidator:v=>v===title||'标题不匹配'})}else{await ElMessageBox.confirm('删除后分享链接将立即失效，是否继续？','强制删除')};const text=await reason('确认删除并填写理由');await api.delete(`/admin/knowledge/${props.kind}/${row.id}`,{data:{reason:text}});ElMessage.success('删除成功');await load()}catch(e){if(e!=='cancel'&&e!=='close')console.error(e)}}
async function detail(row){selected.value=(await api.get(`/admin/knowledge/${props.kind}/${row.id}`)).data;detailVisible.value=true}
onMounted(load)
</script>
<style scoped>.el-pagination{margin-top:20px}.footnote{font-size:13px;color:#909399}.preview{white-space:pre-wrap;overflow-wrap:anywhere;max-height:50vh;overflow:auto}.el-alert{margin-top:20px}:deep(.over-quota){--el-table-tr-bg-color:var(--el-color-danger-light-9)}</style>
