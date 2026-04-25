<template>
  <div class="image-crop-upload">
    <div class="upload-area" @click="openFileDialog">
      <img
        v-if="modelValue && !imageError"
        :src="modelValue"
        :style="{ maxWidth: previewSize, maxHeight: '120px' }"
        class="preview-image"
        @error="imageError = true"
      />
      <div v-if="modelValue && imageError" class="image-error" :style="{ maxWidth: previewSize, maxHeight: '120px' }">
        <ImageIcon class="h-10 w-10" />
      </div>
      <div v-if="!modelValue" class="upload-placeholder" :style="{ width: previewSize, minHeight: '80px' }">
        <Plus class="h-7 w-7 text-muted-foreground mb-2" />
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

    <Dialog :open="cropDialogVisible" @update:open="(val) => { cropDialogVisible = val; if (!val) handleDialogClosed() }">
      <DialogContent class="max-w-[800px]">
        <DialogHeader>
          <DialogTitle>裁剪图片</DialogTitle>
        </DialogHeader>
        <div class="crop-container">
          <img ref="cropperImage" :src="tempImageUrl" style="max-width: 100%" />
        </div>
        <DialogFooter>
          <div class="crop-footer">
            <div class="flex items-center gap-2">
              <Button variant="outline" size="sm" @click="rotateCrop(-90)">
                <RotateCcw class="h-4 w-4 mr-1" /> 左旋
              </Button>
              <Button variant="outline" size="sm" @click="rotateCrop(90)">
                <RotateCw class="h-4 w-4 mr-1" /> 右旋
              </Button>
              <Button variant="outline" size="sm" @click="resetCrop">
                <RefreshCw class="h-4 w-4 mr-1" /> 重置
              </Button>
              <Button variant="outline" size="sm" @click="zoomCrop(0.1)">
                <ZoomIn class="h-4 w-4 mr-1" /> 放大
              </Button>
              <Button variant="outline" size="sm" @click="zoomCrop(-0.1)">
                <ZoomOut class="h-4 w-4 mr-1" /> 缩小
              </Button>
            </div>
            <div class="flex items-center gap-2">
              <Button variant="outline" @click="cropDialogVisible = false">取消</Button>
              <Button :disabled="uploading" @click="handleCropConfirm">
                确定上传
              </Button>
            </div>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { toast } from 'vue-sonner'
import { Plus, Image as ImageIcon, RotateCcw, RotateCw, RefreshCw, ZoomIn, ZoomOut } from 'lucide-vue-next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
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
    default: 5
  },
  aspectRatio: {
    type: Number,
    default: NaN
  }
})

const emit = defineEmits(['update:modelValue'])

const fileInput = ref(null)
const cropperImage = ref(null)
const cropDialogVisible = ref(false)
const tempImageUrl = ref('')
const uploading = ref(false)
const imageError = ref(false)
let cropper = null

const openFileDialog = () => {
  fileInput.value?.click()
}

const handleFileChange = (e) => {
  const file = e.target.files[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    toast.error('请选择图片文件')
    return
  }

  if (file.size / 1024 / 1024 > props.maxSize) {
    toast.error(`图片大小不能超过 ${props.maxSize}MB`)
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    tempImageUrl.value = e.target.result
    cropDialogVisible.value = true
    nextTick(() => {
      initCropper()
    })
  }
  reader.readAsDataURL(file)

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
    const cropData = cropper.getData()
    const maxWidth = 1920
    const scale = cropData.width > maxWidth ? maxWidth / cropData.width : 1

    const canvas = cropper.getCroppedCanvas({
      maxWidth: maxWidth,
      maxHeight: Math.round(cropData.height * scale),
      imageSmoothingEnabled: true,
      imageSmoothingQuality: 'high'
    })

    canvas.toBlob(async (blob) => {
      if (!blob) {
        toast.error('图片处理失败')
        uploading.value = false
        return
      }

      const formData = new FormData()
      formData.append('file', blob, 'cover.jpg')

      try {
        const response = await api.post('/upload/image', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })

        if (response.code === 0 && response.data) {
          emit('update:modelValue', response.data.url)
          toast.success('上传成功')
          cropDialogVisible.value = false
        } else {
          toast.error(response.message || '上传失败')
        }
      } catch (error) {
        console.error('Upload error:', error)
        toast.error('上传失败')
      } finally {
        uploading.value = false
      }
    }, 'image/jpeg', 0.9)
  } catch (error) {
    console.error('Crop error:', error)
    toast.error('图片处理失败')
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
  border: 2px dashed var(--theme-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: all var(--transition-slow);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 80px;
}

.upload-area:hover {
  border-color: var(--theme-primary);
}

.preview-image {
  display: block;
  border-radius: var(--radius-sm);
  object-fit: contain;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: var(--theme-bg-secondary);
  transition: all var(--transition-slow);
}

.upload-placeholder:hover {
  background-color: var(--theme-bg-hover);
}

.upload-text {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.image-error {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background-color: var(--theme-bg-hover);
  color: var(--theme-text-tertiary);
}

.upload-tip {
  font-size: var(--font-size-xs);
  color: var(--theme-text-tertiary);
  margin-top: var(--spacing-sm);
  line-height: var(--line-height-base);
}

.crop-container {
  width: 100%;
  height: 500px;
  background-color: var(--theme-bg-primary);
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
  width: 100%;
}
</style>
