<template>
  <div id="app">
    <el-container class="app-container">
      <!-- 头部 -->
      <el-header class="app-header">
        <div class="header-content">
          <h1 class="app-title">
            <el-icon class="title-icon"><ChatDotRound /></el-icon>
            RAG问答助手
          </h1>
          <p class="app-subtitle">基于文档的智能问答系统</p>
        </div>
      </el-header>

      <!-- 主要内容区域 -->
      <el-main class="app-main">
        <el-row :gutter="20" class="main-content">
          <!-- 左侧：文档上传区域 -->
          <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
            <el-card class="upload-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><Upload /></el-icon>
                  <span>文档上传</span>
                </div>
              </template>
              
              <el-form :model="documentForm" :rules="documentRules" ref="documentFormRef" label-width="80px">
                <el-form-item label="文档标题" prop="title">
                  <el-input
                    v-model="documentForm.title"
                    placeholder="请输入文档标题"
                    clearable
                  />
                </el-form-item>
                
                <el-form-item label="文档内容" prop="content">
                  <el-input
                    v-model="documentForm.content"
                    type="textarea"
                    :rows="8"
                    placeholder="请输入或粘贴文档内容..."
                    resize="vertical"
                  />
                </el-form-item>
                
                <el-form-item>
                  <el-button 
                    type="primary" 
                    @click="uploadDocument"
                    :loading="uploading"
                    :disabled="!documentForm.title || !documentForm.content"
                    class="upload-btn"
                  >
                    <el-icon><Upload /></el-icon>
                    上传文档
                  </el-button>
                </el-form-item>
              </el-form>
            </el-card>
          </el-col>

          <!-- 右侧：问答区域 -->
          <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
            <el-card class="chat-card" shadow="hover">
              <template #header>
                <div class="card-header">
                  <el-icon><ChatDotRound /></el-icon>
                  <span>智能问答</span>
                </div>
              </template>

              <!-- 对话历史 -->
              <div class="chat-history" ref="chatHistoryRef">
                <div 
                  v-for="(message, index) in chatHistory" 
                  :key="index"
                  :class="['message', message.type]"
                >
                  <div class="message-content">
                    <div class="message-header">
                      <el-icon v-if="message.type === 'user'"><User /></el-icon>
                      <el-icon v-else><Robot /></el-icon>
                      <span class="message-label">
                        {{ message.type === 'user' ? '您' : '助手' }}
                      </span>
                    </div>
                    <div class="message-text">{{ message.content }}</div>
                    <div class="message-time">{{ message.time }}</div>
                  </div>
                </div>
                
                <!-- 加载状态 -->
                <div v-if="isLoading" class="message assistant">
                  <div class="message-content">
                    <div class="message-header">
                      <el-icon><Robot /></el-icon>
                      <span class="message-label">助手</span>
                    </div>
                    <div class="loading-content">
                      <el-icon class="is-loading"><Loading /></el-icon>
                      <span>正在思考中...</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 输入区域 -->
              <div class="chat-input">
                <el-form @submit.prevent="askQuestion">
                  <el-form-item>
                    <el-input
                      v-model="questionInput"
                      type="textarea"
                      :rows="3"
                      placeholder="请输入您的问题..."
                      resize="none"
                      @keydown.ctrl.enter="askQuestion"
                      :disabled="isLoading"
                    />
                  </el-form-item>
                  <el-form-item>
                    <el-button 
                      type="primary" 
                      @click="askQuestion"
                      :loading="isLoading"
                      :disabled="!questionInput.trim()"
                      class="ask-btn"
                    >
                      <el-icon><Position /></el-icon>
                      发送问题
                    </el-button>
                    <el-button @click="clearHistory" :disabled="chatHistory.length === 0">
                      <el-icon><Delete /></el-icon>
                      清空对话
                    </el-button>
                    <span class="input-tip">按 Ctrl+Enter 快速发送</span>
                  </el-form-item>
                </el-form>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { ref, reactive, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { apiService } from './services/api'

export default {
  name: 'App',
  setup() {
    // 文档上传相关
    const documentFormRef = ref()
    const documentForm = reactive({
      title: '',
      content: ''
    })
    const documentRules = {
      title: [
        { required: true, message: '请输入文档标题', trigger: 'blur' }
      ],
      content: [
        { required: true, message: '请输入文档内容', trigger: 'blur' }
      ]
    }
    const uploading = ref(false)

    // 问答相关
    const questionInput = ref('')
    const isLoading = ref(false)
    const chatHistory = ref([])
    const chatHistoryRef = ref()

    // 上传文档
    const uploadDocument = async () => {
      if (!documentFormRef.value) return
      
      try {
        await documentFormRef.value.validate()
        uploading.value = true
        
        const response = await apiService.uploadDocument(documentForm.title, documentForm.content)
        
        ElMessage.success(`文档上传成功！共处理 ${response.chunks} 个文档片段`)
        
        // 清空表单
        documentForm.title = ''
        documentForm.content = ''
        documentFormRef.value.resetFields()
        
      } catch (error) {
        console.error('上传文档失败:', error)
        ElMessage.error(error.message || '文档上传失败，请重试')
      } finally {
        uploading.value = false
      }
    }

    // 发送问题
    const askQuestion = async () => {
      if (!questionInput.value.trim() || isLoading.value) return
      
      const question = questionInput.value.trim()
      questionInput.value = ''
      
      // 添加用户消息到历史
      const userMessage = {
        type: 'user',
        content: question,
        time: new Date().toLocaleTimeString()
      }
      chatHistory.value.push(userMessage)
      
      // 滚动到底部
      await nextTick()
      scrollToBottom()
      
      // 发送请求
      isLoading.value = true
      
      try {
        const response = await apiService.askQuestion(question)
        
        // 添加助手回复到历史
        const assistantMessage = {
          type: 'assistant',
          content: response.answer,
          time: new Date().toLocaleTimeString()
        }
        chatHistory.value.push(assistantMessage)
        
        // 滚动到底部
        await nextTick()
        scrollToBottom()
        
      } catch (error) {
        console.error('问答失败:', error)
        ElMessage.error(error.message || '问答失败，请重试')
        
        // 添加错误消息
        const errorMessage = {
          type: 'assistant',
          content: '抱歉，我暂时无法回答您的问题，请稍后再试。',
          time: new Date().toLocaleTimeString()
        }
        chatHistory.value.push(errorMessage)
      } finally {
        isLoading.value = false
      }
    }

    // 清空对话历史
    const clearHistory = async () => {
      try {
        await ElMessageBox.confirm('确定要清空所有对话记录吗？', '确认清空', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        chatHistory.value = []
        ElMessage.success('对话记录已清空')
      } catch {
        // 用户取消
      }
    }

    // 滚动到底部
    const scrollToBottom = () => {
      if (chatHistoryRef.value) {
        chatHistoryRef.value.scrollTop = chatHistoryRef.value.scrollHeight
      }
    }

    return {
      documentFormRef,
      documentForm,
      documentRules,
      uploading,
      questionInput,
      isLoading,
      chatHistory,
      chatHistoryRef,
      uploadDocument,
      askQuestion,
      clearHistory
    }
  }
}
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.app-header {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-content {
  text-align: center;
  color: white;
}

.app-title {
  margin: 0;
  font-size: 2.5rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.title-icon {
  font-size: 2.5rem;
}

.app-subtitle {
  margin: 0.5rem 0 0 0;
  font-size: 1.1rem;
  opacity: 0.9;
}

.app-main {
  padding: 2rem;
}

.main-content {
  max-width: 1400px;
  margin: 0 auto;
}

.upload-card,
.chat-card {
  height: 100%;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  font-size: 1.1rem;
  color: #303133;
}

.upload-btn,
.ask-btn {
  width: 100%;
  height: 40px;
  font-size: 1rem;
  font-weight: 600;
}

.chat-history {
  height: 500px;
  overflow-y: auto;
  padding: 1rem 0;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 1rem;
}

.message {
  margin-bottom: 1.5rem;
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 80%;
  padding: 1rem;
  border-radius: 12px;
  position: relative;
}

.message.user .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  margin-right: 0;
}

.message.assistant .message-content {
  background: #f5f7fa;
  color: #303133;
  margin-left: 0;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
  font-weight: 600;
  font-size: 0.9rem;
}

.message-text {
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.message-time {
  font-size: 0.8rem;
  opacity: 0.7;
  margin-top: 0.5rem;
}

.loading-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #909399;
}

.chat-input {
  padding-top: 1rem;
}

.input-tip {
  font-size: 0.8rem;
  color: #909399;
  margin-left: 1rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-main {
    padding: 1rem;
  }
  
  .app-title {
    font-size: 2rem;
  }
  
  .main-content {
    gap: 1rem;
  }
  
  .message-content {
    max-width: 90%;
  }
  
  .chat-history {
    height: 400px;
  }
}

/* 滚动条样式 */
.chat-history::-webkit-scrollbar {
  width: 6px;
}

.chat-history::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.chat-history::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.chat-history::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
