// internal/handler/file.go
package handler

import (
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "tech-test/backend/internal/domain"
    fileInterface "tech-test/backend/internal/service/interfaces/file"
    "tech-test/backend/internal/utils"
    "tech-test/backend/internal/middleware"
    "os"
    "fmt"
    "io"
    "math"
    "log"
    "path/filepath"
    "mime"
    "github.com/google/uuid"
    "tech-test/backend/internal/config"
    "go.uber.org/zap"
)

type FileHandler struct {
    fileService fileInterface.Service
    uploadDir   string
    config      config.FileConfig
    logger      *zap.Logger
}

func NewFileHandler(fileService fileInterface.Service, config config.FileConfig) *FileHandler {
    absUploadDir, err := filepath.Abs(config.UploadDir)
    if err != nil {
        log.Printf("Error getting absolute path: %v", err)
        absUploadDir = config.UploadDir
    }

    return &FileHandler{
        fileService: fileService,
        uploadDir:   absUploadDir,
        config:      config,
        logger:      zap.NewExample(),
    }
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
    userIDValue := r.Context().Value(middleware.UserIDKey)
    userID, ok := userIDValue.(uint)
    if !ok {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusUnauthorized,
            domain.ErrCodeAuthentication,
            "User ID not found in context",
            nil,
        ))
        return
    }

    if err := r.ParseMultipartForm(h.config.MaxSize); err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "File too large",
            err,
        ))
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "No file provided",
            err,
        ))
        return
    }
    defer file.Close()

    filename := uuid.New().String() + "_" + header.Filename
    filePath := filepath.Join(h.uploadDir, filename)

    if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
        h.logger.Error("Failed to create upload directory", 
            zap.String("dir", h.uploadDir),
            zap.Error(err))
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to create upload directory",
            err,
        ))
        return
    }

    dst, err := os.Create(filePath)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to create file",
            err,
        ))
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to save file",
            err,
        ))
        return
    }

    fileRecord := &domain.File{
        UserID:   userID,
        Name:     header.Filename,
        MimeType: header.Header.Get("Content-Type"),
        Size:     header.Size,
        Path:     filePath,
    }

    if err := h.fileService.Upload(r.Context(), fileRecord); err != nil {
        os.Remove(filePath)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to record file information",
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusCreated, fileRecord)
}

func (h *FileHandler) List(w http.ResponseWriter, r *http.Request) {
    files, err := h.fileService.List(r.Context())
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to list files",
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, files)
}

func (h *FileHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid file ID",
            err,
        ))
        return
    }

    file, err := h.fileService.GetByID(r.Context(), uint(id))
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            "File not found",
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, file)
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid file ID",
            err,
        ))
        return
    }

    file, err := h.fileService.GetByID(r.Context(), uint(id))
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            "File not found",
            err,
        ))
        return
    }

    os.Remove(file.Path)

    if err := h.fileService.Delete(r.Context(), uint(id)); err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to delete file record",
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{
        "message": "File deleted successfully",
    })
}

func (h *FileHandler) GetUserFiles(w http.ResponseWriter, r *http.Request) {
    userIDValue := r.Context().Value(middleware.UserIDKey)
    userID, ok := userIDValue.(uint)
    if !ok {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusUnauthorized,
            domain.ErrCodeAuthentication,
            "Invalid user ID",
            nil,
        ))
        return
    }

    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

    files, total, err := h.fileService.GetUserFilesPaginated(r.Context(), userID, page, pageSize)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to fetch files",
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
        "data":        files,
        "total":       total,           
        "page":        page,            
        "page_size":   pageSize,        
        "total_pages": math.Ceil(float64(total) / float64(pageSize)), 
    })
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileID, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid file ID",
            err,
        ))
        return
    }

    file, err := h.fileService.GetByID(r.Context(), uint(fileID))
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            "File not found",
            err,
        ))
        return
    }

    absPath := filepath.Join(h.uploadDir, filepath.Base(file.Path))
    log.Printf("Downloading file: ID=%d, RelativePath=%s, AbsolutePath=%s", fileID, file.Path, absPath)

    fileContent, err := os.Open(absPath)
    if err != nil {
        log.Printf("Error opening file: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Error opening file",
            err,
        ))
        return
    }
    defer fileContent.Close()

    mimeType := mime.TypeByExtension(filepath.Ext(file.Name))
    if mimeType == "" {
        mimeType = "application/octet-stream"
    }
    w.Header().Set("Content-Type", mimeType)

    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))

    _, err = io.Copy(w, fileContent)
    if err != nil {
        log.Printf("Error streaming file: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Error streaming file",
            err,
        ))
        return
    }

    log.Printf("File %s served successfully", absPath)
}

func (h *FileHandler) View(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileIDStr, exists := vars["id"]
    if !exists {
        http.Error(w, "File ID is missing", http.StatusBadRequest)
        return
    }

    fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
    if err != nil {
        log.Printf("Invalid file ID %s: %v", fileIDStr, err)
        http.Error(w, "Invalid file ID", http.StatusBadRequest)
        return
    }

    file, err := h.fileService.GetByID(r.Context(), uint(fileID))
    if err != nil {
        log.Printf("File with ID %d not found: %v", fileID, err)
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    filePath := file.Path

    log.Printf("Attempting to serve file: %s", filePath)

    http.ServeFile(w, r, filePath) 

    log.Printf("File %s served successfully", filePath)
}

func (h *FileHandler) ServeFiles(w http.ResponseWriter, r *http.Request) {
    http.StripPrefix("/files", http.FileServer(http.Dir(h.uploadDir))).ServeHTTP(w, r)
}

func (h *FileHandler) SearchFiles(w http.ResponseWriter, r *http.Request) {
    userIDValue := r.Context().Value(middleware.UserIDKey)
    userID, ok := userIDValue.(uint)
    if !ok {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusUnauthorized,
            domain.ErrCodeAuthentication,
            "User ID not found in context",
            nil,
        ))
        return
    }

    searchTerm := r.URL.Query().Get("q")

    log.Printf("Searching files for userID: %d with term: %s", userID, searchTerm)

    files, err := h.fileService.SearchFiles(r.Context(), userID, searchTerm)
    if err != nil {
        log.Printf("Error searching files: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to search files",
            err,
        ))
        return
    }

    log.Printf("Found %d files matching search term", len(files))
    utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
        "data": files,
    })
}

func (h *FileHandler) GenerateShareableLink(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileIDStr := vars["id"]

    fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid file ID",
            err,
        ))
        return
    }

    shareableID := uuid.New().String()

    err = h.fileService.UpdateShareableID(r.Context(), strconv.FormatUint(fileID, 10), shareableID)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to update file: " + err.Error(),
            err,
        ))
        return
    }

    shareableLink := fmt.Sprintf("%s/shared/%s", h.config.BaseURL, shareableID)

    response := map[string]interface{}{
        "shareableId":   shareableID,
        "shareableLink": shareableLink,
        "message":       "File shared successfully",
    }

    utils.RespondWithJSON(w, http.StatusOK, response)
}

func (h *FileHandler) GetSharedFile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    shareID := vars["shareId"]

    h.logger.Debug("Getting shared file", zap.String("shareId", shareID))

    
    file, err := h.fileService.GetByShareID(r.Context(), shareID)
    if err != nil {
        h.logger.Error("Failed to get file by share ID", 
            zap.String("shareId", shareID),
            zap.Error(err))
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            "File not found",
            err,
        ))
        return
    }

    filePath := file.Path  

    h.logger.Debug("Attempting to access file", 
        zap.String("path", filePath),
        zap.String("name", file.Name))

    
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        h.logger.Error("File does not exist on disk", 
            zap.String("path", filePath),
            zap.String("name", file.Name),
            zap.Error(err))
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            "File not found on disk",
            err,
        ))
        return
    }

    
    fileContent, err := os.Open(filePath)
    if err != nil {
        h.logger.Error("Failed to open file", 
            zap.String("path", filePath),
            zap.String("name", file.Name),
            zap.Error(err))
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to read file",
            err,
        ))
        return
    }
    defer fileContent.Close()

    contentType := file.MimeType  
    if contentType == "" {
         buffer := make([]byte, 512)
        _, err := fileContent.Read(buffer)
        if err != nil && err != io.EOF {
            h.logger.Error("Failed to read file header", 
                zap.String("name", file.Name),
                zap.Error(err))
            contentType = "application/octet-stream"
        } else {
            contentType = http.DetectContentType(buffer)
            _, err = fileContent.Seek(0, 0) 
            if err != nil {
                h.logger.Error("Failed to reset file pointer", 
                    zap.String("name", file.Name),
                    zap.Error(err))
            }
        }
    }

    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.Name))

     if _, err := io.Copy(w, fileContent); err != nil {
        h.logger.Error("Failed to stream file", 
            zap.String("path", filePath),
            zap.String("name", file.Name),
            zap.Error(err))
        return
    }

    h.logger.Info("Successfully served shared file", 
        zap.String("shareId", shareID),
        zap.String("name", file.Name),
        zap.String("contentType", contentType))
}

func (h *FileHandler) ShareFile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileIDStr := vars["id"]

    fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid file ID",
            err,
        ))
        return
    }

    shareableID := uuid.New().String()

    err = h.fileService.UpdateShareableID(r.Context(), strconv.FormatUint(fileID, 10), shareableID)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Failed to update file: " + err.Error(),
            err,
        ))
        return
    }

    response := map[string]interface{}{
        "shareableId": shareableID,
        "shareableLink": fmt.Sprintf("%s/shared/%s", "http://localhost:3000", shareableID),
        "message": "File shared successfully",
    }

    utils.RespondWithJSON(w, http.StatusOK, response)
}

