package main

import (
	"golang.org/x/exp/errors/fmt"
	"time"
)

type Base struct {
	ID   int
	Name string
}

type Blogger struct {
	Base
	Weibos   []*PostContent
	Comments map[int][]*Comment
	Fans     []FansInterface
}
type PostContent struct {
	Id       int
	Content  string
	PostTime time.Time
	Type     int
	PostMan  string
}

type Comment struct {
	PostContent
	ToName      string
	CommentTime time.Time
}

type BloggerInterface interface {
	Attach(bFans FansInterface)
	Detach(bFans FansInterface)
	Notify(wbid int)
}

func NewBlogger(name string) *Blogger {
	blg := new(Blogger)

	blg.Name = name
	blg.Weibos = make([]*PostContent, 0)
	blg.Comments = make(map[int][]*Comment)

	return blg
}

func (b *Blogger) PostWeibo(content string, wbType int) {
	weibo := new(PostContent)
	weibo.Content = content
	weibo.Id = b.GetID()
	weibo.PostMan = b.Name
	weibo.PostTime = time.Now()
	weibo.Type = wbType

	b.Weibos = append(b.Weibos, weibo)
	b.Notify(weibo.Id)

}

func (b *Blogger) GetID() int {
	if len(b.Weibos) == 0 {
		return 0
	} else {
		return b.Weibos[len(b.Weibos)-1].Id + 1
	}
}

func (b *Blogger) Notify(wbid int) {
	for _, fan := range b.Fans {
		fan.Update(b, wbid)
	}
}

func (b *Blogger) Attach(bFans FansInterface) {
	b.Fans = append(b.Fans, bFans)
}

func (b *Blogger) Detach(bFans FansInterface) {
	for i := 0; i < len(b.Fans); i++ {
		if b.Fans[i] == bFans {
			b.Fans = append(b.Fans[:i], b.Fans[i+1:]...)
		}
	}
}

func (b *Blogger) GetWeibo(wbid int) *PostContent {
	for _, blog := range b.Weibos {
		if blog.Id == wbid {
			return blog
		}
	}

	return nil
}

func (b *Blogger) AddComment(comment Comment, wbid int) {
	b.Comments[wbid] = append(b.Comments[wbid], &comment)
}

func (b *Blogger) ShowComments(wbid int) {
	blog := b.GetWeibo(wbid)
	fmt.Println("博主名称:", b.Name)
	fmt.Println("微博内容:", blog.Content)

	for _, msg := range b.Comments[wbid] {
		fmt.Println("粉丝名称:", msg.PostMan)
		fmt.Println("评论内容:", msg.Content)
	}
}

type Fans struct {
	Base
}

type GoodFans struct {
	Fans
}

type BadFans struct {
	Fans
}

type FansInterface interface {
	Update(b BloggerInterface, wbid int)
	Action(b BloggerInterface, wbid int)
}

func (f *GoodFans) Update(b BloggerInterface, wbid int) {
	fmt.Printf("Hello:%s你所关注的博主发布了一个新的微博\n", f.Name)
	f.Action(b, wbid)
}

func (f *GoodFans) Action(b BloggerInterface, wbid int) {
	blogger, ok := b.(*Blogger)
	if ok {
		weibo := blogger.GetWeibo(wbid)

		cType := weibo.Type
		message := ""

		switch cType {
		case 1:
			message = "非常好啊!!"
		case 2:
			message = "加油"
		default:
			message = "未知心情"
		}

		comment := Comment{PostContent{0, message, time.Now(), cType, f.Name}, blogger.Name, time.Now()}
		blogger.AddComment(comment, wbid)

		blogger.ShowComments(wbid)
	}
}

func main() {
	blg := NewBlogger("张三")
	friedFans := new(GoodFans)
	friedFans.ID = 1
	friedFans.Name = "李四"
	blg.Attach(friedFans) // 添加粉丝

	blg.PostWeibo("今天天气很好", 1)
}
