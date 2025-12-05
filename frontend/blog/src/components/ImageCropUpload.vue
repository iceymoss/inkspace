<template>
  <div class="image-crop-upload">
    <div class="upload-area" @click="openFileDialog">
      <el-image 
        v-if="modelValue" 
        :src="modelValue" 
        :style="{ maxWidth: previewSize, maxHeight: '120px' }"
        fit="contain"
        class="preview-image"
      >
        <template #error>
          <div class="image-error">
            <el-icon><Picture /></el-icon>
          </div>
        </template>
      </el-image>
      <div v-else class="upload-placeholder" :style="{ width: previewSize, minHeight: '80px' }">
        <el-icon class="upload-icon"><Plus /></el-icon>
        <div class="upload-text">{{ placeholder }}</div>
      </div>
    </div>
    <div class="upload-tip">{{ tip }}</div>

    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleFileChange"
    />

    <!-- 裁剪对话框 -->
    <el-dialog
      v-model="cropDialogVisible"
      title="裁剪图片"
      width="800px"
      :close-on-click-modal="false"
      @closed="handleDialogClosed"
    >
      <div class="crop-container">
        <img ref="cropperImage" :src="tempImageUrl" style="max-width: 100%" />
      </div>
      <template #footer>
        <div class="crop-footer">
          <el-space>
            <el-button @click="rotateCrop(-90)">
              <el-icon><RefreshLeft /></el-icon> 左旋
            </el-button>
            <el-button @click="rotateCrop(90)">
              <el-icon><RefreshRight /></el-icon> 右旋
            </el-button>
            <el-button @click="resetCrop">
              <el-icon><Refresh /></el-icon> 重置
            </el-button>
            <el-button @click="zoomCrop(0.1)">
              <el-icon><ZoomIn /></el-icon> 放大
            </el-button>
            <el-button @click="zoomCrop(-0.1)">
              <el-icon><ZoomOut /></el-icon> 缩小
            </el-button>
          </el-space>
          <el-space>
            <el-button @click="cropDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="handleCropConfirm" :loading="uploading">
              确定上传
            </el-button>
          </el-space>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Picture, RefreshLeft, RefreshRight, Refresh, ZoomIn, ZoomOut } from '@element-plus/icons-vue'
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.css'
import api from '@/utils/api'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  previewSize: {
    type: String,
    default: '120px'
  },
  placeholder: {
    type: String,
    default: '点击上传图片'
  },
  tip: {
    type: String,
    default: '可自由裁切任意比例'
  },
  maxSize: {
    type: Number,
    default: 5 // MB
  },
  aspectRatio: {
    type: Number,
    default: NaN // NaN 表示不限制比例，自由裁切
  }
})

const emit = defineEmits(['update:modelValue'])

const fileInput = ref(null)
const cropperImage = ref(null)
const cropDialogVisible = ref(false)
const tempImageUrl = ref('')
const uploading = ref(false)
let cropper = null

const openFileDialog = () => {
  fileInput.value?.click()
}

const handleFileChange = (e) => {
  const file = e.target.files[0]
  if (!file) return

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }

  // 验证文件大小
  if (file.size / 1024 / 1024 > props.maxSize) {
    ElMessage.error(`图片大小不能超过 ${props.maxSize}MB`)
    return
  }

  // 读取图片并打开裁剪对话框
  const reader = new FileReader()
  reader.onload = (e) => {
    tempImageUrl.value = e.target.result
    cropDialogVisible.value = true
    // 初始化裁剪器
    nextTick(() => {
      initCropper()
    })
  }
  reader.readAsDataURL(file)

  // 清空input，允许重复选择同一文件
  e.target.value = ''
}

const initCropper = () => {
  if (cropper) {
    cropper.destroy()
  }
  
  if (cropperImage.value) {
    cropper = new Cropper(cropperImage.value, {
      aspectRatio: props.aspectRatio,
      viewMode: 1,
      dragMode: 'move',
      autoCropArea: 0.8,
      restore: false,
      guides: true,
      center: true,
      highlight: false,
      cropBoxMovable: true,
      cropBoxResizable: true,
      toggleDragModeOnDblclick: false,
    })
  }
}

const rotateCrop = (degree) => {
  cropper?.rotate(degree)
}

const zoomCrop = (ratio) => {
  cropper?.zoom(ratio)
}

const resetCrop = () => {
  cropper?.reset()
}

const handleDialogClosed = () => {
  if (cropper) {
    cropper.destroy()
    cropper = null
  }
  tempImageUrl.value = ''
}

const handleCropConfirm = async () => {
  if (!cropper) return

  uploading.value = true
  try {
    // 获取裁剪后的canvas，保持原始裁切比例
    // 限制最大宽度为 1920px，高度按比例缩放
    const cropData = cropper.getData()
    const maxWidth = 1920
    const scale = cropData.width > maxWidth ? maxWidth / cropData.width : 1
    
    const canvas = cropper.getCroppedCanvas({
      maxWidth: maxWidth,
      maxHeight: Math.round(cropData.height * scale),
      imageSmoothingEnabled: true,
      imageSmoothingQuality: 'high'
    })

    // 转换为blob
    canvas.toBlob(async (blob) => {
      if (!blob) {
        ElMessage.error('图片处理失败')
        uploading.value = false
        return
      }

      // 上传裁剪后的图片
      const formData = new FormData()
      formData.append('file', blob, 'cover.jpg')

      try {
        const response = await api.post('/upload/image', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })

        if (response.code === 0 && response.data) {
          // 直接使用相对路径
          emit('update:modelValue', response.data.url)
          ElMessage.success('上传成功')
          cropDialogVisible.value = false
        } else {
          ElMessage.error(response.message || '上传失败')
        }
      } catch (error) {
        console.error('Upload error:', error)
        ElMessage.error('上传失败')
      } finally {
        uploading.value = false
      }
    }, 'image/jpeg', 0.9)
  } catch (error) {
    console.error('Crop error:', error)
    ElMessage.error('图片处理失败')
    uploading.value = false
  }
}
</script>

<style scoped>
.image-crop-upload {
  display: inline-block;
}

.upload-area {
  cursor: pointer;
  border: 2px dashed #dcdfe6;
  border-radius: 6px;
  overflow: hidden;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 80px;
}

.upload-area:hover {
  border-color: #409eff;
}

.preview-image {
  display: block;
  border-radius: 4px;
  object-fit: contain;
}

.preview-image :deep(img) {
  object-fit: contain;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
  transition: all 0.3s;
}

.upload-placeholder:hover {
  background-color: #f0f2f5;
}

.upload-icon {
  font-size: 28px;
  color: #c0c4cc;
  margin-bottom: 8px;
}

.upload-text {
  font-size: 13px;
  color: #606266;
}

.image-error {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background-color: #f5f7fa;
  color: #c0c4cc;
  font-size: 40px;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
  line-height: 1.5;
}

.crop-container {
  width: 100%;
  height: 500px;
  background-color: #000;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.crop-container img {
  max-width: 100%;
  max-height: 100%;
}

.crop-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 全屏模式下的padding */
:deep(.el-dialog__body) {
  padding: 20px;
}

:deep(.el-dialog.is-fullscreen .el-dialog__body) {
  padding: 30px;
}
</style>

