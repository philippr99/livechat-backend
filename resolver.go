package livechat_gqlgen
//go:generate go run github.com/99designs/gqlgen
import (
	"context"
	"math/rand"
	"strconv"
	"symflower/livechat_gqlgen/auth"
	"symflower/livechat_gqlgen/management"
	"symflower/livechat_gqlgen/models"
	"sync"
	"time"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{
	mutex sync.Mutex
	observers map[string]chan *models.ChatMessage
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SendMessage(ctx context.Context, content string) (*models.ChatMessage, error) {
	decodedToken, err := auth.ValidateJWT(ctx.Value("token").(string))
	if err != nil {
		return nil, err
	}

	var user models.Account = decodedToken.(models.Account)
	chatMessage := models.ChatMessage {
		ID: int(management.MessageCount),
		Author: user.Username,
		Content: content,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}
	management.AddMessage(chatMessage)

	// send to subscribers
	r.mutex.Lock()
	for _, observer := range r.observers {
		observer <- &chatMessage
	}
	r.mutex.Unlock()

	return &chatMessage, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Account(ctx context.Context) (*models.Account, error) {
	decodedToken, err := auth.ValidateJWT(ctx.Value("token").(string))
	if err != nil {
		return nil, err
	}

	var account models.Account;
	account.Username = decodedToken.(models.Account).Username
	account.ID = decodedToken.(models.Account).ID

	return &account, nil
}
func (r *queryResolver) Messages(ctx context.Context) ([]*models.ChatMessage, error) {
	_, err := auth.ValidateJWT(ctx.Value("token").(string))
	if err != nil {
		return nil, err
	}

	var messages []*models.ChatMessage
	for index, _ := range management.GetMessages() {
		messages = append(messages, &management.GetMessages()[index])
	}

	return messages, nil
}

type subscriptionResolver struct{
	*Resolver
}

func (r *subscriptionResolver) MessageReceived(ctx context.Context) (<-chan *models.ChatMessage, error) {
	// No Token validation, should be accessible for everyone atm

	event := make(chan *models.ChatMessage, 1)
	id := strconv.FormatUint(rand.Uint64(),10)

	r.mutex.Lock()
	if r.observers == nil {
		r.observers = make(map[string]chan *models.ChatMessage)
	}
	r.observers[id] = event
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.observers, id)
		r.mutex.Unlock()
	}()
	return event, nil
}
