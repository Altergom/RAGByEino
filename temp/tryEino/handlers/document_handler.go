package handlers

import (
	"github.com/Altergom/tryEino/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DocumentHandler struct {
	documentService *services.DocumentService
	milvusService   *services.MilvusService
}

type docReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func NewDocumentHandler(documentService *services.DocumentService, milvusService *services.MilvusService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
		milvusService:   milvusService,
	}
}

// UploadDocument 上传文档
func (dh *DocumentHandler) UploadDocument(c *gin.Context) {
	var req docReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 处理文档
	documents, err := dh.documentService.ProcessDocument(req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process document"})
		return
	}

	// 插入到Milvus
	if err := dh.milvusService.InsertDocument(documents); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document uploaded successfully",
		"chunks":  len(documents),
		"title":   req.Title,
	})
}
