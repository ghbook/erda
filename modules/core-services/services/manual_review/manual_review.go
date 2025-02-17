// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package manual_review

import (
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/core-services/dao"
	"github.com/erda-project/erda/modules/core-services/model"
)

// ManualReview 成员操作封装
type ManualReview struct {
	db *dao.DBClient
}

// Option 定义 ManualReview 对象配置选项
type Option func(*ManualReview)

// New 新建 ManualReview 实例
func New(options ...Option) *ManualReview {
	mem := &ManualReview{}
	for _, op := range options {
		op(mem)
	}
	return mem
}

// WithDBClient 配置 db client
func WithDBClient(db *dao.DBClient) Option {
	return func(m *ManualReview) {
		m.db = db
	}
}

func (m *ManualReview) CreateReviewUser(param *apistructs.CreateReviewUser) (error, int64) {
	var user = &model.ReviewUser{
		Operator: param.Operator,
		OrgId:    param.OrgId,
		TaskId:   param.TaskId,
	}
	err := m.db.CreateReviewUser(user)
	return err, user.TaskId
}

func (m *ManualReview) GetReviewByTaskId(param *apistructs.GetReviewByTaskIdIdRequest) (apistructs.GetReviewByTaskIdIdResponse, error) {
	return m.db.GetReviewByTaskId(param)
}

func (m *ManualReview) CreateReview(param *apistructs.CreateReviewRequest) (error, int64) {
	var review = &model.ManualReview{
		BuildId:         param.BuildId,
		ProjectId:       param.ProjectId,
		ApplicationId:   param.ApplicationId,
		ApplicationName: param.ApplicationName,
		SponsorId:       param.SponsorId,
		CommitID:        param.CommitID,
		OrgId:           param.OrgId,
		TaskId:          param.TaskId,
		ProjectName:     param.ProjectName,
		BranchName:      param.BranchName,
		ApprovalStatus:  param.ApprovalStatus,
	}
	err := m.db.CreateReview(review)
	return err, review.ID
}

func (m *ManualReview) GetReviewsByUserId(param *apistructs.GetReviewsByUserIdRequest) (int, []apistructs.GetReviewsByUserIdResponse, error) {
	tasks, err := m.db.GetTaskIDByOperator(param)
	if err != nil {
		return 0, nil, err
	}
	total, reviews, err := m.db.GetReviewsByUserId(param, tasks)
	var list []apistructs.GetReviewsByUserIdResponse

	for _, v := range reviews {
		list = append(list, apistructs.GetReviewsByUserIdResponse{
			Id:              v.ID,
			BuildId:         v.BuildId,
			BranchName:      v.BranchName,
			ProjectName:     v.ProjectName,
			ProjectId:       v.ProjectId,
			ApplicationId:   v.ApplicationId,
			ApplicationName: v.ApplicationName,
			Operator:        v.SponsorId,
			CommitMessage:   v.CommitMessage,
			CommitId:        v.CommitID,
			ApprovalStatus:  v.ApprovalStatus,
			ApprovalReason:  v.ApprovalReason,
			ApprovalContent: "pipeline",
		})
	}

	return total, list, nil
}
func (m *ManualReview) GetAuthorityByUserId(param *apistructs.GetAuthorityByUserIdRequest) (apistructs.GetAuthorityByUserIdResponse, error) {
	return m.db.GetAuthorityByUserId(param)
}

func (m *ManualReview) GetReviewsBySponsorId(param *apistructs.GetReviewsBySponsorIdRequest) (int, []apistructs.GetReviewsBySponsorIdResponse, error) {
	total, manualReviews, err := m.db.GetReviewsBySponsorId(param)
	if err != nil {
		return 0, nil, err
	}
	var reviews []int
	for _, v := range manualReviews {
		reviews = append(reviews, v.TaskId)
	}
	reviewusers, _ := m.db.GetOperatorByTaskID(reviews)

	operatorMap := make(map[int][]string)

	for _, v := range reviewusers {
		operatorMap[int(v.TaskId)] = append(operatorMap[int(v.TaskId)], v.Operator)
	}
	var list []apistructs.GetReviewsBySponsorIdResponse
	for _, v := range manualReviews {

		list = append(list, apistructs.GetReviewsBySponsorIdResponse{
			Id:              v.ID,
			BuildId:         v.BuildId,
			BranchName:      v.BranchName,
			ProjectName:     v.ProjectName,
			ProjectId:       v.ProjectId,
			ApplicationName: v.ApplicationName,
			ApplicationId:   v.ApplicationId,
			CommitMessage:   v.CommitMessage,
			CommitId:        v.CommitID,
			Approver:        operatorMap[v.TaskId],
			ApprovalReason:  v.ApprovalReason,
			ApprovalContent: "pipeline",
		})
	}
	return total, list, nil
}

func (m *ManualReview) UpdateApproval(param *apistructs.UpdateApproval) error {
	return m.db.UpdateApproval(param)
}
