package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maksroxx/LFact/db"
	"github.com/maksroxx/LFact/types"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	exists, err := h.store.UserStore.CheckUserExists(c.Context(), params.Email)
	if err != nil {
		return err
	}
	if exists {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"error": "user already exists",
			})
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.store.UserStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.UserStore.GetUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"error": "users resource not found",
			})
	}
	return c.JSON(users)
}
