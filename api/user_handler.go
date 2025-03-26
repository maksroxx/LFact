package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/maksroxx/LFact/db"
	"github.com/maksroxx/LFact/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if params.Balance != 0 {
		balance, err := h.store.UserStore.GetUserById(c.Context(), userID)
		if err != nil {
			return err
		}
		params.Balance += balance.Balance
	}
	filter := db.Map{"_id": oid}
	user, err := h.store.UserStore.UpdateUser(c.Context(), filter, db.Map(params.ToBson()))
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUserById(c *fiber.Ctx) error {
	userId := c.Params("id")
	user, err := h.store.UserStore.GetUserById(c.Context(), userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.store.UserStore.DeleteUser(c.Context(), userId); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userId})
}
