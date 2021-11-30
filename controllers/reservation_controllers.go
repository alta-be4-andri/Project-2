package controllers

import (
	"net/http"
	"project2/lib/databases"
	"project2/middlewares"
	"project2/models"
	"project2/response"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var format_date string = "2006-01-02"

// Controller untuk membuat reservasi baru
func CreateReservationControllers(c echo.Context) error {
	body := models.ReservationBody{}
	c.Bind(&body)
	logged := middlewares.ExtractTokenUserId(c)

	input := models.Reservation{}
	input.Check_In, _ = time.Parse(format_date, body.Check_In)
	input.Check_Out, _ = time.Parse(format_date, body.Check_Out)
	input.UsersID = uint(logged)
	input.RoomsID = body.RoomsID
	input.KartuKreditID = body.KartuKreditID

	reservation, err := databases.CreateReservation(&input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse())
	}

	databases.AddJumlahMalam(input.Check_In, input.Check_Out, reservation.ID)
	databases.AddHargaToReservation(input.RoomsID, reservation.ID)
	return c.JSON(http.StatusOK, response.ReservationSuccessResponse(reservation.ID))
}

func GetReservationControllers(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.FalseParamResponse())
	}
	userId, _ := databases.GetReservationOwner(id)
	if err != nil || userId == 0 {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse())
	}
	logged := middlewares.ExtractTokenUserId(c)
	if uint(logged) != userId {
		return c.JSON(http.StatusBadRequest, response.AccessForbiddenResponse())
	}
	getReservation, _ := databases.GetReservation(id)
	reservation := models.GetReservation{}
	reservation.RoomsID = getReservation.RoomsID
	reservation.Check_In = getReservation.Check_In.Format(format_date)
	reservation.Check_Out = getReservation.Check_Out.Format(format_date)
	reservation.KartuKreditID = getReservation.KartuKreditID
	reservation.Total_Harga = getReservation.Total_Harga

	return c.JSON(http.StatusOK, response.SuccessResponseData(reservation))
}

func CancelReservationController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.FalseParamResponse())
	}
	userId, _ := databases.GetReservationOwner(id)
	if err != nil || userId == 0 {
		return c.JSON(http.StatusBadRequest, response.BadRequestResponse())
	}
	logged := middlewares.ExtractTokenUserId(c)
	if uint(logged) != userId {
		return c.JSON(http.StatusBadRequest, response.AccessForbiddenResponse())
	}
	databases.CancelReservation(id)
	return c.JSON(http.StatusOK, response.SuccessResponseNonData())
}