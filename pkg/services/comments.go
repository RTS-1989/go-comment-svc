package services

import (
	"context"
	"github.com/RTS-1989/go-comment-svc/pkg/db"
	"github.com/RTS-1989/go-comment-svc/pkg/models"
	"github.com/RTS-1989/go-comment-svc/pkg/pb"
	"net/http"
)

type Server struct {
	H db.Handler
}

func NewsPbComment(c *models.CommentUnion) *pb.Comment {
	return &pb.Comment{
		Id:       c.ID,
		Text:     c.Text,
		Censored: c.Censored,
		ParentId: c.ParentId,
		Children: make([]*pb.Comment, 0),
	}
}

func CommentsPb(rawComments []*models.CommentUnion) []*pb.Comment {
	commentsMap := make(map[uint64]*pb.Comment)
	comments := make([]*pb.Comment, 0)

	for _, c := range rawComments {
		// Проверяем вносилась ли запись ранее в мапу, если да, то заполняем ее новыми данными.
		val, ok := commentsMap[c.ID]
		if ok && len(val.Children) > 0 {
			val.Text = c.Text
			val.Censored = c.Censored
			val.ParentId = c.ParentId
		}

		if c.ParentId != 0 {
			v, check := commentsMap[c.ParentId]
			if check {
				// Здесь проверяем есть ли запись в мапе, так как, если она вносилась ранее, то имеет вложенные комментарии
				comment, checkCurrent := commentsMap[c.ID]
				if checkCurrent {
					v.Children = append(v.Children, comment)
					commentsMap[c.ParentId] = v
					continue
				}
				// В случае, если записей нет в мапе создаем новую структуре коммментария
				commentPb := NewsPbComment(c)
				v.Children = append(v.Children, commentPb)
				continue
			}
			// Если родительской записи не существует, то ее создаем
			parentComment := &pb.Comment{Id: c.ParentId}
			// Здесь проверяем есть ли запись в мапе, так как, если она вносилась ранее, то имеет вложенные комментарии
			comment, checkCurrent := commentsMap[c.ID]
			if checkCurrent {
				parentComment.Children = append(parentComment.Children, comment)
				commentsMap[c.ParentId] = parentComment
				continue
			}
			// В случае, если записей нет в мапе создаем новую структуре коммментария
			commentPb := NewsPbComment(c)
			parentComment.Children = append(parentComment.Children, commentPb)
			commentsMap[c.ParentId] = parentComment
			continue
		}
		// Здесь создаем комменты, которые не имеют родительских комментов, к ним прикрепляем значение 1-й детской записи
		// коммента, которая уже содержит ранее вложенные комментарии
		root, ok := commentsMap[c.ID]
		if ok {
			comments = append(comments, root)
			continue
		}
		root = NewsPbComment(c)
		comments = append(comments, root)
	}

	return comments
}

func (s *Server) CreateChild(parentId, childId, newsId uint64) error {
	child := &models.CommentChild{
		ParentId: parentId,
		ChildId:  childId,
		NewsId:   newsId,
	}
	if result := s.H.DB.Create(child); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {

	c := &models.Comment{
		NewsId:   req.NewsId,
		Text:     req.Text,
		Censored: req.Censored,
	}

	if result := s.H.DB.Create(c); result.Error != nil {
		return &pb.CreateCommentResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}

	s.CreateChildEntry(req.ParentId, c.ID, c.NewsId)

	return &pb.CreateCommentResponse{
		Status: http.StatusCreated,
		Id:     c.ID,
	}, nil
}

func (s *Server) CreateChildEntry(parentId, childId, newsId uint64) (*pb.CreateCommentResponse, error) {
	err := s.CreateChild(parentId, childId, newsId)
	if err != nil {
		return &pb.CreateCommentResponse{
			Status: http.StatusConflict,
			Error:  err.Error(),
		}, nil
	}
	return &pb.CreateCommentResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) CommentsByNews(ctx context.Context, req *pb.CommentsByNewsRequest) (*pb.CommentsByNewsResponse, error) {
	rawComments := make([]*models.CommentUnion, 0)
	query := `SELECT c.id as id, text, censored, cc.parent_id as parent_id, 
       cc.child_id as child_id, c.news_id as news_id, c.created_at as created_at  FROM comments c
           INNER JOIN comment_children cc on c.id = cc.child_id 
       WHERE c.news_id = $1
       ORDER BY c.created_at DESC, cc.parent_id DESC;`
	if result := s.H.DB.Raw(query, req.NewsId).Find(&rawComments); result.Error != nil {
		return &pb.CommentsByNewsResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	commentsSlice := CommentsPb(rawComments)

	return &pb.CommentsByNewsResponse{
		Status:   http.StatusOK,
		Comments: commentsSlice,
	}, nil
}
