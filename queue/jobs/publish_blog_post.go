package jobs

import "time"

type PublishBlogPostArgs struct {
	BlogPostID int32     `json:"blog_post_id"`
	PublishAt  time.Time `json:"publish_at"`
}

func (PublishBlogPostArgs) Kind() string { return "publish_blog_post" }
