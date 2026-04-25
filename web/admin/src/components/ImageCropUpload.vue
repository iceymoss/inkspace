<template>
  <div class="image-crop-upload">
    <div class="upload-area" @click="openFileDialog">
      <img
        v-if="modelValue"
        :src="modelValue"
        :style="previewStyle"
        class="preview-image object-cover"
      />
      <div v-else class="upload-placeholder" :style="previewStyle">
        <Plus class="upload-icon h-10 w-10 text-muted-foreground mb-1" />
        <div class="upload-text text-sm text-muted-foreground">{{ placeholder }}</div>
      </div>
    </div>
    <div class="upload-tip text-xs text-muted-foreground mt-1 leading-relaxed">{{ tip }}</div>

    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleFileChange"
    />

    <Dialog v-model:open="cropDialogVisible">
      <DialogContent class="sm:max-w-[800px]" @after-close="handleDialogClosed">
        <DialogHeader>
          <DialogTitle>裁剪图片</DialogTitle>
        </DialogHeader>
        <div class="crop-container w-full h-[500px] bg-black flex items-center justify-center overflow-hidden">
          <img ref="cropperImage" :src="tempImageUrl" style="max-width: 100%" />
        </div>
        <DialogFooter class="flex justify-between items-center sm:justify-between">
          <div class="flex gap-2">
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
          <div class="flex gap-2">
            <Button variant="outline" @click="cropDialogVisible = false">取消</Button>
            <Button @click="handleCropConfirm" :disabled="uploading">
              {{ uploading ? '上传中...' : '确定上传' }}
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, nextTick, computed } from 'vue'
import { toast } from 'vue-sonner'
import { Plus, RotateCcw, RotateCw, RefreshCw, ZoomIn, ZoomOut } from 'lucide-vue-next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.css'
import adminApi from '@/utils/adminApi'

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
    default: '上传正方形图片，系统会自动裁剪'
  },
  maxSize: {
    type: Number,
    default: 2
  },
  aspectRatio: {
    type: Number,
    default: 1
  },
  outputWidth: {
    type: Number,
    default: 800
  },
  outputHeight: {
    type: Number,
    default: 800
  }
})

const emit = defineEmits(['update:modelValue'])

const fileInput = ref(null)
const cropperImage = ref(null)
const cropDialogVisible = ref(false)
const tempImageUrl = ref('')
const uploading = ref(false)
let cropper = null

const previewStyle = computed(() => {
  if (props.aspectRatio && props.aspectRatio !== 1) {
    const width = parseInt(props.previewSize) || 120
    const height = width / props.aspectRatio
    return {
      width: `${width}px`,
      height: `${height}px`
    }
  }
  return {
    width: props.previewSize,
    height: props.previewSize
  }
})

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
    const cropperOptions = {
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
    }
    
    if (props.aspectRatio && props.aspectRatio !== 0) {
      cropperOptions.aspectRatio = props.aspectRatio
    }
    
    cropper = new Cropper(cropperImage.value, cropperOptions)
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
    const canvas = cropper.getCroppedCanvas({
      width: props.outputWidth,
      height: props.outputHeight,
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
      formData.append('file', blob, 'cropped.jpg')

      try {
        const response = await adminApi.post('/upload/image', formData, {
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
        toast.error('上传失败')
      } finally {
        uploading.value = false
      }
    }, 'image/jpeg', 0.9)
  } catch (error) {
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
  border: 2px dashed hsl(var(--border));
  border-radius: var(--radius);
  overflow: hidden;
  transition: all 0.2s;
  display: inline-block;
}

.upload-area:hover {
  border-color: hsl(var(--primary));
}

.preview-image {
  display: block;
  border-radius: calc(var(--radius) - 2px);
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: hsl(var(--muted));
  transition: all 0.2s;
}

.upload-placeholder:hover {
  background-color: hsl(var(--accent));
}

.upload-tip {
  font-size: 0.75rem;
  color: hsl(var(--muted-foreground));
  margin-top: 0.25rem;
  line-height: 1.625;
}
</style>
