package v1

import (
	"encoding/json"
	"fmt"

	"github.com/aliakbariaa1996/mk-test-one/internal/common/errorx"
	"github.com/aliakbariaa1996/mk-test-one/internal/common/tracing"
	"github.com/aliakbariaa1996/mk-test-one/internal/services/delivery"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) makeGetDeliveryHandler(
	deliveryService delivery.UseService,
) func(_ echo.Context) error {
	return func(c echo.Context) error {
		var err error
		sp, _ := tracing.CreateSpan(c.Request().Context(), fmt.Sprintf("Push Order To Delivery"))
		defer sp.Finish()
		defer func() {
			if err != nil {
				tracing.LogSpanError(sp, "", err)
			}
		}()
		sou := delivery.SourceLocation{
			Lat: 55.545454,
			Lng: 12.5465465,
		}
		var listLoc []delivery.DeliverManLocation

		err = json.Unmarshal([]byte(h.cfg.DeliverManLoc), &listLoc)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errorx.Success{Code: errorx.CodeError(errorx.ErrCalculate), Message: errorx.ErrCalculate.Error()})
		}
		res := deliveryService.GetDistance(sou, listLoc)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errorx.Success{Code: errorx.CodeError(errorx.ErrCalculate), Message: errorx.ErrCalculate.Error()})
		}
		return c.JSON(http.StatusOK, errorx.Success{Code: errorx.CodeError(err), Message: "Success Message", Details: res})
	}
}
