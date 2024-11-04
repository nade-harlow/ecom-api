package product

import (
	"github.com/gin-gonic/gin"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/dto"
	"github.com/nade-harlow/ecom-api/internal/adapter/api/http/response"
	"github.com/nade-harlow/ecom-api/internal/app/domain/services/product"
	"github.com/nade-harlow/ecom-api/internal/app/utils/apperrors"
	"github.com/nade-harlow/ecom-api/internal/app/utils/helper"
)

type ProductHandler struct {
	productService product.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: product.NewProductService(),
	}
}

func (p *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(helper.ValidatorFormatErrors(err).Error()))
		return
	}

	if err := helper.ValidateRequestBody(req); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(err.Error()))
		return
	}

	product, err := p.productService.CreateProduct(req)
	if err != nil {
		response.JsonError(ctx, apperrors.InternalServerError("something went wrong"))
		return
	}

	response.JsonCreated(ctx, product, "product")
}

func (p *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var req dto.UpdateProductRequest
	paramID := ctx.Param("productID")

	if paramID == "" {
		response.JsonError(ctx, apperrors.BadRequestError("invalid product id"))
		return
	}

	productID, errS := helper.StringToUUID(paramID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	if err := ctx.ShouldBind(&req); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(helper.ValidatorFormatErrors(err).Error()))
		return
	}

	if err := helper.ValidateRequestBody(req); err != nil {
		response.JsonError(ctx, apperrors.BadRequestError(err.Error()))
		return
	}

	product, err := p.productService.UpdateProduct(productID, &req)
	if err != nil {
		response.JsonError(ctx, apperrors.InternalServerError("something went wrong"))
		return
	}

	response.JsonModified(ctx, product, "product")
}

func (p *ProductHandler) DeleteProduct(ctx *gin.Context) {
	paramID := ctx.Param("productID")
	if paramID == "" {
		response.JsonError(ctx, apperrors.BadRequestError("invalid product id"))
		return
	}

	productID, errS := helper.StringToUUID(paramID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	if err := p.productService.DeleteProduct(productID); err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonDelete(ctx, nil, "product")
}

func (p *ProductHandler) GetProductByID(ctx *gin.Context) {
	paramID := ctx.Param("productID")
	if paramID == "" {
		response.JsonError(ctx, apperrors.BadRequestError("invalid product id"))
		return
	}

	productID, errS := helper.StringToUUID(paramID)
	if errS != nil {
		response.JsonError(ctx, apperrors.BadRequestError("invalid id format"))
		return
	}

	product, err := p.productService.GetProductByID(productID)
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	response.JsonOk(ctx, product)
}

func (p *ProductHandler) GetProducts(ctx *gin.Context) {
	products, err := p.productService.GetProducts()
	if err != nil {
		response.JsonError(ctx, err)
	}

	response.JsonOk(ctx, products)
}
