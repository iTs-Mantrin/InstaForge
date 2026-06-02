package controllers

import (
	"strconv"

	"instaforge/internal/dto"
	"instaforge/internal/metrics"
	"instaforge/internal/models"
	"instaforge/internal/services"
	"instaforge/internal/validators"
	pkgresponse "instaforge/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// InstagramController handles Instagram content endpoints.
type InstagramController struct {
	instaService *services.InstagramService
	userService  *services.UserService
}

// NewInstagramController creates a new Instagram controller.
func NewInstagramController(instaService *services.InstagramService, userService *services.UserService) *InstagramController {
	return &InstagramController{
		instaService: instaService,
		userService:  userService,
	}
}

// PreviewURL handles POST /api/preview (process an Instagram URL).
func (ctrl *InstagramController) PreviewURL(c *fiber.Ctx) error {
	var req dto.PreviewRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.BadRequest(c, "invalid request body")
	}

	urlType, shortcode, err := validators.ValidateInstagramURL(req.URL)
	if err != nil {
		return pkgresponse.ValidationError(c, err.Error())
	}

	var preview dto.PreviewResponse

	switch urlType {
	case validators.URLTypePost:
		post, err := ctrl.instaService.GetPostByURL(shortcode)
		if err != nil {
			metrics.RecordError("instagram", "preview")
			return pkgresponse.InternalError(c, "failed to fetch post details, please try again later")
		}
		preview = dto.PreviewResponse{
			Type:       "post",
			Username:   post.Username,
			Caption:    post.Caption,
			MediaCount: len(post.MediaItems),
			Thumbnail:  post.Thumbnail,
			Likes:      post.Likes,
			Comments:   post.Comments,
			Timestamp:  post.Timestamp,
		}
		for _, item := range post.MediaItems {
			preview.MediaItems = append(preview.MediaItems, dto.MediaItemDTO{
				ID:          item.ID,
				Type:        string(item.Type),
				URL:         item.URL,
				Thumbnail:   item.Thumbnail,
				Width:       item.Width,
				Height:      item.Height,
				FileSize:    item.FileSize,
				ContentType: item.ContentType,
			})
		}

	case validators.URLTypeReel:
		post, err := ctrl.instaService.GetReelByURL(shortcode)
		if err != nil {
			metrics.RecordError("instagram", "preview")
			return pkgresponse.InternalError(c, "failed to fetch reel details, please try again later")
		}
		preview = dto.PreviewResponse{
			Type:       "reel",
			Username:   post.Username,
			Caption:    post.Caption,
			MediaCount: len(post.MediaItems),
			Thumbnail:  post.Thumbnail,
			Likes:      post.Likes,
			Comments:   post.Comments,
			Timestamp:  post.Timestamp,
		}
		for _, item := range post.MediaItems {
			mediaDTO := dto.MediaItemDTO{
				ID:          item.ID,
				Type:        string(item.Type),
				URL:         item.URL,
				Thumbnail:   item.Thumbnail,
				Width:       item.Width,
				Height:      item.Height,
				FileSize:    item.FileSize,
				ContentType: item.ContentType,
			}
			preview.MediaItems = append(preview.MediaItems, mediaDTO)
		}

	case validators.URLTypeStory:
		// For stories we need username and story ID
		stories, err := ctrl.instaService.GetStories(shortcode)
		if err != nil {
			metrics.RecordError("instagram", "preview")
			return pkgresponse.InternalError(c, "failed to fetch story details, please try again later")
		}
		if len(stories) == 0 {
			return pkgresponse.NotFound(c, "no stories found")
		}
		story := stories[0]
		preview = dto.PreviewResponse{
			Type:       "story",
			Username:   story.Username,
			MediaCount: len(stories),
			Timestamp:  story.Timestamp,
		}
		for _, s := range stories {
			preview.MediaItems = append(preview.MediaItems, dto.MediaItemDTO{
				ID:        s.ID,
				Type:      string(s.MediaType),
				URL:       s.URL,
				Thumbnail: s.Thumbnail,
				Duration:  s.Duration,
			})
		}

	default:
		return pkgresponse.BadRequest(c, "unsupported URL type")
	}

	return pkgresponse.Success(c, preview)
}

// SearchUser handles GET /api/user/:username
func (ctrl *InstagramController) SearchUser(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return pkgresponse.ValidationError(c, "username is required")
	}

	if err := validators.ValidateInstagramUsername(username); err != nil {
		return pkgresponse.ValidationError(c, err.Error())
	}

	result, err := ctrl.userService.SearchUsername(username)
	if err != nil {
		metrics.RecordSearch(false)
		return pkgresponse.InternalError(c, "user search failed, please try again later")
	}

	metrics.RecordSearch(true)

	return pkgresponse.Success(c, result)
}

// GetUserFeed handles GET /api/user/:username/feed
func (ctrl *InstagramController) GetUserFeed(c *fiber.Ctx) error {
	username := c.Params("username")
	cursor := c.Query("cursor", "")
	limitStr := c.Query("limit", "12")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 12
	}

	feed, err := ctrl.instaService.GetUserFeed(username, cursor, limit)
	if err != nil {
		return pkgresponse.InternalError(c, "failed to fetch user feed, please try again later")
	}

	hasMore := feed.HasMore
	if !hasMore && feed.Cursor != "" {
		hasMore = true
	}

	return pkgresponse.SuccessWithMeta(c, feed.Items, feed.Cursor, hasMore, int64(len(feed.Items)))
}

// GetUserStories handles GET /api/user/:username/stories
func (ctrl *InstagramController) GetUserStories(c *fiber.Ctx) error {
	username := c.Params("username")

	stories, err := ctrl.instaService.GetStories(username)
	if err != nil {
		return pkgresponse.InternalError(c, "failed to fetch stories, please try again later")
	}

	return pkgresponse.Success(c, stories)
}

// ============================================================
// Frontend-compatible handlers (camelCase JSON, /instagram/ prefix)
// ============================================================

// postToFE converts a backend InstagramPost to frontend PostPreviewResponse.
func postToFE(post *models.InstagramPost) dto.PostPreviewResponse {
	items := make([]dto.MediaItemFEDTO, len(post.MediaItems))
	displayURL := post.Thumbnail
	for i, item := range post.MediaItems {
		if i == 0 && displayURL == "" {
			displayURL = item.URL
		}
		items[i] = dto.MediaItemFEDTO{
			ID:        item.ID,
			Type:      string(item.Type),
			URL:       item.URL,
			Thumbnail: item.Thumbnail,
			Width:     item.Width,
			Height:    item.Height,
		}
	}

	return dto.PostPreviewResponse{
		ID:             post.ID,
		Shortcode:      post.Shortcode,
		MediaType:      string(post.MediaType),
		MediaItems:     items,
		Caption:        post.Caption,
		LikesCount:     post.Likes,
		CommentsCount:  post.Comments,
		Timestamp:      post.Timestamp.UnixMilli(),
		OwnerUsername:  post.Username,
		DisplayUrl:     displayURL,
	}
}

// UserInfo handles GET /instagram/user/:username/info (frontend-facing).
func (ctrl *InstagramController) UserInfo(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return pkgresponse.ValidationError(c, "username is required")
	}

	profile, err := ctrl.instaService.GetUserProfile(username)
	if err != nil {
		return pkgresponse.InternalError(c, "failed to fetch user profile, please try again later")
	}

	return pkgresponse.Success(c, dto.UserInfoResponse{
		ID:            profile.ID,
		Username:      profile.Username,
		FullName:      profile.FullName,
		ProfilePicUrl: profile.ProfilePicURL,
		FollowerCount: profile.FollowerCount,
		FollowingCount: profile.FollowingCount,
		PostCount:     profile.PostCount,
		IsVerified:    profile.IsVerified,
		Biography:     profile.Biography,
		ExternalUrl:   profile.ExternalURL,
		IsPrivate:     profile.IsPrivate,
	})
}

// UserFeedFe handles GET /instagram/user/:username/feed (frontend-facing).
func (ctrl *InstagramController) UserFeedFe(c *fiber.Ctx) error {
	username := c.Params("username")
	cursor := c.Query("cursor", "")
	limitStr := c.Query("limit", "12")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 12
	}

	feed, err := ctrl.instaService.GetUserFeed(username, cursor, limit)
	if err != nil {
		return pkgresponse.InternalError(c, "failed to fetch user feed, please try again later")
	}

	hasMore := feed.HasMore
	if !hasMore && feed.Cursor != "" {
		hasMore = true
	}

	items := make([]dto.PostPreviewResponse, len(feed.Items))
	for i, post := range feed.Items {
		items[i] = postToFE(&post)
	}

	return pkgresponse.Success(c, dto.FeedResponse{
		Items:   items,
		HasMore: hasMore,
		Cursor:  feed.Cursor,
	})
}

// PreviewFe handles POST /instagram/preview (frontend-facing).
func (ctrl *InstagramController) PreviewFe(c *fiber.Ctx) error {
	var req dto.PreviewRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.BadRequest(c, "invalid request body")
	}

	urlType, shortcode, err := validators.ValidateInstagramURL(req.URL)
	if err != nil {
		return pkgresponse.ValidationError(c, err.Error())
	}

	var post *models.InstagramPost

	switch urlType {
	case validators.URLTypePost:
		post, err = ctrl.instaService.GetPostByURL(shortcode)
	case validators.URLTypeReel:
		post, err = ctrl.instaService.GetReelByURL(shortcode)
	case validators.URLTypeStory:
		// For stories, fetch the user's stories and use first as preview
		stories, sErr := ctrl.instaService.GetStories(shortcode)
		if sErr != nil {
			metrics.RecordError("instagram", "preview")
			return pkgresponse.InternalError(c, "failed to fetch story details, please try again later")
		}
		if len(stories) == 0 {
			return pkgresponse.NotFound(c, "no stories found")
		}
		story := stories[0]
		feResp := dto.PostPreviewResponse{
			ID:            story.ID,
			Shortcode:     shortcode,
			MediaType:     string(story.MediaType),
			MediaItems: []dto.MediaItemFEDTO{{
				ID:        story.ID,
				Type:      string(story.MediaType),
				URL:       story.URL,
				Thumbnail: story.Thumbnail,
			}},
			Timestamp:     story.Timestamp.UnixMilli(),
			OwnerUsername: story.Username,
		}
		return pkgresponse.Success(c, feResp)
	default:
		return pkgresponse.BadRequest(c, "unsupported URL type")
	}

	if err != nil {
		metrics.RecordError("instagram", "preview")
		return pkgresponse.InternalError(c, "failed to fetch post details, please try again later")
	}

	return pkgresponse.Success(c, postToFE(post))
}

// StoriesFe handles GET /instagram/user/:username/stories (frontend-facing).
func (ctrl *InstagramController) StoriesFe(c *fiber.Ctx) error {
	username := c.Params("username")

	stories, err := ctrl.instaService.GetStories(username)
	if err != nil {
		return pkgresponse.InternalError(c, "failed to fetch stories, please try again later")
	}

	// Get user profile for the owner profile pic
	ownerPic := ""
	profile, profileErr := ctrl.instaService.GetUserProfile(username)
	if profileErr == nil {
		ownerPic = profile.ProfilePicURL
	}

	type StoryFE struct {
		ID             string `json:"id"`
		Type           string `json:"type"`
		URL            string `json:"url"`
		Thumbnail      string `json:"thumbnail,omitempty"`
		Timestamp      int64  `json:"timestamp"`
		ExpiresAt      int64  `json:"expiresAt"`
		OwnerUsername  string `json:"ownerUsername"`
		OwnerProfilePic string `json:"ownerProfilePic,omitempty"`
	}

	items := make([]StoryFE, len(stories))
	for i, s := range stories {
		expiresAt := int64(0)
		if !s.ExpiresAt.IsZero() {
			expiresAt = s.ExpiresAt.UnixMilli()
		}
		items[i] = StoryFE{
			ID:             s.ID,
			Type:           string(s.MediaType),
			URL:            s.URL,
			Thumbnail:      s.Thumbnail,
			Timestamp:      s.Timestamp.UnixMilli(),
			ExpiresAt:      expiresAt,
			OwnerUsername:  s.Username,
			OwnerProfilePic: ownerPic,
		}
	}

	return pkgresponse.Success(c, fiber.Map{
		"stories": items,
		"owner": fiber.Map{
			"username":       username,
			"profilePicUrl":  ownerPic,
		},
	})
}
