package mysql

import (
	"github.com/jmoiron/sqlx"
	"strings"
	"web_app/models"
)

// CreatePost 创建帖子数据
func CreatePost(p *models.Post) (err error) {
	//执行SQL语句入库
	sqlstr := `insert into post(
                 post_id,
                 author_id,
                 community_id,
                 title,
                 content) 
                 values (?,?,?,?,?)`
	_, err = db.Exec(sqlstr, p.ID, p.AuthorID, p.CommunityID, p.Title, p.Content)
	return
}

// GetPostById 根据post id查询 单个贴子详情数据
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlstr := `select
              post_id,title,content,author_id,community_id,create_time
              from post
              where post_id = ?
               `
	db.Get(post, sqlstr, pid)
	return post, err
}

// GetPostList()根据page、size 获取帖子 列表list[]
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select 
	post_id, title, content, author_id, community_id, create_time
	from post
	ORDER BY create_time
	DESC
	limit  ?,?
	`
	//第一个？ 代表从哪里开始取出，第二个？代表取多少条
	posts = make([]*models.Post, 0, 2) // 不要写成make([]*models.Post, 2),长度为0,容量为2
	db.Select(&posts, sqlStr, (page-1)*size, size)
	return posts, err
}

// 根据给定的id列表查询帖子数据
// sqlx.In的查询示例，在sqlx查询语句中实现In查询和FIND_IN_SET函数。 参考链接：https://www.liwenzhou.com/posts/Go/sqlx/

func GetPostListByIDs(ids []string) (posts []*models.Post, err error) {
	sqlStr := `select 
    post_id, title, content, author_id, community_id, create_time
	from post
   where post_id in (?)
   order by  FIND_IN_SET(post_id,?)
`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...)
	return
}
