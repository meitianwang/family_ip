package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/meitianwang/fast-frame/internal/domain"
	"github.com/meitianwang/fast-frame/internal/pkg/pagination"
	"github.com/meitianwang/fast-frame/internal/pkg/response"
	middleware2 "github.com/meitianwang/fast-frame/internal/server/middleware"
	"github.com/meitianwang/fast-frame/internal/service"

	"github.com/gin-gonic/gin"
)

// ProxyHandler handles user-facing proxy endpoints.
type ProxyHandler struct {
	proxyService   *service.ProxyService
	paymentService *service.PaymentOrderService
}

// NewProxyHandler creates a new ProxyHandler.
func NewProxyHandler(proxyService *service.ProxyService, paymentService *service.PaymentOrderService) *ProxyHandler {
	return &ProxyHandler{proxyService: proxyService, paymentService: paymentService}
}

// ListNodes GET /api/v1/proxy/nodes — list available nodes for purchase.
func (h *ProxyHandler) ListNodes(c *gin.Context) {
	filter := service.ProxyNodeFilter{
		CountryCode: c.Query("country_code"),
		Tag:         c.Query("tag"),
	}
	nodes, err := h.proxyService.ListAvailableNodes(c.Request.Context(), filter)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]proxyNodeDTO, len(nodes))
	for i, n := range nodes {
		out[i] = toProxyNodeDTO(&n)
	}
	response.Success(c, out)
}

// GetNode GET /api/v1/proxy/nodes/:id
func (h *ProxyHandler) GetNode(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	node, err := h.proxyService.GetNode(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if node.Status != domain.ProxyNodeStatusAvailable {
		response.NotFound(c, "node not available")
		return
	}
	response.Success(c, toProxyNodeDTO(node))
}

// ListProducts GET /api/v1/proxy/products
func (h *ProxyHandler) ListProducts(c *gin.Context) {
	products, err := h.proxyService.ListProducts(c.Request.Context(), true)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]proxyProductDTO, len(products))
	for i, p := range products {
		out[i] = toProxyProductDTO(&p)
	}
	response.Success(c, out)
}

// CreateRental POST /api/v1/proxy/rentals — initiate a rental + payment order.
func (h *ProxyHandler) CreateRental(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}

	var req struct {
		NodeID    int64  `json:"node_id" binding:"required"`
		ProductID int64  `json:"product_id" binding:"required"`
		PayType   string `json:"pay_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	ctx := c.Request.Context()

	rental, err := h.proxyService.CreateRental(ctx, service.CreateRentalRequest{
		UserID:    subject.UserID,
		NodeID:    req.NodeID,
		ProductID: req.ProductID,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Create the payment order
	price, err := h.proxyService.PriceForProduct(ctx, req.ProductID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	orderReq := service.CreateOrderRequest{
		UserID:      subject.UserID,
		Amount:      price,
		PaymentType: req.PayType,
		OrderType:   domain.PaymentOrderTypeProxyRental,
		ClientIP:    c.ClientIP(),
		SrcHost:     c.Request.Host,
		SrcURL:      c.Request.RequestURI,
	}
	result, err := h.paymentService.CreateOrder(ctx, orderReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	if err := h.proxyService.SetRentalPaymentOrder(ctx, rental.ID, result.Order.ID); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"rental_id":  rental.ID,
			"order_id":   result.Order.ID,
			"pay_url":    result.PayURL,
			"qr_code":    result.QrCode,
			"amount":     result.Order.Amount,
			"expires_at": result.Order.ExpiresAt,
		},
	})
}

// ListRentals GET /api/v1/proxy/rentals
func (h *ProxyHandler) ListRentals(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}
	page, pageSize := response.ParsePagination(c)
	filter := service.ProxyRentalFilter{Status: c.Query("status"), UserID: &subject.UserID}
	rentals, result, err := h.proxyService.ListUserRentals(c.Request.Context(), subject.UserID, filter,
		pagination.PaginationParams{Page: page, PageSize: pageSize})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]proxyRentalDTO, len(rentals))
	for i, r := range rentals {
		out[i] = toProxyRentalDTO(&r, nil)
	}
	response.Paginated(c, out, result.Total, result.Page, result.PageSize)
}

// GetRental GET /api/v1/proxy/rentals/:id
func (h *ProxyHandler) GetRental(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	rental, cred, err := h.proxyService.GetRentalWithCredential(c.Request.Context(), id, subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toProxyRentalDTO(rental, cred))
}

// CancelRental POST /api/v1/proxy/rentals/:id/cancel
func (h *ProxyHandler) CancelRental(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "not authenticated")
		return
	}
	id, err := parseID(c, "id")
	if err != nil {
		return
	}
	if err := h.proxyService.CancelRental(c.Request.Context(), id, subject.UserID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, nil)
}

// --- DTOs ---

type proxyNodeDTO struct {
	ID           int64    `json:"id"`
	IPAddress    string   `json:"ip_address"`
	Country      string   `json:"country"`
	CountryCode  string   `json:"country_code"`
	City         string   `json:"city"`
	ISP          string   `json:"isp"`
	HTTPPort     int      `json:"http_port"`
	VlessPort    int      `json:"vless_port"`
	Tags         []string `json:"tags"`
	Status       string   `json:"status"`
	Description  string   `json:"description"`
}

func toProxyNodeDTO(n *service.ProxyNode) proxyNodeDTO {
	return proxyNodeDTO{
		ID:          n.ID,
		IPAddress:   n.IPAddress,
		Country:     n.Country,
		CountryCode: n.CountryCode,
		City:        n.City,
		ISP:         n.ISP,
		HTTPPort:    n.HTTPPort,
		VlessPort:   n.VlessPort,
		Tags:        n.Tags,
		Status:      n.Status,
		Description: n.Description,
	}
}

type proxyProductDTO struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DurationDays   int    `json:"duration_days"`
	TrafficLimitGB int    `json:"traffic_limit_gb"`
	Price          string `json:"price"`
	SortOrder      int    `json:"sort_order"`
}

func toProxyProductDTO(p *service.ProxyProduct) proxyProductDTO {
	return proxyProductDTO{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		DurationDays:   p.DurationDays,
		TrafficLimitGB: p.TrafficLimitGB,
		Price:          p.Price.String(),
		SortOrder:      p.SortOrder,
	}
}

type proxyCredentialDTO struct {
	HTTPHost     string `json:"http_host"`
	HTTPPort     int    `json:"http_port"`
	HTTPUsername string `json:"http_username"`
	HTTPPassword string `json:"http_password"`
	VlessLink    string `json:"vless_link"`
}

type proxyRentalDTO struct {
	ID                int64               `json:"id"`
	NodeID            int64               `json:"node_id"`
	ProductID         int64               `json:"product_id"`
	Status            string              `json:"status"`
	StartedAt         *time.Time          `json:"started_at,omitempty"`
	ExpiresAt         *time.Time          `json:"expires_at,omitempty"`
	TrafficUsedBytes  int64               `json:"traffic_used_bytes"`
	TrafficLimitBytes int64               `json:"traffic_limit_bytes"`
	CreatedAt         time.Time           `json:"created_at"`
	Node              *proxyNodeDTO       `json:"node,omitempty"`
	Product           *proxyProductDTO    `json:"product,omitempty"`
	Credential        *proxyCredentialDTO `json:"credential,omitempty"`
}

func toProxyRentalDTO(r *service.ProxyRental, cred *service.ProxyCredential) proxyRentalDTO {
	dto := proxyRentalDTO{
		ID:                r.ID,
		NodeID:            r.NodeID,
		ProductID:         r.ProductID,
		Status:            r.Status,
		StartedAt:         r.StartedAt,
		ExpiresAt:         r.ExpiresAt,
		TrafficUsedBytes:  r.TrafficUsedBytes,
		TrafficLimitBytes: r.TrafficLimitBytes,
		CreatedAt:         r.CreatedAt,
	}
	if r.Node != nil {
		n := toProxyNodeDTO(r.Node)
		dto.Node = &n
	}
	if r.Product != nil {
		p := toProxyProductDTO(r.Product)
		dto.Product = &p
	}
	if cred != nil {
		httpPort := 3128
		if r.Node != nil {
			httpPort = r.Node.HTTPPort
		}
		httpHost := ""
		if r.Node != nil {
			httpHost = r.Node.IPAddress
		}
		dto.Credential = &proxyCredentialDTO{
			HTTPHost:     httpHost,
			HTTPPort:     httpPort,
			HTTPUsername: cred.HTTPUsername,
			HTTPPassword: cred.HTTPPassword,
			VlessLink:    cred.VlessLink,
		}
	}
	return dto
}

// parseID is a small helper to parse int64 path params.
func parseID(c *gin.Context, param string) (int64, error) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return 0, err
	}
	return id, nil
}
