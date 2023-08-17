package logic

import (
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"
)

// 1. 创建帖子 logic函数
func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()
	//前面需要 获得当前用户的ID，即是AuthorID
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	//存储信息到 redis 创建redis 帖子的两个keytime\keyscore
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
	//3. 返回
}

// 2.GetPostById 根据帖子id, 查询帖子详情数据data,//从返回的数据data 入手，新增改造一个 ApiPostDetail结构体 帖子详情接口的结构体
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询  并拼接组合我们接口想用的数据

	//第一步: 查询帖子id 查询post
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid), //同时记录到pid
			zap.Error(err))        //记录err
		return
	}
	// 第二步：根据作者id 查询作者的信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	//第三步：根据社区id 获取社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	//logic 业务逻辑处理

	//接口数据拼接data
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	//data是指针类型，没有初始化直接使用，经典错误
	//data.AuthorName = user.Username 这个形式存在
	//data.CommunityDetail = community
	return
}

// 3.GetPostList() 获取帖子列表 详情 ApiPostDetail []列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed",
			zap.Error(err))
		return nil, err
	}
	//先对data进行一个初始化
	data = make([]*models.ApiPostDetail, 0, len(posts))

	//找到帖子之后，需要查找每一个作者信息，社区信息(上面的接口是根据帖子id来获取作者、社区详情信息)，这里只需要搞一个循环posts 列表即可！
	for _, post := range posts {
		// 第一步：根据作者id 查询作者的信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//第二步：根据社区id 获取社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		PostDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		//将个作者信息，社区信息，PostDetail追加到 data数据返回。
		data = append(data, PostDetail)
	}

	return

}

// 4 .GetPostList2() 获取帖子列表详情 升级版
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//关键点在于 2.去redis查询id列表
	ids, err := redis.GetPostIDInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDInOrder(p)  return 0 failed")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))

	//3.根据id去mysql数据库查询帖子的详情信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs() failed",
			zap.Error(err))
		return nil, err
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))

	// 提前 查询到每篇帖子获得的分数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	//先对data进行一个初始化
	data = make([]*models.ApiPostDetail, 0, len(posts))
	//找到帖子之后，需要查找每一个作者信息，社区信息 填充到帖子中
	//idx 拿到索引
	for idx, post := range posts {
		// 第一步：根据作者id 查询作者的信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//第二步：根据社区id 获取社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		PostDetail := &models.ApiPostDetail{
			VoteNum:         voteData[idx],
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		//将个作者信息，社区信息，PostDetail追加到 data数据返回。
		data = append(data, PostDetail)
	}
	return

}

// 5. 根据社区
func GetCommunityPostList(p *models.ParamCommunityPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetCommunityPostIDsInOrder", zap.Any("ids", ids))
	// 3. 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
