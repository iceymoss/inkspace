<template>
  <div class="work-detail">
    <div class="container">
      <Card v-if="work" class="detail-card">
        <CardContent class="p-6">
          <div class="card-header">
            <div v-if="work.status !== 1" class="work-status-banner">
              <Alert :class="getStatusAlertClass(work.status)">
                <AlertDescription>
                  <div>
                    <div class="font-medium">{{ getStatusText(work.status) }}</div>
                    <div v-if="work.audit_message" class="audit-message-text">
                      {{ work.audit_message }}
                    </div>
                  </div>
                </AlertDescription>
              </Alert>
            </div>
            <div v-if="isWorkOwner" class="card-header-actions">
              <Button
                variant="default"
                size="default"
                @click="handleEdit"
              >
                <Pencil class="mr-1 h-4 w-4" />
                编辑作品
              </Button>
            </div>
          </div>

          <template v-if="work.type === 'photography'">
            <div class="photography-layout">
              <div class="photo-section">
                <div class="main-photo-container" :style="{ height: carouselHeight + 'px' }">
                  <div class="carousel-wrapper" :style="{ transform: `translateX(-${currentPhotoIndex * 100}%)` }">
                    <div
                      v-for="(photo, index) in photos"
                      :key="index"
                      class="carousel-slide"
                    >
                      <img
                        :src="photo.url"
                        :alt="photo.description || work.title"
                        class="main-photo"
                        @click="openPhotoPreview(index)"
                        @load="handleImageLoad"
                      />
                    </div>
                  </div>

                  <div v-if="photos.length > 1" class="carousel-nav">
                    <Button variant="ghost" size="icon" class="carousel-nav-btn" @click="prevPhoto">
                      <ChevronLeft class="h-5 w-5" />
                    </Button>
                    <Button variant="ghost" size="icon" class="carousel-nav-btn" @click="nextPhoto">
                      <ChevronRight class="h-5 w-5" />
                    </Button>
                  </div>

                  <div class="photo-counter">
                    {{ currentPhotoIndex + 1 }} / {{ photos.length }}
                  </div>
                </div>

                <div class="thumbnail-nav" v-if="photos.length > 1">
                  <div
                    v-for="(photo, index) in photos"
                    :key="index"
                    class="thumbnail-item"
                    :class="{ active: index === currentPhotoIndex }"
                    @click="setActivePhoto(index)"
                  >
                    <img :src="photo.url" :alt="photo.description || ''" class="thumbnail-img" />
                  </div>
                </div>

                <div class="photo-description" v-if="currentPhoto.description">
                  <p>{{ currentPhoto.description }}</p>
                </div>

                <div class="album-info">
                  <h2>{{ work.title }}</h2>
                  <div id="work-description-preview-photography" v-if="work.description"></div>
                  <p v-else class="album-description-empty">暂无描述</p>

                  <div class="album-meta" v-if="work.metadata">
                    <div class="meta-item" v-if="work.metadata.location">
                      <MapPin class="h-4 w-4" />
                      <span>{{ work.metadata.location }}</span>
                    </div>
                    <div class="meta-item" v-if="work.metadata.shooting_date">
                      <Calendar class="h-4 w-4" />
                      <span>{{ work.metadata.shooting_date }}</span>
                    </div>
                    <div class="meta-item">
                      <Image class="h-4 w-4" />
                      <span>共 {{ photos.length }} 张照片</span>
                    </div>
                  </div>
                </div>
              </div>

              <div class="info-sidebar">
                <div class="author-card" v-if="work.author">
                  <div class="author-header" @click="goToUserProfile(work.author.id)">
                    <Avatar class="h-[60px] w-[60px]">
                      <AvatarImage :src="work.author.avatar" />
                      <AvatarFallback>{{ (work.author.nickname || work.author.username || '?')[0] }}</AvatarFallback>
                    </Avatar>
                    <div class="author-info">
                      <h3>{{ work.author.nickname || work.author.username }}</h3>
                      <p>{{ work.author.bio || '这个人很懒，什么都没留下' }}</p>
                    </div>
                  </div>
                  <Button
                    v-if="isWorkOwner"
                    variant="default"
                    class="full-width-btn"
                    @click="goToMyProfile"
                  >
                    我的主页
                  </Button>
                  <Button
                    v-else-if="userStore.isLoggedIn"
                    :variant="isFollowing ? 'outline' : 'default'"
                    class="full-width-btn"
                    :disabled="followLoading"
                    @click="handleFollow"
                  >
                    <Loader2 v-if="followLoading" class="mr-1 h-4 w-4 animate-spin" />
                    <template v-else>
                      <Plus v-if="!isFollowing" class="mr-1 h-4 w-4" />
                      <Check v-else class="mr-1 h-4 w-4" />
                    </template>
                    {{ isFollowing ? '已关注' : '关注' }}
                  </Button>
                </div>

                <div class="photo-params-card" v-if="currentPhoto.metadata">
                  <h3>📷 拍摄参数</h3>
                  <div class="params-list">
                    <div class="param-item" v-if="currentPhoto.metadata.camera">
                      <span class="param-label">相机：</span>
                      <span class="param-value">{{ currentPhoto.metadata.camera }}</span>
                    </div>
                    <div class="param-item" v-if="currentPhoto.metadata.lens">
                      <span class="param-label">镜头：</span>
                      <span class="param-value">{{ currentPhoto.metadata.lens }}</span>
                    </div>
                    <div class="param-item" v-if="currentPhoto.metadata.focal_length">
                      <span class="param-label">焦段：</span>
                      <span class="param-value">{{ currentPhoto.metadata.focal_length }}</span>
                    </div>
                    <div class="param-item" v-if="currentPhoto.metadata.aperture">
                      <span class="param-label">光圈：</span>
                      <span class="param-value">{{ currentPhoto.metadata.aperture }}</span>
                    </div>
                    <div class="param-item" v-if="currentPhoto.metadata.shutter_speed">
                      <span class="param-label">快门：</span>
                      <span class="param-value">{{ currentPhoto.metadata.shutter_speed }}</span>
                    </div>
                    <div class="param-item" v-if="currentPhoto.metadata.iso">
                      <span class="param-label">ISO：</span>
                      <span class="param-value">{{ currentPhoto.metadata.iso }}</span>
                    </div>
                  </div>
                </div>

                <div class="stats-card">
                  <h3>互动数据</h3>
                  <div class="stats-list">
                    <div class="stat-item">
                      <Eye class="h-4 w-4" />
                      <span>{{ work.view_count }} 浏览</span>
                    </div>
                    <div class="stat-item">
                      <MessageCircle class="h-4 w-4" />
                      <span>{{ work.comment_count }} 评论</span>
                    </div>
                  </div>
                  <div class="action-buttons" v-if="work && work.is_published">
                    <Button
                      :variant="isLiked ? 'default' : 'outline'"
                      @click="handleLike"
                      :disabled="liking"
                    >
                      <Loader2 v-if="liking" class="mr-1 h-4 w-4 animate-spin" />
                      <Star v-else class="mr-1 h-4 w-4" />
                      {{ work.like_count || 0 }} {{ isLiked ? '已点赞' : '点赞' }}
                    </Button>
                    <Button
                      :variant="isFavorited ? 'accent' : 'outline'"
                      @click="handleFavorite"
                      :disabled="favoriting"
                    >
                      <Loader2 v-if="favoriting" class="mr-1 h-4 w-4 animate-spin" />
                      <Star v-else class="mr-1 h-4 w-4" />
                      {{ work.favorite_count || 0 }} {{ isFavorited ? '已收藏' : '收藏' }}
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else>
            <div class="project-layout">
              <div class="work-header">
                <div class="work-title-section">
                  <h1>{{ work.title }}</h1>
                  <Badge :variant="work.type === 'photography' ? 'secondary' : 'default'">
                    {{ work.type === 'photography' ? '📷 摄影作品' : '💻 开源项目' }}
                  </Badge>
                </div>

                <div class="work-author" v-if="work.author">
                  <Avatar
                    class="clickable-avatar h-10 w-10"
                    @click="goToUserProfile(work.author.id)"
                  >
                    <AvatarImage :src="work.author.avatar" />
                    <AvatarFallback>{{ (work.author.nickname || work.author.username || '?')[0] }}</AvatarFallback>
                  </Avatar>
                  <div class="author-info">
                    <span class="author-name" @click="goToUserProfile(work.author.id)">
                      {{ work.author.nickname || work.author.username }}
                    </span>
                    <span class="publish-time">{{ formatDate(work.created_at) }}</span>
                  </div>
                </div>
              </div>

              <div class="work-content">
                <div id="work-description-preview-project" v-if="work.description"></div>
                <p v-else class="work-description-empty">暂无描述</p>

                <div class="tech-stack" v-if="work.tech_stack">
                  <h3>技术栈</h3>
                  <p>{{ work.tech_stack }}</p>
                </div>

                <div class="work-links" v-if="work.github_url || work.demo_url || work.link">
                  <Button
                    v-if="work.github_url"
                    size="lg"
                    class="work-link-btn github-btn"
                    @click="openLink(work.github_url)"
                  >
                    <ExternalLink class="mr-2 h-4 w-4" />
                    GitHub
                  </Button>
                  <Button
                    v-if="work.link"
                    size="lg"
                    class="work-link-btn project-btn"
                    @click="openLink(work.link)"
                  >
                    <ExternalLink class="mr-2 h-4 w-4" />
                    项目主页
                  </Button>
                  <Button
                    v-if="work.demo_url"
                    variant="outline"
                    size="lg"
                    class="work-link-btn demo-btn"
                    @click="openLink(work.demo_url)"
                  >
                    <ExternalLink class="mr-2 h-4 w-4" />
                    在线演示
                  </Button>
                </div>

                <div class="work-stats">
                  <span class="inline-flex items-center gap-1"><Eye class="h-4 w-4" /> {{ work.view_count }}</span>
                  <span class="inline-flex items-center gap-1"><MessageCircle class="h-4 w-4" /> {{ work.comment_count }}</span>
                </div>

                <div class="work-actions" v-if="work && work.is_published">
                  <Button
                    :variant="isLiked ? 'default' : 'outline'"
                    @click="handleLike"
                    :disabled="liking"
                    size="default"
                    class="action-btn"
                  >
                    <Loader2 v-if="liking" class="mr-1 h-4 w-4 animate-spin" />
                    <Star v-else class="mr-1 h-4 w-4" />
                    {{ work.like_count || 0 }} {{ isLiked ? '已点赞' : '点赞' }}
                  </Button>
                  <Button
                    :variant="isFavorited ? 'accent' : 'outline'"
                    @click="handleFavorite"
                    :disabled="favoriting"
                    size="default"
                    class="action-btn"
                  >
                    <Loader2 v-if="favoriting" class="mr-1 h-4 w-4 animate-spin" />
                    <Star v-else class="mr-1 h-4 w-4" />
                    {{ work.favorite_count || 0 }} {{ isFavorited ? '已收藏' : '收藏' }}
                  </Button>
                </div>
              </div>
            </div>
          </template>

          <div v-if="isWorkOwner" class="edit-actions-bottom">
            <Button
              variant="default"
              size="lg"
              @click="handleEdit"
            >
              <Pencil class="mr-1 h-4 w-4" />
              编辑作品
            </Button>
          </div>

          <Separator v-if="workCommentEnabled && work && work.is_published" />

          <div v-if="workCommentEnabled && work && work.is_published" class="comment-section">
            <h3>评论 ({{ work.comment_count }})</h3>

            <form v-if="userStore.isLoggedIn" @submit.prevent="submitComment" class="comment-form">
              <Textarea
                v-model="commentContent"
                :rows="4"
                placeholder="写下你的评论..."
                maxlength="500"
              />
              <Button variant="default" @click="submitComment" :disabled="submittingComment" class="mt-2">
                <Loader2 v-if="submittingComment" class="mr-1 h-4 w-4 animate-spin" />
                发表评论
              </Button>
            </form>
            <Alert v-else>
              <AlertDescription>
                请<router-link to="/login" class="text-primary underline hover:no-underline">登录</router-link>后发表评论
              </AlertDescription>
            </Alert>

            <div class="comment-list">
              <div v-for="comment in comments" :key="comment.id" class="comment-item">
                <Avatar
                  class="clickable-avatar h-10 w-10"
                  @click="handleUserClick(comment.user_id)"
                >
                  <AvatarImage :src="comment.user?.avatar" />
                  <AvatarFallback>{{ (comment.user?.nickname || comment.nickname || '?')[0] }}</AvatarFallback>
                </Avatar>
                <div class="comment-content">
                  <div class="comment-header">
                    <div class="comment-author-section">
                      <span class="comment-author">
                        {{ comment.user?.nickname || comment.nickname }}
                      </span>
                      <Badge v-if="isCommentAuthor(comment)" variant="accent" class="text-xs">
                        作者
                      </Badge>
                    </div>
                    <div class="comment-meta-section">
                      <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                    </div>
                  </div>
                  <p>{{ comment.content }}</p>
                  <div class="comment-actions">
                    <Button
                      variant="ghost"
                      size="sm"
                      :class="comment.isLiked ? 'text-primary' : ''"
                      @click="handleCommentLike(comment)"
                      :disabled="comment.likeLoading"
                    >
                      <Loader2 v-if="comment.likeLoading" class="mr-1 h-3 w-3 animate-spin" />
                      <Star v-else class="mr-1 h-3 w-3" />
                      {{ comment.like_count || 0 }}
                    </Button>
                    <Button
                      v-if="userStore.isLoggedIn"
                      variant="ghost"
                      size="sm"
                      @click="showReplyInput(comment)"
                    >
                      <MessageCircle class="mr-1 h-3 w-3" /> 回复
                    </Button>
                    <Button
                      v-if="canDeleteComment(comment)"
                      variant="ghost"
                      size="sm"
                      class="text-destructive hover:text-destructive"
                      @click="handleDeleteComment(comment)"
                    >
                      删除
                    </Button>
                  </div>

                  <div v-if="comment.showReply" class="reply-input">
                    <Textarea
                      v-model="comment.replyContent"
                      :rows="3"
                      placeholder="写下你的回复..."
                      maxlength="500"
                    />
                    <div class="reply-actions">
                      <Button variant="outline" size="sm" @click="cancelReply(comment)">取消</Button>
                      <Button variant="default" size="sm" @click="submitReply(comment)" :disabled="comment.replying">
                        <Loader2 v-if="comment.replying" class="mr-1 h-3 w-3 animate-spin" />
                        发表回复
                      </Button>
                    </div>
                  </div>

                  <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
                    <div v-for="reply in getDisplayedReplies(comment)" :key="reply.id" class="reply-item">
                      <Avatar
                        class="clickable-avatar h-8 w-8"
                        @click="handleUserClick(reply.user_id)"
                      >
                        <AvatarImage :src="reply.user?.avatar" />
                        <AvatarFallback>{{ (reply.user?.nickname || reply.nickname || '?')[0] }}</AvatarFallback>
                      </Avatar>
                      <div class="reply-content">
                        <div class="reply-header">
                          <span class="reply-author">
                            {{ reply.user?.nickname || reply.nickname }}
                            <Badge v-if="isCommentAuthor(reply)" variant="accent" class="text-xs">
                              作者
                            </Badge>
                            <span v-if="reply.parent_id" class="reply-to">
                              <template v-if="reply.parent_id === comment.id">
                                回复 @{{ comment.user?.nickname || comment.nickname }}
                              </template>
                              <template v-else>
                                回复 @{{ getReplyTargetName(comment, reply) }}
                              </template>
                            </span>
                          </span>
                          <span class="reply-time">{{ formatDate(reply.created_at) }}</span>
                        </div>
                        <p>{{ reply.content }}</p>
                        <div class="reply-actions">
                          <Button
                            variant="ghost"
                            size="sm"
                            :class="reply.isLiked ? 'text-primary' : ''"
                            @click="handleCommentLike(reply)"
                            :disabled="reply.likeLoading"
                          >
                            <Loader2 v-if="reply.likeLoading" class="mr-1 h-3 w-3 animate-spin" />
                            <Star v-else class="mr-1 h-3 w-3" />
                            {{ reply.like_count || 0 }}
                          </Button>
                          <Button
                            v-if="userStore.isLoggedIn"
                            variant="ghost"
                            size="sm"
                            @click="showReplyToReplyInput(comment, reply)"
                          >
                            <MessageCircle class="mr-1 h-3 w-3" /> 回复
                          </Button>
                          <Button
                            v-if="canDeleteComment(reply)"
                            variant="ghost"
                            size="sm"
                            class="text-destructive hover:text-destructive"
                            @click="handleDeleteComment(reply)"
                          >
                            删除
                          </Button>
                        </div>

                        <div v-if="reply.showReply" class="reply-input">
                          <Textarea
                            v-model="reply.replyContent"
                            :rows="3"
                            placeholder="写下你的回复..."
                            maxlength="500"
                          />
                          <div class="reply-actions">
                            <Button variant="outline" size="sm" @click="cancelReplyToReply(reply)">取消</Button>
                            <Button variant="default" size="sm" @click="submitReplyToReply(comment, reply)" :disabled="reply.replying">
                              <Loader2 v-if="reply.replying" class="mr-1 h-3 w-3 animate-spin" />
                              发表回复
                            </Button>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div v-if="hasMoreReplies(comment)" class="load-more-replies">
                      <Button
                        variant="ghost"
                        size="sm"
                        @click="loadMoreReplies(comment)"
                        :disabled="comment.loadingReplies"
                      >
                        <Loader2 v-if="comment.loadingReplies" class="mr-1 h-3 w-3 animate-spin" />
                        <ChevronDown v-else class="mr-1 h-3 w-3" />
                        {{ comment.loadingReplies ? '加载中...' : `展开更多回复 (${comment.reply_count - getDisplayedReplies(comment).length} 条)` }}
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="load-more-container" v-if="hasMoreComments">
              <Button
                variant="outline"
                @click="loadMoreComments"
                :disabled="loadingMoreComments"
                class="load-more-btn"
              >
                <Loader2 v-if="loadingMoreComments" class="mr-1 h-4 w-4 animate-spin" />
                <ChevronDown v-else class="mr-1 h-4 w-4" />
                {{ loadingMoreComments ? '加载中...' : '加载更多评论' }}
              </Button>
              <div class="comment-count-info">
                已显示 {{ comments.length }} / {{ commentsTotal }} 条评论
              </div>
            </div>
            <div v-else-if="comments.length > 0" class="no-more-comments">
              已显示全部 {{ comments.length }} 条评论
            </div>
          </div>
        </CardContent>
      </Card>

      <div v-else-if="loading" class="loading-container">
        <div class="flex items-center justify-center min-h-[600px]">
          <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
          <span class="ml-2 text-muted-foreground">加载中...</span>
        </div>
      </div>
      <EmptyState v-else description="作品不存在" />

      <!-- Photo Preview Dialog -->
      <Dialog :open="showPhotoPreview" @update:open="showPhotoPreview = false">
        <DialogContent class="max-w-4xl p-0 overflow-hidden">
          <div class="relative">
            <img :src="photos[previewPhotoIndex]?.url" :alt="photos[previewPhotoIndex]?.description || work?.title || ''" class="w-full object-contain max-h-[80vh]" />
            <div v-if="photos.length > 1" class="absolute inset-0 flex items-center justify-between px-2">
              <Button variant="ghost" size="icon" class="bg-black/30 hover:bg-black/50 text-white" @click="previewPhotoIndex = (previewPhotoIndex - 1 + photos.length) % photos.length">
                <ChevronLeft class="h-6 w-6" />
              </Button>
              <Button variant="ghost" size="icon" class="bg-black/30 hover:bg-black/50 text-white" @click="previewPhotoIndex = (previewPhotoIndex + 1) % photos.length">
                <ChevronRight class="h-6 w-6" />
              </Button>
            </div>
            <div class="absolute bottom-4 right-4 bg-black/60 text-white px-3 py-1 rounded-full text-sm">
              {{ previewPhotoIndex + 1 }} / {{ photos.length }}
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>

    <!-- Confirmation Dialog -->
    <Dialog :open="showConfirmDialog" @update:open="onConfirmDialogUpdateOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>确认</DialogTitle>
          <DialogDescription>{{ confirmDialogMessage }}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="cancelConfirmDialog">取消</Button>
          <Button @click="confirmDialogCallback?.()">确认</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import {
  Star, Bookmark, Eye, Clock, Pencil, MessageCircle, ChevronDown,
  MapPin, Calendar, Image, ExternalLink, Plus, Check, Loader2,
  ChevronLeft, ChevronRight
} from 'lucide-vue-next'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { loadCodeTheme, loadHighlightTheme, getMarkdownTheme } from '@/utils/codeTheme'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Textarea } from '@/components/ui/textarea'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Separator } from '@/components/ui/separator'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'
import { EmptyState } from '@/components/ui/empty-state'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const work = ref(null)
const loading = ref(true)
const comments = ref([])
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentContent = ref('')
const submittingComment = ref(false)
const loadingMoreComments = ref(false)
const hasMoreComments = ref(false)

const currentPhotoIndex = ref(0)
const carouselHeight = ref(600)

const isLiked = ref(false)
const isFavorited = ref(false)
const liking = ref(false)
const favoriting = ref(false)

const isFollowing = ref(false)
const followLoading = ref(false)
const workCommentEnabled = ref(true)

const showPhotoPreview = ref(false)
const previewPhotoIndex = ref(0)

const showConfirmDialog = ref(false)
const confirmDialogMessage = ref('')
const confirmDialogCallback = ref(null)
let _confirmDialogReject = null

const confirmDialog = (message) => {
  return new Promise((resolve, reject) => {
    confirmDialogMessage.value = message
    _confirmDialogReject = reject
    confirmDialogCallback.value = () => {
      _confirmDialogReject = null
      showConfirmDialog.value = false
      resolve()
    }
    showConfirmDialog.value = true
  })
}

const onConfirmDialogUpdateOpen = (open) => {
  showConfirmDialog.value = open
  if (!open && _confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    rejectFn('cancel')
  }
}

const cancelConfirmDialog = () => {
  if (_confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    showConfirmDialog.value = false
    rejectFn('cancel')
  }
}

const photos = computed(() => {
  if (!work.value || work.value.type !== 'photography') return []
  return work.value.images || []
})

const photoUrlList = computed(() => {
  return photos.value.map(photo => photo.url)
})

const currentPhoto = computed(() => {
  return photos.value[currentPhotoIndex.value] || { metadata: {} }
})

const isWorkOwner = computed(() => {
  if (!userStore.isLoggedIn || !work.value || !userStore.user) {
    return false
  }
  
  const authorId = work.value.author_id || work.value.author?.id
  const userId = userStore.user.id
  
  if (!authorId || !userId) {
    return false
  }
  
  return Number(authorId) === Number(userId)
})

const getStatusText = (status) => {
  switch (status) {
    case 0:
      return '草稿'
    case 1:
      return '已发布'
    case 2:
      return '审核中'
    case 3:
      return '审核不通过'
    default:
      return '未知状态'
  }
}

const getStatusAlertType = (status) => {
  switch (status) {
    case 0:
      return 'info'
    case 2:
      return 'warning'
    case 3:
      return 'error'
    default:
      return 'info'
  }
}

const getStatusAlertClass = (status) => {
  switch (status) {
    case 0:
      return 'border-blue-200 bg-blue-50 text-blue-800 dark:border-blue-800 dark:bg-blue-950 dark:text-blue-200'
    case 2:
      return 'border-yellow-200 bg-yellow-50 text-yellow-800 dark:border-yellow-800 dark:bg-yellow-950 dark:text-yellow-200'
    case 3:
      return 'border-red-200 bg-red-50 text-red-800 dark:border-red-800 dark:bg-red-950 dark:text-red-200'
    default:
      return 'border-blue-200 bg-blue-50 text-blue-800 dark:border-blue-800 dark:bg-blue-950 dark:text-blue-200'
  }
}

const isCommentAuthor = (comment) => {
  if (!work.value || !comment) return false
  
  const workAuthorId = work.value.author_id || work.value.author?.id
  if (!workAuthorId) return false
  
  const commentAuthorId = comment.user_id || comment.user?.id
  if (!commentAuthorId) return false
  
  return Number(workAuthorId) === Number(commentAuthorId)
}

const handleEdit = () => {
  if (!work.value) return
  router.push(`/dashboard/works/${work.value.id}/edit`)
}

const openPhotoPreview = (index) => {
  previewPhotoIndex.value = index
  showPhotoPreview.value = true
}

const prevPhoto = () => {
  currentPhotoIndex.value = (currentPhotoIndex.value - 1 + photos.value.length) % photos.value.length
}

const nextPhoto = () => {
  currentPhotoIndex.value = (currentPhotoIndex.value + 1) % photos.value.length
}

const renderDescription = async () => {
  if (!work.value || !work.value.description) {
    return
  }
  
  const codeThemeValue = await loadCodeTheme()
  await loadHighlightTheme(codeThemeValue)
  const mdTheme = await getMarkdownTheme()
  
  await new Promise(resolve => setTimeout(resolve, 150))
  
  nextTick(() => {
    const previewId = work.value.type === 'photography' 
      ? 'work-description-preview-photography' 
      : 'work-description-preview-project'
    
    const previewDiv = document.getElementById(previewId)
    if (!previewDiv) {
      console.error(`Description preview element not found: ${previewId}`)
      return
    }
    
    previewDiv.innerHTML = ''
    
    console.log('Rendering work description with Vditor.preview, code theme:', codeThemeValue, 'markdown theme:', mdTheme)
    
    Vditor.preview(previewDiv, work.value.description, {
      mode: mdTheme || 'light',
      markdown: {
        toc: true,
        mark: true,
        footnotes: true,
        autoSpace: true,
      },
      hljs: {
        lineNumber: true,
        style: codeThemeValue || 'github',
        enable: true
      },
      speech: {
        enable: false
      },
      anchor: 1,
      after: () => {
        console.log('Work description rendered successfully!')
      }
    })
  })
}

const loadWork = async (skipView = false) => {
  if (!skipView) {
    loading.value = true
  }
  try {
    const url = skipView 
      ? `/works/${route.params.id}?skip_view=true`
      : `/works/${route.params.id}`
    const response = await api.get(url)
    work.value = response.data
    renderDescription()
  } catch (error) {
    work.value = null
    toast.error('加载作品失败')
  } finally {
    loading.value = false
  }
}

const loadComments = async (append = false) => {
  try {
    const response = await api.get('/comments', {
      params: {
        work_id: route.params.id,
        page: commentsPage.value,
        page_size: 10
      }
    })
    const commentsList = response.data.list || []
    console.debug('Loaded comments:', commentsList.map(c => ({ id: c.id, like_count: c.like_count })))
    
    const existingComments = append ? comments.value : []
    const newComments = commentsList.map(comment => {
      const existingComment = existingComments.find(c => c.id === comment.id)
      return {
        ...comment,
        like_count: comment.like_count !== undefined ? comment.like_count : 0,
        isLiked: false,
        likeLoading: false,
        showReply: existingComment?.showReply || false,
        replyContent: existingComment?.replyContent || '',
        replying: existingComment?.replying || false,
        replyTo: existingComment?.replyTo || null,
        loadingReplies: existingComment?.loadingReplies || false,
        repliesPage: existingComment?.repliesPage || 1,
        replies: (comment.replies || []).map(reply => {
          const existingReply = existingComment?.replies?.find(r => r.id === reply.id)
          return {
            ...reply,
            like_count: reply.like_count !== undefined ? reply.like_count : 0,
            isLiked: false,
            likeLoading: false,
            showReply: existingReply?.showReply || false,
            replyContent: existingReply?.replyContent || '',
            replying: existingReply?.replying || false,
            replyTo: existingReply?.replyTo || null
          }
        })
      }
    })
    
    if (append) {
      comments.value = [...comments.value, ...newComments]
    } else {
      comments.value = newComments
    }
    
    commentsTotal.value = response.data.total || 0
    
    hasMoreComments.value = comments.value.length < commentsTotal.value
    
    const checkPromises = []
    for (const comment of newComments) {
      checkPromises.push(checkCommentLiked(comment))
      if (comment.replies && comment.replies.length > 0) {
        for (const reply of comment.replies) {
          checkPromises.push(checkCommentLiked(reply))
        }
      }
    }
    await Promise.all(checkPromises).catch(err => {
      console.error('Failed to check comment like status:', err)
    })
    console.debug('Comments loaded and like status checked:', comments.value.map(c => ({ id: c.id, like_count: c.like_count, isLiked: c.isLiked })))
  } catch (error) {
    console.error('Failed to load comments:', error)
  }
}

const loadMoreComments = async () => {
  if (loadingMoreComments.value || !hasMoreComments.value) return
  
  loadingMoreComments.value = true
  try {
    commentsPage.value++
    await loadComments(true)
  } finally {
    loadingMoreComments.value = false
  }
}

const canDeleteComment = (comment) => {
  if (!userStore.isLoggedIn || !userStore.user) return false
  
  if (userStore.user.role === 'admin') return true
  
  if (isWorkOwner.value) return true
  
  const commentUserId = comment.user_id || comment.user?.id || comment.userId
  
  if (!commentUserId || commentUserId === 0) {
    return false
  }
  
  return Number(userStore.user.id) === Number(commentUserId)
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    toast.warning('请输入评论内容')
    return
  }

  submittingComment.value = true
  try {
    const response = await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: commentContent.value
    })
    if (response.message && response.message !== 'success') {
      toast.success(response.message)
    }
    commentContent.value = ''
    commentsPage.value = 1
    await loadComments(false)
    await loadWork(true)
  } catch (error) {
    toast.error('评论发表失败')
  } finally {
    submittingComment.value = false
  }
}

const showReplyInput = (comment, replyTo = null) => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }
  
  comments.value.forEach(c => {
    if (c.id !== comment.id) {
      c.showReply = false
    }
  })
  
  comment.showReply = true
  comment.replyTo = replyTo || comment
  comment.replyContent = replyTo ? `@${replyTo.user?.nickname || replyTo.nickname} ` : ''
}

const showReplyToReplyInput = (parentComment, replyTo) => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }
  
  comments.value.forEach(c => {
    c.showReply = false
    if (c.replies) {
      c.replies.forEach(r => {
        r.showReply = false
      })
    }
  })
  
  replyTo.showReply = true
  replyTo.replyTo = replyTo
  replyTo.replyContent = `@${replyTo.user?.nickname || replyTo.nickname} `
}

const cancelReply = (comment) => {
  comment.showReply = false
  comment.replyContent = ''
  comment.replyTo = null
}

const submitReply = async (comment) => {
  if (!comment.replyContent || !comment.replyContent.trim()) {
    toast.warning('请输入回复内容')
    return
  }

  let content = comment.replyContent.trim()
  if (comment.replyTo) {
    const mentionPattern = new RegExp(`^@${comment.replyTo.user?.nickname || comment.replyTo.nickname}\\s*`, 'i')
    content = content.replace(mentionPattern, '')
  }

  if (!content) {
    toast.warning('请输入回复内容')
    return
  }

  comment.replying = true
  try {
    const response = await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: content,
      parent_id: comment.replyTo ? comment.replyTo.id : comment.id
    })
    if (response.message && response.message !== 'success') {
      toast.success(response.message)
    }
    comment.showReply = false
    comment.replyContent = ''
    comment.replyTo = null
    commentsPage.value = 1
    await loadComments(false)
    await loadWork(true)
  } catch (error) {
    toast.error('回复失败')
  } finally {
    comment.replying = false
  }
}

const submitReplyToReply = async (parentComment, reply) => {
  if (!reply.replyContent || !reply.replyContent.trim()) {
    toast.warning('请输入回复内容')
    return
  }

  let content = reply.replyContent.trim()
  if (reply.replyTo) {
    const mentionPattern = new RegExp(`^@${reply.replyTo.user?.nickname || reply.replyTo.nickname}\\s*`, 'i')
    content = content.replace(mentionPattern, '')
  }

  if (!content) {
    toast.warning('请输入回复内容')
    return
  }

  reply.replying = true
  try {
    const response = await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: content,
      parent_id: reply.replyTo ? reply.replyTo.id : reply.id
    })
    if (response.message && response.message !== 'success') {
      toast.success(response.message)
    }
    reply.showReply = false
    reply.replyContent = ''
    reply.replyTo = null
    commentsPage.value = 1
    await loadComments(false)
    await loadWork(true)
  } catch (error) {
    toast.error('回复失败')
  } finally {
    reply.replying = false
  }
}

const cancelReplyToReply = (reply) => {
  reply.showReply = false
  reply.replyContent = ''
  reply.replyTo = null
}

const getDisplayedReplies = (comment) => {
  if (!comment.replies) return []
  return comment.replies
}

const hasMoreReplies = (comment) => {
  if (!comment.reply_count) return false
  const displayedCount = getDisplayedReplies(comment).length
  return comment.reply_count > displayedCount
}

const loadMoreReplies = async (comment) => {
  if (!comment || comment.loadingReplies) return
  
  comment.loadingReplies = true
  try {
    const nextPage = (comment.repliesPage || 1) + 1
    const response = await api.get(`/comments/replies/${comment.id}`, {
      params: {
        page: nextPage,
        page_size: 10
      }
    })
    
    const newReplies = (response.data.list || []).map(reply => ({
      ...reply,
      like_count: reply.like_count !== undefined ? reply.like_count : 0,
      isLiked: false,
      likeLoading: false,
      showReply: false,
      replyContent: '',
      replying: false,
      replyTo: null
    }))
    
    comment.replies = [...(comment.replies || []), ...newReplies]
    comment.repliesPage = nextPage
    
    for (const reply of newReplies) {
      await checkCommentLiked(reply)
    }
  } catch (error) {
    console.error('Failed to load more replies:', error)
    toast.error('加载更多回复失败')
  } finally {
    comment.loadingReplies = false
  }
}

const getReplyTargetName = (parentComment, reply) => {
  if (!reply.parent_id || reply.parent_id === parentComment.id) {
    return parentComment.user?.nickname || parentComment.nickname
  }
  const targetReply = parentComment.replies?.find(r => r.id === reply.parent_id)
  return targetReply ? (targetReply.user?.nickname || targetReply.nickname) : ''
}

const checkCommentLiked = async (comment) => {
  if (!comment || !comment.id) return
  
  try {
    const response = await api.get(`/comments/${comment.id}/is-liked`)
    const isLiked = response.data?.is_liked ?? false
    comment.isLiked = isLiked
    console.debug(`Comment ${comment.id} like status:`, isLiked)
  } catch (error) {
    comment.isLiked = false
    console.debug(`Failed to check comment ${comment.id} like status:`, error)
  }
}

const handleCommentLike = async (comment) => {
  if (!comment || !comment.id || comment.likeLoading) return
  
  comment.likeLoading = true
  const oldIsLiked = comment.isLiked
  const commentId = comment.id
  
  try {
    if (comment.isLiked) {
      console.debug(`Unliking comment ${comment.id}, current count: ${comment.like_count}`)
      const response = await api.delete(`/comments/${comment.id}/like`)
      comment.isLiked = false
      commentsPage.value = 1
      await loadComments(false)
      const updatedComment = comments.value.find(c => c.id === commentId) || 
                            comments.value.flatMap(c => c.replies || []).find(r => r.id === commentId)
      if (updatedComment) {
        await checkCommentLiked(updatedComment)
      }
      const message = response?.message || '取消点赞成功'
      toast.success(message)
      console.debug(`Comment ${comment.id} unliked, new count:`, updatedComment?.like_count)
    } else {
      console.debug(`Liking comment ${comment.id}, current count: ${comment.like_count}`)
      const response = await api.post(`/comments/${comment.id}/like`)
      comment.isLiked = true
      commentsPage.value = 1
      await loadComments(false)
      const updatedComment = comments.value.find(c => c.id === commentId) || 
                            comments.value.flatMap(c => c.replies || []).find(r => r.id === commentId)
      if (updatedComment) {
        await checkCommentLiked(updatedComment)
      }
      const message = response?.message || '点赞成功'
      toast.success(message)
      console.debug(`Comment ${comment.id} liked, new count:`, updatedComment?.like_count)
    }
  } catch (error) {
    comment.isLiked = oldIsLiked
    console.error(`Failed to ${oldIsLiked ? 'unlike' : 'like'} comment ${comment.id}:`, error)
    await checkCommentLiked(comment)
    commentsPage.value = 1
    await loadComments(false)
  } finally {
    comment.likeLoading = false
  }
}

const handleDeleteComment = async (comment) => {
  try {
    await confirmDialog('确定要删除这条评论吗？')
    await api.delete(`/comments/${comment.id}`)
    toast.success('删除成功')
    commentsPage.value = 1
    await loadComments(false)
    await loadWork(true)
  } catch (error) {
    if (error !== 'cancel') {
      const errorMessage = error.response?.data?.message || error.message || '删除失败'
      toast.error(errorMessage)
    }
  }
}

const handleUserClick = (userId) => {
  if (!userId) return
  router.push(`/users/${userId}`)
}

const handlePhotoChange = (index) => {
  currentPhotoIndex.value = index
}

const setActivePhoto = (index) => {
  currentPhotoIndex.value = index
}

const handleImageLoad = () => {
}

const goToUserProfile = (userId) => {
  if (userId) {
    router.push(`/users/${userId}`)
  }
}

const goToMyProfile = () => {
  if (userStore.user?.id) {
    router.push(`/users/${userStore.user.id}`)
  }
}

const checkFollowStatus = async () => {
  if (!userStore.isLoggedIn || !work.value || !work.value.author) {
    isFollowing.value = false
    return
  }

  if (isWorkOwner.value) {
    isFollowing.value = false
    return
  }

  try {
    const authorId = work.value.author_id || work.value.author?.id
    if (!authorId) return

    const response = await api.get(`/users/${authorId}/follow-stats`)
    isFollowing.value = response.data?.is_following || false
  } catch (error) {
    console.error('Failed to check follow status:', error)
    isFollowing.value = false
  }
}

const handleFollow = async () => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }

  if (!work.value || !work.value.author) {
    return
  }

  const authorId = work.value.author_id || work.value.author?.id
  if (!authorId) {
    return
  }

  if (Number(authorId) === Number(userStore.user?.id)) {
    toast.warning('不能关注自己')
    return
  }

  followLoading.value = true
  try {
    if (isFollowing.value) {
      await api.delete(`/users/${authorId}/follow`)
      toast.success('已取消关注')
      isFollowing.value = false
    } else {
      await api.post(`/users/${authorId}/follow`)
      toast.success('关注成功')
      isFollowing.value = true
    }
  } catch (error) {
    const errorMsg = error.response?.data?.message || error.message || '操作失败'
    toast.error(errorMsg)
    if (errorMsg.includes('已经关注过')) {
      isFollowing.value = true
    } else if (errorMsg.includes('未关注')) {
      isFollowing.value = false
    }
  } finally {
    followLoading.value = false
  }
}

const openLink = (url) => {
  if (url) {
    window.open(url, '_blank')
  }
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now - d
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`
  
  return d.toLocaleDateString('zh-CN')
}

const updateCarouselHeight = () => {
  const width = window.innerWidth
  if (width < 768) {
    carouselHeight.value = 300
  } else if (width < 1200) {
    carouselHeight.value = 450
  } else {
    carouselHeight.value = 600
  }
}

const checkLikedStatus = async () => {
  if (!userStore.isLoggedIn) {
    isLiked.value = false
    return
  }
  try {
    const response = await api.get(`/works/${route.params.id}/liked`)
    isLiked.value = response.data.liked || response.data.is_liked || false
  } catch (error) {
    console.error('Failed to check liked status:', error)
    isLiked.value = false
  }
}

const checkFavoritedStatus = async () => {
  if (!userStore.isLoggedIn) {
    isFavorited.value = false
    return
  }
  try {
    const response = await api.get(`/works/${route.params.id}/favorited`)
    isFavorited.value = response.data.favorited || response.data.is_favorited || false
  } catch (error) {
    console.error('Failed to check favorited status:', error)
    isFavorited.value = false
  }
}

const handleLike = async () => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }

  if (liking.value) {
    return
  }

  liking.value = true
  try {
    await api.post(`/works/${route.params.id}/like`)
    
    await loadWork(true)
    await checkLikedStatus()
    
    toast.success(isLiked.value ? '点赞成功' : '取消点赞')
  } catch (error) {
    await checkLikedStatus()
    await loadWork(true)
  } finally {
    liking.value = false
  }
}

const handleFavorite = async () => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }

  if (favoriting.value) {
    return
  }

  favoriting.value = true
  try {
    if (isFavorited.value) {
      await api.delete(`/works/${route.params.id}/favorite`)
      toast.success('取消收藏')
    } else {
      await api.post(`/works/${route.params.id}/favorite`)
      toast.success('收藏成功')
    }
    
    await loadWork(true)
    await checkFavoritedStatus()
  } catch (error) {
    await checkFavoritedStatus()
    await loadWork(true)
  } finally {
    favoriting.value = false
  }
}

watch(() => work.value?.description, () => {
  renderDescription()
})

watch(() => work.value?.author?.id, () => {
  if (work.value) {
    checkFollowStatus()
  }
}, { immediate: false })

const loadCommentSettings = async () => {
  try {
    const response = await api.get('/settings/public')
    const settings = response.data || {}
    workCommentEnabled.value = settings.work_comment_enabled !== '0' && settings.work_comment_enabled !== 'false'
  } catch (error) {
    console.error('Failed to load comment settings:', error)
    workCommentEnabled.value = true
  }
}

onMounted(async () => {
  if (userStore.isLoggedIn && !userStore.user) {
    await userStore.fetchProfile()
  }
  
  const workId = route.params.id
  const preloadedWorkStr = sessionStorage.getItem(`preloaded_work_${workId}`)
  
  if (preloadedWorkStr) {
    try {
      const preloadedWork = JSON.parse(preloadedWorkStr)
      work.value = preloadedWork
      loading.value = false
      sessionStorage.removeItem(`preloaded_work_${workId}`)
      renderDescription()
      await loadCommentSettings()
      if (workCommentEnabled.value) {
        await loadComments()
      }
      checkLikedStatus()
      checkFavoritedStatus()
      checkFollowStatus()
    } catch (error) {
      console.error('Failed to parse preloaded work data:', error)
      await loadCommentSettings()
      await loadWork()
      if (workCommentEnabled.value) {
        await loadComments()
      }
      checkLikedStatus()
      checkFavoritedStatus()
      checkFollowStatus()
    }
  } else {
    await loadCommentSettings()
    
    await loadWork()
    if (workCommentEnabled.value) {
      await loadComments()
    }
    checkLikedStatus()
    checkFavoritedStatus()
    checkFollowStatus()
  }
  
  updateCarouselHeight()
  window.addEventListener('resize', updateCarouselHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateCarouselHeight)
})
</script>

<style>
@import 'vditor/dist/index.css';
</style>

<style scoped>
.work-detail {
  padding: var(--spacing-xl) 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.detail-card {
  max-width: 1400px;
  margin: 0 auto;
  box-shadow: var(--shadow-md);
  background-color: var(--theme-bg-card);
  border: 1px solid var(--theme-border-light);
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.work-status-banner {
  margin-bottom: var(--spacing-sm);
}

.audit-message-text {
  margin-top: var(--spacing-sm);
  font-size: var(--font-size-sm);
  font-weight: normal;
}

.loading-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: var(--spacing-md);
  min-height: 600px;
  position: relative;
}

.card-header-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.edit-actions-bottom {
  display: flex;
  justify-content: center;
  padding: var(--spacing-md) 0;
  margin-top: var(--spacing-md);
  border-top: 1px solid var(--theme-border-light);
}

.photography-layout {
  display: flex;
  gap: var(--spacing-lg);
}

.photo-section {
  flex: 7;
}

.main-photo-container {
  position: relative;
  background: var(--theme-bg-primary);
  border-radius: var(--radius-lg);
  overflow: hidden;
  margin-bottom: var(--spacing-md);
  box-shadow: var(--shadow-lg);
}

.carousel-wrapper {
  display: flex;
  transition: transform 0.3s ease;
  height: 100%;
}

.carousel-slide {
  flex: 0 0 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-photo {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  cursor: pointer;
}

.carousel-nav {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  pointer-events: none;
  padding: 0 0.5rem;
}

.carousel-nav-btn {
  pointer-events: auto;
  background: rgba(0, 0, 0, 0.3);
  color: white;
  border-radius: 50%;
  opacity: 0;
  transition: opacity 0.2s;
}

.main-photo-container:hover .carousel-nav-btn {
  opacity: 1;
}

.photo-counter {
  position: absolute;
  bottom: var(--spacing-lg);
  right: var(--spacing-lg);
  background: rgba(0, 0, 0, 0.75);
  color: var(--color-text-inverse);
  padding: 10px 20px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-base);
  font-weight: 500;
  backdrop-filter: blur(10px);
  box-shadow: var(--shadow-md);
}

.thumbnail-nav {
  display: flex;
  gap: var(--spacing-sm);
  overflow-x: auto;
  padding: var(--spacing-sm) 0;
}

.thumbnail-item {
  flex-shrink: 0;
  width: 80px;
  height: 80px;
  border-radius: var(--radius-md);
  overflow: hidden;
  cursor: pointer;
  border: 3px solid transparent;
  transition: all var(--transition-slow);
  box-shadow: var(--shadow-sm);
}

.thumbnail-item:hover {
  border-color: var(--theme-primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px color-mix(in srgb, var(--theme-primary) 30%, transparent);
}

.thumbnail-item.active {
  border-color: var(--theme-primary);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--theme-primary) 20%, transparent);
  transform: scale(1.05);
}

.thumbnail-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-description {
  padding: var(--spacing-md);
  background: var(--theme-bg-secondary);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  border-left: 4px solid var(--theme-primary);
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
  color: var(--theme-text-secondary);
}

.album-info {
  padding: var(--spacing-md) 0;
}

.album-info h2 {
  margin-bottom: var(--spacing-md);
  color: var(--theme-text-primary);
}

.album-description-empty {
  color: var(--theme-text-secondary);
  margin-bottom: var(--spacing-md);
  font-style: italic;
}

.album-meta {
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  color: var(--theme-text-secondary);
}

.info-sidebar {
  flex: 3;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.author-card,
.photo-params-card,
.stats-card {
  padding: var(--spacing-lg);
  background: var(--theme-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--theme-border-light);
}

.author-header {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
  cursor: pointer;
  transition: opacity var(--transition-base);
}

.author-header:hover {
  opacity: 0.8;
}

.author-info h3 {
  margin: 0 0 5px 0;
  font-size: var(--font-size-lg);
  color: var(--theme-text-primary);
}

.author-info p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.full-width-btn {
  width: 100%;
}

.photo-params-card h3,
.stats-card h3 {
  margin: 0 0 var(--spacing-md) 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--theme-text-primary);
  padding-bottom: var(--spacing-sm);
  border-bottom: 2px solid var(--theme-border-light);
}

.params-list,
.stats-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.param-item {
  display: flex;
  justify-content: space-between;
  font-size: var(--font-size-sm);
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--theme-border-light);
}

.param-item:last-child {
  border-bottom: none;
}

.param-label {
  color: var(--theme-text-tertiary);
  font-weight: 500;
}

.param-value {
  font-weight: 600;
  color: var(--theme-text-primary);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.action-buttons {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-md);
}

.action-buttons button {
  flex: 1;
  height: 44px;
  font-size: var(--font-size-base);
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--transition-slow);
  cursor: pointer;
}

.action-buttons button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.project-layout {
  max-width: 900px;
  margin: 0 auto;
}

.work-header {
  margin-bottom: var(--spacing-lg);
}

.work-title-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.work-title-section h1 {
  margin: 0;
  font-size: var(--font-size-3xl);
  color: var(--theme-text-primary);
}

.work-author {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 500;
  cursor: pointer;
  color: var(--theme-primary);
  transition: color var(--transition-fast);
}

.author-name:hover {
  color: var(--theme-primary-hover);
  text-decoration: underline;
}

.publish-time {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.work-content {
  line-height: var(--line-height-relaxed);
  padding: var(--spacing-lg);
  background-color: var(--theme-content-bg);
  border-radius: var(--radius-md);
  border: 1px solid var(--theme-border-light);
  box-shadow: var(--shadow-sm);
  margin: var(--spacing-md) 0;
}

.work-content :deep(pre) {
  line-height: var(--line-height-tight);
}

.work-description-empty {
  font-size: var(--font-size-lg);
  margin-bottom: var(--spacing-lg);
  color: var(--theme-text-secondary);
  font-style: italic;
}

#work-description-preview-photography {
  padding: var(--spacing-lg);
  background-color: var(--theme-content-bg);
  border-radius: var(--radius-md);
  border: 1px solid var(--theme-border-light);
  box-shadow: var(--shadow-sm);
  margin: var(--spacing-md) 0;
}

#work-description-preview-project {
  padding: 0;
  margin: 0;
}

#work-description-preview-photography :deep(code:not(pre code)),
#work-description-preview-project :deep(code:not(pre code)) {
  background-color: var(--theme-bg-hover);
  color: var(--theme-text-primary);
  padding: 0.2em 0.4em;
  border-radius: var(--radius-sm);
  font-size: 85%;
}

#work-description-preview-photography :deep(pre),
#work-description-preview-project :deep(pre) {
  line-height: var(--line-height-tight);
}

#work-description-preview-photography :deep(pre code),
#work-description-preview-project :deep(pre code) {
  line-height: inherit;
}

.work-links {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  flex-wrap: wrap;
}

.work-link-btn {
  min-width: 140px;
  height: 48px;
  font-size: var(--font-size-lg);
  font-weight: 600;
  border-radius: var(--radius-md);
  transition: all var(--transition-slow);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
}

.work-link-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.github-btn {
  background: linear-gradient(135deg, var(--theme-text-primary) 0%, var(--theme-text-secondary) 100%);
  border: none;
  color: var(--color-text-inverse);
}

.github-btn:hover {
  background: linear-gradient(135deg, var(--theme-text-secondary) 0%, var(--theme-text-primary) 100%);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-md);
}

.project-btn {
  background: linear-gradient(135deg, var(--theme-accent) 0%, var(--theme-primary) 100%);
  border: none;
  color: var(--color-text-inverse);
}

.project-btn:hover {
  background: linear-gradient(135deg, var(--theme-primary) 0%, var(--theme-accent) 100%);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-md);
}

.demo-btn {
  border: 2px solid var(--theme-primary);
  color: var(--theme-primary);
  background: var(--theme-bg-card);
}

.demo-btn:hover {
  background: var(--theme-primary);
  color: var(--color-text-inverse);
  border-color: var(--theme-primary);
  box-shadow: var(--shadow-md);
}

.tech-stack {
  margin-bottom: var(--spacing-lg);
}

.tech-stack h3 {
  margin-bottom: var(--spacing-sm);
  color: var(--theme-text-primary);
}

.work-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-md);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--theme-border-light);
  justify-content: center;
}

.action-btn {
  min-width: 120px;
  height: 40px;
  font-size: var(--font-size-sm);
  font-weight: 500;
  border-radius: var(--radius-md);
  transition: all var(--transition-slow);
  cursor: pointer;
}

.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.work-stats {
  display: flex;
  gap: var(--spacing-md);
  color: var(--theme-text-secondary);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--theme-border-light);
}

.work-stats span {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.comment-section {
  margin-top: var(--spacing-xl);
}

.comment-section h3 {
  margin-bottom: var(--spacing-md);
  color: var(--theme-text-primary);
}

.comment-form {
  margin-bottom: var(--spacing-lg);
}

.comment-list {
  margin-top: var(--spacing-lg);
}

.comment-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid var(--theme-border-light);
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.comment-author-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.comment-author {
  font-weight: 600;
  color: var(--theme-text-primary);
  font-size: var(--font-size-base);
}

.comment-meta-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.comment-time {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
}

.comment-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.reply-input {
  margin-top: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--theme-bg-secondary);
  border-radius: var(--radius-md);
}

.reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.replies-list {
  margin-top: var(--spacing-md);
  padding-left: var(--spacing-md);
  border-left: 2px solid var(--theme-border-light);
}

.reply-item {
  display: flex;
  gap: var(--spacing-sm);
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid var(--theme-border-light);
}

.reply-item:last-child {
  border-bottom: none;
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--spacing-sm);
}

.reply-author {
  font-weight: 600;
  color: var(--theme-text-primary);
  font-size: var(--font-size-sm);
}

.reply-to {
  color: var(--theme-text-secondary);
  font-weight: normal;
  font-size: var(--font-size-xs);
  margin-left: var(--spacing-xs);
}

.reply-time {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-xs);
}

.load-more-container {
  margin-top: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
}

.load-more-btn {
  min-width: 200px;
  cursor: pointer;
}

.comment-count-info {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.no-more-comments {
  margin-top: var(--spacing-lg);
  text-align: center;
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  padding: var(--spacing-md);
  background-color: var(--theme-bg-secondary);
  border-radius: var(--radius-sm);
}

.load-more-replies {
  margin-top: var(--spacing-sm);
  padding: var(--spacing-sm) 0;
  text-align: center;
  border-top: 1px solid var(--theme-border-light);
}

.clickable-avatar {
  cursor: pointer;
  transition: all var(--transition-slow);
}

.clickable-avatar:hover {
  transform: scale(1.1);
  box-shadow: var(--shadow-sm);
}

@media (max-width: 1200px) {
  .photography-layout {
    flex-direction: column;
  }

  .photo-section,
  .info-sidebar {
    flex: 1;
  }
}

@media (max-width: 768px) {
  .work-title-section {
    flex-direction: column;
    align-items: flex-start;
  }

  .thumbnail-nav {
    gap: var(--spacing-xs);
  }

  .thumbnail-item {
    width: 60px;
    height: 60px;
  }
}
</style>
