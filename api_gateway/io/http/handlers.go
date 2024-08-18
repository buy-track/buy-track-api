package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"my-stocks/api-gateway/app"
	"my-stocks/domains"
	"net/http"
	"strconv"
)

type AuthController struct {
	authService *app.AuthService
	userService *app.UserService
}

func NewAuthController(authService *app.AuthService, userService *app.UserService) *AuthController {
	return &AuthController{authService: authService, userService: userService}
}

func (a AuthController) Register(ctx echo.Context) error {
	dto := new(UserRegisterDto)
	if err := ctx.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if exists := a.userService.EmailExists(dto.Email); exists {
		return echo.NewHTTPError(http.StatusBadRequest, "email exists.")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	user, err := a.userService.Create(&domains.User{
		Password: string(hashed),
		Email:    dto.Email,
		Name:     dto.Name,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error during create user")
	}
	token, err := a.authService.Login(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user created. error during login in. please login again")
	}
	return ctx.JSON(http.StatusCreated, &UserRegisterResponseDto{
		Token: token,
		User:  user,
	})

}

func (a AuthController) Login(ctx echo.Context) error {
	dto := new(UserLoginDto)
	if err := ctx.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := a.userService.CheckPassword(dto.Email, dto.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	token, err := a.authService.Login(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error during login user")
	}
	return ctx.JSON(http.StatusCreated, &UserLoginResponseDto{
		Token: token,
		User:  user,
	})
}

func (a AuthController) Logout(ctx echo.Context) error {
	user := ctx.Get("auth_user").(*domains.User)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	token := ctx.Get("auth_token").(string)
	err := a.authService.Logout(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error during logout")
	}

	return echo.NewHTTPError(http.StatusNoContent)
}

type UserController struct {
	userService *app.UserService
}

func NewUserController(userService *app.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u UserController) Profile(ctx echo.Context) error {
	user := ctx.Get("auth_user").(*domains.User)
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}
	return ctx.JSON(http.StatusOK, user)
}

type CoinController struct {
	coinService *app.CoinService
}

func NewCoinController(coinService *app.CoinService) *CoinController {
	return &CoinController{coinService: coinService}
}

func (c CoinController) Paginate(ctx echo.Context) error {
	fmt.Println(c.coinService)
	var page int64 = 1
	var perPage int64 = 15
	if tmp, err := strconv.ParseInt(ctx.QueryParam("page"), 10, 64); err == nil {
		page = tmp
	}
	if tmp, err := strconv.ParseInt(ctx.QueryParam("per_page"), 10, 64); err == nil {
		perPage = tmp
	}
	list := c.coinService.PaginateList(perPage, page)
	return ctx.JSON(http.StatusOK, list)
}

func (c CoinController) Show(ctx echo.Context) error {
	var symbol string
	if symbol = ctx.Param("symbol"); symbol == "" {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	coin, err := c.coinService.GetBySymbol(symbol)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return ctx.JSON(http.StatusOK, coin)
}
