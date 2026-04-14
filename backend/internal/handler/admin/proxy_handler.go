package admin

import (
	"strconv"

	"github.com/meitianwang/fast-frame/internal/pkg/pagination"
	"github.com/meitianwang/fast-frame/internal/pkg/response"
	middleware2 "github.com/meitianwang/fast-frame/internal/server/middleware"
	"github.com/meitianwang/fast-frame/internal/service"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// ProxyAdminHandler handles admin proxy management endpoints.
type ProxyAdminHandler struct {
	proxyService *service.ProxyService
}

// NewProxyAdminHandler creates a new ProxyAdminHandler.
func NewProxyAdminHandler(proxyService *service.ProxyService) *ProxyAdminHandler {
	return &ProxyAdminHandler{proxyService: proxyService}
}

// --- Node management ---

// ListNodes GET /api/v1/admin/proxy/nodes
func (h *ProxyAdminHandler) ListNodes(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	filter := service.ProxyNodeFilter{
		CountryCode: c.Query("country_code"),
		Status:      c.Query("status"),
	}
	nodes, result, err := h.proxyService.ListNodes(c.Request.Context(), filter,
		pagination.PaginationParams{Page: page, PageSize: pageSize})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]adminProxyNodeDTO, len(nodes))
	for i, n := range nodes {
		out[i] = toAdminProxyNodeDTO(&n)
	}
	response.Paginated(c, out, result.Total, result.Page, result.PageSize)
}

// CreateNode POST /api/v1/admin/proxy/nodes
func (h *ProxyAdminHandler) CreateNode(c *gin.Context) {
	var req adminProxyNodeDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	node := fromAdminProxyNodeDTO(&req)
	if err := h.proxyService.CreateNode(c.Request.Context(), node); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, toAdminProxyNodeDTO(node))
}

// UpdateNode PUT /api/v1/admin/proxy/nodes/:id
func (h *ProxyAdminHandler) UpdateNode(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	var req adminProxyNodeDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	req.ID = id
	node := fromAdminProxyNodeDTO(&req)
	if err := h.proxyService.UpdateNode(c.Request.Context(), node); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toAdminProxyNodeDTO(node))
}

// DeleteNode DELETE /api/v1/admin/proxy/nodes/:id
func (h *ProxyAdminHandler) DeleteNode(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	if err := h.proxyService.DeleteNode(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, nil)
}

// --- Product management ---

// ListProducts GET /api/v1/admin/proxy/products
func (h *ProxyAdminHandler) ListProducts(c *gin.Context) {
	products, err := h.proxyService.ListProducts(c.Request.Context(), false)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]adminProxyProductDTO, len(products))
	for i, p := range products {
		out[i] = toAdminProxyProductDTO(&p)
	}
	response.Success(c, out)
}

// CreateProduct POST /api/v1/admin/proxy/products
func (h *ProxyAdminHandler) CreateProduct(c *gin.Context) {
	var req adminProxyProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	p, err := fromAdminProxyProductDTO(&req)
	if err != nil {
		response.BadRequest(c, "invalid price: "+err.Error())
		return
	}
	if err := h.proxyService.CreateProduct(c.Request.Context(), p); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, toAdminProxyProductDTO(p))
}

// UpdateProduct PUT /api/v1/admin/proxy/products/:id
func (h *ProxyAdminHandler) UpdateProduct(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	var req adminProxyProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	req.ID = id
	p, err := fromAdminProxyProductDTO(&req)
	if err != nil {
		response.BadRequest(c, "invalid price: "+err.Error())
		return
	}
	if err := h.proxyService.UpdateProduct(c.Request.Context(), p); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toAdminProxyProductDTO(p))
}

// DeleteProduct DELETE /api/v1/admin/proxy/products/:id
func (h *ProxyAdminHandler) DeleteProduct(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	if err := h.proxyService.DeleteProduct(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, nil)
}

// --- Rental management ---

// ListRentals GET /api/v1/admin/proxy/rentals
func (h *ProxyAdminHandler) ListRentals(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	filter := service.ProxyRentalFilter{Status: c.Query("status")}
	if uid := c.Query("user_id"); uid != "" {
		if id, err := strconv.ParseInt(uid, 10, 64); err == nil {
			filter.UserID = &id
		}
	}
	rentals, result, err := h.proxyService.ListRentals(c.Request.Context(), filter,
		pagination.PaginationParams{Page: page, PageSize: pageSize})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]adminProxyRentalDTO, len(rentals))
	for i, r := range rentals {
		out[i] = toAdminProxyRentalDTO(&r, nil)
	}
	response.Paginated(c, out, result.Total, result.Page, result.PageSize)
}

// GetRental GET /api/v1/admin/proxy/rentals/:id
func (h *ProxyAdminHandler) GetRental(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	rental, cred, err := h.proxyService.AdminGetRental(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toAdminProxyRentalDTO(rental, cred))
}

// UpdateTraffic POST /api/v1/admin/proxy/rentals/:id/traffic
func (h *ProxyAdminHandler) UpdateTraffic(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}
	var req struct {
		DeltaGB float64 `json:"delta_gb" binding:"required"`
		Note    string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}
	deltaBytes := int64(req.DeltaGB * 1024 * 1024 * 1024)
	if err := h.proxyService.UpdateTraffic(c.Request.Context(), id, subject.UserID, deltaBytes, req.Note); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, nil)
}

// ForceExpire POST /api/v1/admin/proxy/rentals/:id/expire
func (h *ProxyAdminHandler) ForceExpire(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	if err := h.proxyService.ForceExpireRental(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, nil)
}

// GetTrafficLogs GET /api/v1/admin/proxy/rentals/:id/traffic
func (h *ProxyAdminHandler) GetTrafficLogs(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil {
		return
	}
	logs, err := h.proxyService.GetTrafficLogs(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, logs)
}

// --- DTOs ---

type adminProxyNodeDTO struct {
	ID           int64    `json:"id,omitempty"`
	IPAddress    string   `json:"ip_address" binding:"required"`
	Country      string   `json:"country"`
	CountryCode  string   `json:"country_code"`
	City         string   `json:"city"`
	ISP          string   `json:"isp"`
	HTTPPort     int      `json:"http_port"`
	VlessPort    int      `json:"vless_port"`
	VlessNetwork string   `json:"vless_network"`
	VlessTLS     bool     `json:"vless_tls"`
	VlessSNI     string   `json:"vless_sni"`
	VlessWSPath  string   `json:"vless_ws_path"`
	Tags         []string `json:"tags"`
	Status       string   `json:"status"`
	Description  string   `json:"description"`
}

func toAdminProxyNodeDTO(n *service.ProxyNode) adminProxyNodeDTO {
	return adminProxyNodeDTO{
		ID:           n.ID,
		IPAddress:    n.IPAddress,
		Country:      n.Country,
		CountryCode:  n.CountryCode,
		City:         n.City,
		ISP:          n.ISP,
		HTTPPort:     n.HTTPPort,
		VlessPort:    n.VlessPort,
		VlessNetwork: n.VlessNetwork,
		VlessTLS:     n.VlessTLS,
		VlessSNI:     n.VlessSNI,
		VlessWSPath:  n.VlessWSPath,
		Tags:         n.Tags,
		Status:       n.Status,
		Description:  n.Description,
	}
}

func fromAdminProxyNodeDTO(req *adminProxyNodeDTO) *service.ProxyNode {
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}
	return &service.ProxyNode{
		ID:           req.ID,
		IPAddress:    req.IPAddress,
		Country:      req.Country,
		CountryCode:  req.CountryCode,
		City:         req.City,
		ISP:          req.ISP,
		HTTPPort:     req.HTTPPort,
		VlessPort:    req.VlessPort,
		VlessNetwork: req.VlessNetwork,
		VlessTLS:     req.VlessTLS,
		VlessSNI:     req.VlessSNI,
		VlessWSPath:  req.VlessWSPath,
		Tags:         tags,
		Status:       req.Status,
		Description:  req.Description,
	}
}

type adminProxyProductDTO struct {
	ID             int64  `json:"id,omitempty"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	DurationDays   int    `json:"duration_days" binding:"required"`
	TrafficLimitGB int    `json:"traffic_limit_gb"`
	Price          string `json:"price" binding:"required"`
	SortOrder      int    `json:"sort_order"`
	IsActive       bool   `json:"is_active"`
}

func toAdminProxyProductDTO(p *service.ProxyProduct) adminProxyProductDTO {
	return adminProxyProductDTO{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		DurationDays:   p.DurationDays,
		TrafficLimitGB: p.TrafficLimitGB,
		Price:          p.Price.String(),
		SortOrder:      p.SortOrder,
		IsActive:       p.IsActive,
	}
}

func fromAdminProxyProductDTO(req *adminProxyProductDTO) (*service.ProxyProduct, error) {
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, err
	}
	return &service.ProxyProduct{
		ID:             req.ID,
		Name:           req.Name,
		Description:    req.Description,
		DurationDays:   req.DurationDays,
		TrafficLimitGB: req.TrafficLimitGB,
		Price:          price,
		SortOrder:      req.SortOrder,
		IsActive:       req.IsActive,
	}, nil
}

type adminProxyRentalDTO struct {
	ID                int64                    `json:"id"`
	UserID            int64                    `json:"user_id"`
	NodeID            int64                    `json:"node_id"`
	ProductID         int64                    `json:"product_id"`
	Status            string                   `json:"status"`
	StartedAt         interface{}              `json:"started_at"`
	ExpiresAt         interface{}              `json:"expires_at"`
	TrafficUsedBytes  int64                    `json:"traffic_used_bytes"`
	TrafficLimitBytes int64                    `json:"traffic_limit_bytes"`
	Node              *adminProxyNodeDTO       `json:"node,omitempty"`
	Product           *adminProxyProductDTO    `json:"product,omitempty"`
	Credential        *adminProxyCredentialDTO `json:"credential,omitempty"`
}

type adminProxyCredentialDTO struct {
	HTTPUsername string `json:"http_username"`
	HTTPPassword string `json:"http_password"`
	VlessUUID    string `json:"vless_uuid"`
	VlessLink    string `json:"vless_link"`
}

func toAdminProxyRentalDTO(r *service.ProxyRental, cred *service.ProxyCredential) adminProxyRentalDTO {
	dto := adminProxyRentalDTO{
		ID:                r.ID,
		UserID:            r.UserID,
		NodeID:            r.NodeID,
		ProductID:         r.ProductID,
		Status:            r.Status,
		StartedAt:         r.StartedAt,
		ExpiresAt:         r.ExpiresAt,
		TrafficUsedBytes:  r.TrafficUsedBytes,
		TrafficLimitBytes: r.TrafficLimitBytes,
	}
	if r.Node != nil {
		n := toAdminProxyNodeDTO(r.Node)
		dto.Node = &n
	}
	if r.Product != nil {
		p := toAdminProxyProductDTO(r.Product)
		dto.Product = &p
	}
	if cred != nil {
		dto.Credential = &adminProxyCredentialDTO{
			HTTPUsername: cred.HTTPUsername,
			HTTPPassword: cred.HTTPPassword,
			VlessUUID:    cred.VlessUUID,
			VlessLink:    cred.VlessLink,
		}
	}
	return dto
}

func parseAdminID(c *gin.Context, param string) (int64, error) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return 0, err
	}
	return id, nil
}
