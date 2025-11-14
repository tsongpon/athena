package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/tsongpon/athena/internal/logger"
	"github.com/tsongpon/athena/internal/model"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

const usersCollection = "users"

// UserFirestoreRepository implements UserRepository interface using GCP Firestore
type UserFirestoreRepository struct {
	client *firestore.Client
	ctx    context.Context
}

// NewUserFirestoreRepository creates a new instance of UserFirestoreRepository
func NewUserFirestoreRepository(ctx context.Context, client *firestore.Client) *UserFirestoreRepository {
	return &UserFirestoreRepository{
		client: client,
		ctx:    ctx,
	}
}

// firestoreUser is the structure used to store/retrieve users in Firestore
type firestoreUser struct {
	ID        string    `firestore:"id"`
	Name      string    `firestore:"name"`
	Email     string    `firestore:"email"`
	Password  string    `firestore:"password"`
	Tier      string    `firestore:"tier"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

// toFirestoreUser converts model.User to firestoreUser
func toFirestoreUser(user model.User) firestoreUser {
	return firestoreUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Tier:      user.Tier,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// toModelUser converts firestoreUser to model.User
func toModelUser(fsUser firestoreUser) model.User {
	return model.User{
		ID:        fsUser.ID,
		Name:      fsUser.Name,
		Email:     fsUser.Email,
		Password:  fsUser.Password,
		Tier:      fsUser.Tier,
		CreatedAt: fsUser.CreatedAt,
		UpdatedAt: fsUser.UpdatedAt,
	}
}

// CreateUser creates a new user in Firestore
func (r *UserFirestoreRepository) CreateUser(user model.User) (model.User, error) {
	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set creation and update times
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	// Convert to Firestore structure
	fsUser := toFirestoreUser(user)

	// Store in Firestore using user ID as document ID
	_, err := r.client.Collection(usersCollection).Doc(user.ID).Set(r.ctx, fsUser)
	if err != nil {
		logger.Error("Failed to create user in Firestore",
			zap.String("user_id", user.ID),
			zap.String("email", user.Email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Debug("Created user in Firestore", zap.String("id", user.ID))
	return user, nil
}

// GetUserByID retrieves a user by their ID from Firestore
func (r *UserFirestoreRepository) GetUserByID(id string) (model.User, error) {
	logger.Debug("Getting user from Firestore", zap.String("id", id))

	docSnap, err := r.client.Collection(usersCollection).Doc(id).Get(r.ctx)
	if err != nil {
		logger.Error("Failed to get user from Firestore",
			zap.String("id", id),
			zap.Error(err))
		return model.User{}, fmt.Errorf("user with ID %s not found: %w", id, err)
	}

	var fsUser firestoreUser
	if err := docSnap.DataTo(&fsUser); err != nil {
		logger.Error("Failed to parse user data from Firestore",
			zap.String("id", id),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to parse user data: %w", err)
	}

	return toModelUser(fsUser), nil
}

// GetUserByEmail retrieves a user by their email address from Firestore
func (r *UserFirestoreRepository) GetUserByEmail(email string) (model.User, error) {
	logger.Debug("Getting user by email from Firestore", zap.String("email", email))

	// Query Firestore for user with matching email
	iter := r.client.Collection(usersCollection).
		Where("email", "==", email).
		Limit(1).
		Documents(r.ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		return model.User{}, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		logger.Error("Failed to get user by email from Firestore",
			zap.String("email", email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	var fsUser firestoreUser
	if err := doc.DataTo(&fsUser); err != nil {
		logger.Error("Failed to parse user data from Firestore",
			zap.String("email", email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to parse user data: %w", err)
	}

	return toModelUser(fsUser), nil
}

// GetUserByEmailAndPassword retrieves a user by email and password from Firestore
func (r *UserFirestoreRepository) GetUserByEmailAndPassword(email, hashedPassword string) (model.User, error) {
	logger.Debug("Getting user by email and password from Firestore", zap.String("email", email))

	// Query Firestore for user with matching email and password
	iter := r.client.Collection(usersCollection).
		Where("email", "==", email).
		Where("password", "==", hashedPassword).
		Limit(1).
		Documents(r.ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		return model.User{}, fmt.Errorf("user not found with provided credentials")
	}
	if err != nil {
		logger.Error("Failed to get user by credentials from Firestore",
			zap.String("email", email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	var fsUser firestoreUser
	if err := doc.DataTo(&fsUser); err != nil {
		logger.Error("Failed to parse user data from Firestore",
			zap.String("email", email),
			zap.Error(err))
		return model.User{}, fmt.Errorf("failed to parse user data: %w", err)
	}

	return toModelUser(fsUser), nil
}
